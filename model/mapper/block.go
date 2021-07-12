package mapper

import (
	"time"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

const superchargedReward = 1440000000000

func BlockFromArchive(input *archive.Block) (*model.Block, error) {
	block := &model.Block{
		Canonical:         false,
		Height:            input.Height,
		Time:              time.Unix(input.Timestamp/1000, 0),
		Creator:           input.Creator,
		Hash:              input.StateHash,
		ParentHash:        input.ParentHash,
		LedgerHash:        input.LedgerHash,
		SnarkedLedgerHash: input.SnarkedLedgerHash,
		Epoch:             int(input.GlobalSlot) / 7140,
		Slot:              int(input.GlobalSlot),
		TransactionsCount: len(input.UserCommands) + len(input.InternalCommands),
	}

	for _, cmd := range input.InternalCommands {
		if cmd.Type == model.TxTypeCoinbase {
			block.Coinbase = types.NewInt64Amount(cmd.Fee)
			break
		}
	}

	if block.Coinbase.Int != nil {
		// The coinbase reward for producing a block is 720 tokens. for supercharged 2x
		block.Supercharged = block.Coinbase.Int64() == superchargedReward
	}

	return block, block.Validate()
}
