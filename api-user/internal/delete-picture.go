package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/api-user/queries"
)

func DeletePicture(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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

	// -------------------- delete file
	err = os.Remove(StoragePath + fmt.Sprint(user.ID) + "." + user.Extension.String)
	if err != nil {
		// Only log the error, there is valid reasons that the delete could fail
		log.Println("Error delete user", user.ID, "picture", err)
	}

	// -------------------- delete the reference
	if err := queries.SqlcQueries.DeleteUserPicture(ctx, token_info.Id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Failed to delete file\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\":true}"))
}
