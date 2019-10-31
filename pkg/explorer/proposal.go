package explorer

type Proposal struct {
	Version             uint32 `json:"version"`
	Hash                string `json:"hash"`
	BlockHash           string `json:"blockHash"`
	Description         string `json:"description"`
	RequestedAmount     uint64 `json:"requestedAmount"`
	NotPaidYet          uint64 `json:"notPaidYet"`
	UserPaidFee         uint64 `json:"userPaidFee"`
	PaymentAddress      string `json:"paymentAddress"`
	ProposalDuration    uint64 `json:"proposalDuration"`
	ExpiresOn           uint64 `json:"expiresOn"`
	Status              string `json:"status"`
	State               uint   `json:"state"`
	StateChangedOnBlock string `json:"stateChangedOnBlock,omitempty"`

	// Custom
	Height uint64         `json:"height"`
	Cycles ProposalCycles `json:"cycles"`
}

type ProposalCycles []ProposalCycle

type ProposalCycle struct {
	VotingCycle uint `json:"votingCycle"`
	VotesYes    uint `json:"votesYes"`
	VotesNo     uint `json:"votesNo"`
}
