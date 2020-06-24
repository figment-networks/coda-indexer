-- +goose Up
CREATE TYPE e_tx_type AS ENUM (
  'payment',
  'delegation',
  'snark_fee',
  'block_reward'
);

CREATE TABLE IF NOT EXISTS transactions (
  id         SERIAL NOT NULL,
  type       e_tx_type NOT NULL,
  hash       TEXT NOT NULL,
  block_hash TEXT NOT NULL,
  height     DOUBLE PRECISION NOT NULL,
  time       TIMESTAMP WITH TIME ZONE NOT NULL,
  nonce      NUMERIC,
  sender     TEXT NOT,
  receiver   TEXT NOT NULL,
  amount     DECIMAL(65, 0) NOT NULL,
  fee        DECIMAL(65, 0),
  memo       TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

CREATE INDEX idx_transactions_type
  ON transactions(type);

CREATE INDEX idx_transactions_block_hash
  ON transactions(block_hash);

CREATE UNIQUE INDEX idx_transactions_hash
  ON transactions(hash);

CREATE INDEX idx_transactions_height
  ON transactions(height);

CREATE INDEX idx_transactions_time
  ON transactions(time);

CREATE INDEX idx_transactions_sender
  ON transactions(sender);

CREATE INDEX idx_transactions_receiver
  ON transactions(receiver);

-- +goose Down
DROP TABLE transactions;
