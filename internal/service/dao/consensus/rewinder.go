package consensus

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Rewinder struct {
	navcoin *navcoind.Navcoind
	elastic *elastic_cache.Index
	repo    *Repository
}

func NewRewinder(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, repo *Repository) *Rewinder {
	return &Rewinder{navcoin, elastic, repo}
}

func (r *Rewinder) Rewind() error {
	log.Infof("Rewinding DAO consensus")

	cfundStats, err := r.navcoin.CfundStats()
	if err != nil {
		log.WithError(err).Error("Failed to get CfundStats from Navcoind")
		return err
	}

	consensus, _ := r.repo.GetConsensus()
	if consensus == nil {
		consensus = new(explorer.Consensus)
		UpdateConsensus(&cfundStats, consensus)
		_, err := r.elastic.Client.Index().Index(elastic_cache.ConsensusIndex.Get()).BodyJson(consensus).Do(context.Background())
		if err != nil {
			log.WithError(err).Fatalf("Failed to persist consensus")
		}
	} else {
		UpdateConsensus(&cfundStats, consensus)
		log.Info("Index Update Cfund Consensus")
		r.elastic.AddUpdateRequest(elastic_cache.ConsensusIndex.Get(), "consensus", consensus, consensus.MetaData.Id)
	}

	return nil
}