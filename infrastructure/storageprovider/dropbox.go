package storageprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

var (
	errNotEnoughScope = errors.New("Not enough scope given from the Dropbox access token. Please grant the required scope 'files.content.write' and reset the access token.")
)

const dropboxProviderID uint = 12345678

type dropbox struct {
	remoteDirectory string
}

type dropboxError struct {
	HttpCode int
	Message  string
	Json     dropboxErrorJson
}

type dropboxErrorJson struct {
	ErrorSummary    string                 `json:"error_summary"`
	ErrorStructured map[string]interface{} `json:"error"`
	UserMessage     string                 `json:"user_message"`
}

// NewDropboxStorageProvider returns new StorageProviderService
func NewDropboxStorageProvider(remoteDirectory string) domain.StorageProviderService {
	return &dropbox{
		remoteDirectory: remoteDirectory,
	}
}

// ID returns provider ID
func (d *dropbox) ID() uint {
	return dropboxProviderID
}

// AccountInfo fetches Dropbox account information
func (d *dropbox) AccountInfo(cred domain.StorageProviderCredential) (domain.StorageProviderAccountInfo, error) {
	var accountInfo domain.StorageProviderAccountInfo

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.dropboxapi.com/2/users/get_current_account",
		nil,
	)
	if err != nil {
		return accountInfo, err
	}

	// prepare header (no need to set content-type)
	req.Header.Set("Authorization", "Bearer "+cred.UserAccessToken)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// do http request
	resp, err := client.Do(req)
	if err != nil {
		return accountInfo, err
	}

	defer resp.Body.Close()

	// read body
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accountInfo, err
	}

	var respBody map[string]interface{}

	err = json.Unmarshal(respBodyBytes, &respBody)
	if err != nil {
		return accountInfo, err
	}

	accountInfo.Email, _ = respBody["email"].(string)
	accountInfo.Photo, _ = respBody["profile_photo_url"].(string)

	return accountInfo, nil
}

// Upload sends the file to Dropbox server
func (d *dropbox) Upload(cred domain.StorageProviderCredential, file io.Reader, fileName, slug string) error {

	req, err := d.prepareRequest(cred.UserAccessToken, file, fileName, slug)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// do the request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		dropboxError, err := d.mapToDropboxError(res.Body, res.StatusCode)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		regularError := d.mapToRegularError(dropboxError)

		return regularError
	}

	return nil
}

func (d *dropbox) prepareRequest(accessToken string, file io.Reader, fileName, slug string) (*http.Request, error) {

	req, err := http.NewRequest(
		http.MethodPost,
		"https://content.dropboxapi.com/2/files/upload",
		file,
	)
	if err != nil {
		return nil, err
	}

	// construct Dropbox API arguments
	dropboxAPIArg := fmt.Sprintf(
		`{"path": "/%s/%s/%s","mode": "add","autorename": true,"mute": false}`,
		d.remoteDirectory,
		slug,
		fileName,
	)

	// prepare header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", dropboxAPIArg)
	return req, nil
}

func (d *dropbox) mapToDropboxError(responseReader io.Reader, httpStatusCode int) (dropboxError, error) {
	byteResponse, err := io.ReadAll(responseReader)
	if err != nil {
		return dropboxError{}, err
	}
	errorRes := dropboxError{HttpCode: httpStatusCode}
	errorRes.Message = string(byteResponse)
	if httpStatusCode == http.StatusUnauthorized {
		errorRes.Json = dropboxErrorJson{}
		json.Unmarshal(byteResponse, &errorRes.Json)
	}

	return errorRes, nil
}

func (d *dropbox) mapToRegularError(dropboxError dropboxError) error {
	if dropboxError.HttpCode == http.StatusBadRequest {
		if strings.Contains(dropboxError.Message, "files.content.write") {
			return errNotEnoughScope
		}
	} else if dropboxError.HttpCode == http.StatusUnauthorized {
		if dropboxError.Json.ErrorStructured[".tag"].(string) == "missing_scope" && dropboxError.Json.ErrorStructured["required_scope"].(string) == "files.content.write" {
			return errNotEnoughScope
		}
	}

	return errors.New("Unknown dropbox error.\n" + dropboxError.Message)
}
