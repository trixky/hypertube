package importer

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/trixky/hypertube/scripts/import-medias/postgres"
	"github.com/trixky/hypertube/scripts/import-medias/sqlc"
)

const batch_size = 1000

func makeNullInt32(value *int32) (null_int32 sql.NullInt32) {
	if value == nil || *value == 0 {
		return
	}
	null_int32.Int32 = *value
	null_int32.Valid = true
	return
}

func makeNullString(value *string) (null_string sql.NullString) {
	if value == nil {
		return
	}
	null_string.String = *value
	null_string.Valid = true
	return
}

// types: [short movie tvEpisode tvSeries tvShort tvMovie tvMiniSeries video tvSpecial videoGame tvPilot]
func Import() error {
	ctx := context.Background()
	// TODO Read and parse all tsv files

	// Read data from CSV file
	csvFile, err := os.Open("./medias/title.basics.tsv")
	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var record Title
	checked := 0
	added := 0

	// Add all medias from the file
	tmp_medias := make([]sqlc.CreateMediasParams, 0, batch_size)
	for index, each := range csvData {
		if index > 0 {
			checked += 1
			if len(each) == 9 {
				record.TConst = each[0]
				record.TitleType = each[1]
				record.PrimaryTitle = each[2]
				record.OriginalTitle = each[3]
				isAdult, _ := strconv.Atoi(each[4])
				record.IsAdult = int8(isAdult)
				startYear, _ := strconv.Atoi(each[5])
				record.StartYear = int32(startYear)
				endYear, _ := strconv.Atoi(each[6])
				record.EndYear = int32(endYear)
				runtimeMinutes, _ := strconv.Atoi(each[7])
				record.RuntimeMinutes = int32(runtimeMinutes)
				record.Genres = each[8]
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
					tmp_medias = append(tmp_medias, sqlc.CreateMediasParams{
						ImdbID:   record.TConst,
						Year:     record.StartYear,
						Duration: record.RuntimeMinutes,
						Genres:   record.Genres,
						// Rating: makeNullInt32(&record.EndYear), // TODO titles.rating.tsv
					})
					if len(tmp_medias) == batch_size {
						_, err = postgres.DB.SqlcQueries.CreateMedias(ctx, tmp_medias)
						if err != nil {
							fmt.Println("after insert", err)
							return err
						}
						tmp_medias = make([]sqlc.CreateMediasParams, 0, batch_size)
					}
					added += 1
				}
			}
			if checked%10000 == 0 {
				fmt.Printf("checked %v medias\n", checked)
			}
		}
	}

	// Add remaining medias
	if len(tmp_medias) > 0 {
		_, err = postgres.DB.SqlcQueries.CreateMedias(ctx, tmp_medias)
		if err != nil {
			fmt.Println("after insert", err)
			return err
		}
		tmp_medias = make([]sqlc.CreateMediasParams, 0, batch_size)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Added %v medias\n", added)

	return nil
}
