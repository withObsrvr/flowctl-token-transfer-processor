package convert

import (
	"testing"
	"time"

	"github.com/stellar/go-stellar-sdk/asset"
	tokensdk "github.com/stellar/go-stellar-sdk/processors/token_transfer"
	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEvent_Transfer(t *testing.T) {
	const issuer = "GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN"

	opIndex := uint32(2)
	sdkEvent := &tokensdk.TokenTransferEvent{
		Meta: &tokensdk.EventMeta{
			LedgerSequence:   123,
			ClosedAt:         timestamppb.New(time.Unix(1700000000, 0)),
			TxHash:           "abc123",
			TransactionIndex: 1,
			OperationIndex:   &opIndex,
			ContractAddress:  "CA123",
		},
		Event: &tokensdk.TokenTransferEvent_Transfer{Transfer: &tokensdk.Transfer{
			From:   "GAFROM",
			To:     "GATO",
			Asset:  asset.NewProtoAsset(xdr.MustNewCreditAsset("USDC", issuer)),
			Amount: "1000000",
		}},
	}

	converted := Event(sdkEvent, true)
	require.NotNil(t, converted)
	require.NotNil(t, converted.Meta)
	require.True(t, converted.Meta.InSuccessfulTx)
	require.Equal(t, uint32(123), converted.Meta.LedgerSequence)
	require.Equal(t, "abc123", converted.Meta.TxHash)
	require.Equal(t, uint32(2), converted.Meta.GetOperationIndex())

	transfer := converted.GetTransfer()
	require.NotNil(t, transfer)
	require.Equal(t, "GAFROM", transfer.From)
	require.Equal(t, "GATO", transfer.To)
	require.Equal(t, "1000000", transfer.Amount)
	require.Equal(t, "USDC", transfer.Asset.GetIssued().AssetCode)
	require.Equal(t, issuer, transfer.Asset.GetIssued().AssetIssuer)
}

func TestEvent_NativeFee(t *testing.T) {
	sdkEvent := &tokensdk.TokenTransferEvent{
		Meta: &tokensdk.EventMeta{
			LedgerSequence:   124,
			ClosedAt:         timestamppb.New(time.Unix(1700000001, 0)),
			TxHash:           "def456",
			TransactionIndex: 3,
		},
		Event: &tokensdk.TokenTransferEvent_Fee{Fee: &tokensdk.Fee{
			From:   "GAFROM",
			Asset:  asset.NewNativeAsset(),
			Amount: "100",
		}},
	}

	converted := Event(sdkEvent, false)
	require.NotNil(t, converted)
	require.False(t, converted.Meta.InSuccessfulTx)

	fee := converted.GetFee()
	require.NotNil(t, fee)
	require.Equal(t, "GAFROM", fee.From)
	require.Equal(t, "100", fee.Amount)
	require.True(t, fee.Asset.GetNative())
}

func TestEvent_Nil(t *testing.T) {
	require.Nil(t, Event(nil, true))
	require.Nil(t, Event(&tokensdk.TokenTransferEvent{}, true))
}
