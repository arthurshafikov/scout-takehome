FROM golang:1.25.5 AS build

ENV CODE_DIR /go/src/
ENV BIN_FILE /go/src/app
ENV GOOSE_FILE /usr/bin/goose

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD ./cmd ${CODE_DIR}cmd
ADD ./internal ${CODE_DIR}internal
# ADD ./pkg ${CODE_DIR}pkg
ADD ./migrations ${CODE_DIR}migrations

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/app/main.go

COPY . ${CODE_DIR}

ADD https://github.com/pressly/goose/releases/download/v3.20.0/goose_linux_x86_64 ${GOOSE_FILE}
RUN chmod +x ${GOOSE_FILE}

FROM alpine:3.9 AS base

LABEL MAINTAINERS="arthurshafikov"

RUN apk add --no-cache tzdata

ENV BIN_FILE /go/src/app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV MIGRATIONS_FOLDER /go/src/migrations
COPY --from=build ${MIGRATIONS_FOLDER} ${MIGRATIONS_FOLDER}

ENV ENV_FILE /go/src/main.env
COPY --from=build ${ENV_FILE} ${ENV_FILE}

ENV GOOSE_FILE /usr/bin/goose
COPY --from=build ${GOOSE_FILE} ${GOOSE_FILE}

ENV CONFIGS_FOLDER /go/src/configs
COPY --from=build ${CONFIGS_FOLDER} ${CONFIGS_FOLDER}

FROM base AS dev
# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait
CMD /wait && \
    ${BIN_FILE} -cfgFolder ${CONFIGS_FOLDER} -env /go/src/

FROM base AS prod

CMD ${BIN_FILE} -cfgFolder ${CONFIGS_FOLDER} -env /go/src/
