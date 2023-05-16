package postgres

const getCmd = `SELECT id, title, description, media_source, n_likes, author_id
				FROM pins
				WHERE id = $1;`

const listByUserCmd = `SELECT id, title, description, media_source, n_likes, author_id
						FROM pins 
						WHERE author_id = $1
						ORDER BY created_at DESC 
						LIMIT $2 OFFSET $3;`

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
