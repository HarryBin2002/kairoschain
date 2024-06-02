#!/bin/bash

KEYS[0]="validator"
KEYS[1]="dev1"
KEYS[2]="dev2"
KEYS[3]="dev3"
# 0xc032bfb0a7f4d79f8bd0d4d6c6169f58e702817a
MNEMONICS[0]="camera foster skate whisper faith opera axis false van urban clean pet shove census surface injury phone alley cup school pet edge trial pony"
# 0x89760f514DCfCCCf1E4c5eDC6Bf6041931c4c183
MNEMONICS[1]="curtain hat remain song receive tower stereo hope frog cheap brown plate raccoon post reflect wool sail salmon game salon group glimpse adult shift"
# 0x21b661c8A270ed83D2826aD49b1E3B78F515E25C
MNEMONICS[2]="coral drink glow assist canyon ankle hole buffalo vendor foster void clip welcome slush cherry omit member legal account lunar often hen winter culture"
# 0x6479D25261A74B1b058778d3F69Ad7cC557341A8
MNEMONICS[3]="depth skull anxiety weasel pulp interest seek junk trumpet orbit glance drink comfort much alarm during lady strong matrix enable write pledge alcohol buzz"

CHAINID="${CHAIN_ID:-kairoschain_80808-1}"
MONIKER="localtestnet"
KEYRING="test" # remember to change to other types of keyring like 'file' in-case exposing to outside world, otherwise your balance will be wiped quickly. The keyring test does not require private key to steal tokens from you
BINARY="kairosd"
MIN_DENOM="ukai"
KEYALGO="eth_secp256k1" #gitleaks:allow
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""
PRUNING="default"
#PRUNING="custom"

CHAINDIR="$HOME/.kairoschain"
GENESIS="$CHAINDIR/config/genesis.json"
TMP_GENESIS="$CHAINDIR/config/tmp_genesis.json"
APP_TOML="$CHAINDIR/config/app.toml"
CONFIG_TOML="$CHAINDIR/config/config.toml"

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# used to exit on first error (any non-zero exit code)
set -e

# Set client config
"$BINARY" config keyring-backend "$KEYRING"
"$BINARY" config chain-id "$CHAINID"

# Recover keys from mnemonics
echo "${MNEMONICS[0]}" | "$BINARY" keys add "${KEYS[0]}" --recover --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
echo "${MNEMONICS[1]}" | "$BINARY" keys add "${KEYS[1]}" --recover --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
echo "${MNEMONICS[2]}" | "$BINARY" keys add "${KEYS[2]}" --recover --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
echo "${MNEMONICS[3]}" | "$BINARY" keys add "${KEYS[3]}" --recover --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"

# Set moniker for this node
"$BINARY" init "$MONIKER" --chain-id "$CHAINID"

# Change parameter token denominations
jq '.app_state.staking.params.bond_denom="'$MIN_DENOM'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state.crisis.constant_fee.denom="'$MIN_DENOM'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state.gov.deposit_params.min_deposit[0].denom="'$MIN_DENOM'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS" # legacy
jq '.app_state.gov.params.min_deposit[0].denom="'$MIN_DENOM'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS" # v0.47
jq '.app_state.evm.params.evm_denom="'$MIN_DENOM'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state.mint.params.mint_denom="'$MIN_DENOM'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# set gov proposing && voting period
jq '.app_state.gov.deposit_params.max_deposit_period="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS" # legacy
jq '.app_state.gov.params.max_deposit_period="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS" # v0.47
jq '.app_state.gov.voting_params.voting_period="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS" # legacy
jq '.app_state.gov.params.voting_period="30s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS" # v0.47

# Set gas limit in genesis
jq '.consensus_params.block.max_gas="10000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# disable produce empty block
sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' "$CONFIG_TOML"

# Allocate genesis accounts (cosmos formatted addresses)
GENESIS_BALANCE="100000000000000000000000000"
"$BINARY" add-genesis-account "${KEYS[0]}" "$GENESIS_BALANCE$MIN_DENOM" --keyring-backend $KEYRING --home "$HOMEDIR"
"$BINARY" add-genesis-account "${KEYS[1]}" "$GENESIS_BALANCE$MIN_DENOM" --keyring-backend $KEYRING --home "$HOMEDIR"
"$BINARY" add-genesis-account "${KEYS[2]}" "$GENESIS_BALANCE$MIN_DENOM" --keyring-backend $KEYRING --home "$HOMEDIR"
"$BINARY" add-genesis-account "${KEYS[3]}" "$GENESIS_BALANCE$MIN_DENOM" --keyring-backend $KEYRING --home "$HOMEDIR"

# Bc is required to add this big numbers
# total_supply=$(bc <<< "$validators_supply")
total_supply=$(echo "${#KEYS[@]} * $GENESIS_BALANCE" | bc)
jq -r --arg total_supply "$total_supply" '.app_state.bank.supply[0].amount=$total_supply' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# set custom pruning settings
if [ "$PRUNING" = "custom" ]; then
  sed -i 's/pruning = "default"/pruning = "custom"/g' "$APP_TOML"
  sed -i 's/pruning-keep-recent = "0"/pruning-keep-recent = "2"/g' "$APP_TOML"
  sed -i 's/pruning-interval = "0"/pruning-interval = "10"/g' "$APP_TOML"
fi

# make sure the localhost IP is 0.0.0.0
sed -i 's/localhost/0.0.0.0/g' "$CONFIG_TOML"
sed -i 's/127.0.0.1/0.0.0.0/g' "$APP_TOML"

# use timeout_commit 1s to make test faster
sed -i 's/timeout_commit = "3s"/timeout_commit = "1s"/g' "$CONFIG_TOML"

# Sign genesis transaction
"$BINARY" gentx "${KEYS[0]}" "1000000000000000000000$MIN_DENOM" --keyring-backend $KEYRING --chain-id "$CHAINID"
## In case you want to create multiple validators at genesis
## 1. Back to `"$BINARY" keys add` step, init more keys
## 2. Back to `"$BINARY" add-genesis-account` step, add balance for those
## 3. Clone this ~/.kairoschain home directory into some others, let's say `~/.clonedHome`
## 4. Run `gentx` in each of those folders
## 5. Copy the `gentx-*` folders under `~/.clonedHome/config/gentx/` folders into the original `~/.kairoschain/config/gentx`

# Collect genesis tx
"$BINARY" collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
"$BINARY" validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
"$BINARY" start \
  --metrics "$TRACE" --log_level "$LOGLEVEL" \
  --minimum-gas-prices="0.0001$MIN_DENOM" \
  --json-rpc.api eth,txpool,personal,net,debug,web3 \
  --api.enable \
  --grpc.enable true \
  --chain-id "$CHAINID"
