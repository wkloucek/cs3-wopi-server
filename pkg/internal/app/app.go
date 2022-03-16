package app

import (
	"context"
	"errors"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/wkloucek/cs3-wopi-server/pkg/internal/logging"

	registryv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/mime"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
	"google.golang.org/grpc"
)

type GRPC struct {
	Addr     string `env:"WOPI_GRPC_ADDR"`
	BindAddr string `env:"WOPI_GRPC_BIND_ADDR"`
}

type HTTP struct {
	Addr     string `env:"WOPI_HTTP_ADDR"`
	BindAddr string `env:"WOPI_HTTP_BIND_ADDR"`
	Scheme   string `env:"WOPI_HTTP_SCHEME"`
}

type WopiApp struct {
	Addr     string `env:"WOPI_APP_ADDR"`
	Insecure bool   `env:"WOPI_APP_INSECURE"`
}

type CS3api struct {
	Addr                   string `env:"WOPI_CS3API_ADDR"`
	CS3DataGatewayInsecure bool   `env:"WOPI_CS3API_DATA_GATEWAY_INSECURE"`
}

type Config struct {
	GRPC
	HTTP
	WopiApp
	CS3api

	JWTSecret      string `env:"WOPI_JWT_SECRET"`
	AppName        string `env:"WOPI_APP_NAME"`
	AppDescription string `env:"WOPI_APP_DESCRIPTION"`
	AppIcon        string `env:"WOPI_APP_ICON"`
}

type demoApp struct {
	gwc        gatewayv1beta1.GatewayAPIClient
	grpcServer *grpc.Server

	appURLs map[string]map[string]string

	Config Config

	Logger log.Logger
}

func New() (*demoApp, error) {
	app := &demoApp{
		Config: Config{
			AppName:        "WOPI app",
			AppDescription: "Open office documents with a WOPI app",
			AppIcon:        "image-edit",
			JWTSecret:      "test",
			CS3api: CS3api{
				Addr:                   "127.0.0.1:9142",
				CS3DataGatewayInsecure: true,
			},
			GRPC: GRPC{
				Addr:     "127.0.0.1:5678",
				BindAddr: "127.0.0.1:5678",
			},
			HTTP: HTTP{
				Addr:     "172.17.0.1:6789",
				BindAddr: "0.0.0.0:6789",
				Scheme:   "http",
			},
			WopiApp: WopiApp{
				Addr:     "https://localhost:8080",
				Insecure: true,
			},
		},
	}

	err := envdecode.Decode(app)
	if err != nil {
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return nil, err
		}
	}

	app.Logger = logging.Configure("wopiserver")

	return app, nil
}

func (app *demoApp) GetCS3apiClient() error {
	// establish a connection to the cs3 api endpoint
	// in this case a REVA gateway, started by oCIS
	gwc, err := pool.GetGatewayServiceClient(app.Config.CS3api.Addr)
	if err != nil {
		return err
	}
	app.gwc = gwc

	return nil
}

func (app *demoApp) RegisterDemoApp(ctx context.Context) error {
	mimeTypesMap := make(map[string]bool)
	for _, extensions := range app.appURLs {
		for ext := range extensions {
			m := mime.Detect(false, ext)
			mimeTypesMap[m] = true
		}
	}

	mimeTypes := make([]string, 0, len(mimeTypesMap))
	for m := range mimeTypesMap {
		mimeTypes = append(mimeTypes, m)
	}

	// TODO: REVA has way to filter supported mimetypes (do we need to implement it here or is it in the registry?)

	// TODO: an added app provider shouldn't last forever. Instead the registry should use a TTL
	// and delete providers that didn't register again. If an app provider dies or get's disconnected,
	// the users will be no longer available to choose to open a file with it (currently, opening a file just fails)
	req := &registryv1beta1.AddAppProviderRequest{
		Provider: &registryv1beta1.ProviderInfo{
			Name:        app.Config.AppName,
			Description: app.Config.AppDescription,
			Icon:        app.Config.AppIcon,
			Address:     app.Config.GRPC.Addr, // address of the grpc server we start in this demo app
			MimeTypes:   mimeTypes,
		},
	}

	resp, err := app.gwc.AddAppProvider(ctx, req)
	if err != nil {
		app.Logger.Error().Err(err).Msg("AddAppProvider failed")
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().Str("status code", resp.Status.Code.String()).Msg("AddAppProvider failed")
		return errors.New("status code != CODE_OK")
	}

	return nil
}
