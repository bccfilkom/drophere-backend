package main

import (
	"encoding/json"
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
	storageProviderPool domain.StorageProviderPool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// get file
		f, fileHeader, err := r.FormFile("file")
		if err != nil {
			if debug {
				log.Println("read file: ", err)
			}
			writeError(w, "Invalid File")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// get linkID
		linkID, err := strconv.Atoi(r.FormValue("linkId"))
		if err != nil {
			if debug {
				log.Println("parsing link ID: ", err)
			}
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
				if debug {
					log.Println("file upload: ", err)
				}
				writeError(w, "Server Error")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// check if the link is connected to a Storage Provider
		if l.UserStorageCredentialID == nil || *l.UserStorageCredentialID < 1 || l.UserStorageCredential == nil {
			writeError(w, "The link is unavailable")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

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

		storageProviderService, err := storageProviderPool.Get(l.UserStorageCredential.ProviderID)
		if err != nil {
			if debug {
				log.Println("get storage provider service: ", err)
			}
			writeError(w, "Sorry, but the Storage Provider is unavailable at the time")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		err = storageProviderService.Upload(
			domain.StorageProviderCredential{
				UserAccessToken: l.UserStorageCredential.ProviderCredential,
			},
			f,
			fileHeader.Filename,
			l.Slug,
		)
		if err != nil {
			if debug {
				log.Println("file upload: ", err)
			}
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
