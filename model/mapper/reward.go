package mapper

import (
	"strings"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

// ValidatorBlockReward returns validator reward models references from the block data
func ValidatorBlockReward(block model.Block) (*model.BlockReward, error) {
	t, err := RewardTimeBucket(block.Time)
	if err != nil {
		return nil, err
	}
	result := model.BlockReward{
		OwnerAccount: block.Creator,
		OwnerType:    string(model.RewardOwnerTypeValidator),
		Epoch:        block.Epoch,
	}
	result.TimeBucket = t
	return &result, nil
}

// DelegatorBlockRewards returns delegator reward models references from the block data
func DelegatorBlockRewards(accounts []model.Delegation, block model.Block) ([]model.BlockReward, error) {
	result := []model.BlockReward{}
	for _, a := range accounts {
		t, err := RewardTimeBucket(block.Time)
		if err != nil {
			return nil, err
		}
		// reward to be calculated next step
		dbr := model.BlockReward{
			OwnerAccount: a.PublicKey,
			Delegate:     a.Delegate,
			Epoch:        block.Epoch,
			OwnerType:    string(model.RewardOwnerTypeDelegator),
		}
		dbr.TimeBucket = t
		result = append(result, dbr)
	}
	return result, nil
}

func TransactionFees(block *graph.Block) types.Amount {
	res := types.NewInt64Amount(0)
	if block.Transactions != nil && block.Transactions.FeeTransfer != nil {
		for _, transfer := range block.Transactions.FeeTransfer {
			if strings.ToLower(transfer.Type) == model.TxTypeFeeTransfer {
				res = res.Add(types.NewAmount(transfer.Fee))
			} else if strings.ToLower(transfer.Type) == model.TxTypeCoinbaseFeeTransfer {
				res = res.Sub(types.NewAmount(transfer.Fee))
			}
		}
	}
	return res
}
