package MultisignatureEscrowAccount

import (
	"account"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"log"
)

type DestinationAccount struct {
	account.Account
}

func (da *DestinationAccount) RetainFromEscrow(ea EscrowAccount) error {
	tx, err := build.Transaction(
		build.SourceAccount{ea.Address()},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.Payment(
			build.Destination{da.Address()},
			build.NativeAmount{fundingBalance}))
	account.PanicIfError(err)

	txHash, err := da.SignAndSubmit(tx)
	log.Print(txHash)
	return err
}
