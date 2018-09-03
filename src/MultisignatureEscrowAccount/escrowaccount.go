package MultisignatureEscrowAccount

import (
	"github.com/stellar/go/build"
	"log"
	"time"
)

type EscrowAccount struct {
	Account
}

func (ea *EscrowAccount) AddSigner(da DestinationAccount, seq build.Sequence) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ea.Address()},
		seq,
		build.TestNetwork,
		build.SetOptions(build.Signer{da.Address(), 1}),
		build.SetOptions(
			build.MasterWeight(1),
			build.SetThresholds(2, 2, 2),
		),
		build.MemoText{"Add signer"},
	)
	if err != nil {
		log.Print(err)
		return
	}

	_, err = ea.SignAndSubmit(tx)
	if err != nil {
		log.Print(err)
		return
	}
}

func YearLater(tm time.Time) uint64 {
	return uint64(tm.Add(24 * 365 * time.Hour).Unix())
}

func MinuteLater(tm time.Time) uint64 {
	return uint64(tm.Add(time.Minute).Unix())
}

func (ea *EscrowAccount) UnlockPreBuild(tm uint64, seq build.Sequence) (tx *build.TransactionBuilder) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ea.Address()},
		seq,
		build.TestNetwork,
		build.Timebounds{MaxTime: tm, MinTime: 0},
		build.SetOptions(build.MasterWeight(0), build.SetThresholds(1, 1, 1)),
		build.MemoText{"Unlock Pre Build"},
	)
	PanicIfError(err)
	return
}

func (ea *EscrowAccount) RecoveryPreBuild(da DestinationAccount, seq build.Sequence) (tx *build.TransactionBuilder) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ea.Address()},
		seq,
		build.TestNetwork,
		build.SetOptions(build.RemoveSigner(da.Address())),
		build.SetOptions(
			build.MasterWeight(1),
			build.SetThresholds(1, 1, 1),
		),
	)
	PanicIfError(err)
	return
}
