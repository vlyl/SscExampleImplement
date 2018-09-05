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

func (ha *HoldingAccount) RemoveSelfAddSigners(seq build.Sequence, signers ...string) (err error) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ha.Address()},
		seq,
		build.TestNetwork,
		build.RemoveSigner(ha.Address()),
	)
	if err != nil {
		return err
	}
	for _, signer := range signers {
		err = tx.Mutate(build.SetOptions(build.AddSigner(signer, 1)))
		if err != nil {
			return
		}
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
			build.Rate{Selling: build.Asset{Code: "IMF", Issuer: ha.Address()},
				Buying: build.NativeAsset(),
				Price:  "1",
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
}
