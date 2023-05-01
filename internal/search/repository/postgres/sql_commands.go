package postgres

const getPinsCmd = `SELECT id, title, description, media_source, n_likes, author_id FROM pins
                        WHERE to_tsquery($1) @@ to_tsvector(pins.title || pins.description);`

const getBoardsCmd = `SELECT * FROM boards
                        WHERE to_tsquery($1) @@ to_tsvector(boards.name);`

const getUsersCmd = `SELECT id, username, name, profile_image, website_url 
					    FROM users
                        WHERE to_tsquery($1) @@ to_tsvector(users.username)`
