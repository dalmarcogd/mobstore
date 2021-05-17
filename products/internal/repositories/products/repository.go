package products

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/products/internal/repositories/filters"
	"github.com/dalmarcogd/mobstore/products/internal/repositories/projections"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

type productRepository struct {
	services.NoopHealth
	serviceManager services.ServiceManager
	ctx            context.Context
}

func New() *productRepository {
	return &productRepository{}
}

func (r *productRepository) ServiceManager() services.ServiceManager {
	return r.serviceManager
}

func (r *productRepository) WithServiceManager(s services.ServiceManager) services.ProductsRepository {
	r.serviceManager = s
	return r
}

func (r *productRepository) Init(ctx context.Context) error {
	r.ctx = ctx
	return nil
}

func (r *productRepository) Close() error {
	return nil
}

func (r *productRepository) Search(ctx context.Context, search domains.ProductSearch) ([]domains.Product, error) {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := r.ServiceManager().Validator().Validate(ctx, search); err != nil {
		span.Error(err)
		return nil, err
	}

	gotProjections := projections.GetProjections(search.Projection)
	query := strings.ReplaceAll(searchProductQuery, ":projections", strings.Join(gotProjections, ", "))

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

	products := make([]domains.Product, 0)
	err := r.ServiceManager().ProductDatabase().TransactionReplica(ctx, func(tx services.DatabaseTransaction) error {
		rows, err := tx.Get(query, argsFilter...)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error select products error=%v", err))
			return err
		}

		defer func() {
			if err := rows.Close(); err != nil {
				r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error closing rows err=%v", err))
			}
		}()

		cols, err := rows.Columns()
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error select products error=%v", err))
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
			product := new(domains.Product)
			err := projections.SetProjections(product, m)
			if err != nil {
				r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error assign projects to Product error=%v", err))
				return err
			}

			products = append(products, *product)
		}
		return nil
	})
	if err != nil {
		span.Error(err)
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Create(ctx context.Context, product *domains.Product) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := r.ServiceManager().Validator().Validate(ctx, product); err != nil {
		span.Error(err)
		return err
	}

	product.Id = ptrs.String(uuid.NewString())
	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()
	r.ServiceManager().Logger().Info(ctx, fmt.Sprintf("Executing insert on products for productId=%v", product.Id))

	err := r.ServiceManager().ProductDatabase().TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		_, err := tx.Insert(insertProduct,
			product.Id,
			product.PriceInCents,
			product.Title,
			product.Description,
			product.CreatedAt,
			product.UpdatedAt,
		)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error trying to insert productId=%v on products err=%s", product.Id, err))
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

func (r *productRepository) Update(ctx context.Context, product *domains.Product) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := r.ServiceManager().Validator().Validate(ctx, product); err != nil {
		span.Error(err)
		return err
	}
	if ptrs.StringValue(product.Id) == "" {
		span.Error(errors.RepositoryIdIsRequiredError)
		return errors.RepositoryIdIsRequiredError
	}

	product.UpdatedAt = time.Now().UTC()

	projection := domains.ProductProjection{
		PriceInCents: product.PriceInCents != nil,
		Title:        product.Title != nil,
		Description:  product.Description != nil,
		DeletedAt:    product.DeletedAt != nil,
	}
	getProjections := projections.GetProjections(projection)
	if len(getProjections) == 0 {
		span.Error(errors.RepositoryProjectionsIsRequiredError)
		return errors.RepositoryProjectionsIsRequiredError
	}
	var queryUpdates []string
	var queryValues []interface{}
	for _, proj := range getProjections {
		if val := projections.GetProjectionValue(*product, proj); val != nil {
			queryUpdates = append(queryUpdates, fmt.Sprintf("%v = ?", proj))
			queryValues = append(queryValues, val)
		}
	}
	queryUpdates = append(queryUpdates, "updated_at = ?")
	queryValues = append(queryValues, product.UpdatedAt, product.Id)
	query := strings.ReplaceAll(updateProduct, ":updates", strings.Join(queryUpdates, ","))

	err := r.ServiceManager().ProductDatabase().TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		_, err := tx.Update(query, queryValues...)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error trying to update productId=%v on products err=%s", product.Id, err))
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

func (r *productRepository) Delete(ctx context.Context, product *domains.Product) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	utc := time.Now().UTC()
	product.DeletedAt = &utc
	return r.Update(ctx, product)
}

func (r *productRepository) Migrate(ctx context.Context) error {
	ctx, span := r.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	err := r.ServiceManager().ProductDatabase().TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		_, err := tx.Migrate(createTableProduct)
		if err != nil {
			r.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error trying to migrate on products err=%s", err))
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
