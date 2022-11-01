// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package version

import (
	"fmt"

	sdkVersion "github.com/ava-labs/avalanchego/version"
)

var (
	// String is displayed when CLI arg --version is used
	String string

	// GitCommit is set in the build script at compile time
	GitCommit string
)

func init() {
	format := "core: %s [database: %s"
	args := []interface{}{
		sdkVersion.Current,
		sdkVersion.CurrentDatabase,
	}
	if GitCommit != "" {
		format += ", commit=%s"
		args = append(args, GitCommit)
	}
	format += "]\n"
	String = fmt.Sprintf(format, args...)
}
