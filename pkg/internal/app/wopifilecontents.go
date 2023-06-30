package app

import (
	"bytes"
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
		app.gwc,
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
	file, err := io.ReadAll(resp.Body)
	if err != nil {
		app.Logger.Error().Str("FileReference", wopiContext.FileReference.String()).Msg("GetFile: reading from the download body failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// but just return the file here
	_, err = w.Write(file)
	if err != nil {
		app.Logger.Error().Str("FileReference", wopiContext.FileReference.String()).Msg("GetFile: writing to the response body failed")
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
	file, err := io.ReadAll(r.Body)
	if err != nil {
		app.Logger.Error().Str("FileReference", wopiContext.FileReference.String()).Msg("PutFile: reading from the body failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// upload the file
	err = helpers.UploadFile(
		ctx,
		bytes.NewReader(file),
		&wopiContext.FileReference,
		app.gwc,
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
