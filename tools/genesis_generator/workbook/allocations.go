package workbook

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/units"
)

const WB_ALLOCATIONS_NAME = "Camino Allocation"

type Allocation struct {
	RowNo               int
	Bucket              string
	Kyc                 string
	Amount              uint64
	Address             ids.ShortID
	ConsortiumMember    string
	ControlGroup        string
	MsigThreshold       uint32
	NodeId              ids.NodeID
	ValidatorPeriodDays uint32
	UnbondingStart      float64
	UnbondingPeriod     float64
	Additional1Percent  string
	RewardPercent       int
}

func (a *Allocation) FromRow(row []string) error {
	// COLUMNS
	const (
		ROW_NO = iota
		COMAPNY_NAME
		FIRST_NAME
		LAST_NAME
		KNOWN_BY
		KYC
		STREET
		STREET2
		ZIP
		CITY
		COUNTRY
		BUCKET
		CAM_PURCHASE_PRICE
		CAM_AMOUNT
		P_CHAIN_ADDRESS
		CONSORTIUM_MEMBER
		CONTROL_GROUP
		MSIG_THRESHOLD
		NODE_ID
		VALIDATION_PERIOD_DAYS
		UNBONDING_START_YEARS
		UNBONDING_PERIOD_YEARS
		ADDITIONAL_1PERCENT
		REWARD_PERCENT
	)

	var err error
	a.Bucket = row[BUCKET]
	a.Kyc = row[KYC]
	a.ConsortiumMember = row[CONSORTIUM_MEMBER]
	a.ControlGroup = row[CONTROL_GROUP]
	a.Additional1Percent = row[ADDITIONAL_1PERCENT]

	a.RowNo, err = strconv.Atoi(row[ROW_NO])
	if err != nil {
		return fmt.Errorf("could not parse row number %s", row[ROW_NO])
	}

	a.Amount, err = strconv.ParseUint(row[CAM_AMOUNT], 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse amount %s", row[CAM_AMOUNT])
	}
	a.Amount *= units.Avax

	if row[MSIG_THRESHOLD] != "" {
		threshold, err := strconv.ParseUint(row[MSIG_THRESHOLD], 10, 32)
		a.MsigThreshold = uint32(threshold)
		if err != nil {
			return fmt.Errorf("could not parse msig threshold %s", row[MSIG_THRESHOLD])
		}
	}

	row[NODE_ID] = strings.TrimSpace(row[NODE_ID])
	if row[NODE_ID] != "" && row[NODE_ID] != "X" {
		a.NodeId, err = ids.NodeIDFromString(row[NODE_ID])
		if err != nil {
			fmt.Println("could not parse node id", row[NODE_ID])
		}
	}

	if row[VALIDATION_PERIOD_DAYS] != "" {
		vpd, err := strconv.ParseUint(row[VALIDATION_PERIOD_DAYS], 10, 32)
		a.ValidatorPeriodDays = uint32(vpd)
		if err != nil {
			fmt.Println("could not parse Validator Period: ", row[VALIDATION_PERIOD_DAYS])
		}
	}

	if len(row[P_CHAIN_ADDRESS]) >= 47 {
		_, _, addrBytes, err := address.Parse(row[P_CHAIN_ADDRESS])
		if err != nil {
			return fmt.Errorf("could not parse address %s", row[P_CHAIN_ADDRESS])
		}
		a.Address, _ = ids.ToShortID(addrBytes)
	}

	a.UnbondingStart, err = strconv.ParseFloat(row[UNBONDING_START_YEARS], 64)
	if err != nil {
		a.UnbondingStart = 0
	}
	a.UnbondingPeriod, err = strconv.ParseFloat(row[UNBONDING_PERIOD_YEARS], 64)
	if err != nil {
		a.UnbondingPeriod = 0
	}

	if strings.HasSuffix(row[REWARD_PERCENT], "%") {
		a.RewardPercent, err = strconv.Atoi(row[REWARD_PERCENT][:len(row[REWARD_PERCENT])-1])
		if err != nil {
			a.RewardPercent = -1
		}
	} else {
		a.RewardPercent = -1
	}

	return nil
}
