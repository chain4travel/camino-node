package workbook

import (
	"fmt"
	"log"
	"sort"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/maps"
)

type TabName string

const (
	MultisigDefinitions TabName = "MultiSig Addresses"
	DepositOffers       TabName = "depositOffers"
	Allocations         TabName = "Camino Allocation"
)

// ParseAllocations Reads all rows from xls file "Allocations" workbook
func ParseAllocations(xls *excelize.File) []*AllocationRow {
	var (
		err  error
		rows = []*AllocationRow{}
	)
	unparsedRows := loadRows(xls, Allocations)

	for i, urow := range unparsedRows {
		row := &AllocationRow{}
		if detectHeaderRow(i, row.Header(), urow) {
			continue
		}

		if err = row.FromRow(i, urow); err != nil {
			log.Panic("could not parse row: ", i+1, err)
		}
		rows = append(rows, row)
	}

	return rows
}

// ParseMultiSigGroups Reads all rows from xls file "Multisig" workbook
func ParseMultiSigGroups(xls *excelize.File) []*MultiSigGroup {
	var (
		err  error
		rows = []*MultiSigRow{}
	)
	unparsedRows := loadRows(xls, MultisigDefinitions)

	for i, urow := range unparsedRows {
		row := &MultiSigRow{}
		if detectHeaderRow(i, row.Header(), urow) {
			continue
		}

		if err = row.FromRow(i, urow); err != nil {
			log.Panic("could not parse row", i+1, err)
		}
		rows = append(rows, row)
	}

	multis := map[string]*MultiSigGroup{}

	for _, ms := range rows {
		// skip header urow
		if ms.ControlGroup == "Control Group" && ms.Threshold == 0 {
			continue
		}

		group, ok := multis[ms.ControlGroup]
		if ok {
			if group.Threshold != ms.Threshold {
				log.Panic("ctrl group which differs by threshold found ", ms.ControlGroup, ": ", group.Threshold, " vs ", ms.Threshold)
			}
		} else {
			group = &MultiSigGroup{ControlGroup: ms.ControlGroup, Threshold: ms.Threshold, Addrs: []ids.ShortID{}}
		}
		group.Addrs = append(group.Addrs, ms.Addr)
		multis[ms.ControlGroup] = group
	}

	// also lets have MSig ordered by CtrlGroup
	cgroups := maps.Keys(multis)
	sort.Strings(cgroups)
	sortedMultis := make([]*MultiSigGroup, len(cgroups))
	for i, cgroup := range cgroups {
		sortedMultis[i] = multis[cgroup]
	}

	return sortedMultis
}

// DepositOffersWithOrder helps to populate offers into json in the same order as in xls
type DepositOffersWithOrder struct {
	Offers map[string]*genesis.UnparsedDepositOffer
	Order  []string
}

func ParseDepositOfferRows(xls *excelize.File) DepositOffersWithOrder {
	var (
		err  error
		rows = []*DepositOfferRow{}
	)
	unparsedRows := loadRows(xls, DepositOffers)

	for i, urow := range unparsedRows {
		row := &DepositOfferRow{}
		if detectHeaderRow(i, row.Header(), urow) {
			continue
		}

		if err = row.FromRow(i, urow); err != nil {
			log.Panic("could not parse row", i+1, err)
		}
		rows = append(rows, row)
	}

	orderedOffers := DepositOffersWithOrder{
		Offers: make(map[string]*genesis.UnparsedDepositOffer),
	}
	for _, row := range rows {
		offerID, offer := RowToOffer(row)
		if err != nil {
			fmt.Println("could not parse row", offerID, err)
			continue
		}
		orderedOffers.Offers[offerID] = offer
		orderedOffers.Order = append(orderedOffers.Order, offerID)
	}

	return orderedOffers
}

func loadRows(xls *excelize.File, workbook TabName) [][]string {
	rows, err := xls.GetRows(string(workbook))
	if err != nil {
		log.Panic("Could not load workbook", workbook, err)
	}

	return rows
}

func detectHeaderRow(idx int, headerTitles, row []string) bool {
	startsWith := func(expected, row []string) bool {
		for i, expValue := range expected {
			if expValue != row[i] {
				return false
			}
		}
		return true
	}

	if idx == 0 && startsWith(headerTitles, row) {
		return true
	}
	return false
}
