package workbook

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
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
	)

	msig.Threshold = threshold
	msig.Company = rowGroup[0][Company]
	msig.ControlGroup = rowGroup[0][ControlGroup]

	for _, row := range rowGroup {
		if row[ControlGroup] != msig.ControlGroup {
			return fmt.Errorf("control group mismatch")
		}
		_, _, addrBytes, err := address.Parse(row[PChainAddress])
		if err != nil {
			return fmt.Errorf("could not parse address %s for ctrl group %s - err: %s", row[PChainAddress], msig.ControlGroup, err)
		}
		addr, _ := ids.ToShortID(addrBytes)
		msig.Addrs = append(msig.Addrs, addr)
	}
	return nil
}
