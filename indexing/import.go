package indexing

import (
	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store"
)

// Import creates new database records for the chain data
func Import(db *store.Store, data *Data) error {
	log.Debug("creating block")

	existing, err := db.Blocks.FindByHash(data.Block.Hash)
	if err != nil {
		if err != store.ErrNotFound {
			return err
		}
		existing = nil
	}

	if existing != nil {
		data.Block.ID = existing.ID
		err = db.Blocks.Update(data.Block)
	} else {
		err = db.Blocks.Create(data.Block)
	}
	if err != nil {
		return err
	}

	log.WithField("count", len(data.Accounts)).Debug("creating accounts")
	if err := db.Accounts.Import(data.Accounts); err != nil {
		return err
	}

	log.WithField("count", 1).Debug("creating validators")
	if err := db.Validators.Import([]model.Validator{*data.Validator}); err != nil {
		return err
	}

	log.WithField("count", len(data.ValidatorEpochs)).Debug("creating validator epochs")
	if err := db.ValidatorsEpochs.Import(data.ValidatorEpochs); err != nil {
		return err
	}

	log.WithField("count", len(data.Transactions)).Debug("creating transactions")
	if err := db.Transactions.Import(data.Transactions); err != nil {
		return err
	}

	log.WithField("count", len(data.Snarkers)).Debug("creating snarkers")
	if err := db.Snarkers.Import(data.Snarkers); err != nil {
		return err
	}

	log.WithField("count", len(data.SnarkJobs)).Debug("creating snarkjobs")
	if err := db.Jobs.Import(data.SnarkJobs); err != nil {
		return err
	}

	return nil
}
