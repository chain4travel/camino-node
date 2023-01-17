package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/chain4travel/camino-node/tools/genesis/workbook"
	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) < 5 {
		usage := fmt.Sprintf("Usage: %s <workbook> <genesis_json> <network> <output_dir>", os.Args[0])
		log.Panic(usage)
	}

	spreadsheetFile := os.Args[1]
	genesisFile := os.Args[2]
	networkName := os.Args[3]
	outputPath := os.Args[4]

	// Allows to choose where additional unlocked funds (e.g. 1%) will be sent
	unlockedFunds := TransferToPChain

	networkID := uint32(0)
	switch networkName {
	case "camino":
		networkID = constants.CaminoID
	case "columbus":
		networkID = constants.ColumbusID
	case "kopernikus":
		networkID = constants.KopernikusID
	default:
		log.Panic("Need to provide a valid network name (camino|columbus|kopernikus)")
	}

	genesisConfig, err := readGenesisConfig(genesisFile)
	if err != nil {
		log.Panic("Could not open the genesis template file", err)
	}
	fmt.Println("Read genesis template with NetworkID", genesisConfig.NetworkID, " overwriting with ", networkID)
	genesisConfig.NetworkID = networkID

	// Set ID on DepositOffers
	offerValuesToID := make(map[string]ids.ID)
	depositOffers := make([]genesis.UnparsedDepositOffer, len(genesisConfig.Camino.DepositOffers))
	for i, offer := range genesisConfig.Camino.DepositOffers {
		parsedOffer, err := offer.Parse(genesisConfig.StartTime)
		if err != nil {
			log.Panic("Error parsing offer at", i, err)
		}
		parsedOffer.End += workbook.TierTimeDelayInSeconds
		id, _ := parsedOffer.ID()
		index := valueIndex(offer)

		offerValuesToID[index] = id
		fmt.Println("DepositOffer index ", index, " ID:", id)

		if err = depositOffers[i].Unparse(parsedOffer, genesisConfig.StartTime); err != nil {
			log.Panic("Error unparsing offer ", i, "after modifications", err)
		}
	}
	genesisConfig.Camino.DepositOffers = depositOffers

	fmt.Println("Loadingspreadsheet", spreadsheetFile)
	xls, err := excelize.OpenFile(spreadsheetFile)
	if err != nil {
		log.Panic("Could not open the file", err)
	}
	defer xls.Close()
	allocationRows := loadRows(xls, workbook.Allocations)
	multisigRows := loadRows(xls, workbook.MultisigDefinitions)

	allocations, err := parseAllocations(allocationRows)
	fmt.Println("Loaded allocations", len(allocations), "err", err)
	multisigs, err := parseMultiSigGroups(multisigRows, allocations)
	fmt.Println("Loaded multisigs", len(multisigs), "err", err)

	msigGroups, _ := generateMSigDefinitions(genesisConfig.NetworkID, multisigs)
	genesisConfig.Camino.InitialMultisigAddresses = msigGroups.MultisigAliaseas
	// create Genesis allocation records
	genAlloc, adminAddr := generateAllocations(genesisConfig.NetworkID, allocations, offerValuesToID, msigGroups.ControlGroupToAlias, unlockedFunds)
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

	outputFileName := fmt.Sprintf("%s/genesis_%s.json", outputPath, constants.NetworkIDToHRP[networkID])
	fmt.Println("Saving genesis to", outputFileName)
	err = os.WriteFile(outputFileName, bytes, 0o600)
	if err != nil {
		log.Panic("Could not write the output file: ", outputFileName, err)
	}

	fmt.Println("DONE")
}

func readGenesisConfig(genesisFile string) (genesis.UnparsedConfig, error) {
	genesisConfig := genesis.UnparsedConfig{}
	file, err := os.Open(genesisFile)
	if err != nil {
		log.Panic("unable to read genesis file", genesisFile, err)
	}
	fileBytes, _ := io.ReadAll(file)
	err = json.Unmarshal(fileBytes, &genesisConfig)
	if err != nil {
		log.Panic("error while parsing genesis json", err)
	}

	return genesisConfig, err
}
