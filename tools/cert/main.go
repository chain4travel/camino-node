// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/ava-labs/avalanchego/network/peer"
	"github.com/ava-labs/avalanchego/staking"
)

var (
	keyFile  = "staker%s.key"
	certFile = "staker%s.crt"
	destPath = "./"
)

func main() {
	count := 1
	flag.StringVar(&destPath, "destPath", destPath, "Destination path")
	flag.IntVar(&count, "count", 1, "Number of certificates")
	flag.Parse()

	num := ""
	for i := 1; i <= count; i++ {
		if count > 1 {
			num = fmt.Sprintf("%d", i)
		}

		keyPath := path.Join(destPath, fmt.Sprintf(keyFile, num))
		certPath := path.Join(destPath, fmt.Sprintf(certFile, num))

		err := staking.InitNodeStakingKeyPair(keyPath, certPath)
		if err != nil {
			fmt.Printf("couldn't create certificate files: %s\n", err)
			os.Exit(1)
		}

		cert, err := staking.LoadTLSCertFromFiles(keyPath, certPath)
		if err != nil {
			fmt.Printf("couldn't read staking certificate: %s\n", err)
			os.Exit(1)
		}

		id, err := peer.CertToID(cert.Leaf)
		if err != nil {
			fmt.Printf("cannot extract nodeID from certificate: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("NodeID%s: %s\n", num, id.String())
	}
}
