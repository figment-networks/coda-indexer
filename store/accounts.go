package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// AccountsStore handles operations on accounts
type AccountsStore struct {
	baseStore
}

func (s AccountsStore) Count() (int, error) {
	var n int
	err := s.db.Table("accounts").Count(&n).Error
	return n, err
}

// FindBy returns an account for a matching attribute
func (s AccountsStore) FindBy(key string, value interface{}) (*model.Account, error) {
	result := &model.Account{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns an account for the ID
func (s AccountsStore) FindByID(id int64) (*model.Account, error) {
	return s.FindBy("id", id)
}

// FindByPublicKey returns an account for the public key
func (s AccountsStore) FindByPublicKey(key string) (*model.Account, error) {
	return s.FindBy("public_key", key)
}

// AllByDelegator returns all accounts delegated to another account
func (s AccountsStore) AllByDelegator(account string) ([]model.Account, error) {
	result := []model.Account{}
	err := s.db.
		Where("delegate = ?", account).
		Find(&result).
		Error
	return result, checkErr(err)
}

// ByHeight returns all accounts that were created at a given height
func (s AccountsStore) ByHeight(height int64) ([]model.Account, error) {
	result := []model.Account{}

	err := s.db.
		Where("start_height <= ?", height).
		Order("id DESC").
		Find(&result).
		Error

	return result, checkErr(err)
}

// All returns all accounts
func (s AccountsStore) All() ([]model.Account, error) {
	result := []model.Account{}

	err := s.db.
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s AccountsStore) UpdateStaking() error {
	return s.db.Exec(queries.AccountsUpdateStaking).Error
}

func (s AccountsStore) Import(records []model.Account) error {
	n := len(records)
	if n == 0 {
		return nil
	}

	batchSize := 250

	for idx := 0; idx < n; idx += batchSize {
		endIdx := idx + batchSize
		if endIdx > n {
			endIdx = n
		}

		batch := records[idx:endIdx]

		err := bulk.Import(s.db, queries.AccountsImport, len(batch), func(rowIdx int) bulk.Row {
			acc := batch[rowIdx]
			now := time.Now()

			return bulk.Row{
				acc.PublicKey,
				acc.Delegate,
				acc.Balance,
				acc.BalanceUnknown,
				acc.Nonce,
				acc.StartHeight,
				acc.StartTime,
				acc.LastHeight,
				acc.LastTime,
				now,
				now,
			}
		})

		if err != nil {
			return err
		}
	}

	return nil
}
