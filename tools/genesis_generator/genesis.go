package main

import (
	"fmt"
	"genesis_generator/workbook"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/set"
	platform "github.com/ava-labs/avalanchego/vms/platformvm/genesis"
	"golang.org/x/exp/maps"
	"strconv"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
)

func generateMSigDefinitions(networkID uint32, msigs []*workbook.MultiSig) (MultisigDefs, error) {
	var (
		GenesisTxId ids.ID = ids.Empty
		msDefs             = map[ids.ShortID]platform.MultisigAlias{}
		cgToMSig           = map[string]ids.ShortID{}
	)

	for _, ms := range msigs {
		ma, err := platform.NewMultisigAlias(GenesisTxId, ms.Addrs, ms.Threshold)
		if err != nil {
			fmt.Println("Could not create multisig definition for ", ms.ControlGroup, err)
		}
		msDefs[ma.Alias] = ma
		cgToMSig[ms.ControlGroup] = ma.Alias
	}

	aliases := set.NewSet[ids.ShortID](len(msDefs))
	aliases.Add(maps.Keys(msDefs)...)

	uniqAliases := aliases.List()
	utils.Sort(uniqAliases)

	defs := MultisigDefs{
		ControlGroupToAlias: cgToMSig,
		MultisigAliaseas:    make([]genesis.UnparsedMultisigAlias, 0, len(uniqAliases)),
	}

	strAliases := map[ids.ShortID]genesis.UnparsedMultisigAlias{}
	for _, ali := range uniqAliases {
		ma, _ := msDefs[ali]
		uma := genesis.UnparsedMultisigAlias{}
		err := uma.Unparse(ma, networkID)
		if err != nil {
			fmt.Println("Could not unparse multisig definition for ", ali.String(), err)
		}
		strAliases[ali] = uma
		defs.MultisigAliaseas = append(defs.MultisigAliaseas, uma)
	}

	return defs, nil
}

func generateAllocations(
	allocations []*workbook.Allocation,
	offerValueToID map[string]ids.ID,
	msigCtrlGrpToAlias map[string]ids.ShortID,
) []*genesis.CaminoAllocation {
	var parsedAlloc []*genesis.CaminoAllocation
	skippedRows := 0
	for _, al := range allocations {

		msigAlias, ok := msigCtrlGrpToAlias[al.ControlGroup]
		if ok {
			al.Address = msigAlias
			fmt.Println("replaced address with its control group alias for row", al.RowNo)
		}

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

type MultisigDefs struct {
	ControlGroupToAlias map[string]ids.ShortID
	MultisigAliaseas    []genesis.UnparsedMultisigAlias
}
