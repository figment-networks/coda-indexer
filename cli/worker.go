package cli

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/client/staketab"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/store"
	"github.com/figment-networks/mina-indexer/worker"
)

func startSyncWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	client := graph.NewDefaultClient(cfg.MinaEndpoint)
	archiveClient := archive.NewDefaultClient(cfg.ArchiveEndpoint)
	staketabClient := staketab.NewDefaultClient(cfg.StaketabEndpoint)
	syncWorker := worker.NewSyncWorker(cfg, db, client, archiveClient, staketabClient)
	timer := time.NewTimer(cfg.SyncDuration())

	wg.Add(1)

	go func() {
		defer func() {
			timer.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-timer.C:
				lag, err := syncWorker.Run()
				if err != nil {
					log.WithError(err).Error("sync failed")
				}
				if lag > 10 {
					timer.Reset(time.Second)
				} else {
					timer.Reset(cfg.SyncDuration())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startCleanupWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(cfg.CleanupDuration())

	go func() {
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ticker.C:
				worker.RunCleanup(cfg, db)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startWorker(cfg *config.Config) error {
	log.Info("using mina graph endpoint: ", cfg.MinaEndpoint)
	log.Info("using mina archive endpoint: ", cfg.ArchiveEndpoint)
	log.Info("sync will run every: ", cfg.SyncInterval)
	log.Info("cleanup will run every: ", cfg.CleanupInterval)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}

	cancelSync := startSyncWorker(wg, cfg, db)
	cancelCleanup := startCleanupWorker(wg, cfg, db)

	s := <-initSignals()

	log.Info("received signal: ", s)
	cancelSync()
	cancelCleanup()

	wg.Wait()
	return nil
}
