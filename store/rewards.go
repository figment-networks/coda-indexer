package store

import (
	"strings"
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// RewardStore handles operations on rewards
type RewardStore struct {
	baseStore
}

// FetchRewardsByInterval fetches rewards by interval
func (s *RewardStore) FetchRewardsByInterval(ownerAccount string, delegate string, from time.Time, to time.Time, timeInterval model.TimeInterval, rewardOwnerType model.RewardOwnerType) ([]model.RewardsSummary, error) {
	var res []model.RewardsSummary
	q := strings.Replace(queries.BlockRewards, "$INTERVAL", timeInterval.String(), -1)
	var err error
	if delegate == "" {
		q = strings.Replace(q, "AND delegate = ?", "", -1)
		err = s.db.Raw(q, ownerAccount, from, to, rewardOwnerType).Scan(&res).Error
	} else {
		err = s.db.Raw(q, ownerAccount, delegate, from, to, rewardOwnerType).Scan(&res).Error
	}
	return res, err
}

// Import creates new rewards
func (s RewardStore) Import(records []model.BlockReward) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.BlockRewardsImport, len(records), func(idx int) bulk.Row {
		tx := records[idx]
		return bulk.Row{
			tx.OwnerAccount,
			tx.Delegate,
			tx.Epoch,
			tx.TimeBucket,
			tx.Reward,
			tx.OwnerType,
		}
	})
}
