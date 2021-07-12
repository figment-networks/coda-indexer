package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

const batchSize = 100

// StakingStore handles operations on staking data
type StakingStore struct {
	baseStore
}

// CreateLedger creates a new ledger record
func (s StakingStore) CreateLedger(ledger *model.Ledger) error {
	return s.Create(ledger)
}

// CreateLedgerEntries create a batch of ledger entries
func (s StakingStore) CreateLedgerEntries(records []model.LedgerEntry) error {
	var err error
	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = bulk.Import(s.db, queries.LedgerImportEntries, j-i, func(k int) bulk.Row {
			r := records[i+k]

			return bulk.Row{
				r.LedgerID,
				r.PublicKey,
				r.Delegate,
				r.Delegation,
				r.Balance,
				r.TimingInitialMinimumBalance,
				r.TimingCliffTime,
				r.TimingCliffAmount,
				r.TimingVestingPeriod,
				r.TimingVestingIncrement,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// FindLedger returns the most recent ledger of an epoch
func (s StakingStore) FindLedger(epoch int) (*model.Ledger, error) {
	ledger := &model.Ledger{}

	err := s.db.
		Model(ledger).
		Where("epoch = ?", epoch).
		First(ledger).
		Error

	if err != nil {
		ledger = nil
	}

	return ledger, checkErr(err)
}

// AllLedgers returns all existing ledgers
func (s StakingStore) AllLedgers() ([]model.Ledger, error) {
	result := []model.Ledger{}

	err := s.db.
		Model(&model.Ledger{}).
		Order("epoch ASC").
		Find(&result).
		Error

	return result, err
}

// LastLedger returns the most recent ledger record
func (s StakingStore) LastLedger() (*model.Ledger, error) {
	ledger := &model.Ledger{}

	err := s.db.
		Model(ledger).
		Order("id DESC").
		First(ledger).
		Error

	if err == ErrNotFound {
		ledger = nil
	}

	return ledger, checkErr(err)
}

type FindDelegationsParams struct {
	LedgerID  *int
	PublicKey string
	Delegate  string
}

// LedgerRecords returns all ledger records from current epoch
func (s StakingStore) LedgerRecords(ledgerID int) ([]model.LedgerEntry, error) {
	result := []model.LedgerEntry{}

	err := s.db.
		Model(&model.LedgerEntry{}).
		Where("ledger_id = ?", ledgerID).
		Find(&result).
		Error

	return result, checkErr(err)
}

// LedgerRecordsOfDelegate returns delegations' ledger info
func (s StakingStore) LedgerRecordsOfDelegate(ledgerID int, delegate string) ([]model.LedgerEntry, error) {
	result := []model.LedgerEntry{}

	err := s.db.
		Model(&model.LedgerEntry{}).
		Where("ledger_id = ? AND delegate = ? AND delegation = ?", ledgerID, delegate, true).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindDelegations returns delegations for a given ledger ID
func (s StakingStore) FindDelegations(params FindDelegationsParams) ([]model.Delegation, error) {
	result := []model.Delegation{}

	if params.LedgerID == nil {
		ledger, err := s.LastLedger()
		if err != nil && err != ErrNotFound {
			return result, err
		}
		params.LedgerID = &ledger.ID
	}

	scope := s.db.
		Table("ledger_entries").
		Where("ledger_id = ?", params.LedgerID).
		Where("delegation = ?", true)

	if params.Delegate != "" {
		scope = scope.Where("delegate = ?", params.Delegate)
	}
	if params.PublicKey != "" {
		scope = scope.Where("public_key = ?", params.PublicKey)
	}

	err := scope.Find(&result).Error

	return result, checkErr(err)
}
