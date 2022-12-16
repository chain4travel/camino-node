package workbook

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
)

const WB_MSIG_NAME = "MultiSig Addresses"

type MultiSig struct {
	ControlGroup string
	Company      string
	Threshold    uint32
	Addrs        []ids.ShortID
}

func (msig *MultiSig) FromRow(threshold uint32, row_group [][]string) error {
	// COLUMNS
	const (
		CONTROL_GROUP = iota
		COMPANY
		FIRST_NAME
		LAST_NAME
		KYC
		P_CHAIN_ADDRESS
	)

	msig.Threshold = threshold
	msig.Company = row_group[0][COMPANY]
	msig.ControlGroup = row_group[0][CONTROL_GROUP]

	for _, row := range row_group {
		if row[CONTROL_GROUP] != msig.ControlGroup {
			return fmt.Errorf("control group mismatch")
		}
		_, _, addrBytes, err := address.Parse(row[P_CHAIN_ADDRESS])
		if err != nil {
			return fmt.Errorf("could not parse address %s for ctrl group %s - err: %s", row[P_CHAIN_ADDRESS], msig.ControlGroup, err)
		}
		addr, _ := ids.ToShortID(addrBytes)
		msig.Addrs = append(msig.Addrs, addr)
	}
	return nil
}
