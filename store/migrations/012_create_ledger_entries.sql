-- +goose Up
CREATE TABLE ledger_entries (
  id                             SERIAL PRIMARY KEY,
  ledger_id                      INTEGER NOT NULL,
  public_key                     TEXT NOT NULL,
  delegate                       TEXT NOT NULL,
  delegation                     BOOLEAN DEFAULT FALSE,
  balance                        CHAIN_CURRENCY,
  timing_initial_minimum_balance CHAIN_CURRENCY,
  timing_cliff_time              INTEGER,
  timing_cliff_amount            CHAIN_CURRENCY,
  timing_vesting_period          INTEGER,
  timing_vesting_increment       INTEGER
);

CREATE INDEX idx_ledger_entries_ledger_id
  ON ledger_entries(ledger_id);

CREATE INDEX idx_ledger_entries_public_key
  ON ledger_entries(public_key);

CREATE INDEX idx_ledger_entries_delegate
  ON ledger_entries(delegate);

CREATE INDEX idx_ledger_entries_delegation
  ON ledger_entries(delegation);

-- +goose Down
DROP TABLE ledger_entries;
