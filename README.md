# CS3api WOPI server

This project has been moved to https://github.com/owncloud/cs3-wopi-server

## development

1. start one office from `./dev-dependencies/Collabora` or `./dev-dependencies/OnlyOffice` by running `docker compose up -d` in that directory
2. run oCIS, please see https://github.com/owncloud/ocis. Please note that oCIS needs to use a shared service registry with the WOPi server. You can eg. use `"MICRO_REGISTRY": "mdns",` (see `./.vscode/launch.json`).
3. run WOPI server in debug mode by using the VSCode run target or by executing `go run ./cmd/cs3-wopi-server`

4. check that the office is running: http://localhost:8080/hosting/discovery
5. log in into oCIS
6. open / create office files
