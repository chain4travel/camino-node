package workbook

import (
	"fmt"
	"strconv"

	"github.com/ava-labs/avalanchego/genesis"
)

const (
	DepositOffers = "depositOffers"
)

type DepositOfferRow int

func RowToOffer(row []string) (string, *genesis.UnparsedDepositOffer, error) {
	const (
		OfferID DepositOfferRow = iota
		InterestRateNominator
		StartOffset
		EndOffset
		MinAmount
		MinDuration
		MaxDuration
		UnlockPeriodDuration
		NoRewardsPeriodDuration
		Locked
		Comment
	)
	var err error
	offerID := row[OfferID]
	offer := &genesis.UnparsedDepositOffer{}

	offer.InterestRateNominator, err = strconv.ParseUint(row[InterestRateNominator], 10, 64)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse interest rate nominator %s, err %w", row[InterestRateNominator], err)
	}

	offer.StartOffset, err = strconv.ParseUint(row[StartOffset], 10, 64)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse start offset %s, err %w", row[StartOffset], err)
	}

	offer.EndOffset, err = strconv.ParseUint(row[EndOffset], 10, 64)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse end offset %s, err %w", row[EndOffset], err)
	}

	offer.MinAmount, err = strconv.ParseUint(row[MinAmount], 10, 64)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse min amount %s, err %w", row[MinAmount], err)
	}

	minDuration, err := strconv.ParseUint(row[MinDuration], 10, 32)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse min duration %s, err %w", row[MinDuration], err)
	}
	offer.MinDuration = uint32(minDuration)

	maxDuration, err := strconv.ParseUint(row[MaxDuration], 10, 32)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse max duration %s, err %w", row[MaxDuration], err)
	}
	offer.MaxDuration = uint32(maxDuration)

	upd, err := strconv.ParseUint(row[UnlockPeriodDuration], 10, 32)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse unlock period duration %s, err %w", row[UnlockPeriodDuration], err)
	}
	offer.UnlockPeriodDuration = uint32(upd)

	nrpd, err := strconv.ParseUint(row[NoRewardsPeriodDuration], 10, 32)
	if err != nil {
		return offerID, offer, fmt.Errorf("could not parse no rewards period duration %s, err %w", row[NoRewardsPeriodDuration], err)
	}
	offer.NoRewardsPeriodDuration = uint32(nrpd)

	if row[Locked] != TrueValue && row[Locked] != FalseValue {
		return offerID, offer, fmt.Errorf("locked value must be either TRUE or FALSE, got %s", row[Locked])
	}
	offer.Flags = genesis.UnparsedDepositOfferFlags{Locked: row[Locked] == TrueValue}

	offer.Memo = row[Comment]

	return offerID, offer, nil
}
