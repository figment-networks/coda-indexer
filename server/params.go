package server

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model"
)

type rewardsParams struct {
	From            time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To              time.Time `form:"to" binding:"required" time_format:"2006-01-02"`
	AccountId       string    `form:"account_id" binding:"-" `
	RewardOwnerType string    `form:"owner_type" binding:"required" `
	Interval        string    `form:"interval" binding:"required" `
}

func (p *rewardsParams) Validate() error {
	if p.From.IsZero() && p.To.IsZero() {
		return errors.New("invalid time range: " + "")
	}

	var ok bool
	if _, ok = model.GetTypeForTimeInterval(p.Interval); !ok {
		return errors.New("time interval type is wrong")
	}

	if _, ok = model.GetTypeForRewardOwnerType(p.RewardOwnerType); !ok {
		return errors.New("owner type is wrong")
	}

	return nil
}
