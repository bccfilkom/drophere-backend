package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

func writeError(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": []map[string]string{
			{
				"message": msg,
			},
		},
	})
}

func fileUploadHandler(
	userSvc domain.UserService,
	linkSvc domain.LinkService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// get file
		f, fileHeader, err := r.FormFile("file")
		if err != nil {
			writeError(w, "Invalid File")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// get linkID
		linkID, err := strconv.Atoi(r.FormValue("linkId"))
		if err != nil {
			writeError(w, "Invalid Link ID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// fetch link from database
		l, err := linkSvc.FetchLink(uint(linkID))
		if err != nil {
			if err == domain.ErrLinkNotFound {
				writeError(w, err.Error())
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Println("file upload: ", err)
				writeError(w, "Server Error")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// check if the link creator has Dropbox Access Token
		if l.User.DropboxToken == nil {
			writeError(w, "User Dropbox is not connected")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		dropboxAccessToken := *l.User.DropboxToken

		// check for password
		if l.IsProtected() {
			password := r.FormValue("password")

			if !linkSvc.CheckLinkPassword(l, password) {
				writeError(w, "Invalid Password")
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		// check for deadline
		if l.Deadline != nil && l.Deadline.Before(time.Now()) {
			writeError(w, "Link is Expired")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// prepare request to Dropbox
		req, err := http.NewRequest(http.MethodPost, "https://content.dropboxapi.com/2/files/upload", f)
		if err != nil {
			log.Println("file upload: ", err)
			writeError(w, "Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// construct Dropbox API arguments
		dropboxAPIArg := fmt.Sprintf(
			`{"path": "/%s/%s/%s","mode": "add","autorename": true,"mute": false}`,
			"drophere-dev",
			l.Slug,
			fileHeader.Filename,
		)

		// prepare header
		req.Header.Set("Authorization", "Bearer "+dropboxAccessToken)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Dropbox-API-Arg", dropboxAPIArg)

		// make the request
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Println("file upload: ", err)
			writeError(w, "Server Error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "File is successfully uploaded",
		})
		w.WriteHeader(http.StatusOK)
	}
}
