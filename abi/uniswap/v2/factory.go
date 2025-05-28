package v2

import (
	"abchain_scan/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"strings"
)

const (
	FactoryAbiJson       = `[{"type":"constructor","stateMutability":"nonpayable","payable":false,"inputs":[{"type":"address","name":"_feeToSetter","internalType":"address"}]},{"type":"event","name":"PairCreated","inputs":[{"type":"address","name":"token0","internalType":"address","indexed":true},{"type":"address","name":"token1","internalType":"address","indexed":true},{"type":"address","name":"pair","internalType":"address","indexed":false},{"type":"uint256","name":"pairsLength","internalType":"uint256","indexed":false}],"anonymous":false},{"type":"function","stateMutability":"view","payable":false,"outputs":[{"type":"address","name":"","internalType":"address"}],"name":"allPairs","inputs":[{"type":"uint256","name":"","internalType":"uint256"}],"constant":true},{"type":"function","stateMutability":"view","payable":false,"outputs":[{"type":"uint256","name":"","internalType":"uint256"}],"name":"allPairsLength","inputs":[],"constant":true},{"type":"function","stateMutability":"nonpayable","payable":false,"outputs":[{"type":"address","name":"pair","internalType":"address"}],"name":"createPair","inputs":[{"type":"address","name":"tokenA","internalType":"address"},{"type":"address","name":"tokenB","internalType":"address"}],"constant":false},{"type":"function","stateMutability":"view","payable":false,"outputs":[{"type":"address","name":"","internalType":"address"}],"name":"feeTo","inputs":[],"constant":true},{"type":"function","stateMutability":"view","payable":false,"outputs":[{"type":"address","name":"","internalType":"address"}],"name":"feeToSetter","inputs":[],"constant":true},{"type":"function","stateMutability":"view","payable":false,"outputs":[{"type":"address","name":"","internalType":"address"}],"name":"getPair","inputs":[{"type":"address","name":"pair_address","internalType":"address"},{"type":"address","name":"","internalType":"address"}],"constant":true},{"type":"function","stateMutability":"nonpayable","payable":false,"outputs":[],"name":"setFeeTo","inputs":[{"type":"address","name":"_feeTo","internalType":"address"}],"constant":false},{"type":"function","stateMutability":"nonpayable","payable":false,"outputs":[],"name":"setFeeToSetter","inputs":[{"type":"address","name":"_feeToSetter","internalType":"address"}],"constant":false}]`
	FactoryAddressHex    = "0x723913136a42684B5e3657e3cD2f67ee3e83A82D"
	PairCreatedTopic0Hex = "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9"
)

var (
	FactoryAbi        *abi.ABI
	FactoryAddress    = common.HexToAddress(FactoryAddressHex)
	PairCreatedTopic0 = common.HexToHash(PairCreatedTopic0Hex)
	PairCreatedEvent  *abi.Event
)

func init() {
	factoryAbi, err := abi.JSON(strings.NewReader(FactoryAbiJson))
	if err != nil {
		log.Logger.Fatal("Failed to parse factory ABI", zap.Error(err))
	}
	FactoryAbi = &factoryAbi

	pairCreatedEvent, err := factoryAbi.EventByID(PairCreatedTopic0)
	if err != nil {
		log.Logger.Fatal("Failed to find PairCreatedTopic0", zap.Error(err))
	}
	PairCreatedEvent = pairCreatedEvent
}
