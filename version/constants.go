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

// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package version

import (
	"github.com/chain4travel/caminogo/utils/constants"
	sdkVersion "github.com/chain4travel/caminogo/version"
)

// These are globals that describe network upgrades and node versions
var (
	Current                      = sdkVersion.NewDefaultVersion(0, 2, 0)
	CurrentApp                   = sdkVersion.NewDefaultApplication(constants.PlatformName, Current.Major(), Current.Minor(), Current.Patch())
	MinimumCompatibleVersion     = sdkVersion.NewDefaultApplication(constants.PlatformName, 0, 2, 0)
	PrevMinimumCompatibleVersion = sdkVersion.NewDefaultApplication(constants.PlatformName, 0, 1, 0)
	MinimumUnmaskedVersion       = sdkVersion.NewDefaultApplication(constants.PlatformName, 0, 0, 0)
	PrevMinimumUnmaskedVersion   = sdkVersion.NewDefaultApplication(constants.PlatformName, 0, 0, 0)
	VersionParser                = sdkVersion.NewDefaultApplicationParser()
)

func GetCompatibility(networkID uint32) sdkVersion.Compatibility {
	return sdkVersion.NewCompatibility(
		CurrentApp,
		MinimumCompatibleVersion,
		sdkVersion.GetSunrisePhase0Time(networkID),
		PrevMinimumCompatibleVersion,
		MinimumUnmaskedVersion,
		sdkVersion.GetApricotPhase0Time(networkID),
		PrevMinimumUnmaskedVersion,
	)
}
