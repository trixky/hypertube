package sites

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/trixky/hypertube/api-scrapper/postgres"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	ut "github.com/trixky/hypertube/api-scrapper/utils"
)

type TMDBMovieResult struct {
	Id               int32   `json:"id"`
	Adult            bool    `json:"adult"`
	BackdropPath     *string `json:"backdrop_path"`
	GenreIds         []int32 `json:"genre_ids"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       *string `json:"poster_path"`
	Popularity       float64 `json:"popularity"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int32   `json:"vote_count"`
}

type TMDBKnownFor struct {
	Id               int32   `json:"id"`
	MediaType        string  `json:"media_type"`
	Overview         string  `json:"overview"`
	PosterPath       *string `json:"poster_path"`
	BackdropPath     *string `json:"backdrop_path"`
	OriginalLanguage string  `json:"original_language"`
	GenreIds         []int32 `json:"genre_ids"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int32   `json:"vote_count"`
	Title            *string `json:"title"`
	OriginalTitle    *string `json:"original_title"`
}

type TMDBPersonResult struct {
	Id          int32        `json:"id"`
	ProfilePath *string      `json:"profile_path"`
	Adult       bool         `json:"adult"`
	Name        string       `json:"name"`
	Popularity  int32        `json:"popularity"`
	KnownFor    TMDBKnownFor `json:"known_for"`
}

type TMDBFindResponse struct {
	MovieResults  []TMDBMovieResult  `json:"movie_results"`
	PersonResults []TMDBPersonResult `json:"person_results"`
}

type TMDBCast struct {
	Id                 int32   `json:"id"`
	Adult              bool    `json:"adult"`
	Gender             int32   `json:"gender"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
	CastId             int32   `json:"cast_id"`
	Character          string  `json:"character"`
	CreditId           string  `json:"credit_id"`
	Order              int32   `json:"order"`
}

type TMDBCrew struct {
	Id                 int32   `json:"id"`
	Adult              bool    `json:"adult"`
	Gender             int32   `json:"gender"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
	CreditId           string  `json:"credit_id"`
	Department         string  `json:"department"`
	Job                string  `json:"job"`
}

type TMDBMovieResponse struct {
	Id                  int32   `json:"id"`
	Adult               bool    `json:"adult"`
	PosterPath          *string `json:"poster_path"`
	BackdropPath        *string `json:"backdrop_path"`
	BelongsToCollection *struct {
		Id           int32   `json:"id"`
		Name         string  `json:"name"`
		PosterPath   *string `json:"poster_path"`
		BackdropPath *string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int64 `json:"budget"`
	Genres []struct {
		Id   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	HomePage            string  `json:"homepage"`
	IMDBId              *string `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	ProductionCompanies []struct {
		Id            int32   `json:"id"`
		Name          string  `json:"name"`
		LogoPath      *string `json:"logo_path"`
		OriginCountry string  `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		ISO_3166_1 string `json:"iso_3166_1"`
		Name       string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int64  `json:"revenue"`
	Runtime         int32  `json:"runtime"`
	SpokenLanguages []struct {
		ISO_639_1 string `json:"iso_639_1"`
		Name      string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     *string `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int32   `json:"vote_count"`
	Credits     struct {
		Cast []TMDBCast `json:"cast"`
		Crew []TMDBCrew `json:"crew"`
	} `json:"credits"`
	Translations struct {
		Translations []struct {
			ISO_3166_1  string `json:"iso_3166_1"`
			ISO_639_1   string `json:"iso_639_1"`
			Name        string `json:"name"`
			EnglishName string `json:"english_name"`
			Data        struct {
				Title    string `json:"title"`
				Overview string `json:"overview"`
				HomePage string `json:"homepage"`
				Runtime  int32  `json:"runtime"`
				Tagline  string `json:"tagline"`
			}
		} `json:"translations"`
	} `json:"translations"`
}

type TMDBSearchResponse struct {
	Page    int32 `json:"page"`
	Results []struct {
		Id               int32   `json:"id"`
		Adult            bool    `json:"adult"`
		PosterPath       *string `json:"poster_path"`
		BackdropPath     *string `json:"backdrop_path"`
		Overview         string  `json:"overview"`
		Genres           []int32 `json:"genre_ids"`
		ReleaseDate      string  `json:"release_date"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Title            string  `json:"title"`
		Popularity       float64 `json:"popularity"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int32   `json:"vote_count"`
	} `json:"results"`
	TotalResults int32 `json:"total_results"`
	TotalPages   int32 `json:"total_pages"`
}

type MediaCrew struct {
	Id        int32
	Name      string
	Thumbnail *string
	Job       string
}

type MediaActor struct {
	Id        int32
	Name      string
	Thumbnail *string
	Character string
}

type MediaName struct {
	Lang  string
	Title string
}

type MediaInformations struct {
	ImdbId      *string
	TmdbId      int32
	Title       string
	Poster      *string
	Background  *string
	Description string
	Year        *int32
	Duration    *int32
	Rating      *float64
	Genres      []string
	Crew        []MediaCrew
	Actors      []MediaActor
	Names       []MediaName
}

var year_extractor = regexp.MustCompile("(\\d{4})")
var match_name = regexp.MustCompile("(?i)(.+?)\\s*(?:\\(?(\\d{4})\\)?)(?:\\s*-\\s*|\\s*)(?:4k|2k|2160p|1080p|720p|proper|hqcam|cam|ts|blu-?ray|dvdrip|brrip|hdrip|x24|h264|x265|h265|web|hmax|imax|\\(|\\[|\\{)?")

var api_key = os.Getenv("TMDB_API_KEY")

func GenerateImage(size string, path string) string {
	return "https://image.tmdb.org/t/p/" + size + path
}

func GetTMDBMedia(tmdb_id int32) (informations *MediaInformations, err error) {
	resp, err := http.Get("https://api.themoviedb.org/3/movie/" + fmt.Sprint(tmdb_id) + "?api_key=" + api_key + "&append_to_response=credits,translations")
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var media TMDBMovieResponse
	err = json.Unmarshal(body, &media)
	if err != nil {
		return
	}

	// Convert all informations to MediaInformations
	informations = &MediaInformations{
		TmdbId: tmdb_id,
		Crew:   make([]MediaCrew, 0),
		Actors: make([]MediaActor, 0),
		Names:  make([]MediaName, 0),
	}
	informations.TmdbId = media.Id
	informations.Title = media.Title
	if media.PosterPath != nil {
		poster := GenerateImage("w500", *media.PosterPath)
		informations.Poster = &poster
	}
	if media.BackdropPath != nil {
		background := GenerateImage("original", *media.BackdropPath)
		informations.Background = &background
	}
	informations.Genres = make([]string, 0)
	for _, genre := range media.Genres {
		informations.Genres = append(informations.Genres, genre.Name)
	}
	informations.Description = media.Overview
	if media.ReleaseDate != "" {
		var year_extractor = regexp.MustCompile("(\\d{4})")
		matches := year_extractor.FindStringSubmatch(media.ReleaseDate)
		if len(matches) == 2 {
			if year, _ := strconv.Atoi(matches[1]); year > 0 {
				year_int32 := int32(year)
				informations.Year = &year_int32
			}
		}
	}
	if media.Runtime > 0 {
		informations.Duration = &media.Runtime
	}
	if media.VoteAverage > 0 {
		informations.Rating = &media.VoteAverage
	}

	// Add relations
	for _, crew := range media.Credits.Crew {
		var picture *string
		if crew.ProfilePath != nil {
			picture_path := GenerateImage("w500", *crew.ProfilePath)
			picture = &picture_path
		}
		informations.Crew = append(informations.Crew, MediaCrew{
			Id:        crew.Id,
			Name:      crew.Name,
			Thumbnail: picture,
			Job:       crew.Job,
		})
	}
	for _, actor := range media.Credits.Cast {
		var picture *string
		if actor.ProfilePath != nil {
			picture_path := GenerateImage("w500", *actor.ProfilePath)
			picture = &picture_path
		}
		informations.Actors = append(informations.Actors, MediaActor{
			Id:        actor.Id,
			Name:      actor.Name,
			Thumbnail: picture,
			Character: actor.Character,
		})
	}
	for _, translation := range media.Translations.Translations {
		if translation.Data.Title != "" {
			informations.Names = append(informations.Names, MediaName{
				Lang:  translation.ISO_3166_1,
				Title: translation.Data.Title,
			})
		}
	}

	return
}

func SearchTMDBMedia(query string, year int32) (tmdb_id int32, err error) {
	args := map[string]string{
		"api_key":              api_key,
		"query":                url.QueryEscape(query),
		"primary_release_year": fmt.Sprint(year),
		"page":                 "1",
		"include_adult":        "false",
	}

	query_args := ""
	for arg, arg_value := range args {
		if query_args != "" {
			query_args = query_args + "&"
		}
		query_args = query_args + arg + "=" + arg_value
	}
	log.Println("searching", "https://api.themoviedb.org/3/search/movie", args["query"], args["primary_release_year"])
	resp, err := http.Get("https://api.themoviedb.org/3/search/movie?" + query_args)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result TMDBSearchResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	// Select the first matching media
	// -- and ignore the search if there is more than 20 results
	if result.TotalResults == 0 || result.TotalResults > 20 {
		return
	}
	media := result.Results[0]

	return media.Id, nil
}

func GetIMDBMedia(imdb_id string) (informations *MediaInformations, err error) {
	resp, err := http.Get("https://api.themoviedb.org/3/find/" + imdb_id + "?api_key=" + api_key + "&external_source=imdb_id")
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var find_response TMDBFindResponse
	err = json.Unmarshal(body, &find_response)
	if err != nil {
		return
	}

	if len(find_response.MovieResults) == 1 {
		tmdb_id := find_response.MovieResults[0].Id
		time.Sleep(time.Second)
		informations, err = GetTMDBMedia(tmdb_id)
	}

	if informations != nil {
		informations.ImdbId = &imdb_id
	}

	return
}

func InsertMediaInformations(informations *MediaInformations) (media_id int32, err error) {
	ctx := context.Background()

	// Insert Media
	created_media, err := postgres.DB.SqlcQueries.CreateMedia(ctx, sqlc.CreateMediaParams{
		ImdbID:      ut.MakeNullString(informations.ImdbId),
		TmdbID:      informations.TmdbId,
		Description: ut.MakeNullString(&informations.Description),
		Duration:    ut.MakeNullInt32(informations.Duration),
		Thumbnail:   ut.MakeNullString(informations.Poster),
		Background:  ut.MakeNullString(informations.Background),
		Year:        ut.MakeNullInt32(informations.Year),
		Rating:      ut.MakeNullFloat64(informations.Rating),
	})
	if err != nil {
		return
	}
	media_id = int32(created_media.ID)

	// Insert the main media name
	err = postgres.DB.SqlcQueries.CreateMediaName(ctx, sqlc.CreateMediaNameParams{
		MediaID: media_id,
		Name:    informations.Title,
		Lang:    "__",
	})
	if err != nil {
		return
	}
	// ... and insert other media names
	for _, media_name := range informations.Names {
		err = postgres.DB.SqlcQueries.CreateMediaName(ctx, sqlc.CreateMediaNameParams{
			MediaID: media_id,
			Name:    media_name.Title,
			Lang:    media_name.Lang,
		})
		if err != nil {
			return
		}
	}

	// Insert media genres
	for _, genre_name := range informations.Genres {
		// Check if the genre exist
		genre, err := postgres.DB.SqlcQueries.GetGenre(ctx, genre_name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
				created_genre, err := postgres.DB.SqlcQueries.CreateGenre(ctx, genre_name)
				if err != nil {
					return media_id, err
				}
				genre.ID = created_genre.ID
			} else {
				return media_id, err
			}
		}

		// Add the relation
		err = postgres.DB.SqlcQueries.CreateMediaGenre(ctx, sqlc.CreateMediaGenreParams{
			MediaID: media_id,
			GenreID: int32(genre.ID),
		})
		if err != nil {
			return media_id, err
		}
	}

	// Insert related persons
	for _, crew := range informations.Crew {
		// Check if the name exist
		name, err := postgres.DB.SqlcQueries.GetNameByTMDB(ctx, crew.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
				created_name, err := postgres.DB.SqlcQueries.CreateName(ctx, sqlc.CreateNameParams{
					TmdbID:    crew.Id,
					Name:      crew.Name,
					Thumbnail: ut.MakeNullString(crew.Thumbnail),
				})
				if err != nil {
					return media_id, err
				}
				name.ID = created_name.ID
			} else {
				return media_id, err
			}
		}

		// Add the relation
		err = postgres.DB.SqlcQueries.CreateMediaStaff(ctx, sqlc.CreateMediaStaffParams{
			MediaID: media_id,
			NameID:  int32(name.ID),
			Role:    ut.MakeNullString(&crew.Job),
		})
		if err != nil {
			return media_id, err
		}
	}
	for _, actor := range informations.Actors {
		// Check if the name exist
		name, err := postgres.DB.SqlcQueries.GetNameByTMDB(ctx, actor.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
				created_name, err := postgres.DB.SqlcQueries.CreateName(ctx, sqlc.CreateNameParams{
					TmdbID:    actor.Id,
					Name:      actor.Name,
					Thumbnail: ut.MakeNullString(actor.Thumbnail),
				})
				if err != nil {
					return media_id, err
				}
				name.ID = created_name.ID
			} else {
				return media_id, err
			}
		}

		// Add the relation
		err = postgres.DB.SqlcQueries.CreateMediaActor(ctx, sqlc.CreateMediaActorParams{
			MediaID:   media_id,
			NameID:    int32(name.ID),
			Character: ut.MakeNullString(&actor.Character),
		})
		if err != nil {
			return media_id, err
		}
	}

	return
}

func TryInsertOrGetMedia(name string) (media_id int32, err error) {
	ctx := context.Background()

	// Try to fix the name for the query, and find a Year
	clean_name := strings.ReplaceAll(name, ".", " ")
	matches := match_name.FindStringSubmatch(clean_name)
	if len(matches) != 3 {
		return
	}
	query := matches[1]
	year_int, err := strconv.Atoi(matches[2])
	year := int32(year_int)

	// Search for a Media
	result, err := SearchTMDBMedia(query, year)
	if result == 0 || err != nil {
		log.Println("no match for", query, year)
		return
	} else {
		log.Println("found TMDB ID", result, "from", query, year)
	}
	existing_media, err := postgres.DB.SqlcQueries.GetMediaByTMDBID(ctx, result)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	err = nil

	// Check if the Media already exists and return it
	if existing_media.ID > 0 {
		media_id = int32(existing_media.ID)
		return media_id, err
	}
	// ... or else find all informations and save them

	// Get all informations from TMDB
	informations, err := GetTMDBMedia(result)
	if informations == nil || err != nil {
		return
	}
	media_id, err = InsertMediaInformations(informations)

	return

}

func InsertOrGetMedia(imdb_id string) (media_id int32, err error) {
	ctx := context.Background()

	existing_media, err := postgres.DB.SqlcQueries.GetMediaByIMDB(ctx, ut.MakeNullString(&imdb_id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	err = nil

	// Check if the Media already exists and return it
	if existing_media.ID > 0 {
		media_id = int32(existing_media.ID)
		return media_id, err
	}
	// ... or else find all informations and save them

	// Get all informations from TMDB
	informations, err := GetIMDBMedia(imdb_id)
	if informations == nil || err != nil {
		return
	}
	media_id, err = InsertMediaInformations(informations)

	return
}
