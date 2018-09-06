package JointEntityCrowdfunding

import (
	"account"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"log"
)

type HoldingAccount struct {
	account.Account
}

func (ha *HoldingAccount) RemoveSelfAddSigners(seq build.Sequence, signer1, signer2 string) (err error) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ha.Address()},
		seq,
		build.BaseFee{Amount: build.DefaultBaseFee * 10},
		build.TestNetwork,
		build.SetOptions(
			build.MasterWeight(0),
			build.AddSigner(signer1, 1),
			build.AddSigner(signer2, 1),
			build.SetThresholds(2, 2, 2)),
	)
	if err != nil {
		return err
	}

	txHash, err := ha.SignAndSubmit(tx)
	log.Println("RemoveSelfAddSigners txid:", txHash)
	if err != nil {
		return
	}
	return
}

func (ha *HoldingAccount) CrowdFunding(seq build.Sequence, signers ...string) (err error) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ha.Address()},
		seq,
		build.TestNetwork,
		build.CreateOffer(
			build.Rate{
				Selling: build.Asset{Code: "IMF", Issuer: ha.Address()},
				Buying:  build.NativeAsset(),
				Price:   "1",
			},
			build.Amount("100")),
	)
	if err != nil {
		return
	}

	txe, err := tx.Sign(signers...)
	if err != nil {
		return
	}

	txeB64, _ := txe.Base64()
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	log.Print(resp)
	account.LogIfErrorMsg(err, account.GetResultCodeFromError(err))
	return
}

func (ha *HoldingAccount) PayAsset(d build.Destination, asset build.Asset, amount build.Amount, seq build.Sequence, tm uint64) (err error) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ha.Address()},
		seq,
		build.TestNetwork,
		build.Timebounds{MaxTime: 0, MinTime: tm},
		build.Payment(d, asset, amount),
	)
	if err != nil {
		return
	}

	txHash, err := ha.SignAndSubmit(tx)
	log.Println("Pay to ", d.AddressOrSeed, " txid: ", txHash)
	return
}
