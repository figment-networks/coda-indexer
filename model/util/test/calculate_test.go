package test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/model/util"
)

func TestCalculateWeight(t *testing.T) {
	type args struct {
		balance            types.Amount
		totalStakedBalance types.Amount
	}
	tests := []struct {
		name    string
		args    args
		result  types.Percentage
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				balance:            types.NewInt64Amount(10),
				totalStakedBalance: types.NewInt64Amount(10000),
			},
			result: types.NewPercentage("0.001"),
		},
		{
			name: "error case stake value",
			args: args{
				balance:            types.NewInt64Amount(10),
				totalStakedBalance: types.NewInt64Amount(0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CalculateWeight(tt.args.balance, tt.args.totalStakedBalance)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateDelegatorReward(t *testing.T) {
	w, _ := new(big.Float).SetString("0.3")
	type args struct {
		weight          big.Float
		remainingReward types.Percentage
	}
	tests := []struct {
		name    string
		args    args
		result  types.Percentage
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				weight:          *w,
				remainingReward: types.NewPercentage("95"),
			},
			result: types.NewPercentage("28.5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CalculateDelegatorReward(tt.args.weight, tt.args.remainingReward)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateValidatorReward(t *testing.T) {
	type args struct {
		blockReward  types.Amount
		validatorFee types.Percentage
	}
	tests := []struct {
		name    string
		args    args
		result  types.Amount
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				blockReward:  types.NewInt64Amount(100),
				validatorFee: types.NewPercentage("5"),
			},
			result: types.NewAmount("5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CalculateValidatorReward(tt.args.blockReward, tt.args.validatorFee)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateSuperchargedWeighting(t *testing.T) {
	type args struct {
		block model.Block
	}
	tests := []struct {
		name    string
		args    args
		result  types.Percentage
		wantErr bool
	}{
		{
			name: "successful zero transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("0"),
				},
			},
			result: types.NewPercentage("2"),
		},
		{
			name: "successful same amount transaction fee and coinbase",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("200"),
				},
			},
			result: types.NewPercentage("1.5"),
		},
		{
			name: "successful low transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("5"),
				},
			},
			result: types.NewPercentage("1.975609756"),
		},
		{
			name: "successful high transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("10000"),
				},
			},
			result: types.NewPercentage("1.019607843"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CalculateSuperchargedWeighting(tt.args.block)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateWeightsSupercharged(t *testing.T) {
	ct := 86400
	vp := 1
	vi := 0

	records := []model.LedgerEntry{
		{
			PublicKey: "one",
			Balance:   types.NewAmount("20000"),
		},
		{
			PublicKey: "two",
			Balance:   types.NewAmount("50000"),
		},
		{
			PublicKey:                   "three",
			Balance:                     types.NewAmount("30000"),
			TimingInitialMinimumBalance: types.NewAmount("230400"),
			TimingCliffAmount:           types.NewAmount("230400"),
			TimingCliffTime:             &ct,
			TimingVestingPeriod:         &vp,
			TimingVestingIncrement:      &vi,
		},
	}

	delegations := []model.Delegation{
		{
			PublicKey: "one",
			Balance:   types.NewAmount("20000"),
		},
		{
			PublicKey: "two",
			Balance:   types.NewAmount("50000"),
		},
		{
			PublicKey: "three",
			Balance:   types.NewAmount("30000"),
		},
	}

	type args struct {
		superchargedContribution types.Percentage
		delegations              []model.Delegation
		records                  []model.LedgerEntry
		firstSlotOfEpoch         int
	}
	tests := []struct {
		name    string
		args    args
		result  []string
		wantErr bool
	}{
		{
			name: "successful with timing but locked entire epoch",
			args: args{
				superchargedContribution: types.NewPercentage("1.98765"),
				delegations:              delegations,
				records:                  records,
				firstSlotOfEpoch:         70000,
			},
			result: []string{"0.2350364057", "0.5875910143", "0.17737258"},
		},
		{
			name: "successful with timing but locked entire epoch",
			args: args{
				superchargedContribution: types.NewPercentage("1.98765"),
				delegations:              delegations,
				records:                  records,
				firstSlotOfEpoch:         80000,
			},
			result: []string{"0.2308451533", "0.5771128832", "0.1920419636"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegations := tt.args.delegations
			err := util.CalculateWeightsSupercharged(tt.args.superchargedContribution, delegations, tt.args.records, tt.args.firstSlotOfEpoch)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				for i, _ := range records {
					assert.Equal(t, tt.result[i], delegations[i].Weight.String())
				}
			}
		})
	}
}
