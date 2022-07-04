package finder

import (
	"context"
	"fmt"
	"strings"

	"github.com/trixky/hypertube/api-search/postgres"
	"github.com/trixky/hypertube/api-search/sqlc"
)

// ? Imported from sqlc generated code to handle different ORDER BY and genres IN condition

const findMedias = `-- name: FindMedias :many
SELECT DISTINCT medias.id, medias.imdb_id, medias.tmdb_id, medias.description, medias.duration, medias.thumbnail, medias.background, medias.year, medias.rating
FROM medias
WHERE
	(
		CASE WHEN $2::bool
		THEN EXISTS (
			SELECT id FROM media_names
			WHERE media_names.name ILIKE '%' || $3 || '%'
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

func FindMedias(ctx context.Context, arg FindMediasParams) ([]sqlc.Media, error) {
	db := postgres.DB.SqlDatabase

	// Dynamically update the query
	query := strings.Replace(findMedias, "{{sort_column}}", arg.SortColumn, 1)
	query = strings.Replace(query, "{{sort_order}}", arg.SortOrder, 1)
	args := []interface{}{
		arg.Offset,
		arg.SearchQuery,
		arg.Query,
		arg.SearchRating,
		arg.Rating,
		arg.SearchYear,
		arg.Year,
		arg.SearchGenres,
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
