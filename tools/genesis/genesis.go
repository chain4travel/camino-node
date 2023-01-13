package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	platform "github.com/ava-labs/avalanchego/vms/platformvm/genesis"
	"github.com/chain4travel/camino-node/tools/genesis/workbook"
)

func generateMSigDefinitions(networkID uint32, msigs []*workbook.MultiSig) (MultisigDefs, error) {
	var (
		msDefs   = []platform.MultisigAlias{}
		cgToMSig = map[string]ids.ShortID{}
	)

	txID := ids.Empty
	for idx, ms := range msigs {
		// Note: only control_group makes an alias, I'm ignoring unlikely possible hashing collisions
		memo := ms.ControlGroup
		ma, err := platform.NewMultisigAlias(txID, ms.Addrs, ms.Threshold, memo)
		if err != nil {
			fmt.Println("Could not create multisig definition for ", ms.ControlGroup, err)
		}

		if err = memoSanityCheck(&ma, idx); err != nil {
			log.Panic(err)
		}

		msDefs = append(msDefs, ma)
		cgToMSig[memo] = ma.Alias

		fmt.Println("MSig alias generated ", memo, "  Addr:", addrToString(networkID, ma.Alias))
	}

	defs := MultisigDefs{
		ControlGroupToAlias: cgToMSig,
		MultisigAliaseas:    make([]genesis.UnparsedMultisigAlias, 0, len(msDefs)),
	}

	strAliases := map[ids.ShortID]genesis.UnparsedMultisigAlias{}
	for _, ali := range msDefs {
		uma := genesis.UnparsedMultisigAlias{}
		err := uma.Unparse(ali, networkID)
		if err != nil {
			fmt.Println("Could not unparse multisig definition for ", ali.Alias, err)
		}
		strAliases[ali.Alias] = uma
		defs.MultisigAliaseas = append(defs.MultisigAliaseas, uma)
	}

	return defs, nil
}

type UnlockedFunds int

// TransferToPChain for now is the default. At some point we want to have a choice.
const (
	TransferToPChain UnlockedFunds = iota
	// TransferToXChain
)

func generateAllocations(
	networkID uint32,
	allocations []*workbook.Allocation,
	offerValueToID map[string]ids.ID,
	msigCtrlGrpToAlias map[string]ids.ShortID,
	unlockedFundsDestination UnlockedFunds,
) ([]*genesis.CaminoAllocation, ids.ShortID) {
	parsedAlloc := make([]*genesis.CaminoAllocation, 0, len(allocations))
	skippedRows := 0
	adminAddr := ids.ShortEmpty
	for _, al := range allocations {
		msigAlias, ok := msigCtrlGrpToAlias[al.ControlGroup]
		if ok {
			al.Address = msigAlias
			fmt.Printf("replaced row %3d address with its control group alias %s\n", al.RowNo, al.ControlGroup)
		}

		// print addresses generated from public keys
		if !ok && al.PublicKey != "" {
			fmt.Printf("replaced row %3d public key %s resolved to address %s\n", al.RowNo, al.PublicKey[:11], addrToString(networkID, al.Address))
		}

		// early exits
		if al.Address == ids.ShortEmpty {
			fmt.Println("Skipping Row # ", al.RowNo, " Reason: Address Empty")
			skippedRows++
			continue
		}

		if al.FirstName == "ADMIN" {
			adminAddr = al.Address
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

		directAmount := uint64(0)
		if offerValueMinDuration == 0 && offerValueUnlockPeriodDuration == 0 {
			directAmount = al.Amount
		}

		onePercent := uint64(0)
		if al.Additional1Percent == "y" {
			onePercent = al.Amount / 100
		}

		isConsortiumMember := al.ConsortiumMember != ""
		isKycVerified := al.Kyc == "y"

		a := &genesis.CaminoAllocation{
			AVAXAddr:      al.Address,
			AddressStates: genesis.AddressStates{ConsortiumMember: isConsortiumMember, KYCVerified: isKycVerified},
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
				DepositDuration:   offerValueMinDuration,
				NodeID:            al.NodeID,
				ValidatorDuration: uint64(al.ValidatorPeriodDays * 24 * 60 * 60),
				TimestampOffset:   al.TokenDeliveryOffset,
				Memo:              strconv.Itoa(al.RowNo),
			}
			a.PlatformAllocations = append(a.PlatformAllocations, pa)
		}

		unlockedFunds := directAmount + onePercent
		if unlockedFunds > 0 && unlockedFundsDestination == TransferToPChain {
			additionalUnlocked := genesis.PlatformAllocation{
				Amount:          unlockedFunds,
				TimestampOffset: al.TokenDeliveryOffset,
				Memo:            fmt.Sprintf("%d+", al.RowNo),
			}
			a.PlatformAllocations = append(a.PlatformAllocations, additionalUnlocked)
		} else {
			a.XAmount = unlockedFunds
		}

		parsedAlloc = append(parsedAlloc, a)
	}
	fmt.Println("Skipped ", skippedRows, "allocation rows")

	return parsedAlloc, adminAddr
}

func valueIndex(offer genesis.UnparsedDepositOffer) string {
	// offer.MinDuration :: Min-duration in seconds -- is for example 3.5years as seconds for the offer with 2.5 + 1 year unlock
	// offer.UnlockPeriodDuration :: The duration the unlock will last -- it's exactly what it's written in the excel file (in years) for unbonding period in seconds
	// offer.InterestRateNominator :: the percentage given as a reward * 10000
	index := strconv.FormatUint(uint64(offer.MinDuration), 10) + "_" + strconv.FormatUint(uint64(offer.UnlockPeriodDuration), 10) + "_" + strconv.FormatUint(offer.InterestRateNominator/10_000, 10)
	return index
}

func unparseAllocations(genAlloc []*genesis.CaminoAllocation, networkID uint32) []genesis.UnparsedCaminoAllocation {
	confAlloc := make([]genesis.UnparsedCaminoAllocation, 0, len(genAlloc))
	for i, ga := range genAlloc {
		uga, err := ga.Unparse(networkID)
		if err != nil {
			fmt.Println("Could not unparse allocation for ", i, err)
		}
		confAlloc = append(confAlloc, uga)
	}
	return confAlloc
}

func addrToString(networkID uint32, addr ids.ShortID) string {
	fmtAddr, _ := address.Format("X", constants.NetworkIDToHRP[networkID], addr.Bytes())
	return fmtAddr
}

func memoSanityCheck(ma *platform.MultisigAlias, index int) error {
	// Sanity check: Unparse & Parse & Verify
	uma := genesis.UnparsedMultisigAlias{}
	err := uma.Unparse(*ma, 1)
	if err != nil {
		return err
	}
	mm, err := uma.Parse()
	if err != nil {
		return err
	}

	// reverse computation check
	return mm.Verify(ids.Empty)
}

type MultisigDefs struct {
	ControlGroupToAlias map[string]ids.ShortID
	MultisigAliaseas    []genesis.UnparsedMultisigAlias
}
