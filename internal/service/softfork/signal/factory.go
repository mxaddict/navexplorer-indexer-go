package signal

import "github.com/navcoin/navexplorer-indexer-go/v2/pkg/explorer"

func CreateSignal(block *explorer.Block, softForks *explorer.SoftForks) *explorer.Signal {
	sig := &explorer.Signal{Address: block.StakedBy, Height: block.Height}

	for _, softFork := range *softForks {
		if (softFork.State == explorer.SoftForkLockedIn && block.Height <= softFork.LockedInHeight || softFork.IsOpen()) && block.Version>>softFork.SignalBit&1 == 1 {
			sig.SoftForks = append(sig.SoftForks, softFork.Name)
		}
	}

	if len(sig.SoftForks) == 0 {
		return nil
	}

	return sig
}
