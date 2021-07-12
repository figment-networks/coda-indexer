INSERT INTO block_rewards (
  owner_account,
  delegate,
  epoch,
  time_bucket,
  reward,
  owner_type
)
VALUES @values

ON CONFLICT (owner_account, delegate, epoch, time_bucket, owner_type) DO UPDATE
SET
  reward = COALESCE(block_rewards.reward,0) + excluded.reward
