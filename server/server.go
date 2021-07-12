package server

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store"
)

// Server handles HTTP requests
type Server struct {
	*gin.Engine

	graphClient *graph.Client
	db          *store.Store
	log         *logrus.Logger
}

// New returns a new server instance
func New(db *store.Store, cfg *config.Config, logger *logrus.Logger) *Server {
	s := &Server{
		Engine: gin.New(),

		db:          db,
		graphClient: graph.NewDefaultClient(cfg.MinaEndpoint),
		log:         logger,
	}

	s.initMiddleware(cfg)
	s.initRoutes()

	return s
}

func (s *Server) initRoutes() {
	s.GET("/health", s.GetHealth)
	s.GET("/status", s.GetStatus)
	s.GET("/height", s.GetCurrentHeight)
	s.GET("/block", s.GetCurrentBlock)
	s.GET("/blocks", s.GetBlocks)
	s.GET("/blocks/:id", s.GetBlock)
	s.GET("/blocks/:id/transactions", s.GetBlockTransactions)
	s.GET("/block_times", s.GetBlockTimes)
	s.GET("/block_stats", timeBucketMiddleware(), s.GetBlockStats)
	s.GET("/chain_stats", timeBucketMiddleware(), s.GetBlockStats)
	s.GET("/validators", s.GetValidators)
	s.GET("/validators/:id", s.GetValidator)
	s.GET("/validators/:id/stats", timeBucketMiddleware(), s.GetValidatorStats)
	s.GET("/delegations", s.GetDelegations)
	s.GET("/rewards/:id", s.GetRewards)
	s.GET("/snarkers", s.GetSnarkers)
	s.GET("/snarker/:id", s.GetSnarker)
	s.GET("/transactions", s.GetTransactions)
	s.GET("/pending_transactions", s.GetPendingTransactions)
	s.GET("/transactions/:id", s.GetTransaction)
	s.GET("/accounts/:id", s.GetAccount)
	s.GET("/ledgers", s.GetLedgers)
	s.GET("/ledger", s.GetLedger)
}

func (s *Server) initMiddleware(cfg *config.Config) {
	s.Use(gin.Recovery())
	s.Use(requestLoggerMiddleware(logrus.StandardLogger()))

	if cfg.IsDevelopment() {
		s.Use(corsMiddleware())
	}

	if cfg.RollbarToken != "" {
		s.Use(rollbarMiddleware())
	}
}

// GetHealth renders the server health status
func (s Server) GetHealth(c *gin.Context) {
	resp := HealthResponse{Healthy: true}

	if err := s.db.Test(); err != nil {
		s.log.WithError(err).Error("database check error")
		resp.Healthy = false
		jsonResponse(c, 500, resp)
		return
	}

	jsonOk(c, resp)
}

// GetStatus returns the status of the service
func (s Server) GetStatus(c *gin.Context) {
	resp := StatusResponse{
		AppName:    config.AppName,
		AppVersion: config.AppVersion,
		GitCommit:  config.GitCommit,
		GoVersion:  config.GoVersion,
		SyncStatus: "stale",
	}

	// Fetch node status as quickly as possible.
	// We don't care if node's sync status is reported as error at this point.
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
	defer cancel()

	daemonStatus, err := s.graphClient.GetDaemonStatus(ctx)
	if err == nil {
		resp.NodeVersion = daemonStatus.CommitID
		resp.NodeStatus = string(daemonStatus.SyncStatus)
	} else {
		logrus.WithError(err).Error("node status fetch failed")
		resp.NodeError = true
	}

	if block, err := s.db.Blocks.Recent(); err == nil {
		resp.LastBlockTime = block.Time
		resp.LastBlockHeight = block.Height

		if time.Since(block.Time).Minutes() <= 30 {
			resp.SyncStatus = "current"
		}
	} else {
		logrus.WithError(err).Error("recent block fetch failed")
	}

	jsonOk(c, resp)
}

// GetCurrentHeight returns the current blockchain height
func (s *Server) GetCurrentHeight(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, HeightResponse{
		Height: block.Height,
		Time:   block.Time,
	})
}

// GetCurrentBlock returns the current blockchain height
func (s *Server) GetCurrentBlock(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, block)
}

// GetBlock returns a single block
func (s *Server) GetBlock(c *gin.Context) {
	var block *model.Block
	var err error

	id := resourceID(c, "id")
	if id.IsNumeric() {
		if id.UInt64() == 0 {
			badRequest(c, errors.New("height must be greater than 0"))
			return
		}
		block, err = s.db.Blocks.FindByHeight(id.UInt64())
	} else {
		block, err = s.db.Blocks.FindByHash(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	creator, err := s.db.Accounts.FindByPublicKey(block.Creator)
	if err == store.ErrNotFound {
		creator = nil
		err = nil
	}
	if shouldReturn(c, err) {
		return
	}

	transactions, err := s.db.Transactions.ByHeight(block.Height, uint(block.TransactionsCount))
	if shouldReturn(c, err) {
		return
	}

	jobs, err := s.db.Jobs.ByHash(block.Hash)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, BlockResponse{
		Block:        block,
		Creator:      creator,
		Transactions: transactions,
		SnarkJobs:    jobs,
	})
}

func (s *Server) GetBlockTransactions(c *gin.Context) {
	var block *model.Block
	var err error

	id := resourceID(c, "id")
	if id.IsNumeric() {
		if id.UInt64() == 0 {
			badRequest(c, errors.New("height must be greater than 0"))
			return
		}
		block, err = s.db.Blocks.FindByHeight(id.UInt64())
	} else {
		block, err = s.db.Blocks.FindByHash(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	transactions, err := s.db.Transactions.ByHeight(block.Height, uint(block.TransactionsCount))
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, transactions)
}

// GetBlocks returns a list of available blocks matching the filter
func (s *Server) GetBlocks(c *gin.Context) {
	search := &store.BlockSearch{}

	if err := c.BindQuery(search); err != nil {
		badRequest(c, err)
		return
	}

	if err := search.Validate(); err != nil {
		badRequest(c, err)
		return
	}

	blocks, err := s.db.Blocks.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, blocks)
}

// GetBlockTimes returns avg block times info
func (s *Server) GetBlockTimes(c *gin.Context) {
	params := blockTimesParams{}

	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.AvgTimes(params.Limit)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, result)
}

// GetBlockStats returns block stats for an interval
func (s *Server) GetBlockStats(c *gin.Context) {
	tb := c.MustGet("timebucket").(timeBucket)
	result, err := s.db.Blocks.Stats(tb.Period, tb.Interval)
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, result)
}

// GetTransaction returns a single transaction details
func (s *Server) GetTransaction(c *gin.Context) {
	var tran *model.Transaction
	var err error

	id := resourceID(c, "id")
	if id.IsNumeric() {
		tran, err = s.db.Transactions.FindByID(id.Int64())
	} else {
		tran, err = s.db.Transactions.FindByHash(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, tran)
}

// GetValidators rendes all existing validators
func (s *Server) GetValidators(c *gin.Context) {
	validators, err := s.db.Validators.Index()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, validators)
}

// GetValidator renders the validator details
func (s *Server) GetValidator(c *gin.Context) {
	validator, err := s.db.Validators.FindByPublicKey(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	account, err := s.db.Accounts.FindByPublicKey(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}
	if err != store.ErrNotFound && shouldReturn(c, err) {
		return
	}

	delegations, err := s.db.Staking.FindDelegations(store.FindDelegationsParams{
		Delegate: validator.PublicKey,
	})
	if err != store.ErrNotFound && shouldReturn(c, err) {
		return
	}

	stats30d, err := s.db.Stats.ValidatorStats(validator, 30, store.BucketDay)
	if shouldReturn(c, err) {
		return
	}

	stats24h, err := s.db.Stats.ValidatorStats(validator, 48, store.BucketHour)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, ValidatorResponse{
		Validator:   validator,
		Account:     account,
		Delegations: delegations,
		Stats:       stats30d,
		StatsHourly: stats24h,
		StatsDaily:  stats30d,
	})
}

// GetValidatorStats renders validator stats for a given time bucket
func (s *Server) GetValidatorStats(c *gin.Context) {
	tb := c.MustGet("timebucket").(timeBucket)

	validator, err := s.db.Validators.FindByPublicKey(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	stats, err := s.db.Stats.ValidatorStats(validator, tb.Period, tb.Interval)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, stats)
}

// GetDelegations rendes all existing delegations
func (s *Server) GetDelegations(c *gin.Context) {
	delegations, err := s.db.Staking.FindDelegations(store.FindDelegationsParams{
		PublicKey: c.Query("public_key"),
		Delegate:  c.Query("delegate"),
	})
	if err != store.ErrNotFound && shouldReturn(c, err) {
		return
	}

	jsonOk(c, delegations)
}

// GeRewards returns rewards
func (s Server) GetRewards(c *gin.Context) {
	var params rewardsParams
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, errors.New("missing parameter"))
		return
	}
	if err := params.Validate(); err != nil {
		badRequest(c, err)
		return
	}
	interval, _ := model.GetTypeForTimeInterval(params.Interval)
	rewardOwnerType, _ := model.GetTypeForRewardOwnerType(params.RewardOwnerType)
	resp, err := s.db.Rewards.FetchRewardsByInterval(c.Param("id"), params.ValidatorId, params.From, params.To, interval, rewardOwnerType)
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, resp)
}

// GetSnarkers renders all existing snarkers
func (s *Server) GetSnarkers(c *gin.Context) {
	snarkers, err := s.db.Snarkers.All()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, snarkers)
}

// GetSnarker get snarker info for canonical
func (s *Server) GetSnarker(c *gin.Context) {
	snarker, err := s.db.Snarkers.FindSnarker(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	result, err := s.db.Snarkers.SnarkerInfoFromCanonicalBlocks(snarker.Account, snarker.StartHeight, snarker.LastHeight)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, result)
}

// GetTransactions returns transactions by height
func (s *Server) GetTransactions(c *gin.Context) {
	search := store.TransactionSearch{}
	if err := c.BindQuery(&search); err != nil {
		badRequest(c, err)
		return
	}

	if err := search.Validate(); err != nil {
		badRequest(c, err)
		return
	}

	transactions, err := s.db.Transactions.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, transactions)
}

// GetPendingTransactions returns transactions by height
func (s *Server) GetPendingTransactions(c *gin.Context) {
	transactions, err := s.graphClient.GetPendingTransactions()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, transactions)
}

// GetAccount returns account for by hash or ID
func (s *Server) GetAccount(c *gin.Context) {
	var (
		acc *model.Account
		err error
	)

	id := resourceID(c, "id")
	if id.IsNumeric() {
		acc, err = s.db.Accounts.FindByID(id.Int64())
	} else {
		acc, err = s.db.Accounts.FindByPublicKey(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, acc)
}

// GetLedgers returns a list of all existing ledgers
func (s *Server) GetLedgers(c *gin.Context) {
	ledgers, err := s.db.Staking.AllLedgers()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, ledgers)
}

// GetLedger records the current epoch ledger records
func (s *Server) GetLedger(c *gin.Context) {
	var (
		ledger *model.Ledger
		err    error
	)

	input := &LedgerRequest{}
	if err := c.Bind(input); err != nil {
		badRequest(c, err)
		return
	}

	if epoch := input.Epoch; epoch != nil {
		ledger, err = s.db.Staking.FindLedger(*epoch)
	} else {
		ledger, err = s.db.Staking.LastLedger()
	}
	if shouldReturn(c, err) {
		return
	}

	records, err := s.db.Staking.LedgerRecords(ledger.ID)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, LedgerResponse{
		Ledger:  ledger,
		Records: records,
	})
}
