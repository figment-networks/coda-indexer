package model

import (
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

type Ledger struct {
	ID                int          `json:"-"`
	Time              time.Time    `json:"time"`
	Epoch             int          `json:"epoch"`
	EntriesCount      int          `json:"entries_count"`
	StakedAmount      types.Amount `json:"staked_amount"`
	DelegationsCount  int          `json:"delegations_count"`
	DelegationsAmount types.Amount `json:"delegations_amount"`
}

func (Ledger) TableName() string {
	return "ledgers"
}

type LedgerEntry struct {
	ID                          int          `json:"-"`
	LedgerID                    int          `json:"-"`
	PublicKey                   string       `json:"public_key"`
	Delegate                    string       `json:"delegate"`
	Delegation                  bool         `json:"delegation"`
	Balance                     types.Amount `json:"balance"`
	TimingInitialMinimumBalance types.Amount `json:"timing_initial_minimum_balance"`
	TimingCliffTime             *int         `json:"timing_cliff_time"`
	TimingCliffAmount           types.Amount `json:"timing_cliff_amount"`
	TimingVestingPeriod         *int         `json:"timing_vesting_period"`
	TimingVestingIncrement      *int         `json:"timing_vesting_increment"`
}

func (LedgerEntry) TableName() string {
	return "ledger_entries"
}

func (l LedgerEntry) IsUntimed() bool {
	return l.TimingInitialMinimumBalance.Int == nil && l.TimingCliffTime == nil &&
		l.TimingCliffAmount.Int == nil && l.TimingVestingPeriod == nil && l.TimingVestingIncrement == nil
}
