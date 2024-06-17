package indexer

import (
	"github.com/christophercampbell/bridge-connector/types"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
)

var (
	// New Ger event
	updateL1InfoTreeSignatureHash = ethgo.Hash(ethgo.Keccak256([]byte("UpdateL1InfoTree(bytes32,bytes32)")))
	l1InfoTreeEvent               = abi.MustNewEvent(`event UpdateL1InfoTree(
        bytes32 indexed mainnetExitRoot,
        bytes32 indexed rollupExitRoot
	)`)

	// PreLxLy events
	updateGlobalExitRootSignatureHash = ethgo.Hash(ethgo.Keccak256([]byte("UpdateGlobalExitRoot(bytes32,bytes32)")))
	v1GEREvent                        = abi.MustNewEvent(`event UpdateGlobalExitRoot(
        bytes32 indexed mainnetExitRoot,
        bytes32 indexed rollupExitRoot
	)`)

	// New Bridge events
	depositEventSignatureHash = ethgo.Hash(ethgo.Keccak256([]byte("BridgeEvent(uint8,uint32,address,uint32,address,uint256,bytes,uint32)"))) // Used in oldBridge as well
	depositEvent              = abi.MustNewEvent(`event BridgeEvent(
	   uint8 leafType,
	   uint32 originNetwork,
	   address originAddress,
	   uint32 destinationNetwork,
	   address destinationAddress,
	   uint256 amount,
	   bytes metadata,
	   uint32 depositCount
	)`)

	//     * @param globalIndex Global index is defined as:
	//     * | 191 bits |    1 bit     |   32 bits   |     32 bits    |
	//     * |    0     |  mainnetFlag | rollupIndex | localRootIndex |
	//     * note that only the rollup index will be used only in case the mainnet flag is 0
	//     * note that global index do not assert the unused bits to 0.
	//     * This means that when synching the events, the globalIndex must be decoded the same way that in the Smart contract
	//     * to avoid possible synch attacks

	claimEventSignatureHash = ethgo.Hash(ethgo.Keccak256([]byte("ClaimEvent(uint256,uint32,address,address,uint256)")))
	claimEvent              = abi.MustNewEvent(`event ClaimEvent(
        uint256 globalIndex,
        uint32 originNetwork,
        address originAddress,
        address destinationAddress,
        uint256 amount
	)`)

	// Old Bridge events
	oldClaimEventSignatureHash = ethgo.Hash(ethgo.Keccak256([]byte("ClaimEvent(uint32,uint32,address,address,uint256)")))
	oldClaimEvent              = abi.MustNewEvent(`event ClaimEvent(
        uint32 index,
        uint32 originNetwork,
        address originAddress,
        address destinationAddress,
        uint256 amount
	)`)

	verifyBatchesEtrogSignatureHash = ethgo.Hash(ethgo.Keccak256([]byte("VerifyBatches(uint64,bytes32,address)")))
	verifyBatchesEtrogEvent         = abi.MustNewEvent(`event VerifyBatches(
        uint64 indexed numBatch,
        bytes32 stateRoot,
        address indexed aggregator
    )`)

	verifyBatchesTrustedSequencerHash  = ethgo.Hash(ethgo.Keccak256([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)")))
	verifyBatchesTrustedSequencerEvent = abi.MustNewEvent(`event VerifyBatchesTrustedAggregator(
        uint64 indexed numBatch,
        bytes32 stateRoot,
        address indexed aggregator
    )`)
)

const (
	BridgeEventL1InfoTree = iota
	BridgeEventV1GER
	BridgeEventDeposit
	BridgeEventV2Claim
	BridgeEventV1Claim
	BridgeEventVerifyBatchesEtrog
	BridgeEventVerifyTrustedSequencer
)

var (
	bridgeEventTypeMap = map[ethgo.Hash]int{
		l1InfoTreeEvent.ID():                    BridgeEventL1InfoTree,
		v1GEREvent.ID():                         BridgeEventV1GER,
		depositEvent.ID():                       BridgeEventDeposit,
		claimEvent.ID():                         BridgeEventV2Claim,
		oldClaimEvent.ID():                      BridgeEventV1Claim,
		verifyBatchesEtrogEvent.ID():            BridgeEventVerifyBatchesEtrog,
		verifyBatchesTrustedSequencerEvent.ID(): BridgeEventVerifyTrustedSequencer,
	}

	bridgeEventParseMap = map[int]func(log *ethgo.Log) (map[string]interface{}, error){
		BridgeEventL1InfoTree:             l1InfoTreeEvent.ParseLog,
		BridgeEventV1GER:                  v1GEREvent.ParseLog,
		BridgeEventDeposit:                depositEvent.ParseLog,
		BridgeEventV2Claim:                claimEvent.ParseLog,
		BridgeEventV1Claim:                oldClaimEvent.ParseLog,
		BridgeEventVerifyBatchesEtrog:     verifyBatchesEtrogEvent.ParseLog,
		BridgeEventVerifyTrustedSequencer: verifyBatchesTrustedSequencerEvent.ParseLog,
	}
)

func maybeFromLog(l *ethgo.Log) *types.BridgeEvent {
	if et, ok := bridgeEventTypeMap[l.Topics[0]]; !ok {
		return nil
	} else {
		data, _ := bridgeEventParseMap[et](l) // TODO: handle error
		be := types.BridgeEvent{
			Removed:          l.Removed,
			BlockNumber:      l.BlockNumber,
			TransactionIndex: l.TransactionIndex,
			LogIndex:         l.LogIndex,
			TransactionHash:  l.TransactionHash,
			EventType:        uint8(et),
			Data:             data,
		}
		return &be
	}
}
