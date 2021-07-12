package graph

import "fmt"

var (
	// Get the node status
	queryDaemonStatus = `
		query {
			daemonStatus {
				syncStatus
				uptimeSecs
				stateHash
				commitId
				highestBlockLengthReceived
				blockchainLength
				numAccounts
			}
		}`

	// Get block details
	queryBlocks = `
		query {
			blocks(%s) {
				nodes {
					%s
				}
			}
		}`

	queryBestChain = `
		query {
			bestChain(maxLength: 290) {
				stateHash
				protocolState {
					consensusState {
						epoch
						epochCount
						slot
						blockHeight
					}
				}
			}
		}`

	queryBestTip = `
		query {
			bestChain(maxLength: 1) {
				stateHash
				protocolState {
					consensusState {
						epoch
						epochCount
						slot
						blockHeight
					}
				}
			}
		}`

	queryBlock = `
		query {
			block(stateHash: "%s") {
				%s
			}
		}
	`

	// Block details fields
	queryBlockFields = `
		stateHash
		stateHashField
		creator
		creatorAccount {
			publicKey
			delegate
			nonce
			votingFor
			balance {
				blockHeight
				total
				unknown
			}
		}
		protocolState {
			blockchainState {
				date
				utcDate
				stagedLedgerHash
				snarkedLedgerHash
			}
			consensusState {
				blockHeight
				blockchainLength
				epoch
				epochCount
				hasAncestorInSameCheckpointWindow
				lastVrfOutput
				totalCurrency
				minWindowDensity
				slot
				stakingEpochData {
					ledger {
						totalCurrency
					}
					epochLength
					lockCheckpoint
					seed
					startCheckpoint
				}
			}
			previousStateHash
		}
		snarkJobs {
			fee
			prover
			workIds
		}
		transactions {
			coinbase
			feeTransfer {
				fee
				recipient
				type
			}
		}
		winnerAccount {
			publicKey
			locked
		}`

	queryAccount = `
		query {
			account(publicKey: "%s") {
				nonce
				inferredNonce
				receiptChainHash
				delegate
				delegateAccount {
					publicKey
					delegate
					nonce
					votingFor
					balance {
						blockHeight
						total
						unknown
					}
				}
				votingFor
				locked
				balance {
					unknown
					total
					blockHeight
				}
			}
		}`

	queryPendingTx = `
		query {
			pooledUserCommands {
				amount
				fee
				from
				hash
				id
				isDelegation
				nonce
				memo
				kind
				to
			}
		}`
)

func buildBestChainQuery() string {
	return fmt.Sprintf(queryBestChain)
}

func buildBlocksQuery(filter string) string {
	return fmt.Sprintf(queryBlocks, filter, queryBlockFields)
}

func buildAccountQuery(filter string) string {
	return fmt.Sprintf(queryAccount, filter)
}
