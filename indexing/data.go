package indexing

import "github.com/figment-networks/coda-indexer/model"

// Data contains all the records processed for a height
type Data struct {
	Block        *model.Block
	Validator    *model.Validator
	Accounts     []model.Account
	Snarkers     []model.Snarker
	Transactions []model.Transaction
	FeeTransfers []model.FeeTransfer
	SnarkJobs    []model.Job
}
