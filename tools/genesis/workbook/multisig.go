package workbook

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chain4travel/camino-node/tools/genesis/utils"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
)

type MultiSigGroup struct {
	ControlGroup string
	Threshold    uint32
	PublicKeys   []*secp256k1.PublicKey
}

type MultiSigColumn int

type MultiSigRow struct {
	ControlGroup string
	Threshold    uint32
	PublicKey    *secp256k1.PublicKey
}

func (msig *MultiSigRow) Header() []string { return []string{"Control Group", "Threshold", "Company"} }

func (msig *MultiSigRow) FromRow(_ int, msigRow []string) error {
	// COLUMNS
	const (
		ControlGroup MultiSigColumn = iota
		Threshold
		_Company
		_FirstName
		_LastName
		_Kyc
		_PChainAddress
		PublicKey
	)

	msig.ControlGroup = strings.TrimSpace(msigRow[ControlGroup])
	if msigRow[Threshold] != "" {
		threshold, err := strconv.ParseUint(msigRow[Threshold], 10, 32)
		msig.Threshold = uint32(threshold)
		if err != nil {
			return fmt.Errorf("could not parse msig threshold %s", msigRow[Threshold])
		}
	}

	if len(msigRow) > int(PublicKey) && msigRow[PublicKey] != "" {
		msigRow[PublicKey] = strings.TrimPrefix(strings.TrimSpace(msigRow[PublicKey]), "0x")
		pk, err := utils.PublicKeyFromString(msigRow[PublicKey])
		msig.PublicKey = pk
		if err != nil {
			return fmt.Errorf("could not parse public key")
		}
	} else {
		return fmt.Errorf("empty / invalid public key")
	}

	return nil
}
