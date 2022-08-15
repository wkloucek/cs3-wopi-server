package app

import (
	"net/http"
	"time"

	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
)

const (
	// WOPI Locks generally have a lock duration of 30 minutes and will be refreshed before expiration if needed
	// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/concepts#lock
	lockDuration time.Duration = 30 * time.Minute
)

// GetLock returns a lock or an empty string if no lock exists
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/getlock
func GetLock(app *demoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	req := &providerv1beta1.GetLockRequest{
		Ref: &wopiContext.FileReference,
	}
	resp, err := app.gwc.GetLock(
		ctx,
		req,
	)
	if err != nil {
		app.Logger.Error().Err(err).Str("FileReference", wopiContext.FileReference.String()).Msg("GetLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().Str("status_code", resp.Status.Code.String()).Str("status_msg", resp.Status.Message).Str("FileReference", wopiContext.FileReference.String()).Msg("GetLock failed")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	lockID := ""
	if resp.Lock != nil {
		lockID = resp.Lock.LockId
	}
	w.Header().Set(HeaderWopiLock, lockID)
	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}

// Lock returns a WOPI lock or performs an unlock and relock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/lock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlockandrelock
func Lock(app *demoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// TODO: handle un- and relock

	lockID := r.Header.Get(HeaderWopiLock)
	if lockID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := &providerv1beta1.SetLockRequest{
		Ref: &wopiContext.FileReference,
		Lock: &providerv1beta1.Lock{
			LockId:  lockID,
			AppName: app.Config.AppLockName,
			Type:    providerv1beta1.LockType_LOCK_TYPE_WRITE,
			Expiration: &typesv1beta1.Timestamp{
				Seconds: uint64(time.Now().Add(lockDuration).Unix()),
			},
		},
	}

	app.Logger.Debug().Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("Performing SetLock")
	resp, err := app.gwc.SetLock(
		ctx,
		req,
	)
	if err != nil {
		app.Logger.Error().Err(err).Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("SetLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	switch resp.Status.Code {
	case rpcv1beta1.Code_CODE_OK:
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return

	case rpcv1beta1.Code_CODE_FAILED_PRECONDITION, rpcv1beta1.Code_CODE_ABORTED:
		// already locked
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}
		resp, err := app.gwc.GetLock(
			ctx,
			req,
		)
		if err != nil {
			app.Logger.Error().Err(err).Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("SetLock failed, fallback to GetLock failed too")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
			app.Logger.Error().Str("status_code", resp.Status.Code.String()).Str("status_msg", resp.Status.Message).Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("SetLock failed, fallback to GetLock failed too")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if resp.Lock != nil {
			if resp.Lock.LockId != lockID {
				w.Header().Set(HeaderWopiLock, resp.Lock.LockId)
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			// TODO: according to the spec we need to treat this as a RefreshLock

			http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
			return
		}

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return

	default:
		app.Logger.Error().Str("status_code", resp.Status.Code.String()).Str("status_msg", resp.Status.Message).Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("SetLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// RefreshLock refreshes a provided lock for 30 minutes
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/refreshlock
func RefreshLock(app *demoApp, w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// UnLock removes a given lock from a file
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlock
func UnLock(app *demoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	lockID := r.Header.Get(HeaderWopiLock)
	if lockID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := &providerv1beta1.UnlockRequest{
		Ref: &wopiContext.FileReference,
		Lock: &providerv1beta1.Lock{
			LockId:  lockID,
			AppName: app.Config.AppLockName,
		},
	}

	resp, err := app.gwc.Unlock(
		ctx,
		req,
	)
	if err != nil {
		app.Logger.Error().Err(err).Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("UnLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().Str("status_code", resp.Status.Code.String()).Str("status_msg", resp.Status.Message).Str("lock_id", lockID).Str("FileReference", wopiContext.FileReference.String()).Msg("UnLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}
