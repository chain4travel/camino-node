package workbook

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/chain4travel/camino-node/tools/genesis/utils"
)

const MultisigDefinitions = "MultiSig Addresses"

type MultiSig struct {
	ControlGroup string
	Company      string
	Threshold    uint32
	Addrs        []ids.ShortID
}

type MultiSigRow int

func (msig *MultiSig) FromRow(threshold uint32, rowGroup [][]string) error {
	// COLUMNS
	const (
		ControlGroup MultiSigRow = iota
		Company
		FirstName
		LastName
		Kyc
		PChainAddress
		PublicKey
	)

	msig.Threshold = threshold
	msig.Company = rowGroup[0][Company]
	msig.ControlGroup = rowGroup[0][ControlGroup]

	for _, row := range rowGroup {
		if row[ControlGroup] != msig.ControlGroup {
			return fmt.Errorf("control group mismatch")
		}

		keyRead := false
		var addr ids.ShortID
		if row[PublicKey] != "" {
			row[PublicKey] = strings.TrimPrefix(row[PublicKey], "0x")

			pk, err := utils.PublicKeyFromString(row[PublicKey])
			if err != nil {
				return fmt.Errorf("could not parse public key, expected uncompressed bytes %s", row[PublicKey])
			}
			addr, err = utils.ToPAddress(pk)
			if err != nil {
				return fmt.Errorf("[X/P] could not parse public key %s, %w", row[PublicKey], err)
			}

			keyRead = true
		}
		if !keyRead && len(row[PChainAddress]) >= 47 {
			_, _, addrBytes, err := address.Parse(strings.TrimSpace(row[PChainAddress]))
			if err != nil {
				return fmt.Errorf("could not parse address %s for ctrl group %s - err: %s", row[PChainAddress], msig.ControlGroup, err)
			}
			addr, _ = ids.ToShortID(addrBytes)
		}
		msig.Addrs = append(msig.Addrs, addr)
	}
	return nil
}
