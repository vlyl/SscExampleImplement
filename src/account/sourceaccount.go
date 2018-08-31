package account

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"log"
)

type SourceAccount struct {
	Account
}

func (sa *SourceAccount) CreateEscrowAccount(startingBalance string, seq build.Sequence) (ea EscrowAccount) {
	full, err := keypair.Random()
	if err != nil {
		log.Panic(err)
		return
	}

	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: sa.Address()},
		seq,
		build.TestNetwork,
		build.CreateAccount(
			build.Destination{AddressOrSeed: full.Address()},
			build.Amount(startingBalance),
		),
	)
	if err != nil {
		log.Print(err)
		return
	}

	_, err = sa.SignAndSubmit(tx)
	if err != nil {
		log.Print(err)
		return
	}

	return EscrowAccount{NewAccount(full.Seed())}
}

func (sa *SourceAccount) Funding(da DestinationAccount, seq build.Sequence) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: sa.Address()},
		seq,
		build.TestNetwork,
		build.Payment(build.NativeAsset(), build.NativeAmount{"1000"}),
	)
	PanicIfError(err)

	txHash, err := sa.SignAndSubmit(tx)
	PanicIfError(err)
	log.Println("txid:", txHash)
}