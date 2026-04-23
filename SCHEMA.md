# stellar.token.transfer.v1

Schema reference for the `token-transfer-processor` output event type.

## Event type

- `stellar.token.transfer.v1`

## Batch message

- `stellar.v1.TokenTransferBatch`

## Contained event message

- `stellar.v1.TokenTransferEvent`

## Event variants

Each `TokenTransferEvent` contains one of:

- `transfer`
- `mint`
- `burn`
- `clawback`
- `fee`

## Metadata

Each event includes `meta` with:

- `ledger_sequence`
- `closed_at`
- `tx_hash`
- `transaction_index`
- `operation_index` (optional)
- `contract_address`
- `in_successful_tx`

## Asset encoding

Assets are encoded as one of:

- native XLM
- issued asset
- contract id

## Amount encoding

Amounts are encoded as strings to preserve precision.

## Source of truth

See:

- `proto/stellar/v1/token_transfers.proto`
