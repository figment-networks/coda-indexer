# mina-indexer

Blockchain data indexer and API service for Mina blockchain protocol

## Requirements

- PostgreSQL 10.x+
- Go 1.15+

## Installation

*Not available yet*

## Configuration

You can configure the service using either a config file or environment variables.

### Config File

Example:

```json
{
  "app_env": "production",
  "mina_endpoint": "http://YOUR_NODE_IP:3085/graphql",
  "archive_endpoint": "http://YOUR_NODE_IP:3088",
  "server_addr": "127.0.0.1",
  "server_port": 8080,
  "database_url": "postgres://user:pass@host/dbname",
  "sync_interval": "30s",
  "cleanup_interval": "10m"
}
```

### Environment Variables

| Name               | Description             | Default
|--------------------|-------------------------|-------------------
| `DATABASE_URL`     | PostgreSQL database URL
| `MINA_ENDPOINT`    | Mina GraphQL Endpoint
| `ARCHIVE_ENDPOINT` | Mina Archive API Endpoint
| `APP_ENV`          | Application environment | `development`
| `SERVER_ADDR`      | Server listen address   | `0.0.0.0`
| `SERVER_PORT`      | Server listen port      | `8080`
| `SYNC_INTERVAL`    | Data sync interval      | `10s`
| `CLEANUP_INTERVAL` | Data cleanup interval   | `10min`
| `LOG_LEVEL`        | Application log level   | `info`
| `LOG_FORMAT`       | Application log format  | `text`. Available: `text`, `json`

## Running Application

Once you have created a database and specified all configuration options, you
need to migrate the database. You can do that by running the command below:

```bash
mina-indexer -config path/to/config.json -cmd=migrate
```

Start the data indexer:

```bash
mina-indexer -config path/to/config.json -cmd=worker
```

Start the API server:

```bash
mina-indexer -config path/to/config.json -cmd=server
```

## API Reference

| Method | Path                            | Description
|--------|---------------------------------|------------------------------------
| GET    | /health                         | Healthcheck endpoint
| GET    | /height                         | Current indexed blockchain height
| GET    | /blocks                         | Blocks search
| GET    | /blocks/:hash                   | Block details by ID or Hash
| GET    | /block_times                    | Block times stats
| GET    | /block_times_interval           | Block creation stats
| GET    | /transactions                   | Transactions search
| GET    | /pending_transactions           | Pending Transactions
| GET    | /transactions/:id               | Transaction details by ID or Hash
| GET    | /accounts                       | Accounts search
| GET    | /accounts/:id                   | Account details by ID or Key
| GET    | /snarkers                       | All existing snarkers from all blocks(including non-canonical)
| GET    | /snarker/:id                    | Snarker info from canonical blocks
| GET    | /rewards/:id                    | Reward details

## How do we calculate rewards?
We wait a certain number of blocks to calculate rewards, also known as transaction finality.
(Please see https://learn.figment.io/network-documentation/mina/mina-101#transaction-speed-and-finality) so 15, which provides a 99.9% confidence that the transaction will not be reversed, is the threshold for calculations. That means we only show the rewards calculated 45 minute before.

We sum rewards daily based on each block rewards one by one. Here is the formula of a block's reward;

```
block reward = coinbase reward + transaction fees - snark job fees

validator reward = block reward * validator fee / 100

delegators reward = block reward - validator reward
```
For calculating delegator rewards we have 2 different algorithms based on coinbase reward.

1. The coinbase reward for producing a block is 720 tokens. At this situation here is the formula;
```
total staked amount = sum of delegations

foreach delegator 
      weight = delegation amount / total staked amount
      reward = delegators reward / weight
```
2. Supercharged rewards when coinbase reward for producing a block is 1440 tokens. At this case calculating weights different;
```
supercharged weighting = 1 + (1 / (1 + transaction fees / coinbase))

foreach delegator
      timed weighting = 0 for if entry is unlocked entire epoch
                        1 for if entry is locked entire epoch 
                        proportion calculated based on slots if it becomes unlocked during epoch
      supercharged contribution = ((supercharged weighting - 1) * timed weighting factor) + 1
      pool stake = stake * supercharged contribution

foreach delegator
      pool weighting = pool stake / sum(pool stakes)

foreach delegator 
      weight = pool stake / sum(pool stakes)
      reward = delegators reward / weight
```
