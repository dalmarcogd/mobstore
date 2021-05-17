package products

const (
	searchProductQuery = `select :projections from products :filters order by id`

	insertProduct = `
		insert into products 
			(id, price_in_cents, title, description, created_at, updated_at)
		values 
			(?, ?, ?, ?, ?, ?)
	`

	updateProduct = `
		update products
			set
				:updates
		where id = ?
	`

	createTableProduct = `
		create table if not exists products (
			id varchar(60) not null primary key,
			title varchar(100) not null,
			description varchar(500) null,
			price_in_cents bigint not null,
			created_at timestamp not null,
			updated_at timestamp not null,
			deleted_at timestamp null
		);
	`
)
