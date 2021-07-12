package store

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// BlocksStore handles operations on blocks
type BlocksStore struct {
	baseStore
}

// FindBy returns a block for a matching attribute
func (s BlocksStore) FindBy(key string, value interface{}) (*model.Block, error) {
	result := &model.Block{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a block with matching ID
func (s BlocksStore) FindByID(id int64) (*model.Block, error) {
	return s.FindBy("id", id)
}

// FindByHash returns a block with the matching hash
func (s BlocksStore) FindByHash(hash string) (*model.Block, error) {
	return s.FindBy("hash", hash)
}

// FindByHeight returns a canonical block with the matching height
func (s BlocksStore) FindByHeight(height uint64) (*model.Block, error) {
	result := model.Block{}

	scope := s.db.Limit(1)
	scope = scope.Where("height = ? AND canonical = ?", height, true)
	err := scope.Find(&result).Error
	return &result, checkErr(err)
}

// Recent returns the most recent block
func (s BlocksStore) Recent() (*model.Block, error) {
	block := &model.Block{}
	err := s.db.Where("canonical = ?", true).Order("height DESC").Limit(1).Take(block).Error
	return block, checkErr(err)
}

// LastBlock returns the last block
func (s BlocksStore) LastBlock() (*model.Block, error) {
	block := &model.Block{}
	err := s.db.Order("height DESC").Limit(1).Take(block).Error
	return block, checkErr(err)
}

// Search returns blocks that match search filters
func (s BlocksStore) Search(search *BlockSearch) ([]model.Block, error) {
	result := []model.Block{}

	scope := s.db.
		Order(fmt.Sprintf("%s %s", search.Sort, search.Order)).
		Limit(search.Limit)

	if search.MinHeight > 0 {
		scope = scope.Where("height >= ?", search.MinHeight)
	}

	if search.MaxHeight > 0 {
		scope = scope.Where("height <= ?", search.MaxHeight)
	}

	if search.Creator != "" {
		scope = scope.Where("creator = ?", search.Creator)
	}

	return result, scope.Find(&result).Error
}

// AvgTimes returns recent blocks averages
func (s BlocksStore) AvgTimes(limit int64) ([]byte, error) {
	return jsonquery.MustObject(s.db, queries.BlocksTimes, limit)
}

// Stats returns block stats for a given interval
func (s BlocksStore) Stats(period uint, interval string) ([]byte, error) {
	return jsonquery.MustArray(s.db, queries.BlocksStats, period, interval)
}

// MarkBlocksOrphan updates all blocks as non canonical at a height
func (s BlocksStore) MarkBlocksOrphan(height uint64) error {
	return s.db.Exec(queries.MarkBlocksOrphan, height).Error
}

// MarkBlockCanonical updates canonical at a height
func (s BlocksStore) MarkBlockCanonical(hash string) error {
	return s.db.Exec(queries.MarkBlockCanonical, hash).Error
}

// FindUnsafeBlocks returns the last indexed unsafe blocks that may be orphaned
func (s BlocksStore) FindUnsafeBlocks(startingHeight uint64) ([]model.Block, error) {
	result := []model.Block{}

	scope := s.db.
		Where("height >= ?", startingHeight).
		Order("height asc")

	return result, scope.Find(&result).Error
}

// FirstBlockOfEpoch returns the first block of epoch
func (s BlocksStore) FirstBlockOfEpoch(epoch string) (*model.Block, error) {
	block := &model.Block{}
	err := s.db.Where("epoch >= ?", epoch).Order("height ASC").Limit(1).Take(block).Error
	return block, checkErr(err)
}

// FindNonCalculatedBlockRewards returns the non calculated blocks for rewards
func (s BlocksStore) FindNonCalculatedBlockRewards(from, to uint64) ([]model.Block, error) {
	result := []model.Block{}

	scope := s.db.Where("height >= ? AND height < ? AND canonical = ? AND reward_calculated = ?", from, to, true, false).
		Order("height asc")
	return result, scope.Find(&result).Error
}

// FindLastCalculatedBlockReward returns the last calculated block reward
func (s BlocksStore) FindLastCalculatedBlockReward(from uint64) (*model.Block, error) {
	block := &model.Block{}

	err := s.db.Where("height > ? AND canonical = ? AND reward_calculated = ?", from, true, true).
		Order("height desc").Limit(1).Take(block).Error
	return block, checkErr(err)
}
