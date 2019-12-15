package daemon

import (
	"github.com/NavExplorer/navexplorer-indexer-go/generated/dic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var container *dic.Container

func Execute() {
	config.Init()

	container, _ = dic.NewContainer(dingo.App)

	container.GetElastic().InstallMappings()
	container.GetSoftforkService().LoadSoftForks()

	height := getHeight()
	if err := container.GetRewinder().RewindToHeight(height); err != nil {
		log.WithError(err).Fatal("Failed to rewind index")
	}

	if height != 0 {
		if block, err := container.GetBlockRepo().GetBlockByHeight(height); err != nil {
			log.WithError(err).Fatal("Failed to get block")
		} else {
			blockCycle := block.BlockCycle(config.Get().DaoCfundConsensus.BlocksPerVotingCycle, config.Get().DaoCfundConsensus.Quorum)
			container.GetDaoProposalService().LoadVotingProposals(block, blockCycle)
			container.GetDaoPaymentRequestService().LoadVotingPaymentRequests(block, blockCycle)
		}
	}

	// Bulk index the backlog
	container.GetIndexer().BulkIndex()

	// Subscribe to 0MQ
	container.GetSubscriber().Subscribe()
}

func getHeight() uint64 {
	if height, err := container.GetRedis().Start(); err != nil {
		log.WithError(err).Fatal("Failed to start redis")
	} else {
		if height >= uint64(config.Get().BulkIndexSize) {
			return height - uint64(config.Get().BulkIndexSize)
		}
	}

	return 0
}
