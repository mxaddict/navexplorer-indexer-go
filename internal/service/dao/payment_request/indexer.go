package payment_request

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	navcoin *navcoind.Navcoind
	elastic *elastic_cache.Index
}

func NewIndexer(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index) *Indexer {
	return &Indexer{navcoin, elastic}
}

func (i *Indexer) Index(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if tx.Version != 5 {
			continue
		}

		if navP, err := i.navcoin.GetPaymentRequest(tx.Hash); err == nil {
			paymentRequest := CreatePaymentRequest(navP, tx.Height)

			index := elastic_cache.PaymentRequestIndex.Get()
			resp, err := i.elastic.Client.Index().Index(index).BodyJson(paymentRequest).Do(context.Background())
			if err != nil {
				raven.CaptureError(err, nil)
				log.WithError(err).Fatal("Failed to save new payment request")
			}

			paymentRequest.MetaData = explorer.NewMetaData(resp.Id, resp.Index)
			PaymentRequests = append(PaymentRequests, paymentRequest)
		}
	}
}

func (i *Indexer) Update(blockCycle *explorer.BlockCycle, block *explorer.Block) {
	for _, p := range PaymentRequests {
		if p == nil {
			continue
		}

		navP, err := i.navcoin.GetPaymentRequest(p.Hash)
		if err != nil {
			raven.CaptureError(err, nil)
			log.WithError(err).Fatalf("Failed to find active proposal: %s", p.Hash)
		}

		UpdatePaymentRequest(navP, block.Height, p)
		if p.UpdatedOnBlock == block.Height {
			i.elastic.AddUpdateRequest(elastic_cache.PaymentRequestIndex.Get(), p.Hash, p, p.MetaData.Id)
		}

		if p.Status == explorer.PaymentRequestExpired || p.Status == explorer.PaymentRequestRejected {
			if block.Height-p.UpdatedOnBlock >= uint64(blockCycle.Size) {
				PaymentRequests.Delete(p.Hash)
			}
		}
	}
}
