package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// TransactionsStore handles operations on transactions
type TransactionsStore struct {
	baseStore
}

// FindBy returns transactions by a given key and value
func (s TransactionsStore) FindBy(key string, value interface{}) (*model.Transaction, error) {
	result := &model.Transaction{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a transaction for a given ID
func (s TransactionsStore) FindByID(id int64) (*model.Transaction, error) {
	return s.FindBy("id", id)
}

// FindByHash returns a transaction for a given hash
func (s TransactionsStore) FindByHash(hash string) (*model.Transaction, error) {
	return s.FindBy("hash", hash)
}

// ByAccount returns a list of transactions sent or received by the account
func (s TransactionsStore) ByAccount(account string) ([]model.Transaction, error) {
	var canonical = true
	return s.Search(TransactionSearch{Account: account, Canonical: &canonical})
}

// ByHeight returns transactions for a given height
func (s TransactionsStore) ByHeight(height uint64, limit uint) ([]model.Transaction, error) {
	var canonical = true
	return s.Search(TransactionSearch{Height: height, Limit: limit, Canonical: &canonical})
}

// Search returns a list of transactions that matches the filters
func (s TransactionsStore) Search(search TransactionSearch) ([]model.Transaction, error) {
	scope := s.db.
		Order("time DESC").
		Limit(search.Limit)

	if search.BeforeID > 0 {
		scope = scope.Where("id < ?", search.BeforeID)
	}
	if search.AfterID > 0 {
		scope = scope.Where("id > ?", search.AfterID)
	}
	if search.BlockHash != "" {
		scope = scope.Where("block_hash = ?", search.BlockHash)
	}
	if search.Height > 0 {
		scope = scope.Where("block_height = ?", search.Height)
	}
	if search.Type != "" {
		scope = scope.Where("type IN (?)", strings.Split(search.Type, ","))
	}
	if search.Account != "" {
		scope = scope.Where("sender = ? OR receiver = ?", search.Account, search.Account)
	} else {
		if search.Sender != "" {
			scope = scope.Where("sender = ?", search.Sender)
		}
		if search.Receiver != "" {
			scope = scope.Where("receiver = ?", search.Receiver)
		}
	}
	if len(search.Memo) > 2 {
		scope = scope.Where("memo ILIKE ?", fmt.Sprintf("%%%s%%", search.Memo))
	}
	if search.Status != "" {
		scope = scope.Where("status = ?", search.Status)
	}
	if search.startTime != nil {
		scope = scope.Where("time >= ?", search.startTime)
	}
	if search.endTime != nil {
		scope = scope.Where("time <= ?", search.endTime)
	}
	if search.Canonical != nil {
		scope = scope.Where("canonical = ?", *search.Canonical)
	}

	result := []model.Transaction{}
	err := scope.Find(&result).Error

	return result, err
}

func (s TransactionsStore) Import(records []model.Transaction) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.TransactionsImport, len(records), func(idx int) bulk.Row {
		tx := records[idx]
		now := time.Now()

		return bulk.Row{
			tx.Type,
			tx.Hash,
			tx.BlockHash,
			tx.BlockHeight,
			tx.Time,
			tx.Nonce,
			tx.Sender,
			tx.Receiver,
			tx.Amount,
			tx.Fee,
			tx.Memo,
			tx.Status,
			tx.Canonical,
			tx.FailureReason,
			tx.SequenceNumber,
			tx.SecondarySequenceNumber,
			now,
			now,
		}
	})
}

// MarkTransactionsOrphan updates all transactions as non canonical at a height
func (s TransactionsStore) MarkTransactionsOrphan(height uint64) error {
	return s.db.Exec(queries.MarkTransactionsOrphan, height).Error
}

// MarkTransactionsCanonical updates transactions canonical for given block hash
func (s TransactionsStore) MarkTransactionsCanonical(blockHash string) error {
	return s.db.Exec(queries.MarkTransactionsCanonical, blockHash).Error
}
