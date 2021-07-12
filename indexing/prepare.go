package indexing

import (
	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/mapper"
	"github.com/figment-networks/mina-indexer/model/types"
)

// Prepare generates a new models from the graph block data
func Prepare(archiveBlock *archive.Block, graphBlock *graph.Block, validatorEpochs []model.ValidatorEpoch) (*Data, error) {
	block, err := mapper.BlockFromArchive(archiveBlock)
	if err != nil {
		return nil, err
	}

	if graphBlock != nil {
		block.TotalCurrency = types.NewAmount(graphBlock.ProtocolState.ConsensusState.TotalCurrency)
		block.TransactionsFees = mapper.TransactionFees(graphBlock)
	}

	// Prepare validator record
	validator, err := mapper.Validator(archiveBlock)
	if err != nil {
		return nil, err
	}

	// Prepare transaction records
	transactions, err := mapper.TransactionsFromArchive(archiveBlock)
	if err != nil {
		return nil, err
	}
	block.TransactionsCount = len(transactions)

	// Prepare snarkers
	snarkers, err := mapper.Snarkers(graphBlock)
	if err != nil {
		return nil, err
	}
	block.SnarkersCount = len(snarkers)
	var snarkerAccounts []string
	for _, s := range snarkers {
		snarkerAccounts = append(snarkerAccounts, s.Account)
	}
	block.SnarkerAccounts = snarkerAccounts

	// Prepare snarker jobs
	snarkJobs, err := mapper.SnarkJobs(graphBlock)
	if err != nil {
		return nil, err
	}
	block.SnarkJobsCount = len(snarkJobs)
	block.SnarkJobsFees = types.NewInt64Amount(0)
	for _, job := range snarkJobs {
		block.SnarkJobsFees = block.SnarkJobsFees.Add(job.Fee)
	}

	// Prepare accounts
	accounts, err := mapper.Accounts(graphBlock)
	if err != nil {
		return nil, err
	}

	data := &Data{
		Block:           block,
		Validator:       validator,
		ValidatorEpochs: validatorEpochs,
		Accounts:        accounts,
		Transactions:    transactions,
		Snarkers:        snarkers,
		SnarkJobs:       snarkJobs,
	}

	return data, nil
}
