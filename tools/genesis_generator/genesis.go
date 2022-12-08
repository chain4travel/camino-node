package main

import (
	"fmt"
	"genesis_generator/workbook"
	platform "github.com/ava-labs/avalanchego/vms/platformvm/genesis"
	"strconv"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
)

func generateAllocations(allocations []*workbook.Allocation, offerValueToID map[string]ids.ID) []*genesis.CaminoAllocation {
	var parsedAlloc []*genesis.CaminoAllocation
	skippedRows := 0
	for _, al := range allocations {

		// early exits
		if al.Address == ids.ShortEmpty {
			fmt.Println("Skipping Row # ", al.RowNo, " Reason: Address Empty")
			skippedRows++
			continue
		}

		if al.Amount == 0 {
			fmt.Println("Skipping Row # ", al.RowNo, " Reason: No allocation amount given")
			skippedRows++
			continue
		}

		// Computation of the offer value as a key to the map of DepositOffers
		YearToSeconds := float64(365 * 24 * 60 * 60)
		offerValueMinDuration := uint64((al.UnbondingStart + al.UnbondingPeriod) * YearToSeconds)
		offerValueUnlockPeriodDuration := uint64(al.UnbondingPeriod * YearToSeconds)
		offerValueIndex := strconv.FormatUint(offerValueMinDuration, 10) + "_" + strconv.FormatUint(offerValueUnlockPeriodDuration, 10) + "_" + strconv.FormatInt(int64(al.RewardPercent), 10)

		directXAmount := uint64(0)
		if offerValueMinDuration == 0 && offerValueUnlockPeriodDuration == 0 {
			directXAmount = al.Amount
		}

		onePercent := uint64(0)
		if al.Additional1Percent == "y" {
			onePercent = al.Amount / 100
		}
		a := &genesis.CaminoAllocation{
			XAmount:  directXAmount + onePercent,
			AVAXAddr: al.Address,
		}

		if offerValueMinDuration != 0 && offerValueUnlockPeriodDuration != 0 {
			depositOfferID, ok := offerValueToID[offerValueIndex]
			if !ok {
				fmt.Println("Skipping Row # ", al.RowNo, " Reason: No fitting allocation found for values -- index: ", offerValueIndex)
				skippedRows++
				continue
			}

			pa := genesis.PlatformAllocation{
				Amount:            al.Amount,
				DepositOfferID:    depositOfferID,
				NodeID:            al.NodeId,
				ValidatorDuration: uint64(al.ValidatorPeriodDays * 24 * 60 * 60),
			}
			a.PlatformAllocations = append(a.PlatformAllocations, pa)
		}

		parsedAlloc = append(parsedAlloc, a)
	}
	fmt.Println("Skipped ", skippedRows, "allocation rows")

	return parsedAlloc
}

func valueIndex(offer platform.DepositOffer) string {
	// offer.MinDuration :: Min-duration in seconds -- is for example 3.5years as seconds for the offer with 2.5 + 1 year unlock
	// offer.UnlockPeriodDuration :: The duration the unlock will last -- it's exactly what it's written in the excel file (in years) for unbonding period in seconds
	// offer.InterestRateNominator :: the percentage given as a reward * 10000
	index := strconv.FormatUint(uint64(offer.MinDuration), 10) + "_" + strconv.FormatUint(uint64(offer.UnlockPeriodDuration), 10) + "_" + strconv.FormatUint(offer.InterestRateNominator/10_000, 10)
	return index
}

func unparseAllocations(genAlloc []*genesis.CaminoAllocation) []genesis.UnparsedCaminoAllocation {
	var confAlloc []genesis.UnparsedCaminoAllocation
	for i, ga := range genAlloc {
		uga, err := ga.Unparse(constants.KopernikusID)
		if err != nil {
			fmt.Println("Could not unparse allocation for ", i, err)
		}
		confAlloc = append(confAlloc, uga)
	}
	return confAlloc
}

//type ExtMultisig struct {
//	ControlGroup                  string `serialize:"true" json:"controlGroup"`
//	genesis.UnparsedMultisigAlias `serialize:"true"`
//}
//
//var GENESIS_TX_ID ids.ID = ids.Empty
//
//func newFrom(ms *workbook.MultiSig) (*ExtMultisig, error) {
//	ma, err := pchain.NewMultisigAlias(GENESIS_TX_ID, ms.Addrs, ms.Threshold)
//	uma := genesis.UnparsedMultisigAlias{}
//	err = uma.Unparse(ma, constants.CaminoID)
//	if err != nil {
//		return nil, err
//	}
//	return &ExtMultisig{
//		ControlGroup:          ms.ControlGroup,
//		UnparsedMultisigAlias: uma,
//	}, err
//}
//
//func printMultiSig(ms []*workbook.MultiSig) error {
//	mas := make([]*ExtMultisig, len(ms))
//	for i, m := range ms {
//		ma, err := newFrom(m)
//		if err != nil {
//			fmt.Println("could not create multisig", m.ControlGroup, err)
//			continue
//		}
//		mas[i] = ma
//	}
//
//	msDef, err := json.MarshalIndent(mas, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	fmt.Println(string(msDef))
//	return nil
//}
