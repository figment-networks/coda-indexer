package store

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/mina-indexer/model"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Blocks           BlocksStore
	Accounts         AccountsStore
	Validators       ValidatorsStore
	ValidatorsEpochs ValidatorsEpochsStore
	Transactions     TransactionsStore
	Jobs             JobsStore
	Snarkers         SnarkersStore
	Stats            StatsStore
	Staking          StakingStore
	Rewards          RewardStore
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// Conn returns an underlying database connection
func (s *Store) Conn() *sql.DB {
	return s.db.DB()
}

// SetDebugMode enabled detailed query logging
func (s *Store) SetDebugMode(enabled bool) {
	s.db.LogMode(enabled)
}

// New returns a new store from the connection string
func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: conn,

		Blocks:           NewBlocksStore(conn),
		Accounts:         NewAccountsStore(conn),
		Validators:       NewValidatorsStore(conn),
		ValidatorsEpochs: NewValidatorsEpochsStore(conn),
		Transactions:     NewTransactionsStore(conn),
		Snarkers:         NewSnarkersStore(conn),
		Jobs:             NewJobsStore(conn),
		Stats:            NewStatsStore(conn),
		Staking:          NewStakingStore(conn),
		Rewards:          NewRewardStore(conn),
	}, nil
}

func NewBlocksStore(db *gorm.DB) BlocksStore {
	return BlocksStore{scoped(db, model.Block{})}
}

func NewAccountsStore(db *gorm.DB) AccountsStore {
	return AccountsStore{scoped(db, model.Account{})}
}

func NewValidatorsStore(db *gorm.DB) ValidatorsStore {
	return ValidatorsStore{scoped(db, model.Validator{})}
}

func NewValidatorsEpochsStore(db *gorm.DB) ValidatorsEpochsStore {
	return ValidatorsEpochsStore{scoped(db, model.ValidatorEpoch{})}
}

func NewTransactionsStore(db *gorm.DB) TransactionsStore {
	return TransactionsStore{scoped(db, model.Transaction{})}
}

func NewSnarkersStore(db *gorm.DB) SnarkersStore {
	return SnarkersStore{scoped(db, model.Snarker{})}
}

func NewJobsStore(db *gorm.DB) JobsStore {
	return JobsStore{scoped(db, model.SnarkJob{})}
}

func NewStatsStore(db *gorm.DB) StatsStore {
	return StatsStore{baseStore{db: db}}
}

func NewStakingStore(db *gorm.DB) StakingStore {
	return StakingStore{scoped(db, nil)}
}

func NewRewardStore(db *gorm.DB) RewardStore {
	return RewardStore{scoped(db, model.BlockReward{})}
}
