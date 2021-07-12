package util

import (
	"errors"
	"math"
	"math/big"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

const slotsPerEpoch = 7140

var zero = types.NewInt64Amount(0)

// CalculateWeight calculates weight
func CalculateWeight(balance types.Amount, totalStakedAmount types.Amount) (types.Percentage, error) {
	if totalStakedAmount.Compare(zero) == 0 {
		return types.Percentage{}, errors.New("total staked amount can not be zero")
	}
	w := types.NewPercentage(balance.String())
	t := types.NewPercentage(totalStakedAmount.String())
	return w.Quo(t), nil
}

// CalculateValidatorReward calculates validator reward
func CalculateValidatorReward(blockReward types.Amount, validatorFee types.Percentage) (types.Percentage, error) {
	vr := types.NewPercentage(blockReward.String())
	fee := validatorFee.Quo(types.NewPercentage("100"))
	vr = vr.Mul(fee)
	result := new(big.Int)
	vr.Int(result)
	return types.NewPercentage(result.String()), nil
}

// CalculateDelegatorReward calculates delegator reward
func CalculateDelegatorReward(weight big.Float, remainingReward types.Percentage) (types.Percentage, error) {
	w := types.NewPercentage(weight.String())
	return w.Mul(remainingReward), nil
}

// CalculateSuperchargedWeighting calculates supercharged weighting for given block
// supercharged weighting = 1 + (1 / (1 + transaction fees / coinbase))
func CalculateSuperchargedWeighting(block model.Block) (types.Percentage, error) {
	trFees := types.NewPercentage(block.TransactionsFees.String())
	coinbase := types.NewPercentage(block.Coinbase.String())
	denom := trFees.Quo(coinbase)
	denom = denom.Add(types.NewPercentage("1"))
	enum := types.NewPercentage("1")
	enum = enum.Quo(denom)
	res := enum.Add(types.NewPercentage("1"))
	return res, nil
}

// CalculateWeightsNonSupercharged calculates weights when block reward is not doubled for supercharged
func CalculateWeightsNonSupercharged(delegations []model.Delegation) error {
	stakedAmount := new(big.Int)
	bln := new(big.Int)
	var ok bool
	for _, r := range delegations {
		bln, ok = bln.SetString(r.Balance.String(), 10)
		if !ok {
			return errors.New("error with balance")
		}
		stakedAmount.Add(stakedAmount, bln)
	}
	for i, r := range delegations {
		w, err := CalculateWeight(r.Balance, types.NewAmount(stakedAmount.String()))
		if err != nil {
			return err
		}
		delegations[i].Weight = w
	}
	return nil
}

// CalculateWeightsSupercharged calculates weights when block reward is doubled for supercharged
func CalculateWeightsSupercharged(superchargedWeighting types.Percentage, delegations []model.Delegation, records []model.LedgerEntry, firstSlotOfEpoch int) error {
	sum := new(big.Float)
	sc := new(big.Float)
	bln := new(big.Float)
	stk := new(big.Float)
	w := new(big.Float)
	var ok bool

	recordsMap := map[string]model.LedgerEntry{}
	for _, r := range records {
		recordsMap[r.PublicKey] = r
	}

	// pool stakes
	for i, r := range delegations {
		record, ok := recordsMap[r.PublicKey]
		if !ok {
			return errors.New("ledger record not found")
		}

		timedWeighting, err := calculateTimedWeighting(record, firstSlotOfEpoch)
		if err != nil {
			return err
		}

		superchargedContribution, err := calculateSuperchargedContribution(superchargedWeighting, timedWeighting)
		if err != nil {
			return err
		}

		sc, ok = sc.SetString(superchargedContribution.String())
		if !ok {
			return errors.New("error with supercharged contribution")
		}
		bln, ok := bln.SetString(r.Balance.String())
		if !ok {
			return errors.New("error with balance")
		}
		stk = bln.Mul(bln, sc)

		delegations[i].Weight = types.NewPercentage(stk.String())
		sum.Add(sum, stk)
	}

	// pool weights
	for i, r := range delegations {
		w, ok = new(big.Float).SetString(r.Weight.String())
		if !ok {
			return errors.New("error with weight")
		}
		w = w.Quo(w, sum)
		delegations[i].Weight = types.NewPercentage(w.String())
	}

	return nil
}

// calculateTimedWeighting calculates timed weighting
//
// return 0 for if entry is unlocked entire epoch,
//        1 for if entry is locked entire epoch,
//        proportion calculated based on slots if it becomes unlocked during epoch
func calculateTimedWeighting(record model.LedgerEntry, firstSlotOfEpoch int) (types.Percentage, error) {
	if record.IsUntimed() {
		return types.NewFloat64Percentage(1), nil
	}

	// unlockedTime: global slot at which that account's tokens will be fully unlocked
	//
	// vesting_amount = initial_min_balance - cliff_amount
	// vesting_time = vesting_amount * math.ceil(vesting_period / vesting_increment)
	// unlocked_time = cliff_time + vesting_time
	imb := types.NewAmount(record.TimingInitialMinimumBalance.String())
	ca := types.NewAmount(record.TimingCliffAmount.String())
	ct := types.NewInt64Amount(int64(*record.TimingCliffTime))
	vestingAmount := imb.Sub(ca)
	vestingPerc := math.Ceil(float64(*record.TimingVestingPeriod) / float64(*record.TimingVestingIncrement))
	vestingTime := types.NewInt64Amount(int64(vestingPerc))
	vestingTime = vestingTime.Mul(vestingAmount)
	unLockedTime := ct.Add(vestingTime)

	factor := float64(slotsPerEpoch-(unLockedTime.Int64()-int64(firstSlotOfEpoch))) / float64(slotsPerEpoch)
	if factor < 0 {
		// means it will be unlocked after this epoch, so it is locked in this one
		return types.NewFloat64Percentage(0), nil
	}
	return types.NewFloat64Percentage(factor), nil
}

// calculateSuperchargedContribution calculates supercharged contribution
//
// supercharged contribution = ((supercharged weighting - 1) * timed weighting factor) + 1
func calculateSuperchargedContribution(superchargedWeighting, timedWeighting types.Percentage) (types.Percentage, error) {
	res := superchargedWeighting.Sub(types.NewPercentage("1"))
	res = res.Mul(timedWeighting)
	res = res.Add(types.NewPercentage("1"))
	return res, nil
}
