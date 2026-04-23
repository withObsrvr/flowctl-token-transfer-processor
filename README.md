# flowctl-token-transfer-processor

A reference `flowctl` processor that extracts Stellar token transfer events from raw ledger events.

This repo is intended to be one of the first flagship reference processors in the nebu → flowctl promotion path:

```text
Prototype in nebu
  ↓
keep the contract stable
  ↓
run it in production with flowctl
```

## What it does

This processor:

- consumes `stellar.ledger.v1`
- uses Stellar's token transfer processor to extract token events
- emits `stellar.token.transfer.v1`
- preserves useful metadata like ledger sequence, tx hash, operation index, contract address, and transaction success

It emits these event variants:

- transfer
- mint
- burn
- clawback
- fee

## Why this repo exists

This processor is the flowctl companion to the token transfer work already proven in:

- nebu
- nebu-processor-registry
- Stellar's `token_transfer` processor pattern

The goal is to provide a production-oriented reference implementation for:

- flowctl processor authors
- example pipelines
- flagship docs
- future registry/discovery work

## Current contract

### Input type

- `stellar.ledger.v1`

### Output type

- `stellar.token.transfer.v1`

### Output protobuf

The protobuf schema currently lives in this repo at:

- `proto/stellar/v1/token_transfers.proto`

This is intentionally close to the token transfer schema already explored in the wider ecosystem. When the shared `flow-proto` version is ready and generated, this processor should migrate to that shared contract.

## Repository layout

```text
flowctl-token-transfer-processor/
├── cmd/token-transfer-processor/main.go
├── internal/convert/convert.go
├── internal/convert/convert_test.go
├── proto/stellar/v1/token_transfers.proto
├── examples/testnet-duckdb-pipeline.yaml
├── processor.yaml
├── Makefile
└── README.md
```

## Build

This repo now targets the published SDK tag:

- `github.com/withObsrvr/flowctl-sdk v0.1.0`

```bash
make proto
make tidy
make build
```

Binary output:

```bash
./bin/token-transfer-processor
```

## Test

```bash
make test
```

## Run locally

```bash
export NETWORK_PASSPHRASE="Test SDF Network ; September 2015"
./bin/token-transfer-processor
```

Optional flowctl integration:

```bash
export ENABLE_FLOWCTL=true
export FLOWCTL_ENDPOINT=localhost:8080
export PORT=:50051
export HEALTH_PORT=8088
./bin/token-transfer-processor
```

## Example flowctl pipeline

See:

- `examples/testnet-duckdb-pipeline.yaml`

The intended chain is:

```text
stellar-live-source
  → token-transfer-processor
  → duckdb-consumer
```

## Implementation notes

This processor uses:

- `github.com/withObsrvr/flowctl-sdk/pkg/stellar`
- `github.com/stellar/go-stellar-sdk/processors/token_transfer`

The wrapper is intentionally thin:

1. decode `stellar.ledger.v1`
2. turn XDR back into `xdr.LedgerCloseMeta`
3. run Stellar token transfer extraction
4. convert the results into the processor's protobuf output
5. emit a single `stellar.token.transfer.v1` batch event per input ledger when events exist

## Metadata included

Each output event includes metadata fields such as:

- `ledger_sequence`
- `closed_at`
- `tx_hash`
- `transaction_index`
- `operation_index`
- `contract_address`
- `in_successful_tx`

## Relationship to nebu

This repo is part of the documented processor promotion path for flowctl.

Related references:

- `flowctl/docs/BUILDING_PROCESSORS.md`
- `flowctl/docs/REFERENCE_PROCESSORS.md`
- `nebu/docs/BUILDING_PROCESSORS.md`
- `nebu-processor-registry/processors/token-transfer`

## Schema

See:

- `SCHEMA.md`
- `proto/stellar/v1/token_transfers.proto`

## Releasing

This repo now uses a fast binary-first release setup:

- `.github/workflows/release-binaries.yml` — runs on version tags and creates GitHub releases with binary assets
- `.github/workflows/publish-image.yml` — manual Docker publishing workflow

### Release binaries

```bash
git tag v0.1.0
git push origin v0.1.0
```

That creates a GitHub Release with assets like:

```text
token-transfer-processor-linux-amd64.tar.gz
token-transfer-processor-linux-arm64.tar.gz
token-transfer-processor-darwin-amd64.tar.gz
token-transfer-processor-darwin-arm64.tar.gz
checksums.txt
```

### Publish a Docker image manually

Use the GitHub Actions UI and run:

- `Publish Docker Image`

Recommended defaults for now:

- `image_tag`: `v0.1.0`
- `platforms`: `linux/amd64`
- `push_latest`: `false`

## CI

This repo also includes a normal CI workflow:

- `.github/workflows/ci.yml`

It verifies:

- proto generation is up to date
- tests pass
- the processor builds

## Next steps

Planned follow-up work:

- add an integration test with a bounded ledger range
- migrate to shared `flow-proto` types once available
- add a companion filter processor like `amount-filter-processor`
- add a sample end-to-end `flowctl run` verification script
