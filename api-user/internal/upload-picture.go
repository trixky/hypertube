package internal

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	osuser "os/user"
	"strconv"
	"strings"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-user/queries"
	"github.com/trixky/hypertube/api-user/sqlc"
)

func UploadPicture(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	ctx := context.Background()

	// Manually check cookies
	cookies := r.Cookies()
	token := ""
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			token = cookie.Value
			break
		}
	}
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\":\"You need to be logged in to update your picture\"}"))
		return
	}

	// -------------------- cache
	token_info, err := databases.RetrieveToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\":\"You need to be logged in to update your picture\"}"))
		return
	}

	// -------------------- db
	user, err := queries.SqlcQueries.GetUserById(ctx, token_info.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{\"error\":\"No users found for this token\"}"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Failed to find user associated with token\"}"))
		return
	}

	// -------------------- Parse FormData input
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	multipartFile, fileHeader, err := r.FormFile("picture")
	if err != nil {
		panic(err)
	}

	fmt.Println("filename", fileHeader.Filename)
	fmt.Println("size", fileHeader.Size)
	// fmt.Println("file-header", fileHeader.Header)

	// -------------------- Check file validity
	if fileHeader.Size > 2_000_000 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"File too large, the limit is 2Mb\"}"))
		return
	}

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, multipartFile)
	if err != nil {
		panic(err)
	}

	fileParts := strings.Split(fileHeader.Filename, ".")
	if len(fileParts) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Invalid file\"}"))
		return
	}
	extension := fileParts[len(fileParts)-1]
	path := StoragePath + fmt.Sprint(user.ID) + "." + extension
	log.Println("Saving picture to", path)

	// -------------------- Save the picture to disk and in the database
	if err := os.WriteFile(path, buffer.Bytes(), 0); err != nil {
		log.Println("Failed to save image to disk")
		panic(err)
	}
	us, err := osuser.Lookup("hypertube")
	uid, _ := strconv.Atoi(us.Uid)
	gid, _ := strconv.Atoi(us.Gid)
	if err := os.Chown(path, uid, gid); err != nil {
		log.Println("Failed to save image to disk")
		panic(err)
	}

	queries.SqlcQueries.UpdateUserPicture(ctx, sqlc.UpdateUserPictureParams{
		ID:        user.ID,
		Extension: utils.MakeNullString(&extension),
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\":true}"))
}
