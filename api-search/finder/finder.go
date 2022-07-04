package finder

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/trixky/hypertube/api-search/postgres"
	"github.com/trixky/hypertube/api-search/sqlc"
)

// ? Imported from sqlc generated code to handle different ORDER BY and genres IN condition

const countMedias = `-- name: FindMedias :many
SELECT count(DISTINCT medias.id)
FROM medias
WHERE
	(
		CASE WHEN $1::bool
		THEN EXISTS (
			SELECT id FROM media_names
			WHERE media_names.media_id = medias.id AND media_names.name ILIKE '%' || $2 || '%'
			LIMIT 1
		)
		ELSE true
		END
	)
	AND (
		CASE WHEN $3::bool
		THEN medias.rating >= $4
		ELSE true
		END
	) AND (
		CASE WHEN $5::bool
		THEN medias.year = $6
		ELSE true
		END
	) AND (
		CASE WHEN $7::bool
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
SELECT DISTINCT medias.id, medias.imdb_id, medias.tmdb_id, medias.description, medias.duration, medias.thumbnail, medias.background, medias.year, medias.rating
FROM medias
WHERE
	(
		CASE WHEN $2::bool
		THEN EXISTS (
			SELECT id FROM media_names
			WHERE media_names.media_id = medias.id AND media_names.name ILIKE '%' || $3 || '%'
			LIMIT 1
		)
		ELSE true
		END
	)
	AND (
		CASE WHEN $4::bool
		THEN medias.rating >= $5
		ELSE true
		END
	) AND (
		CASE WHEN $6::bool
		THEN medias.year = $7
		ELSE true
		END
	) AND (
		CASE WHEN $8::bool
		THEN EXISTS (
			SELECT id FROM media_genres
			WHERE media_genres.media_id = medias.id AND media_genres.genre_id IN ({{genres}})
			LIMIT 1
		)
		ELSE true
		END
	)
ORDER BY {{sort_column}} {{sort_order}}
LIMIT 5 OFFSET $1
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
	query := strings.Replace(base, "{{sort_column}}", arg.SortColumn, 1)
	query = strings.Replace(query, "{{sort_order}}", arg.SortOrder, 1)
	args := []interface{}{
		arg.SearchQuery,
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
	db := postgres.DB.SqlDatabase

	// Dynamically update the query
	query, args := GenerateQuery("count", &arg)
	log.Println("query", query, args)

	row := db.QueryRowContext(ctx, query, args...)
	var count int64
	err := row.Scan(&count)

	return count, err
}

func FindMedias(ctx context.Context, arg FindMediasParams) ([]sqlc.Media, error) {
	db := postgres.DB.SqlDatabase

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
