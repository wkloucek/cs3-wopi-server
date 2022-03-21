FROM owncloudci/golang as build

COPY ./ /wopiserver

WORKDIR /wopiserver

RUN go build ./cmd/cs3-wopi-server

FROM alpine

LABEL maintainer="ownCloud GmbH <devops@owncloud.com>" \
	org.label-schema.name="ownCloud CS3 WOPI server" \
	org.label-schema.vendor="ownCloud GmbH" \
	org.label-schema.schema-version="1.0"

ENTRYPOINT ["/usr/bin/cs3-wopi-server"]

COPY --from=build /wopiserver/cs3-wopi-server /usr/bin/cs3-wopi-server
