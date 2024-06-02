package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/HarryBin2002/kairoschain/v12/constants"
)

// RootCodespace is the codespace for all errors defined in this package
const RootCodespace = constants.ApplicationName

// ErrInvalidChainID returns an error resulting from an invalid chain ID.
var ErrInvalidChainID = errorsmod.Register(RootCodespace, 3, "invalid chain ID")
