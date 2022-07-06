package finder

import (
	"context"
	"fmt"
	"strings"

	"github.com/trixky/hypertube/api-search/databases"
	"github.com/trixky/hypertube/api-search/sqlc"
)

const PerPage int32 = 20

// ? The query result handle code is the same as the sqlc generated code
// ? Only the queries are different and are updated dynamically to handle different ORDER BY and genres IN condition

const countMedias = `-- name: FindMedias :many
SELECT count(DISTINCT medias.id)
FROM medias
WHERE
	(
		CASE WHEN $1::bool
		THEN
			EXISTS (
				SELECT id FROM media_names
				WHERE media_names.media_id = medias.id AND media_names.name ILIKE '%' || $2 || '%'
				LIMIT 1
			) OR (
				medias.imdb_id ILIKE '%' || $3 || '%'
			) OR (
				medias.tmdb_id::varchar ILIKE '%' || $4 || '%'
			) OR (
				medias.description ILIKE '%' || $5 || '%'
			)
		ELSE true
		END
	)
	AND (
		CASE WHEN $6::bool
		THEN medias.rating >= $7
		ELSE true
		END
	) AND (
		CASE WHEN $8::bool
		THEN medias.year = $9
		ELSE true
		END
	) AND (
		CASE WHEN $10::bool
		THEN EXISTS (
			SELECT id FROM media_genres
			WHERE media_genres.media_id = medias.id AND media_genres.genre_id IN ({{genres}})
			LIMIT 1
		)
		ELSE true
		END
	)
`

const findMedias = `-- name: FindMedias :many
SELECT DISTINCT medias.id, medias.imdb_id, medias.tmdb_id, medias.description, medias.duration, medias.thumbnail, medias.background, medias.year, medias.rating, media_names.name
FROM medias
RIGHT JOIN (
	SELECT media_names.id, media_names.media_id, media_names.name FROM media_names
	WHERE lang = '__'
	ORDER BY name {{sort_order}}
) media_names ON media_names.media_id = medias.id
WHERE
	(
		CASE WHEN $2::bool
		THEN
			EXISTS (
				SELECT id FROM media_names
				WHERE (media_names.media_id = medias.id AND media_names.name ILIKE '%' || $3 || '%')
				LIMIT 1
			) OR (
				medias.imdb_id ILIKE '%' || $4 || '%'
			) OR (
				medias.tmdb_id::varchar ILIKE '%' || $5 || '%'
			) OR (
				medias.description ILIKE '%' || $6 || '%'
			)
		ELSE true
		END
	)
	AND (
		CASE WHEN $7::bool
		THEN medias.rating >= $8
		ELSE true
		END
	) AND (
		CASE WHEN $9::bool
		THEN medias.year = $10
		ELSE true
		END
	) AND (
		CASE WHEN $11::bool
		THEN EXISTS (
			SELECT id FROM media_genres
			WHERE media_genres.media_id = medias.id AND media_genres.genre_id IN ({{genres}})
			LIMIT 1
		)
		ELSE true
		END
	)
ORDER BY {{sort_column}} {{sort_order}}
LIMIT {{per_page}} OFFSET $1
`

type FindMediasParams struct {
	Offset       int32
	SearchQuery  bool
	Query        string
	SearchGenres bool
	Genres       []int32
	SearchRating bool
	Rating       float64
	SearchYear   bool
	Year         int32
	SortColumn   string
	SortOrder    string
}

func GenerateQuery(mode string, arg *FindMediasParams) (string, []interface{}) {
	// Dynamically update the query
	base := countMedias
	if mode == "find" {
		base = findMedias
	}
	query := strings.Replace(base, "{{sort_order}}", arg.SortOrder, 2)
	if arg.SortColumn == "name" {
		query = strings.Replace(query, "{{sort_column}}", "media_names.name", 1)
	} else {
		query = strings.Replace(query, "{{sort_column}}", arg.SortColumn, 1)
	}
	query = strings.Replace(query, "{{per_page}}", fmt.Sprint(PerPage), 1)
	args := []interface{}{
		arg.SearchQuery,
		arg.Query,
		arg.Query,
		arg.Query,
		arg.Query,
		arg.SearchRating,
		arg.Rating,
		arg.SearchYear,
		arg.Year,
		arg.SearchGenres,
	}
	if mode == "find" {
		args = append([]interface{}{arg.Offset}, args...)
	}

	// Update the genres IN condition by generating args at the end
	start := len(args) + 1
	genres := ""
	if !arg.SearchGenres || len(arg.Genres) == 0 {
		arg.SearchGenres = false
		genres = fmt.Sprintf("$%v", start)
		args = append(args, 0)
	} else {
		for _, genre_id := range arg.Genres {
			if genres == "" {
				genres = fmt.Sprintf("$%v", start)
			} else {
				genres = genres + fmt.Sprintf(",$%v", start)
			}
			start += 1
			args = append(args, genre_id)
		}
	}
	query = strings.Replace(query, "{{genres}}", genres, 1)

	return query, args
}

func CountMedias(ctx context.Context, arg FindMediasParams) (int64, error) {
	db := databases.DBs.SqlDatabase

	// Dynamically update the query
	query, args := GenerateQuery("count", &arg)

	row := db.QueryRowContext(ctx, query, args...)
	var count int64
	err := row.Scan(&count)

	return count, err
}

func FindMedias(ctx context.Context, arg FindMediasParams) ([]sqlc.Media, error) {
	db := databases.DBs.SqlDatabase

	// Dynamically update the query
	query, args := GenerateQuery("find", &arg)

	// Send the query and handle the result
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sqlc.Media
	for rows.Next() {
		var i sqlc.Media
		var name string
		if err := rows.Scan(
			&i.ID,
			&i.ImdbID,
			&i.TmdbID,
			&i.Description,
			&i.Duration,
			&i.Thumbnail,
			&i.Background,
			&i.Year,
			&i.Rating,
			&name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
