// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code whose
// original notices appear below.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************
// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/utils/crypto"
)

func TestLoadTestKeys(t *testing.T) {
	keys, err := LoadHexTestKeys("test.insecure.secp256k1.keys")
	require.NoError(t, err)
	for i, k := range keys {
		curAddr := encodeShortAddr(k)
		t.Logf("[%d] loaded %v", i, curAddr)
	}
}

func encodeShortAddr(pk *crypto.PrivateKeySECP256K1R) string {
	return pk.PublicKey().Address().String()
}
