package users

const (
	searchUserQuery = `select :projections from users :filters order by id`

	insertUser = `
		insert into users 
			(id, first_name, last_name, birth_date, created_at, updated_at)
		values 
			(?, ?, ?, ?, ?, ?)
	`

	updateUser = `
		update users
			set
				:updates
		where id = ?
	`

	createTableUser = `
		create table if not exists users (
			id varchar(60) not null primary key,
			first_name varchar(100) not null,
			last_name varchar(500) null,
			birth_date timestamp not null,
			created_at timestamp not null,
			updated_at timestamp not null,
			deleted_at timestamp null
		);
	`
)
