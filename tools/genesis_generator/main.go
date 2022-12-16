package main

import (
	"encoding/json"
	"fmt"
	"genesis_generator/workbook"
	"io"
	"log"
	"os"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) < 3 {
		usage := fmt.Sprintf("Usage: %s <workbook> <genesis_json> <network>", os.Args[0])
		log.Panic(usage)
	}

	spreadsheet_file := os.Args[1]
	genesis_file := os.Args[2]
	network_name := os.Args[3]

	networkID := uint32(0)
	switch network_name {
	case "camino":
		networkID = constants.CaminoID
	case "columbus":
		networkID = constants.ColumbusID
	case "kopernikus":
		networkID = constants.KopernikusID
	default:
		log.Panic("Need to provide a valid network name (camino|columbus|kopernikus)")
	}

	genesisConfig, err := readGenesisConfig(genesis_file)
	if err != nil {
		log.Panic("Could not open the genesis template file", err)
	}
	fmt.Println("Read genesis template with NetworkID", genesisConfig.NetworkID, " overwriting with ", networkID)
	genesisConfig.NetworkID = networkID

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
	genAlloc, adminAddr := generateAllocations(allocations, offerValuesToID, msigGroups.ControlGroupToAlias)
	// Overwrite the admin addr if given
	if adminAddr != ids.ShortEmpty {
		avaxAddr, _ := address.Format(
			"X",
			constants.GetHRP(networkID),
			adminAddr.Bytes(),
		)
		genesisConfig.Camino.InitialAdmin = avaxAddr
		fmt.Println("InitialAdmin address set with:", avaxAddr)
	}

	// Uparse for Kopernikus and fill the allocation to config
	confAlloc := unparseAllocations(genAlloc, networkID)
	genesisConfig.Camino.Allocations = confAlloc

	// saving the json file
	bytes, err := json.MarshalIndent(genesisConfig, "", "  ")
	if err != nil {
		fmt.Println("Could not marshal genesis config - error: ", err)
		return
	}

	outputFN := fmt.Sprintf("genesis_%s.json", constants.NetworkIDToHRP[networkID])

	err = os.WriteFile(outputFN, bytes, 0644)
	if err != nil {
		log.Panic("Could not write the output file: ", outputFN, err)
	}

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
