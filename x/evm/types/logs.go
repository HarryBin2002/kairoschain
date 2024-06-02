package types

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	evertypes "github.com/HarryBin2002/kairoschain/v12/types"
)

// Validate performs a basic validation of an ethereum Log fields.
func (log *Log) Validate() error {
	if err := evertypes.ValidateAddress(log.Address); err != nil {
		return fmt.Errorf("invalid log address %w", err)
	}
	if evertypes.IsEmptyHash(log.BlockHash) {
		return fmt.Errorf("block hash cannot be the empty %s", log.BlockHash)
	}
	if log.BlockNumber == 0 {
		return errors.New("block number cannot be zero")
	}
	if evertypes.IsEmptyHash(log.TxHash) {
		return fmt.Errorf("tx hash cannot be the empty %s", log.TxHash)
	}
	return nil
}

// ToEthereum returns the Ethereum type Log from a Ethermint proto compatible Log.
func (log *Log) ToEthereum() *ethtypes.Log {
	topics := make([]common.Hash, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = common.HexToHash(topic)
	}

	return &ethtypes.Log{
		Address:     common.HexToAddress(log.Address),
		Topics:      topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      common.HexToHash(log.TxHash),
		TxIndex:     uint(log.TxIndex),
		Index:       uint(log.Index),
		BlockHash:   common.HexToHash(log.BlockHash),
		Removed:     log.Removed,
	}
}

func NewLogsFromEth(ethlogs []*ethtypes.Log) []*Log {
	var logs []*Log //nolint: prealloc
	for _, ethlog := range ethlogs {
		logs = append(logs, NewLogFromEth(ethlog))
	}

	return logs
}

// LogsToEthereum casts the Ethermint Logs to a slice of Ethereum Logs.
func LogsToEthereum(logs []*Log) []*ethtypes.Log {
	var ethLogs []*ethtypes.Log //nolint: prealloc
	for i := range logs {
		ethLogs = append(ethLogs, logs[i].ToEthereum())
	}
	return ethLogs
}

// NewLogFromEth creates a new Log instance from a Ethereum type Log.
func NewLogFromEth(log *ethtypes.Log) *Log {
	topics := make([]string, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = topic.String()
	}

	return &Log{
		Address:     log.Address.String(),
		Topics:      topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.String(),
		TxIndex:     uint64(log.TxIndex),
		Index:       uint64(log.Index),
		BlockHash:   log.BlockHash.String(),
		Removed:     log.Removed,
	}
}
