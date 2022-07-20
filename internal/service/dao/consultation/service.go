package consultation

import (
	"github.com/navcoin/navexplorer-indexer-go/v2/pkg/explorer"
	"go.uber.org/zap"
	"math"
)

type Service interface {
	LoadOpenConsultations(block *explorer.Block)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return service{repository}
}

func (s service) LoadOpenConsultations(block *explorer.Block) {
	zap.L().Info("ConsultationService: Load Open Consultations")

	excludeOlderThan := block.Height - (uint64(block.BlockCycle.Size * 2))
	if excludeOlderThan < 0 {
		excludeOlderThan = 0
	}

	consultations, err := s.repository.GetOpenConsultations(excludeOlderThan)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("ConsultationService: Failed to load consultations")
	}

	for _, c := range consultations {
		zap.L().With(zap.String("consultation", c.Hash)).Info("ConsultationService: Loaded consultation")
		Consultations[c.Hash] = c
	}
}

func ConsultationSupportRequired() int {
	//minSupport := float64(consensus.Parameters.Get(consensus.CONSULTATION_MIN_SUPPORT).Value)
	//cycleSize := float64(consensus.Parameters.Get(consensus.VOTING_CYCLE_LENGTH).Value)
	//
	//return int(math.Ceil((minSupport / 10000) * cycleSize))
	return 0
}

func AnswerSupportRequired(minSupport *explorer.ConsensusParameter, votingCycleLength *explorer.ConsensusParameter) int {
	return int(math.Ceil((float64(minSupport.Value) / 10000) * float64(votingCycleLength.Value)))
}
