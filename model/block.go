package model

import (
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/figment-networks/mina-indexer/model/types"
)

// Block model contains block data
type Block struct {
	ID                int            `json:"-"`
	Height            uint64         `json:"height"`
	Hash              string         `json:"hash"`
	ParentHash        string         `json:"parent_hash"`
	Time              time.Time      `json:"time"`
	Canonical         bool           `json:"canonical"`
	LedgerHash        string         `json:"ledger_hash"`
	SnarkedLedgerHash string         `json:"snarked_ledger_hash"`
	Creator           string         `json:"creator"`
	Coinbase          types.Amount   `json:"coinbase"`
	TotalCurrency     types.Amount   `json:"total_currency"`
	Epoch             int            `json:"epoch"`
	Slot              int            `json:"slot"`
	TransactionsCount int            `json:"transactions_count"`
	TransactionsFees  types.Amount   `json:"transactions_fees"`
	SnarkersCount     int            `json:"snarkers_count"`
	SnarkerAccounts   pq.StringArray `json:"snarker_accounts"`
	SnarkJobsCount    int            `json:"snark_jobs_count"`
	SnarkJobsFees     types.Amount   `json:"snark_jobs_fees"`
	RewardCalculated  bool           `json:"reward_calculated"`
	Supercharged      bool           `json:"supercharged"`
}

// BlockIntervalStat contains block count stats for a given time interval
type BlockIntervalStat struct {
	TimeInterval string  `json:"time_interval"`
	Count        int64   `json:"count"`
	Avg          float64 `json:"avg"`
}

// BlockAvgStat contains block averagess
type BlockAvgStat struct {
	StartHeight int64   `json:"start_height"`
	EndHeight   int64   `json:"end_height"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Count       int64   `json:"count"`
	Diff        float64 `json:"diff"`
	Avg         float64 `json:"avg"`
}

// TableName returns the model table name
func (Block) TableName() string {
	return "blocks"
}

// Validate returns an error if block data is invalid
func (b Block) Validate() error {
	if b.Time.IsZero() {
		return errors.New("time is invalid")
	}
	if b.Height == 0 {
		return errors.New("height is invalid")
	}
	if b.Hash == "" {
		return errors.New("hash is required")
	}
	if b.ParentHash == "" {
		return errors.New("parent hash is required")
	}
	if b.LedgerHash == "" {
		return errors.New("ledger hash is required")
	}
	if b.Creator == "" {
		return errors.New("creator is required")
	}
	return nil
}
