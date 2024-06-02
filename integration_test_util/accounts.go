package integration_test_util

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
import (
	"github.com/HarryBin2002/kairoschain/v12/crypto/ethsecp256k1"
	etherminthd "github.com/HarryBin2002/kairoschain/v12/crypto/hd"
	itutiltypes "github.com/HarryBin2002/kairoschain/v12/integration_test_util/types"
	
	cosmoshd "github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"
	"testing"
)

// newValidatorAccounts inits and return predefined validator accounts.
// By defining this, data in test cases will be more consistency.
func newValidatorAccounts(t *testing.T) itutiltypes.TestAccounts {
	val1Account := newTestAccountFromMnemonic(t, IT_VAL_1_MNEMONIC)
	val1Account.Type = itutiltypes.TestAccountTypeValidator
	require.Equal(t, IT_VAL_1_VAL_ADDR, val1Account.GetValidatorAddress().String())
	require.Equal(t, IT_VAL_1_CONS_ADDR, val1Account.GetConsensusAddress().String())
	require.Equal(t, IT_VAL_1_ADDR, val1Account.GetCosmosAddress().String())

	val2Account := newTestAccountFromMnemonic(t, IT_VAL_2_MNEMONIC)
	val2Account.Type = itutiltypes.TestAccountTypeValidator
	require.Equal(t, IT_VAL_2_VAL_ADDR, val2Account.GetValidatorAddress().String())
	require.Equal(t, IT_VAL_2_CONS_ADDR, val2Account.GetConsensusAddress().String())
	require.Equal(t, IT_VAL_2_ADDR, val2Account.GetCosmosAddress().String())

	val3Account := newTestAccountFromMnemonic(t, IT_VAL_3_MNEMONIC)
	val3Account.Type = itutiltypes.TestAccountTypeValidator
	require.Equal(t, IT_VAL_3_VAL_ADDR, val3Account.GetValidatorAddress().String())
	require.Equal(t, IT_VAL_3_CONS_ADDR, val3Account.GetConsensusAddress().String())
	require.Equal(t, IT_VAL_3_ADDR, val3Account.GetCosmosAddress().String())

	val4Account := newTestAccountFromMnemonic(t, IT_VAL_4_MNEMONIC)
	val4Account.Type = itutiltypes.TestAccountTypeValidator
	require.Equal(t, IT_VAL_4_VAL_ADDR, val4Account.GetValidatorAddress().String())
	require.Equal(t, IT_VAL_4_CONS_ADDR, val4Account.GetConsensusAddress().String())
	require.Equal(t, IT_VAL_4_ADDR, val4Account.GetCosmosAddress().String())

	val5Account := newTestAccountFromMnemonic(t, IT_VAL_5_MNEMONIC)
	val5Account.Type = itutiltypes.TestAccountTypeValidator
	require.Equal(t, IT_VAL_5_VAL_ADDR, val5Account.GetValidatorAddress().String())
	require.Equal(t, IT_VAL_5_CONS_ADDR, val5Account.GetConsensusAddress().String())
	require.Equal(t, IT_VAL_5_ADDR, val5Account.GetCosmosAddress().String())

	return []*itutiltypes.TestAccount{
		val1Account,
		val2Account,
		val3Account,
		val4Account,
		val5Account,
	}
}

// newWalletsAccounts inits and return predefined wallet accounts.
// By defining this, data in test cases will be more consistency.
func newWalletsAccounts(t *testing.T) itutiltypes.TestAccounts {
	wal1Account := newTestAccountFromMnemonic(t, IT_WAL_1_MNEMONIC)
	wal1Account.Type = itutiltypes.TestAccountTypeWallet
	require.Equal(t, IT_WAL_1_ETH_ADDR, wal1Account.GetEthAddress().String())
	require.Equal(t, IT_WAL_1_ADDR, wal1Account.GetCosmosAddress().String())

	wal2Account := newTestAccountFromMnemonic(t, IT_WAL_2_MNEMONIC)
	wal2Account.Type = itutiltypes.TestAccountTypeWallet
	require.Equal(t, IT_WAL_2_ETH_ADDR, wal2Account.GetEthAddress().String())
	require.Equal(t, IT_WAL_2_ADDR, wal2Account.GetCosmosAddress().String())

	wal3Account := newTestAccountFromMnemonic(t, IT_WAL_3_MNEMONIC)
	wal3Account.Type = itutiltypes.TestAccountTypeWallet
	require.Equal(t, IT_WAL_3_ETH_ADDR, wal3Account.GetEthAddress().String())
	require.Equal(t, IT_WAL_3_ADDR, wal3Account.GetCosmosAddress().String())

	wal4Account := newTestAccountFromMnemonic(t, IT_WAL_4_MNEMONIC)
	wal4Account.Type = itutiltypes.TestAccountTypeWallet
	require.Equal(t, IT_WAL_4_ETH_ADDR, wal4Account.GetEthAddress().String())
	require.Equal(t, IT_WAL_4_ADDR, wal4Account.GetCosmosAddress().String())

	wal5Account := newTestAccountFromMnemonic(t, IT_WAL_5_MNEMONIC)
	wal5Account.Type = itutiltypes.TestAccountTypeWallet
	require.Equal(t, IT_WAL_5_ETH_ADDR, wal5Account.GetEthAddress().String())
	require.Equal(t, IT_WAL_5_ADDR, wal5Account.GetCosmosAddress().String())

	return []*itutiltypes.TestAccount{
		wal1Account,
		wal2Account,
		wal3Account,
		wal4Account,
		wal5Account,
	}
}

// newTestAccountFromMnemonic creates a new test account from a mnemonic.
func newTestAccountFromMnemonic(t *testing.T, mnemonic string) *itutiltypes.TestAccount {
	var algo keyring.SignatureAlgo
	var err error

	//goland:noinspection SpellCheckingInspection
	algo, err = keyring.NewSigningAlgoFromString("eth_secp256k1", supportedKeyringAlgorithms)
	require.NoError(t, err)

	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath)
	require.NoError(t, err)

	privKey := algo.Generate()(derivedPriv)

	priv := &ethsecp256k1.PrivKey{
		Key: privKey.Bytes(),
	}

	return NewTestAccount(t, priv)
}

// NewTestAccount creates a new test account. If the private key is not provided, a new one will be generated.
func NewTestAccount(t *testing.T, nilAblePrivKey *ethsecp256k1.PrivKey) *itutiltypes.TestAccount {
	testAccount := &itutiltypes.TestAccount{}

	var err error

	if nilAblePrivKey == nil {
		nilAblePrivKey, err = ethsecp256k1.GenerateKey()
		require.NoError(t, err)
		require.NotNil(t, nilAblePrivKey)
	}

	testAccount.PrivateKey = nilAblePrivKey
	testAccount.Signer = itutiltypes.NewSigner(nilAblePrivKey)

	return testAccount
}

var supportedKeyringAlgorithms = keyring.SigningAlgoList{etherminthd.EthSecp256k1, cosmoshd.Secp256k1}
var hdPath = cosmoshd.CreateHDPath(60, 0, 0).String()

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
var (
	IT_VAL_1_ADDR      = "kairos1cqetlv987ntelz7s6ntvv95ltrns9qt6q3mumx"
	IT_VAL_1_VAL_ADDR  = "kairosvaloper1cqetlv987ntelz7s6ntvv95ltrns9qt6szcp6l"
	IT_VAL_1_CONS_ADDR = "kairosvalcons1vv3kjxtrh7jredjehk5xw66r62euensszrjglr"
	IT_VAL_1_MNEMONIC  = "camera foster skate whisper faith opera axis false van urban clean pet shove census surface injury phone alley cup school pet edge trial pony"

	IT_VAL_2_ADDR      = "kairos19k6gu9tkr40uyhf86sjmlgy6hu4lpfx4ssjyf0"
	IT_VAL_2_VAL_ADDR  = "kairosvaloper19k6gu9tkr40uyhf86sjmlgy6hu4lpfx4qr3egk"
	IT_VAL_2_CONS_ADDR = "kairosvalcons19fphsrnm2rx9jk4exfdeq46d6ptwlpy38x43kj"
	IT_VAL_2_MNEMONIC  = "explain captain crucial fault symptom degree divorce beyond path security jewel alien beach finish bridge decide toast scene pelican sorry achieve off denial wall"

	IT_VAL_3_ADDR      = "kairos1rxczyg2x94dqcn77t4pyhcndg3r889dw6j5vlj"
	IT_VAL_3_VAL_ADDR  = "kairosvaloper1rxczyg2x94dqcn77t4pyhcndg3r889dw2ph37t"
	IT_VAL_3_CONS_ADDR = "kairosvalcons1vxky3ld4llhaqk8nl6pw6xkxqy97rwda245fvu"
	IT_VAL_3_MNEMONIC  = "worth talent fire announce file skull acquire ethics injury yard home list clap guard busy describe bag front grass noise index vacuum govern number"

	IT_VAL_4_ADDR      = "kairos1gmjvfd4pr0yd94t0x8xw4uwg2j0cn9g9gzgge8"
	IT_VAL_4_VAL_ADDR  = "kairosvaloper1gmjvfd4pr0yd94t0x8xw4uwg2j0cn9g9c3t4c7"
	IT_VAL_4_CONS_ADDR = "kairosvalcons1yl9a7v952ejxju9fec6hqeuuku4372pn3dpd95"
	IT_VAL_4_MNEMONIC  = "question joke action slice mistake carbon virtual still culture push estate inhale true endless market flip hammer word lecture pen toddler lyrics creek regular"

	IT_VAL_5_ADDR      = "kairos1fpveqajjpt2emsfkr5xwp80074mkn38xp0e6pj"
	IT_VAL_5_VAL_ADDR  = "kairosvaloper1fpveqajjpt2emsfkr5xwp80074mkn38x3u68qt"
	IT_VAL_5_CONS_ADDR = "kairosvalcons1p6n7qpnn5lqyyujzrp344drz228l3wx0dcyzz2"
	IT_VAL_5_MNEMONIC  = "tornado fuel drill critic indicate pool few wheat omit sight stage focus mountain amused neck surge post giant vague nut marine spoon fragile outdoor"
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
var (
	IT_WAL_1_ADDR     = "kairos139mq752delxv78jvtmwxhasyrycufsvrtr3924"
	IT_WAL_1_ETH_ADDR = "0x89760f514DCfCCCf1E4c5eDC6Bf6041931c4c183"
	IT_WAL_1_MNEMONIC = "curtain hat remain song receive tower stereo hope frog cheap brown plate raccoon post reflect wool sail salmon game salon group glimpse adult shift"

	IT_WAL_2_ADDR     = "kairos1yxmxrj9zwrkc855zdt2fk83m0r63tcjun4rxny"
	IT_WAL_2_ETH_ADDR = "0x21b661c8A270ed83D2826aD49b1E3B78F515E25C"
	IT_WAL_2_MNEMONIC = "coral drink glow assist canyon ankle hole buffalo vendor foster void clip welcome slush cherry omit member legal account lunar often hen winter culture"

	IT_WAL_3_ADDR     = "kairos1v3uay5np5a93kpv80rfldxkhe32hxsdgc46yau"
	IT_WAL_3_ETH_ADDR = "0x6479D25261A74B1b058778d3F69Ad7cC557341A8"
	IT_WAL_3_MNEMONIC = "depth skull anxiety weasel pulp interest seek junk trumpet orbit glance drink comfort much alarm during lady strong matrix enable write pledge alcohol buzz"

	IT_WAL_4_ADDR     = "kairos1zsdj9vsw44kk46fmnka7k76smsaxgh6p05cyzt"
	IT_WAL_4_ETH_ADDR = "0x141B22B20ead6d6AE93B9DBBeB7b50DC3A645F41"
	IT_WAL_4_MNEMONIC = "author humble raise whisper allow appear typical release fossil address spy jazz damage runway spy gossip add embark wrap frost toe advice matrix laundry"

	IT_WAL_5_ADDR     = "kairos1862crydur2cpjww66dhfzcc26yglvrcsgkpa70"
	IT_WAL_5_ETH_ADDR = "0x3E958191BC1AB01939DAD36e91630Ad111F60f10"
	IT_WAL_5_MNEMONIC = "museum stumble kingdom impulse replace angle exercise trial spring sphere cube brief foil bridge dish earn practice surprise quantum hunt scale solve october scout"
)
