package mapper

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model/util"
)

var (
	errNoProtocolState   = errors.New("no protocol state")
	errNoConsensusState  = errors.New("no consensus state")
	errNoBlockchainState = errors.New("no blockchain state")
)

// BlockHeight returns a parsed block height
func BlockHeight(input *graph.Block) uint64 {
	return util.MustUInt64(input.ProtocolState.ConsensusState.BlockHeight)
}

// BlockTime returns a parsed block time
func BlockTime(input *graph.Block) time.Time {
	return util.MustTime(input.ProtocolState.BlockchainState.Date)
}

func blockCheck(input *graph.Block) error {
	if input.ProtocolState == nil {
		return errNoProtocolState
	}
	if input.ProtocolState.ConsensusState == nil {
		return errNoConsensusState
	}
	if input.ProtocolState.BlockchainState == nil {
		return errNoBlockchainState
	}
	return nil
}

// RewardTimeBucket daily time bucket for rewards
func RewardTimeBucket(input time.Time) (time.Time, error) {
	t, err := time.Parse("2006-01-02", input.Format("2006-01-02"))
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}
