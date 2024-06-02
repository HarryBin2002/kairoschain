package ante_test

import (
	"github.com/HarryBin2002/kairoschain/v12/constants"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/HarryBin2002/kairoschain/v12/app"
	"github.com/HarryBin2002/kairoschain/v12/crypto/ethsecp256k1"
	"github.com/HarryBin2002/kairoschain/v12/encoding"
	"github.com/HarryBin2002/kairoschain/v12/testutil"
	feemarkettypes "github.com/HarryBin2002/kairoschain/v12/x/feemarket/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var s *AnteTestSuite

type AnteTestSuite struct {
	suite.Suite

	ctx       sdk.Context
	clientCtx client.Context
	app       *app.Kairoschain
	denom     string
}

func (suite *AnteTestSuite) SetupTest() {
	t := suite.T()
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	isCheckTx := false
	chainID := constants.TestnetFullChainId
	suite.app = app.Setup(isCheckTx, feemarkettypes.DefaultGenesisState(), chainID)
	suite.Require().NotNil(suite.app.AppCodec())

	header := testutil.NewHeader(
		1, time.Now().UTC(), chainID, consAddress, nil, nil)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, header)

	suite.denom = constants.BaseDenom
	evmParams := suite.app.EvmKeeper.GetParams(suite.ctx)
	evmParams.EvmDenom = suite.denom
	_ = suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)

	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)
}

func TestAnteTestSuite(t *testing.T) {
	s = new(AnteTestSuite)
	suite.Run(t, s)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Run AnteHandler Integration Tests")
}
