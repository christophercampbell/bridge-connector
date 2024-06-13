package indexer

import (
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
)
