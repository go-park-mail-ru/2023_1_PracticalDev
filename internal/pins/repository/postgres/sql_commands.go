package postgres

const createCmd = `INSERT INTO pins (title, media_source, description, author_id)
				   VALUES ($1, $2, $3, $4)
				   RETURNING id, title, media_source, description, author_id;`

const getCmd = `SELECT id, title, description, media_source, author_id
				FROM pins
				WHERE id = $1;`

const listByUserCmd = `SELECT id, title, description, media_source, author_id
						FROM pins 
						WHERE author_id = $1
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

const listByBoardCmd = `SELECT pins.id, title, description, media_source, author_id 
						FROM pins 
						JOIN boards_pins AS b
						ON b.board_id = $1 AND b.pin_id = pins.id
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

const listCmd = `SELECT id, title, description, media_source, author_id 
					FROM pins 
					ORDER BY created_at DESC 
					LIMIT $1 OFFSET $2;`

const fullUpdateCmd = `UPDATE pins
						SET title = $1::VARCHAR,
						description = $2::TEXT,
						media_source = $3::TEXT
						WHERE id = $4
						RETURNING id, title, description, media_source, author_id;`

const deleteCmd = `DELETE FROM pins 
					WHERE id = $1;`

const addToBoardCmd = `INSERT INTO boards_pins(pin_id, board_id)
						VALUES($1, $2);`

const deleteFromBoardCmd = `DELETE FROM boards_pins
							WHERE pin_id = $1 AND board_id = $2;`

const checkWriteCmd = `SELECT EXISTS(SELECT id
						FROM pins
					 	WHERE id = $1 AND author_id = $2);`
