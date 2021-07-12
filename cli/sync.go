package cli

import (
	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/client/staketab"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/worker"
)

func runSync(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	archiveClient := archive.NewDefaultClient(cfg.ArchiveEndpoint)
	graphClient := graph.NewDefaultClient(cfg.MinaEndpoint)
	graphClient.SetDebug(cfg.LogLevel == "debug")
	staketabClient := staketab.NewDefaultClient(cfg.StaketabEndpoint)
	syncWorker := worker.NewSyncWorker(cfg, db, graphClient, archiveClient, staketabClient)

	_, err = syncWorker.Run()
	return err
}
