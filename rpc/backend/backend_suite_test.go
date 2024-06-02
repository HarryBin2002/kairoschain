package backend

import (
	"bufio"
	"github.com/HarryBin2002/kairoschain/v12/constants"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	dbm "github.com/cometbft/cometbft-db"

	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"

	"github.com/HarryBin2002/kairoschain/v12/app"
	"github.com/HarryBin2002/kairoschain/v12/crypto/hd"
	"github.com/HarryBin2002/kairoschain/v12/encoding"
	"github.com/HarryBin2002/kairoschain/v12/indexer"
	"github.com/HarryBin2002/kairoschain/v12/rpc/backend/mocks"
	rpctypes "github.com/HarryBin2002/kairoschain/v12/rpc/types"
	utiltx "github.com/HarryBin2002/kairoschain/v12/testutil/tx"
	evmtypes "github.com/HarryBin2002/kairoschain/v12/x/evm/types"
)

type BackendTestSuite struct {
	suite.Suite

	backend *Backend
	from    common.Address
	acc     sdk.AccAddress
	signer  keyring.Signer
}

func TestBackendTestSuite(t *testing.T) {
	suite.Run(t, new(BackendTestSuite))
}

const ChainID = constants.TestnetFullChainId

// SetupTest is executed before every BackendTestSuite test
func (suite *BackendTestSuite) SetupTest() {
	ctx := server.NewDefaultContext()
	ctx.Viper.Set("telemetry.global-labels", []interface{}{})

	baseDir := suite.T().TempDir()
	nodeDirName := "node"
	clientDir := filepath.Join(baseDir, nodeDirName, "kairoschaincli")
	keyRing, err := suite.generateTestKeyring(clientDir)
	if err != nil {
		panic(err)
	}

	// Create Account with set sequence
	suite.acc = sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	accounts := map[string]client.TestAccount{}
	accounts[suite.acc.String()] = client.TestAccount{
		Address: suite.acc,
		Num:     uint64(1),
		Seq:     uint64(1),
	}

	from, priv := utiltx.NewAddrKey()
	suite.from = from
	suite.signer = utiltx.NewSigner(priv)
	suite.Require().NoError(err)

	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	clientCtx := client.Context{}.WithChainID(ChainID).
		WithHeight(1).
		WithTxConfig(encodingConfig.TxConfig).
		WithKeyringDir(clientDir).
		WithKeyring(keyRing).
		WithAccountRetriever(client.TestAccountRetriever{Accounts: accounts})

	allowUnprotectedTxs := false
	idxer := indexer.NewKVIndexer(dbm.NewMemDB(), ctx.Logger, clientCtx)

	suite.backend = NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, idxer)
	suite.backend.queryClient.QueryClient = mocks.NewEVMQueryClient(suite.T())
	suite.backend.clientCtx.Client = mocks.NewClient(suite.T())
	suite.backend.queryClient.FeeMarket = mocks.NewFeeMarketQueryClient(suite.T())
	suite.backend.ctx = rpctypes.ContextWithHeight(1)
	suite.backend.indexer = mocks.NewEVMTxIndexer(suite.T())

	// Add codec
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	suite.backend.clientCtx.Codec = encCfg.Codec
}

// buildEthereumTx returns an example legacy Ethereum transaction
func (suite *BackendTestSuite) buildEthereumTx() (*evmtypes.MsgEthereumTx, []byte) {
	ethTxParams := evmtypes.EvmTxArgs{
		ChainID:  suite.backend.chainID,
		Nonce:    uint64(0),
		To:       &common.Address{},
		Amount:   big.NewInt(0),
		GasLimit: 100000,
		GasPrice: big.NewInt(1),
	}
	msgEthereumTx := evmtypes.NewTx(&ethTxParams)

	// A valid msg should have empty `From`
	msgEthereumTx.From = suite.from.Hex()

	txBuilder := suite.backend.clientCtx.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgEthereumTx)
	suite.Require().NoError(err)

	bz, err := suite.backend.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	suite.Require().NoError(err)
	return msgEthereumTx, bz
}

// buildFormattedBlock returns a formatted block for testing
func (suite *BackendTestSuite) buildFormattedBlock(
	blockRes *tmrpctypes.ResultBlockResults,
	resBlock *tmrpctypes.ResultBlock,
	fullTx bool,
	tx *evmtypes.MsgEthereumTx,
	validator sdk.AccAddress,
	baseFee *big.Int,
) map[string]interface{} {
	header := resBlock.Block.Header
	gasLimit := int64(^uint32(0)) // for `MaxGas = -1` (DefaultConsensusParams)
	gasUsed := new(big.Int).SetUint64(uint64(blockRes.TxsResults[0].GasUsed))

	var transactions ethtypes.Transactions
	var receipts ethtypes.Receipts
	if tx != nil {
		transactions = append(transactions, tx.AsTransaction())
		receipt := createTestReceipt(nil, resBlock, tx, false, mockGasUsed)
		receipts = append(receipts, receipt)
	}

	bloom := ethtypes.CreateBloom(receipts)

	return rpctypes.FormatBlock(
		header,
		suite.backend.chainID,
		resBlock.Block.Size(),
		gasLimit, gasUsed, baseFee,
		transactions, fullTx,
		receipts,
		bloom,
		common.BytesToAddress(validator.Bytes()),
		suite.backend.logger,
	)
}

func createTestReceipt(root []byte, resBlock *tmrpctypes.ResultBlock, tx *evmtypes.MsgEthereumTx, failed bool, gasUsed uint64) *ethtypes.Receipt {
	var status uint64
	if failed {
		status = ethtypes.ReceiptStatusFailed
	} else {
		status = ethtypes.ReceiptStatusSuccessful
	}

	transaction := tx.AsTransaction()

	return &ethtypes.Receipt{
		Type:              transaction.Type(),
		PostState:         root,
		Status:            status,
		CumulativeGasUsed: gasUsed,
		Bloom:             ethtypes.BytesToBloom(ethtypes.LogsBloom([]*ethtypes.Log{})),
		Logs:              []*ethtypes.Log{},
		TxHash:            transaction.Hash(),
		ContractAddress:   common.Address{},
		GasUsed:           gasUsed,
		BlockHash:         common.HexToHash(resBlock.Block.Header.Hash().String()),
		BlockNumber:       big.NewInt(resBlock.Block.Height),
		TransactionIndex:  0,
	}
}

func (suite *BackendTestSuite) generateTestKeyring(clientDir string) (keyring.Keyring, error) {
	buf := bufio.NewReader(os.Stdin)
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	return keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, clientDir, buf, encCfg.Codec, []keyring.Option{hd.EthSecp256k1Option()}...)
}

func (suite *BackendTestSuite) signAndEncodeEthTx(msgEthereumTx *evmtypes.MsgEthereumTx) []byte {
	from, priv := utiltx.NewAddrKey()
	signer := utiltx.NewSigner(priv)

	queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
	RegisterParamsWithoutHeader(queryClient, 1)

	ethSigner := ethtypes.LatestSigner(suite.backend.ChainConfig())
	msgEthereumTx.From = from.String()
	err := msgEthereumTx.Sign(ethSigner, signer)
	suite.Require().NoError(err)

	tx, err := msgEthereumTx.BuildTx(suite.backend.clientCtx.TxConfig.NewTxBuilder(), constants.BaseDenom)
	suite.Require().NoError(err)

	txEncoder := suite.backend.clientCtx.TxConfig.TxEncoder()
	txBz, err := txEncoder(tx)
	suite.Require().NoError(err)

	return txBz
}

func (suite *BackendTestSuite) signMsgEthTx(msgEthereumTx *evmtypes.MsgEthereumTx) (*evmtypes.MsgEthereumTx, []byte) {
	queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
	RegisterParamsWithoutHeader(queryClient, 1)

	ethSigner := ethtypes.LatestSigner(suite.backend.ChainConfig())
	msgEthereumTx.From = suite.from.String()
	err := msgEthereumTx.Sign(ethSigner, suite.signer)
	suite.Require().NoError(err)

	tx, err := msgEthereumTx.BuildTx(suite.backend.clientCtx.TxConfig.NewTxBuilder(), constants.BaseDenom)
	suite.Require().NoError(err)

	txEncoder := suite.backend.clientCtx.TxConfig.TxEncoder()
	txBz, err := txEncoder(tx)
	suite.Require().NoError(err)

	return msgEthereumTx, txBz
}