package api

const (
	CREATE_ITEM string = `
		INSERT INTO items (
			id,
			order_id,
			title,
			price,
			quantity
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING id, created
	`

	CREATE_ORDER string = `
		INSERT INTO orders (
			id,
			status,
			total_price,
			total_quantity
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING id, created
	`

	SELECT_ORDER_BY_ID string = `
		SELECT 
			id,
			status,
			total_price,
			total_quantity,
			created,
			updated
		FROM
			orders
		WHERE
			id = $1
	`
)
