package app

import (
	"io"
	"net/http"

	"github.com/wkloucek/cs3-wopi-server/pkg/internal/helpers"
)

// GetFile downloads the file from the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/getfile
func GetFile(app *demoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// download the file
	resp, err := helpers.DownloadFile(
		ctx,
		&wopiContext.FileReference,
		app.GatewayAPIClient,
		wopiContext.AccessToken,
		app.Config.CS3DataGatewayInsecure,
		app.Logger,
	)

	if err != nil || resp.StatusCode != http.StatusOK {
		app.Logger.Error().Err(err).Str("status_code", http.StatusText(resp.StatusCode)).Str("FileReference", wopiContext.FileReference.String()).Msg("GetFile: downloading the file failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// read the file from the body
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		app.Logger.Error().Str("FileReference", wopiContext.FileReference.String()).Msg("GetFile: copying the file content to the response body failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
}

// PutFile uploads the file to the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/putfile
func PutFile(app *demoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// read the file from the body
	defer r.Body.Close()

	// upload the file
	err := helpers.UploadFile(
		ctx,
		r.Body,
		&wopiContext.FileReference,
		app.GatewayAPIClient,
		wopiContext.AccessToken,
		r.Header.Get(HeaderWopiLock),
		app.Config.CS3DataGatewayInsecure,
		app.Logger,
	)

	if err != nil {
		app.Logger.Error().Err(err).Str("FileReference", wopiContext.FileReference.String()).Msg("PutFile: uploading the file failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
}
