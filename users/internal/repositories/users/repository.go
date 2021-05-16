package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/users/internal/repositories/filters"
	"github.com/dalmarcogd/mobstore/users/internal/repositories/projections"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

type userRepository struct {
	services.NoopHealth
	serviceManager services.ServiceManager
	ctx            context.Context
}

func New() *userRepository {
	return &userRepository{}
}

func (r *userRepository) ServiceManager() services.ServiceManager {
	return r.serviceManager
}

func (r *userRepository) WithServiceManager(s services.ServiceManager) services.UsersRepository {
	r.serviceManager = s
	return r
}

func (r *userRepository) Init(ctx context.Context) error {
	r.ctx = ctx
	return nil
}

func (r *userRepository) Close() error {
	return nil
}

func (r *userRepository) Search(ctx context.Context, search domains.UserSearch) ([]domains.User, error) {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := r.ServiceManager().Validator().Validate(ctx, search); err != nil {
		span.Error(err)
		return nil, err
	}

	gotProjections := projections.GetProjections(search.Projection)
	query := strings.ReplaceAll(searchUserQuery, ":projections", strings.Join(gotProjections, ", "))

	gotFilters := filters.GetFilters(search.Filter)

	var gotFilterStr strings.Builder
	var argsFilter []interface{}
	gotFilterStr.WriteString("where deleted_at is null ")
	if len(gotFilters) > 0 {
		for field, value := range gotFilters {
			gotFilterStr.WriteString(fmt.Sprintf("and %v = ? ", field))
			argsFilter = append(argsFilter, value)
		}
	}
	query = strings.ReplaceAll(query, ":filters", gotFilterStr.String())

	users := make([]domains.User, 0)
	err := r.ServiceManager().UserDatabase().TransactionReplica(ctx, func(tx services.DatabaseTransaction) error {
		rows, err := tx.Get(query, argsFilter...)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error select users error=%v", err))
			return err
		}

		defer func() {
			if err := rows.Close(); err != nil {
				r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error closing rows err=%v", err))
			}
		}()

		cols, err := rows.Columns()
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error select users error=%v", err))
			return err
		}

		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = new(sql.RawBytes)
			}

			if err := rows.Scan(columnPointers...); err != nil {
				r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error scanning error=%v", err))
				return err
			}

			m := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*sql.RawBytes)
				m[colName] = *val
			}
			user := new(domains.User)
			err := projections.SetProjections(user, m)
			if err != nil {
				r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error assign projects to User error=%v", err))
				return err
			}

			users = append(users, *user)
		}
		return nil
	})
	if err != nil {
		span.Error(err)
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Create(ctx context.Context, user *domains.User) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := r.ServiceManager().Validator().Validate(ctx, user); err != nil {
		span.Error(err)
		return err
	}

	user.Id = ptrs.String(uuid.NewString())
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	r.ServiceManager().Logger().Info(ctx, fmt.Sprintf("Executing insert on users for userId=%v", user.Id))

	err := r.ServiceManager().UserDatabase().TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		_, err := tx.Insert(insertUser,
			user.Id,
			user.FirstName,
			user.LastName,
			user.BirthDate,
			user.CreatedAt,
			user.UpdatedAt,
		)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error trying to insert userId=%v on users err=%s", user.Id, err))
			return err
		}
		return nil
	})
	if err != nil {
		span.Error(err)
		return err
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domains.User) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := r.ServiceManager().Validator().Validate(ctx, user); err != nil {
		span.Error(err)
		return err
	}
	if ptrs.StringValue(user.Id) == "" {
		span.Error(errors.RepositoryIdIsRequiredError)
		return errors.RepositoryIdIsRequiredError
	}

	user.UpdatedAt = time.Now().UTC()

	projection := domains.UserProjection{
		FirstName: user.FirstName != nil,
		LastName:  user.LastName != nil,
		BirthDate: user.BirthDate != nil,
		DeletedAt: user.DeletedAt != nil,
	}
	getProjections := projections.GetProjections(projection)
	if len(getProjections) == 0 {
		span.Error(errors.RepositoryProjectionsIsRequiredError)
		return errors.RepositoryProjectionsIsRequiredError
	}
	var queryUpdates []string
	var queryValues []interface{}
	for _, proj := range getProjections {
		if val := projections.GetProjectionValue(*user, proj); val != nil {
			queryUpdates = append(queryUpdates, fmt.Sprintf("%v = ?", proj))
			queryValues = append(queryValues, val)
		}
	}
	queryUpdates = append(queryUpdates, "updated_at = ?")
	queryValues = append(queryValues, user.UpdatedAt, user.Id)
	query := strings.ReplaceAll(updateUser, ":updates", strings.Join(queryUpdates, ","))

	err := r.ServiceManager().UserDatabase().TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		_, err := tx.Update(query, queryValues...)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error trying to update userId=%v on users err=%s", user.Id, err))
			return err
		}
		return nil
	})
	if err != nil {
		span.Error(err)
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, user *domains.User) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	utc := time.Now().UTC()
	user.DeletedAt = &utc
	return r.Update(ctx, user)
}

func (r *userRepository) Migrate(ctx context.Context) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	err := r.ServiceManager().UserDatabase().TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		_, err := tx.Migrate(createTableUser)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error trying to migrate on users err=%s", err))
			return err
		}
		return nil
	})
	if err != nil {
		span.Error(err)
		return err
	}

	return nil
}
