package account

import (
	"context"
	"encoding/json"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestMultisignatureEscrowAccount(t *testing.T) {
	// source account
	seed := "SBMSM5WYMB4ULOL72FRLZTGHEJ3SPXUFASKGWXB7DOSFLJMBBKFEOAKI"
	a := NewAccount(seed)
	sa := SourceAccount{a}

	// destination account
	seed = "SDEDX4OS5KFPAYIBYD77IKJRZ6UNV7DFZZ4XFJQH6ENTF5PFNT6AGXYJ"
	a = NewAccount(seed)
	da := DestinationAccount{a}

	// sequence number
	// N, M - sequence number of escrow account and source account
	seqM := sa.GetSequence()
	// T - the lock-up period
	// D - the date upon which the lock-up period starts
	// R - the recovery period

	/// Tx1: Create Escrow Account
	startingBalance := "1"
	ea := sa.CreateEscrowAccount(startingBalance, seqM)
	eaNet, err := horizon.DefaultTestNetClient.LoadAccount(ea.Address())
	if err != nil || (
		func () int {n, _ := strconv.Atoi(eaNet.Balances[0].Balance); return n}() ==
			func () int {n, _ := strconv.Atoi(startingBalance); return n}()) {
		t.Fail()
	}

	seqN := ea.GetSequence()

	/// Tx3: UnlockPreBuild
	txUnlock :=ea.UnlockPreBuild(MinuteLater(time.Now()), SequenceIncrement(seqN))
	/// Tx4: RecoveryPreBuild
	txRecovery := ea.RecoveryPreBuild(da, SequenceIncrement(seqN))



	/// Tx2: Enabling Multi-sig
	ea.EnableMultiSig(da, seqN)

	// Tx5: Funding to destination account
	sa.Funding(da, SequenceIncrement(seqM))


}

func TestAccount_SignAndSubmit(t *testing.T) {
	seed := "SBMSM5WYMB4ULOL72FRLZTGHEJ3SPXUFASKGWXB7DOSFLJMBBKFEOAKI"
	a := NewAccount(seed)
	tx, _ := build.Transaction(
		build.SourceAccount{a.Address()},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.SetOptions(),
		)

	txHash, err := a.SignAndSubmit(tx)
	if err == nil {
		log.Print(txHash)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		horizon.DefaultTestNetClient.StreamTransactions(ctx, a.Address(), nil, func (t horizon.Transaction) {
			log.Println("Tx:")
			content, _ := json.Marshal(t)
			log.Print(string(content))
		},
		)
	}
}