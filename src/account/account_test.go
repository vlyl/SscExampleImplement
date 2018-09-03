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
	seqM := SequenceIncrement(sa.GetSequence())
	// T - the lock-up period
	// D - the date upon which the lock-up period starts
	// R - the recovery period

	/// Tx1: Create Escrow Account
	startingBalance := "10"
	ea := sa.CreateEscrowAccount(startingBalance, seqM)
	eaNet, err := horizon.DefaultTestNetClient.LoadAccount(ea.Address())
	if err != nil || (func() int { n, _ := strconv.Atoi(eaNet.Balances[0].Balance); return n }() ==
		func() int { n, _ := strconv.Atoi(startingBalance); return n }()) {
		t.Fail()
	}

	seqN := SequenceIncrement(ea.GetSequence())

	/// Tx2: Enabling Multi-sig
	ea.AddSigner(da, seqN)

	/// Tx3: UnlockPreBuild
	txUnlock := ea.UnlockPreBuild(MinuteLater(time.Now()), SequenceIncrement(seqN))
	// sign tx3(txUnlock) with escrow account
	txeUnlock := ea.SignTx(txUnlock)
	// sign tx3(txUnlock) with destination account
	da.SignTxe(&txeUnlock)

	/// Tx4: RecoveryPreBuild
	txRecovery := ea.RecoveryPreBuild(da, SequenceIncrement(seqN))
	// sign tx4(txRecovery) with escrow account
	txeRecovery := ea.SignTx(txRecovery)
	// sign tx4(txRecovery) with destination account
	da.SignTxe(&txeRecovery)

	// Tx5: Funding to escrow account
	sa.Funding(ea, SequenceIncrement(seqM))

	escrow, err := horizon.DefaultTestNetClient.LoadAccount(ea.Address())
	log.Print(time.Now())
	log.Println("escrow balance: ")
	eb, _ := json.Marshal(escrow.Balances)
	log.Print(string(eb))

	log.Println("sleeping S...")
	time.Sleep(time.Minute / 2)

	log.Println("submit unlock tx")
	txHash, err := SubmitTxe(txeUnlock)
	if err != nil {
		log.Print(err)
	}
	log.Print("unlock txid:", txHash)
	//resp, err := http.Get("https://horizon-testnet.stellar.org/transactions/"+txHash)
	//if err != nil {
	//	log.Print(err)
	//}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Print(err)
	//}
	//log.Print(string(body))

	escrow, err = horizon.DefaultTestNetClient.LoadAccount(ea.Address())
	log.Println("escrow balance: ")
	escrowJson, _ := json.Marshal(escrow.Balances)
	log.Print(string(escrowJson))

	da.RetainFromEscrow(ea)

	escrow, err = horizon.DefaultTestNetClient.LoadAccount(ea.Address())
	log.Println("escrow balance after target retain: ")
	escrowJson, _ = json.Marshal(escrow.Balances)
	log.Print(string(escrowJson))
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
		horizon.DefaultTestNetClient.StreamTransactions(ctx, a.Address(), nil, func(t horizon.Transaction) {
			log.Println("Tx:")
			content, _ := json.Marshal(t)
			log.Print(string(content))
		},
		)
	}
}
