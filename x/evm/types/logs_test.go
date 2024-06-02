package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	utiltx "github.com/HarryBin2002/kairoschain/v12/testutil/tx"
	"github.com/HarryBin2002/kairoschain/v12/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
)

func TestValidateLog(t *testing.T) {
	addr := utiltx.GenerateAddress().String()

	testCases := []struct {
		name    string
		log     *types.Log
		expPass bool
	}{
		{
			"valid log",
			&types.Log{
				Address:     addr,
				Topics:      []string{common.BytesToHash([]byte("topic")).String()},
				Data:        []byte("data"),
				BlockNumber: 1,
				TxHash:      common.BytesToHash([]byte("tx_hash")).String(),
				TxIndex:     1,
				BlockHash:   common.BytesToHash([]byte("block_hash")).String(),
				Index:       1,
				Removed:     false,
			},
			true,
		},
		{
			"empty log", &types.Log{}, false,
		},
		{
			"zero address",
			&types.Log{
				Address: common.Address{}.String(),
			},
			false,
		},
		{
			"empty block hash",
			&types.Log{
				Address:   addr,
				BlockHash: common.Hash{}.String(),
			},
			false,
		},
		{
			"zero block number",
			&types.Log{
				Address:     addr,
				BlockHash:   common.BytesToHash([]byte("block_hash")).String(),
				BlockNumber: 0,
			},
			false,
		},
		{
			"empty tx hash",
			&types.Log{
				Address:     addr,
				BlockHash:   common.BytesToHash([]byte("block_hash")).String(),
				BlockNumber: 1,
				TxHash:      common.Hash{}.String(),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.log.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
