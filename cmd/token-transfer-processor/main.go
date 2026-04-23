package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stellar/go-stellar-sdk/ingest"
	tokensdk "github.com/stellar/go-stellar-sdk/processors/token_transfer"
	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/withObsrvr/flowctl-sdk/pkg/stellar"
	"github.com/withObsrvr/flowctl-token-transfer-processor/internal/convert"
	stellarv1 "github.com/withObsrvr/flowctl-token-transfer-processor/proto/stellar/v1"
	"google.golang.org/protobuf/proto"
)

var (
	Version   = "dev"
	CommitSHA = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("token-transfer-processor %s (%s)\n", Version, CommitSHA)
		return
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	stellar.Run(stellar.ProcessorConfig{
		ProcessorName: "Token Transfer Processor",
		OutputType:    "stellar.token.transfer.v1",
		ProcessLedger: func(passphrase string, ledger xdr.LedgerCloseMeta) (proto.Message, error) {
			processor := tokensdk.NewEventsProcessor(passphrase)
			events, err := processor.EventsFromLedger(ledger)
			if err != nil {
				return nil, err
			}

			if len(events) == 0 {
				return nil, nil
			}

			txSuccessMap, err := buildTxSuccessMap(passphrase, ledger)
			if err != nil {
				return nil, err
			}

			batch := &stellarv1.TokenTransferBatch{
				Events: make([]*stellarv1.TokenTransferEvent, 0, len(events)),
			}

			for _, event := range events {
				successful := true
				if event.Meta != nil {
					if found, ok := txSuccessMap[event.Meta.TxHash]; ok {
						successful = found
					}
				}

				converted := convert.Event(event, successful)
				if converted != nil {
					batch.Events = append(batch.Events, converted)
				}
			}

			if len(batch.Events) == 0 {
				return nil, nil
			}

			log.Printf("found %d token transfer events in ledger %d", len(batch.Events), ledger.LedgerSequence())
			return batch, nil
		},
	})
}

func buildTxSuccessMap(passphrase string, ledger xdr.LedgerCloseMeta) (map[string]bool, error) {
	reader, err := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(passphrase, ledger)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	result := make(map[string]bool)
	for {
		tx, err := reader.Read()
		if err != nil {
			break
		}
		result[tx.Result.TransactionHash.HexString()] = tx.Result.Successful()
	}

	return result, nil
}
