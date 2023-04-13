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

const listCmd = `SELECT id, title, description, media_source, author_id 
					FROM pins 
					ORDER BY created_at DESC 
					LIMIT $1 OFFSET $2;`

const fullUpdateCmd = `UPDATE pins
						SET title = $1::VARCHAR,
						description = $2::TEXT
						WHERE id = $3
						RETURNING id, title, description, media_source, author_id;`

const deleteCmd = `DELETE FROM pins 
					WHERE id = $1;`

const checkWriteCmd = `SELECT EXISTS(SELECT id
						FROM pins
					 	WHERE id = $1 AND author_id = $2);`
