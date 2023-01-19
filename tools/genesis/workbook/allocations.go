package workbook

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/chain4travel/camino-node/tools/genesis/utils"
)

const (
	Allocations = "Camino Allocation"
)

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
	Additional1Percent  string
	OfferID             string
	FirstName           string
	TokenDeliveryOffset uint64
	DepositDuration     uint32
	PublicKey           string
}

type AllocationRow int

const (
	TrueValue  = "TRUE"
	FalseValue = "FALSE"
)

func (a *Allocation) FromRow(fileRowNo int, row []string) error {
	// COLUMNS
	const (
		_RowNo AllocationRow = iota
		ComapnyName
		FirstName
		_LastName
		_KnownBy
		Kyc
		_Street
		_Street2
		_Zip
		_City
		_Country
		Bucket
		_CamPurchasePrice
		CamAmount
		PChainAddress
		PublicKey
		ConsortiumMember
		ControlGroup
		MultisigThreshold
		NodeID
		ValidationPeriodDays
		Additional1Percent
		OfferID
		AllocationStartOffet
		DepositDuration
	)

	var err error
	a.RowNo = fileRowNo
	a.Bucket = row[Bucket]
	a.Kyc = strings.TrimSpace(row[Kyc])
	a.FirstName = row[FirstName]
	a.ConsortiumMember = strings.TrimSpace(row[ConsortiumMember])
	a.ControlGroup = strings.TrimSpace(row[ControlGroup])
	a.Additional1Percent = strings.TrimSpace(row[Additional1Percent])
	a.OfferID = strings.TrimSpace(row[OfferID])

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

	keyRead := false
	if row[PublicKey] != "" {
		row[PublicKey] = strings.TrimPrefix(row[PublicKey], "0x")

		pk, err := utils.PublicKeyFromString(row[PublicKey])
		if err != nil {
			return fmt.Errorf("could not parse public key, expected uncompressed bytes %s", row[PublicKey])
		}
		addr, err := utils.ToPAddress(pk)
		if err != nil {
			return fmt.Errorf("[X/P] could not parse public key %s, %w", row[PublicKey], err)
		}

		a.Address, keyRead = addr, true
		a.PublicKey = row[PublicKey]
	}
	if !keyRead && len(row[PChainAddress]) >= 47 {
		_, _, addrBytes, err := address.Parse(strings.TrimSpace(row[PChainAddress]))
		if err != nil {
			return fmt.Errorf("could not parse address %s", row[PChainAddress])
		}
		a.Address, _ = ids.ToShortID(addrBytes)
	}

	a.TokenDeliveryOffset, err = strconv.ParseUint(row[AllocationStartOffet], 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse allocation offset %s", row[AllocationStartOffet])
	}

	dd, err := strconv.ParseUint(row[DepositDuration], 10, 32)
	if row[DepositDuration] != "" && err != nil {
		return fmt.Errorf("could not parse deposit duration %s: %w", row[DepositDuration], err)
	}
	a.DepositDuration = uint32(dd)

	return nil
}
