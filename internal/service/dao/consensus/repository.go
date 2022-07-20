package consensus

import (
	"context"
	"encoding/json"
	"github.com/navcoin/navexplorer-indexer-go/v2/internal/elastic_cache"
	"github.com/navcoin/navexplorer-indexer-go/v2/pkg/explorer"
)

type Repository interface {
	GetConsensusParameters() (explorer.ConsensusParameters, error)
}

type repository struct {
	elastic elastic_cache.Index
}

func NewRepo(elastic elastic_cache.Index) Repository {
	return repository{elastic}
}

func (r repository) GetConsensusParameters() (explorer.ConsensusParameters, error) {
	results, err := r.elastic.GetClient().Search(elastic_cache.ConsensusIndex.Get()).
		Sort("id", true).
		Size(10000).
		Do(context.Background())
	if err != nil || results == nil {
		return explorer.ConsensusParameters{}, err
	}

	if len(results.Hits.Hits) == 0 {
		return explorer.ConsensusParameters{}, nil
	}

	consensusParameters := explorer.ConsensusParameters{}
	for _, hit := range results.Hits.Hits {
		var consensusParameter explorer.ConsensusParameter
		if err = json.Unmarshal(hit.Source, &consensusParameter); err != nil {
			return explorer.ConsensusParameters{}, err
		}

		consensusParameters.Add(consensusParameter)
	}

	return consensusParameters, nil
}
