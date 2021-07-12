INSERT INTO validator_epochs (
  account_id,
  account_address,
  epoch,
  validator_fee
)
VALUES
  @values
ON CONFLICT (account_id, epoch) DO UPDATE
SET
  validator_fee   = excluded.validator_fee