package convert

import (
	"github.com/stellar/go-stellar-sdk/asset"
	tokensdk "github.com/stellar/go-stellar-sdk/processors/token_transfer"
	stellarv1 "github.com/withObsrvr/flowctl-token-transfer-processor/proto/stellar/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Event(sdkEvent *tokensdk.TokenTransferEvent, inSuccessfulTx bool) *stellarv1.TokenTransferEvent {
	if sdkEvent == nil || sdkEvent.Meta == nil {
		return nil
	}

	meta := &stellarv1.TokenTransferEventMeta{
		LedgerSequence:   sdkEvent.Meta.LedgerSequence,
		ClosedAt:         timestamppb.New(sdkEvent.Meta.ClosedAt.AsTime()),
		TxHash:           sdkEvent.Meta.TxHash,
		TransactionIndex: sdkEvent.Meta.TransactionIndex,
		ContractAddress:  sdkEvent.Meta.ContractAddress,
		InSuccessfulTx:   inSuccessfulTx,
	}
	if sdkEvent.Meta.OperationIndex != nil {
		meta.OperationIndex = sdkEvent.Meta.OperationIndex
	}

	out := &stellarv1.TokenTransferEvent{Meta: meta}

	switch ev := sdkEvent.Event.(type) {
	case *tokensdk.TokenTransferEvent_Transfer:
		out.Event = &stellarv1.TokenTransferEvent_Transfer{Transfer: &stellarv1.Transfer{
			From:   ev.Transfer.From,
			To:     ev.Transfer.To,
			Asset:  assetToProto(ev.Transfer.Asset),
			Amount: ev.Transfer.Amount,
		}}
	case *tokensdk.TokenTransferEvent_Mint:
		out.Event = &stellarv1.TokenTransferEvent_Mint{Mint: &stellarv1.Mint{
			To:     ev.Mint.To,
			Asset:  assetToProto(ev.Mint.Asset),
			Amount: ev.Mint.Amount,
		}}
	case *tokensdk.TokenTransferEvent_Burn:
		out.Event = &stellarv1.TokenTransferEvent_Burn{Burn: &stellarv1.Burn{
			From:   ev.Burn.From,
			Asset:  assetToProto(ev.Burn.Asset),
			Amount: ev.Burn.Amount,
		}}
	case *tokensdk.TokenTransferEvent_Clawback:
		out.Event = &stellarv1.TokenTransferEvent_Clawback{Clawback: &stellarv1.Clawback{
			From:   ev.Clawback.From,
			Asset:  assetToProto(ev.Clawback.Asset),
			Amount: ev.Clawback.Amount,
		}}
	case *tokensdk.TokenTransferEvent_Fee:
		out.Event = &stellarv1.TokenTransferEvent_Fee{Fee: &stellarv1.Fee{
			From:   ev.Fee.From,
			Asset:  assetToProto(ev.Fee.Asset),
			Amount: ev.Fee.Amount,
		}}
	default:
		return nil
	}

	return out
}

func assetToProto(a *asset.Asset) *stellarv1.Asset {
	if a == nil {
		return nil
	}

	if a.GetNative() {
		return &stellarv1.Asset{Asset: &stellarv1.Asset_Native{Native: true}}
	}

	if issued := a.GetIssuedAsset(); issued != nil {
		return &stellarv1.Asset{Asset: &stellarv1.Asset_Issued{Issued: &stellarv1.IssuedAsset{
			AssetCode:   issued.GetAssetCode(),
			AssetIssuer: issued.GetIssuer(),
		}}}
	}

	return nil
}
