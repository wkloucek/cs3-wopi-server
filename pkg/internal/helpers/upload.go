package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

func UploadFile(
	ctx context.Context,
	content *bytes.Reader,
	ref *providerv1beta1.Reference,
	gwc gatewayv1beta1.GatewayAPIClient,
	token string,
	lockID string,
	insecure bool,
	logger log.Logger,
) error {

	req := &providerv1beta1.InitiateFileUploadRequest{
		Ref:    ref,
		LockId: lockID,
		// TODO: if-match
		//Options: &providerv1beta1.InitiateFileUploadRequest_IfMatch{
		//	IfMatch: "",
		//},
	}

	resp, err := gwc.InitiateFileUpload(ctx, req)
	if err != nil {
		logger.Error().Err(
			err,
		).Str(
			"FileReference", ref.String(),
		).Msg("UploadHelper: InitiateFileUpload failed")
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		return errors.New("status code != CODE_OK")
	}

	uploadEndpoint := ""
	uploadToken := ""

	for _, proto := range resp.Protocols {
		if proto.Protocol == "simple" || proto.Protocol == "spaces" {
			uploadEndpoint = proto.UploadEndpoint
			uploadToken = proto.Token
		}
	}

	if uploadEndpoint == "" {
		return errors.New("upload endpoint or token is missing")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	httpReq, err := http.NewRequest(http.MethodPut, uploadEndpoint, content)
	if err != nil {
		return err
	}

	if uploadToken != "" {
		// public link uploads have the token in the upload endpoint
		httpReq.Header.Add("X-Reva-Transfer", uploadToken)
	}
	// TODO: the access token shouldn't be needed
	httpReq.Header.Add("x-access-token", token)

	// TODO: better mechanism for the upload while locked, relies on patch in REVA
	//if lockID, ok := ctxpkg.ContextGetLockID(ctx); ok {
	//	httpReq.Header.Add("X-Lock-Id", lockID)
	//}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		return errors.New("status code was not 200")
	}

	return nil
}
