package importer

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/trixky/hypertube/scripts/import-medias/postgres"
	"github.com/trixky/hypertube/scripts/import-medias/sqlc"
)

const batch_size = 10000

// Import all names from names.basic.tsv
// Relations are collected and returned as a pair to add them later when Medias are inserted
func ImportNames() (relations map[string][]string, err error) {
	fmt.Println("Loading names...")
	relations = make(map[string][]string)
	ctx := context.Background()
	file, err := openTsvFile("./medias/name.basics.tsv")
	defer file.Close()
	if err != nil {
		return
	}

	// Add all names from the file
	checked := 0
	scanner := bufio.NewScanner(file)
	for index := 0; scanner.Scan(); index++ {
		if index == 0 {
			continue
		}
		raw_data := scanner.Text()
		data := strings.Split(raw_data, "\t")

		checked += 1
		if checked%50000 == 0 {
			fmt.Printf("checked %v names\n", checked)
		}

		if len(data) == 6 {
			nconst := data[0]
			name := data[1]
			birthYear_int, _ := strconv.Atoi(data[2])
			birthYear := int32(birthYear_int)
			deathYear_int, _ := strconv.Atoi(data[3])
			deathYear := int32(deathYear_int)
			related_titles := strings.Split(data[5], ",")
			for _, relation := range related_titles {
				title_relations, ok := relations[relation]
				if ok {
					title_relations = append(title_relations, nconst)
				} else {
					relations[relation] = []string{nconst}
				}
			}

			// Insert the name -- ON CONFLICT will ignore duplicates
			err = postgres.DB.SqlcQueries.CreateName(ctx, sqlc.CreateNameParams{
				ImdbID:    nconst,
				Name:      name,
				BirthYear: makeNullInt32(&birthYear),
				DeathYear: makeNullInt32(&deathYear),
			})
			if err != nil {
				return
			}
		}
	}

	fmt.Printf("Checked %v names\n", checked)

	return
}

func GetRatings() (map[string]float64, error) {
	fmt.Println("Loading ratings...")
	ratings := make(map[string]float64)
	file, err := openTsvFile("./medias/title.ratings.tsv")
	defer file.Close()
	if err != nil {
		return ratings, err
	}

	// Find all ratings in the file
	scanner := bufio.NewScanner(file)
	for index := 0; scanner.Scan(); index++ {
		if index == 0 {
			continue
		}
		raw_data := scanner.Text()
		data := strings.Split(raw_data, "\t")
		if len(data) == 3 {
			rating, _ := strconv.ParseFloat(data[1], 64)
			ratings[data[0]] = rating
		}
	}

	fmt.Printf("Found %v ratings\n", len(ratings))

	return ratings, nil
}

// Import Titles and all related informations
func ImportTitles(ratings *map[string]float64) error {
	fmt.Println("Loading titles...")
	ctx := context.Background()
	file, err := openTsvFile("./medias/title.basics.tsv")
	defer file.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)

	var record Title
	checked := 0
	added := 0

	// Add all medias from the file
	tmp_medias := make([]sqlc.CreateMediasParams, 0, batch_size)
	for index := 0; scanner.Scan(); index++ {
		if index == 0 {
			continue
		}
		raw_data := scanner.Text()
		data := strings.Split(raw_data, "\t")

		checked += 1
		if checked%10000 == 0 {
			fmt.Printf("checked %v medias\n", checked)
		}

		if len(data) == 9 {
			record.TConst = data[0]
			record.TitleType = data[1]
			record.PrimaryTitle = data[2]
			record.OriginalTitle = data[3]
			isAdult, _ := strconv.Atoi(data[4])
			record.IsAdult = int8(isAdult)
			startYear, _ := strconv.Atoi(data[5])
			record.StartYear = int32(startYear)
			endYear, _ := strconv.Atoi(data[6])
			record.EndYear = int32(endYear)
			runtimeMinutes, _ := strconv.Atoi(data[7])
			record.RuntimeMinutes = int32(runtimeMinutes)
			record.Genres = data[8]
			if record.Genres == "\\N" {
				record.Genres = ""
			}

			// Ignore adult shows and TV series
			if record.IsAdult == 0 && record.TitleType != "tvEpisode" && record.TitleType != "tvSeries" && record.TitleType != "tvMiniSeries" {
				count_total, err := postgres.DB.SqlcQueries.CheckMediaExistByIMDB(ctx, record.TConst)
				if err != nil {
					return err
				} else if count_total >= 1 {
					continue
				} else {
					err = nil
				}

				// Accumulate medias to insert them by bigger groups to make it faster
				rating := (*ratings)[record.TConst]
				tmp_medias = append(tmp_medias, sqlc.CreateMediasParams{
					ImdbID:   record.TConst,
					Year:     record.StartYear,
					Duration: record.RuntimeMinutes,
					Genres:   record.Genres,
					Rating:   makeNullFloat64(&rating),
				})
				if len(tmp_medias) == batch_size {
					_, err = postgres.DB.SqlcQueries.CreateMedias(ctx, tmp_medias)
					if err != nil {
						return err
					}
					tmp_medias = make([]sqlc.CreateMediasParams, 0, batch_size)
				}

				added += 1
			}
		}
	}

	// Add remaining medias
	if len(tmp_medias) > 0 {
		_, err = postgres.DB.SqlcQueries.CreateMedias(ctx, tmp_medias)
		if err != nil {
			return err
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Added %v medias\n", added)

	return nil
}

func ImportTitlesNames() error {
	fmt.Println("Loading titles names...")
	ctx := context.Background()
	file, err := openTsvFile("./medias/title.akas.tsv")
	defer file.Close()
	if err != nil {
		return err
	}

	// Keep track of IDs to ignore
	ignore := make(map[string]bool)
	existing := make(map[string]int64)
	var current_names_for int64 = 0
	current_names := make(map[string]bool)

	// Read line by line
	checked := 0
	scanner := bufio.NewScanner(file)
	for index := 0; scanner.Scan(); index++ {
		if index == 0 {
			continue
		}
		raw_data := scanner.Text()
		data := strings.Split(raw_data, "\t")

		// Ignore header
		checked += 1
		if checked%50000 == 0 {
			fmt.Printf("checked %v names\n", checked)
		}

		// Import each names
		if len(data) == 8 {
			// Check if we can ignore the Title
			tconst := data[0]
			if _, ok := ignore[tconst]; ok {
				continue
			}
			lang := data[3]
			// Only save langs in French, English and the original title
			if lang == "\\N" {
				lang = "__"
			} else if lang != "FR" && lang != "US" && lang != "GB" {
				continue
			}

			// Check if the Media exist
			media_id, ok := existing[tconst]
			if !ok {
				existing_media, err := postgres.DB.SqlcQueries.GetMediaByIMDB(ctx, tconst)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return err
					} else {
						ignore[tconst] = true
						continue
					}
				} else {
					err = nil
				}
				media_id = existing_media.ID
				existing[tconst] = media_id
			}

			// Avoid adding duplicate names for the same language
			name := data[2]
			if current_names_for != media_id {
				current_names_for = media_id
				current_names = make(map[string]bool)
			} else if current_names[name+lang] {
				fmt.Println("skipped", name, lang)
				continue
			}
			current_names[name+lang] = true

			// Insert the Media Name -- ON CONFLICT will ignore duplicates
			err_insert := postgres.DB.SqlcQueries.CreateMediaName(ctx, sqlc.CreateMediaNameParams{
				MediaID: int32(media_id),
				Name:    name,
				Lang:    lang,
			})
			if err_insert != nil {
				return err_insert
			}
		}
	}

	fmt.Printf("Checked %v names\n", checked)

	return nil
}

func ImportTitlesPrincipal() error {
	fmt.Println("Loading titles principals...")
	ctx := context.Background()
	file, err := openTsvFile("./medias/title.principals.tsv")
	defer file.Close()
	if err != nil {
		return err
	}

	// Keep track of IDs to ignore
	ignore := make(map[string]bool)
	media_exist := make(map[string]int64)
	name_exist := make(map[string]int64)

	// Read line by line
	checked := 0
	scanner := bufio.NewScanner(file)
	for index := 0; scanner.Scan(); index++ {
		if index == 0 {
			continue
		}
		raw_data := scanner.Text()
		data := strings.Split(raw_data, "\t")

		checked += 1
		if checked%50000 == 0 {
			fmt.Printf("checked %v principals\n", checked)
		}

		// Check if we can ignore the Title
		tconst := data[0]
		if _, ok := ignore[tconst]; ok {
			continue
		}
		nconst := data[2]
		if _, ok := ignore[nconst]; ok {
			continue
		}

		// Check if the Name exist
		name_id, ok := name_exist[nconst]
		if !ok {
			exisiting_name, err := postgres.DB.SqlcQueries.GetNameByIMDB(ctx, nconst)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return err
				} else {
					ignore[nconst] = true
					continue
				}
			} else {
				err = nil
			}
			name_id = exisiting_name.ID
			name_exist[nconst] = name_id
		}

		// Check if the Media exist
		media_id, ok := media_exist[tconst]
		if !ok {
			existing_media, err := postgres.DB.SqlcQueries.GetMediaByIMDB(ctx, tconst)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return err
				} else {
					ignore[tconst] = true
					continue
				}
			} else {
				err = nil
			}
			media_id = existing_media.ID
			media_exist[tconst] = media_id
		}

		// Insert the Media Principal -- ON CONFLICT will ignore duplicates
		category := data[3]
		if category == "actor" {
			// Characters are stored as array of strings
			raw_characters := data[5]
			if raw_characters == "\\N" {
				continue
			}
			var characters []string
			err = json.Unmarshal([]byte(raw_characters), &characters)
			if err != nil {
				return err
			}

			// Add each characters for the Name
			for _, character := range characters {
				err_insert := postgres.DB.SqlcQueries.CreateMediaActor(ctx, sqlc.CreateMediaActorParams{
					MediaID:   int32(media_id),
					NameID:    int32(name_id),
					Character: makeNullString(&character),
				})
				if err_insert != nil {
					return err_insert
				}
			}
		} else {
			role := data[4]
			err_insert := postgres.DB.SqlcQueries.CreateMediaStaff(ctx, sqlc.CreateMediaStaffParams{
				MediaID: int32(media_id),
				NameID:  int32(name_id),
				Role:    makeNullString(&role),
			})
			if err_insert != nil {
				return err_insert
			}
		}
	}

	fmt.Printf("Checked %v principals\n", checked)

	return nil
}

func ImportNameRelations(relations *map[string][]string) error {
	fmt.Println("Loading name relations...")
	ctx := context.Background()

	// Keep track of IDs to ignore
	ignore := make(map[string]bool)
	media_exist := make(map[string]int64)
	name_exist := make(map[string]int64)

	// Add all relations
	for relation_imdb, name_imdbs := range *relations {
		// Check if we can ignore the Media
		if _, ok := ignore[relation_imdb]; ok {
			continue
		}

		// Check if the Media exist
		media_id, ok := media_exist[relation_imdb]
		if !ok {
			existing_media, err := postgres.DB.SqlcQueries.GetMediaByIMDB(ctx, relation_imdb)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return err
				} else {
					ignore[relation_imdb] = true
					continue
				}
			} else {
				err = nil
			}
			media_id = existing_media.ID
			media_exist[relation_imdb] = media_id
		}

		for _, name_imdb := range name_imdbs {
			// Check if we can ignore the Name
			if _, ok := ignore[name_imdb]; ok {
				continue
			}

			// Check if the Name exist
			name_id, ok := name_exist[name_imdb]
			if !ok {
				exisiting_name, err := postgres.DB.SqlcQueries.GetNameByIMDB(ctx, name_imdb)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return err
					} else {
						ignore[name_imdb] = true
						continue
					}
				} else {
					err = nil
				}
				name_id = exisiting_name.ID
				name_exist[name_imdb] = name_id
			}

			// Add the relation
			err_insert := postgres.DB.SqlcQueries.CreateNameRelation(ctx, sqlc.CreateNameRelationParams{
				MediaID: int32(media_id),
				NameID:  int32(name_id),
			})
			if err_insert != nil {
				return err_insert
			}
		}

	}

	fmt.Printf("Checked %v name relations\n", len(*relations))

	return nil
}

// types: [short movie tvEpisode tvSeries tvShort tvMovie tvMiniSeries video tvSpecial videoGame tvPilot]
func Import() error {
	// Import names first
	name_relations, err := ImportNames()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Open Medias related informations
	// ratings, err := GetRatings()
	// if err != nil {
	// fmt.Println(err)
	// 	return err
	// }

	// Import Titles
	// err = ImportTitles(&ratings)
	// if err != nil {
	// fmt.Println(err)
	// 	return err
	// }

	// Add Title names after they exists
	// err := ImportTitlesNames()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// Import Title crew and principals
	// err := ImportTitlesPrincipal()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// Import Name relations to Titles
	err = ImportNameRelations(&name_relations)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Done !")

	return nil
}
