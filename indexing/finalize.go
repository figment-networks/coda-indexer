package indexing

import (
	"github.com/figment-networks/mina-indexer/store"
	log "github.com/sirupsen/logrus"
)

// Finalize generates summary records
func Finalize(db *store.Store, data *Data) error {
	if err := db.Validators.UpdateStaking(); err != nil {
		return err
	}

	if err := db.Accounts.UpdateStaking(); err != nil {
		return err
	}

	ts := data.Block.Time
	buckets := []string{store.BucketHour, store.BucketDay}

	for _, bucket := range buckets {
		log.WithField("bucket", bucket).Debug("creating chain stats")
		if err := db.Stats.CreateChainStats(bucket, ts); err != nil {
			return err
		}

		log.WithField("bucket", bucket).Debug("creating validator stats")
		if err := db.Stats.CreateValidatorStats(data.Validator.PublicKey, bucket, ts); err != nil {
			return err
		}

		validators, err := db.Stats.FindValidatorsForDefaultStats(bucket, ts)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		for _, v := range validators {
			if err := db.Stats.CreateValidatorStats(v.PublicKey, bucket, ts); err != nil {
				return err
			}
		}

	}

	return nil
}
