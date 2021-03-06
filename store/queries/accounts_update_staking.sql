WITH staking AS (
  SELECT delegate, SUM(balance) AS total, COUNT(1) AS delegations
  FROM ledger_entries
  WHERE ledger_id = (SELECT id FROM ledgers ORDER BY id DESC LIMIT 1)
  GROUP BY delegate
)
UPDATE accounts
SET stake = staking.total
FROM staking
WHERE accounts.public_key = staking.delegate