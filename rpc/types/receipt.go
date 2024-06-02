package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func (r *RPCReceipt) AsEthReceipt() *ethtypes.Receipt {
	var contractAddress common.Address
	if r.ContractAddress != nil {
		contractAddress = *r.ContractAddress
	}

	return &ethtypes.Receipt{
		Type:              uint8(r.Type),
		PostState:         nil,
		Status:            uint64(r.Status),
		CumulativeGasUsed: uint64(r.CumulativeGasUsed),
		Bloom:             r.Bloom,
		Logs:              r.Logs,
		TxHash:            r.TransactionHash,
		ContractAddress:   contractAddress,
		GasUsed:           uint64(r.GasUsed),
		BlockHash:         r.BlockHash,
		BlockNumber:       big.NewInt(int64(r.BlockNumber)),
		TransactionIndex:  uint(r.TransactionIndex),
	}
}
