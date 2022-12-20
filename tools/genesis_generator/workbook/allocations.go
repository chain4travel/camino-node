package workbook

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/units"
)

const Allocations = "Camino Allocation"

type Allocation struct {
	RowNo               int
	Bucket              string
	Kyc                 string
	Amount              uint64
	Address             ids.ShortID
	ConsortiumMember    string
	ControlGroup        string
	MsigThreshold       uint32
	NodeID              ids.NodeID
	ValidatorPeriodDays uint32
	UnbondingStart      float64
	UnbondingPeriod     float64
	Additional1Percent  string
	RewardPercent       int
	FirstName           string
}

type AllocationRow int

func (a *Allocation) FromRow(row []string) error {
	// COLUMNS
	const (
		RowNo AllocationRow = iota
		ComapnyName
		FirstName
		LastName
		KnownBy
		Kyc
		Street
		Street2
		Zip
		City
		Country
		Bucket
		CamPurchasePrice
		CamAmount
		PChainAddress
		ConsortiumMember
		ControlGroup
		MultisigThreshold
		NodeID
		ValidationPeriodDays
		UnbondingStartYears
		UnbondingPeriodYears
		Additional1Percent
		RewardPercent
	)

	var err error
	a.Bucket = row[Bucket]
	a.Kyc = row[Kyc]
	a.FirstName = row[FirstName]
	a.ConsortiumMember = row[ConsortiumMember]
	a.ControlGroup = row[ControlGroup]
	a.Additional1Percent = row[Additional1Percent]

	a.RowNo, err = strconv.Atoi(row[RowNo])
	if err != nil {
		return fmt.Errorf("could not parse row number %s", row[RowNo])
	}

	a.Amount, err = strconv.ParseUint(row[CamAmount], 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse amount %s", row[CamAmount])
	}
	a.Amount *= units.Avax

	if row[MultisigThreshold] != "" {
		threshold, err := strconv.ParseUint(row[MultisigThreshold], 10, 32)
		a.MsigThreshold = uint32(threshold)
		if err != nil {
			return fmt.Errorf("could not parse msig threshold %s", row[MultisigThreshold])
		}
	}

	row[NodeID] = strings.TrimSpace(row[NodeID])
	if row[NodeID] != "" && row[NodeID] != "X" {
		a.NodeID, err = ids.NodeIDFromString(row[NodeID])
		if err != nil {
			fmt.Println("could not parse node id", row[NodeID])
		}
	}

	if row[ValidationPeriodDays] != "" {
		vpd, err := strconv.ParseUint(row[ValidationPeriodDays], 10, 32)
		a.ValidatorPeriodDays = uint32(vpd)
		if err != nil {
			fmt.Println("could not parse Validator Period: ", row[ValidationPeriodDays])
		}
	}

	if len(row[PChainAddress]) >= 47 {
		_, _, addrBytes, err := address.Parse(row[PChainAddress])
		if err != nil {
			return fmt.Errorf("could not parse address %s", row[PChainAddress])
		}
		a.Address, _ = ids.ToShortID(addrBytes)
	}

	a.UnbondingStart, err = strconv.ParseFloat(row[UnbondingStartYears], 64)
	if err != nil {
		a.UnbondingStart = 0
	}
	a.UnbondingPeriod, err = strconv.ParseFloat(row[UnbondingPeriodYears], 64)
	if err != nil {
		a.UnbondingPeriod = 0
	}

	if strings.HasSuffix(row[RewardPercent], "%") {
		a.RewardPercent, err = strconv.Atoi(row[RewardPercent][:len(row[RewardPercent])-1])
		if err != nil {
			a.RewardPercent = -1
		}
	} else {
		a.RewardPercent = -1
	}

	return nil
}
