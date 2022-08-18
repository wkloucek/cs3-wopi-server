package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

func DownloadFile(
	ctx context.Context,
	ref *providerv1beta1.Reference,
	gwc gatewayv1beta1.GatewayAPIClient,
	token string,
	insecure bool,
	logger log.Logger,
) (http.Response, error) {

	req := &providerv1beta1.InitiateFileDownloadRequest{
		Ref: ref,
	}

	resp, err := gwc.InitiateFileDownload(ctx, req)
	if err != nil {
		logger.Error().Err(
			err,
		).Str(
			"FileReference", ref.String(),
		).Msg("DownloadHelper: InitiateFileDownload failed")
		return http.Response{}, err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		logger.Error().Str(
			"status_code", resp.Status.Code.String(),
		).Str(
			"status_msg", resp.Status.Message,
		).Str(
			"FileReference", ref.String(),
		).Msg("DownloadHelper: InitiateFileDownload failed")
		return http.Response{}, errors.New("status code != CODE_OK")
	}

	downloadEndpoint := ""
	downloadToken := ""

	for _, proto := range resp.Protocols {
		if proto.Protocol == "simple" || proto.Protocol == "spaces" {
			downloadEndpoint = proto.DownloadEndpoint
			downloadToken = proto.Token
		}
	}

	if downloadEndpoint == "" {
		return http.Response{}, errors.New("download endpoint is missing")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	httpReq, err := http.NewRequest(http.MethodGet, downloadEndpoint, bytes.NewReader([]byte("")))
	if err != nil {
		return http.Response{}, err
	}
	if downloadToken != "" {
		// public link downloads have the token in the download endpoint
		httpReq.Header.Add("X-Reva-Transfer", downloadToken)
	}
	// TODO: the access token shouldn't be needed
	httpReq.Header.Add("x-access-token", token)

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return http.Response{}, err
	}

	return *httpResp, nil
}
