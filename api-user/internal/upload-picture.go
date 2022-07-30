package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-user/queries"
	"github.com/trixky/hypertube/api-user/sqlc"
)

func UploadFile(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	ctx := context.Background()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	multipartFile, fileHeader, err := r.FormFile("picture")
	if err != nil {
		panic(err)
	}

	fmt.Println("filename", fileHeader.Filename)
	fmt.Println("size", fileHeader.Size)
	fmt.Println("file-header", fileHeader.Header)

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, multipartFile)
	if err != nil {
		panic(err)
	}

	fileParts := strings.Split(fileHeader.Filename, ".")
	extension := fileParts[len(fileParts)-1]
	path := StoragePath + fmt.Sprint(1) + "." + extension
	log.Println("Saving picture to", path)
	if err := os.WriteFile(path, buffer.Bytes(), 0); err != nil {
		log.Println("Failed to save image to disk")
		panic(err)
	}

	queries.SqlcQueries.UpdateUserPicture(ctx, sqlc.UpdateUserPictureParams{
		ID:        1,
		Extension: utils.MakeNullString(&extension),
	})

	w.WriteHeader(http.StatusOK)

	// // -------------------- get token
	// sanitized_token, err := utils.ExtractSanitizedTokenFromGrpcGatewayCookies("", ctx)

	// if err != nil {
	// 	return nil, err
	// }

	// // -------------------- cache
	// token_info, err := databases.RetrieveToken(sanitized_token)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	// }

	// // -------------------- db
	// user, err := queries.SqlcQueries.GetUserById(context.Background(), token_info.Id)

	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, status.Errorf(codes.NotFound, "no user found with this id")
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	// }

	// // -------------------- check uploaded file validity
	// if len(in.GetPicture()) == 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "Invalid picture")
	// }

	// // -------------------- Save file in storage (shared volume)
	// log.Println("picture", in.Picture)
	// path := StoragePath + fmt.Sprint(user.ID) + "." + "png"
	// if err := os.WriteFile(path, in.GetPicture(), 0); err != nil {
	// 	log.Println("Failed to save image to disk")
	// 	return nil, status.Errorf(codes.Internal, "failed to save image to disk")
	// }

	// return &pb.UploadPictureResponse{
	// 	Success: true,
	// }, nil
}
