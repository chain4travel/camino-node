package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/chain4travel/camino-node/tools/genesis/workbook"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/maps"
)

// Reads all rows from xls file "Allocations" workbook
func parseAllocations(allocationRows [][]string) ([]*workbook.Allocation, error) {
	var err error
	allocations := make([]*workbook.Allocation, len(allocationRows))
	for i, row := range allocationRows {
		allocations[i] = &workbook.Allocation{}
		err = allocations[i].FromRow(row)
		if err != nil {
			fmt.Println("could not parse row", i, err)
			continue
		}
	}

	return allocations, nil
}

// Reads all rows from xls file "Multisig" workbook
func parseMultiSigGroups(multisigRows [][]string, allocs []*workbook.Allocation) ([]*workbook.MultiSig, error) {
	msMap := make(map[string][][]string)
	colID := 0 // control group
	for _, row := range multisigRows {
		controlGroup := row[colID]
		msMap[controlGroup] = append(msMap[controlGroup], row)
	}

	multis := map[string]*workbook.MultiSig{}
	for _, a := range allocs {
		ctrlGroup, ok := msMap[a.ControlGroup]
		if !ok {
			if a.ControlGroup != "" {
				fmt.Println("could not find control group", a.ControlGroup)
			}
			continue
		}
		ms := &workbook.MultiSig{}
		err := ms.FromRow(a.MsigThreshold, ctrlGroup)
		if err != nil {
			fmt.Println("could not parse multisig for ", a.RowNo, a.ControlGroup, err)
			continue
		}

		if prevMs, ok := multis[ms.ControlGroup]; ok {
			if prevMs.Threshold != ms.Threshold {
				log.Panic("ctrl group which differs by threshold found ", ms.ControlGroup, ": ", prevMs.Threshold, " vs ", ms.Threshold)
			}
			continue
		}
		multis[ms.ControlGroup] = ms
	}

	// also lets have MSig ordered by CtrlGroup
	cgroups := maps.Keys(multis)
	sort.Strings(cgroups)
	sortedMultis := make([]*workbook.MultiSig, len(cgroups))
	for i, cgroup := range cgroups {
		sortedMultis[i] = multis[cgroup]
	}

	return sortedMultis, nil
}

func loadRows(xls *excelize.File, workbook string) [][]string {
	rows, err := xls.GetRows(workbook)
	if err != nil {
		log.Panic("Could not load workbook", workbook, err)
	}

	return rows
}
