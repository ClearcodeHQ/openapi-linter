# Build Stage
FROM golang:1.13.4 AS build-stage

LABEL app="build-openapi-linter"
LABEL REPO="https://github.com/clearcodehq/openapi-linter"

ENV PROJPATH=/go/src/github.com/clearcodehq/openapi-linter

# Because of https://bitbucket.org/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/clearcodehq/openapi-linter
WORKDIR /go/src/github.com/clearcodehq/openapi-linter

RUN make build-alpine

# Final Stage
FROM alpine:3.10.3

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/clearcodehq/openapi-linter"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://bitbucket.org/docker/docker/issues/14914
ENV PATH=$PATH:/opt/openapi-linter/bin

WORKDIR /opt/openapi-linter/bin

COPY --from=build-stage /go/src/github.com/clearcodehq/openapi-linter/bin/openapi-linter /opt/openapi-linter/bin/
RUN chmod +x /opt/openapi-linter/bin/openapi-linter

# Create appuser
RUN adduser -D -g '' openapi-linter
USER openapi-linter

CMD ["/opt/openapi-linter/bin/openapi-linter"]
