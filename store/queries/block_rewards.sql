SELECT
  to_char(time_bucket, '$INTERVAL') AS interval,
  epoch,
  SUM(reward) AS amount
FROM
  block_rewards
WHERE
  owner_account = ?
  AND delegate = ?
  AND time_bucket BETWEEN ? AND ?
  AND owner_type = ?
GROUP BY
  to_char(time_bucket, '$INTERVAL'),
  epoch
ORDER BY
  to_char(time_bucket, '$INTERVAL')
