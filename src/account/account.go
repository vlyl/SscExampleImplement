package account

import (
	"fmt"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"log"
	//protoHorizon "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/clients/horizon"
)

type Account struct{
	seed string
}

func NewAccount(seed string) (a Account) {
	return Account{seed:seed}
}

func (a *Account) Address() string {
	kp, err := keypair.Parse(a.seed)
	PanicIfError(err)
	return kp.Address()
}
func (a *Account) SignMessage(input []byte) []byte {
	if len(a.seed) == 0 {
		return nil
	}

	kp, err := keypair.Parse(a.seed)
	PanicIfError(err)

	signature, err := kp.Sign(input)
	PanicIfError(err)
	return signature
}

func (a *Account) SignTx(tx *build.TransactionBuilder) (txe build.TransactionEnvelopeBuilder) {
	txe, err := tx.Sign(a.seed)
	PanicIfError(err)
	return
}

func (a *Account) SignTxe(txe *build.TransactionEnvelopeBuilder) {
	txe.Mutate(build.Sign{Seed:a.seed})
}

func SubmitTxe(txe build.TransactionEnvelopeBuilder) (txHash string, err error) {
	txeB64, _ := txe.Base64()
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	log.Print(resp)
	LogIfErrorMsg(err, GetResultCodeFromError(err))
	return resp.Hash, nil
}

func (a *Account) SignAndSubmit(tx *build.TransactionBuilder) (txHash string, err error) {
	return SubmitTxe(a.SignTx(tx))
}

func (a *Account) GetSequence() build.Sequence {
	xdrSeq, err := horizon.DefaultTestNetClient.SequenceForAccount(a.Address())
	PanicIfError(err)
	return build.Sequence{Sequence: uint64(xdrSeq)}
}

func SequenceIncrement(seq build.Sequence) build.Sequence {
	ns := seq.Sequence + 1
	return build.Sequence{Sequence: ns}
}

func GetResultCodeFromError(err error) string {
	herr, isHorizonError := err.(*horizon.Error)
	if isHorizonError {
		resultCodes, err := herr.ResultCodes()
		if err != nil {
			fmt.Println("failed to extract result codes from horizon response")
			return ""
		}
		return resultCodes.TransactionCode
	}
	return ""
}

func PanicIfError(e error) {
	PanicErrorMsg(e, "")
}

func PanicErrorMsg(e error, msg string) {
	if e != nil {
		panic(e.Error() + "\n" + msg)
	}
}

func LogError(e error) {
	LogIfErrorMsg(e, "")
}

func LogIfErrorMsg(e error, msg string) {
	if e != nil {
		log.Println(e)
		log.Println(msg)
	}
}
