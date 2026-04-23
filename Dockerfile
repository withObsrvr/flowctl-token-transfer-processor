FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG COMPONENT_NAME=token-transfer-processor
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/${COMPONENT_NAME} \
    ./cmd/${COMPONENT_NAME}

FROM gcr.io/distroless/base-debian12

ARG COMPONENT_NAME=token-transfer-processor
COPY --from=builder /out/${COMPONENT_NAME} /app/${COMPONENT_NAME}
COPY --from=builder /build/processor.yaml /app/processor.yaml

LABEL org.opencontainers.image.title="Token Transfer Processor"
LABEL org.opencontainers.image.description="Extracts Stellar token transfer events from Stellar ledgers"
LABEL org.opencontainers.image.source="https://github.com/withObsrvr/flowctl-token-transfer-processor"
LABEL org.opencontainers.image.documentation="https://github.com/withObsrvr/flowctl-token-transfer-processor/blob/main/README.md"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.authors="withObsrvr"
LABEL io.flowctl.component.type="processor"
LABEL io.flowctl.component.api-version="v1"
LABEL io.flowctl.component.capabilities="streaming,batch"

USER nobody:nobody
WORKDIR /app
ENTRYPOINT ["/app/token-transfer-processor"]
CMD []
