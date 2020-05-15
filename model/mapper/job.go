package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Job returns a job model constructed from the coda input
func Job(block coda.Block, w *coda.CompletedWork) (*model.Job, error) {
	j := &model.Job{
		Height: blockHeight(block),
		Time:   blockTime(block),
		Prover: w.Prover,
		Fee:    util.MustInt64(w.Fee),
	}
	return j, j.Validate()
}

// Jobs returns list of jobs constructed from the coda input
func Jobs(block coda.Block) ([]model.Job, error) {
	result := []model.Job{}

	for _, w := range block.SnarkJobs {
		j, err := Job(block, w)
		if err != nil {
			return nil, err
		}
		result = append(result, *j)
	}

	return result, nil
}