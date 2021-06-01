package util

import (
	"math/big"
	"testing"

	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/stretchr/testify/assert"
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
			res, err := CalculateWeight(tt.args.balance, tt.args.totalStakedBalance)
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
		weight       big.Float
		blockReward  types.Amount
		validatorFee types.Percentage
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
				weight:       *w,
				blockReward:  types.NewInt64Amount(100),
				validatorFee: types.NewPercentage("5"),
			},
			result: types.NewPercentage("28.5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateDelegatorReward(tt.args.weight, tt.args.blockReward, tt.args.validatorFee)
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
			res, err := CalculateValidatorReward(tt.args.blockReward, tt.args.validatorFee)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}
