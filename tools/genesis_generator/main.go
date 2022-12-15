package main

import (
	"encoding/json"
	"fmt"
	"genesis_generator/workbook"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		usage := fmt.Sprintf("Usage: %s <workbook> <genesis_json>", os.Args[0])
		log.Panic(usage)
	}

	spreadsheet_file := os.Args[1]
	genesis_file := os.Args[2]

	genesisConfig, err := readGenesisConfig(genesis_file)
	fmt.Println("genesis config NetworkID", genesisConfig.NetworkID)

	// Set ID on DepositOffers
	offerValuesToID := make(map[string]ids.ID)
	for _, offer := range genesisConfig.Camino.DepositOffers {
		id, _ := offer.ID()
		index := valueIndex(offer)

		offerValuesToID[index] = id
		fmt.Println("DepositOffer index ", index, " ID:", id)
	}

	fmt.Println("Loadingspreadsheet", spreadsheet_file)
	xls, err := excelize.OpenFile(spreadsheet_file)
	if err != nil {
		log.Panic("Could not open the file", err)
	}
	defer xls.Close()
	allocationRows := loadRows(xls, workbook.WB_ALLOCATIONS_NAME)
	multisigRows := loadRows(xls, workbook.WB_MSIG_NAME)

	allocations, err := parseAllocations(allocationRows)
	fmt.Println("Loaded allocations", len(allocations), "err", err)
	multisigs, err := parseMultiSigGroups(multisigRows, allocations)
	fmt.Println("Loaded multisigs", len(multisigs), "err", err)

	msigGroups, _ := generateMSigDefinitions(genesisConfig.NetworkID, multisigs)
	genesisConfig.Camino.InitialMultisigAddresses = msigGroups.MultisigAliaseas
	// create Genesis allocation records
	genAlloc := generateAllocations(allocations, offerValuesToID, msigGroups.ControlGroupToAlias)
	// Uparse for Kopernikus and fill the allocation to config
	confAlloc := unparseAllocations(genAlloc)
	genesisConfig.Camino.Allocations = confAlloc

	// saving the json file
	bytes, err := json.MarshalIndent(genesisConfig, "", "  ")
	if err != nil {
		fmt.Println("Could not marshal genesis config", err)
		return
	}

	err = os.WriteFile("genesis.json", bytes, 0644)
	fmt.Println("DONE")
}

func readGenesisConfig(genesis_file string) (genesis.UnparsedConfig, error) {
	genesisConfig := genesis.UnparsedConfig{}
	file, err := os.Open(genesis_file)
	if err != nil {
		log.Panic("unable to read genesis file", genesis_file, err)
	}
	fileBytes, _ := io.ReadAll(file)
	err = json.Unmarshal(fileBytes, &genesisConfig)
	if err != nil {
		log.Panic("error while parsing genesis json", err)
	}

	return genesisConfig, err
}
