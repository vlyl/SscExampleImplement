package JointEntityCrowdfunding

import (
	"account"
	"github.com/stellar/go/build"
	"testing"
)

func TestJointEntityCrowdfunding(t *testing.T) {
	// GATQAK52W2NQVRLA65U6FERLG6IJ3LLMPSLN5JZ32XX7LTMXCOIOT6VI
	A := account.NewAccount("SDUZU4ERQ5SYV5YYRYC2NTCP3IPJOFBVHEFMRDJIH4PJKCEEGKZ4BPOR")
	// GDAGKGKT5FMNRMAHQ3QBUV66CMZLL3TPGGGFYTGJEMV4VJH3SRKVPO3F
	B := account.NewAccount("SBFP73TFKMD4VNU5HQQEFPQ7RWCU7Y3M6HSCAD2M7IHILS6XMVT2I6UT")
	// buying account
	// GCT7LYQMVOEFZQZNSLIZ2HEIFCBC3OCEB6MCMOF6PZK7WFP2UILCBKDY
	C := account.NewAccount("SA7QRAZWPGDYDN26A3L5ZPQO3DV36QH4TPVZ2AP5XVJLD3Q7CGQK3BE6")
	// target account
	// GAMYTGCOEW25CKACOXSDZZMCJE6GWQIIG43M7TZTSDLB2CPPGJUYKBQG
	D := account.NewAccount("SD5WAPJYQJ5URL4CSU4G736YMSVKU3LAGINKUYXTJHBM7AKPDU47HMVB")

	seqM := account.SequenceIncrement(A.GetSequence())
	// tx1: create holding account
	ha, err := A.CreateNewAccount("10", seqM)
	if err != nil {
		t.Error(err)
	}

	seqN := account.SequenceIncrement(ha.GetSequence())
	// tx2: remove holding account from itself signer, add A and B as holding account singer
	err = HoldingAccount{Account: ha}.RemoveSelfAddSigners(seqN, A.Address(), B.Address())
	if err != nil {
		t.Error(err)
	}

	// tx3: begin crowd funding
	err = HoldingAccount{Account: ha}.CrowdFunding(account.SequenceIncrement(seqN), A.Address(), B.Address())
	if err != nil {
		t.Error(err)
	}

	err = C.TrustAsset(build.Asset{Code: "IMF", Issuer: ha.Address(), Native: false})
	if err != nil {
		t.Error(err)
	}

}
