package types

const (
	ProtocolIdNewSwap = iota + 1
	ProtocolIdUniswapV3
)

const (
	ProtocolNameNewSwap   = "NewSwap"
	ProtocolNameUniswapV3 = "UniswapV3"
)

func GetProtocolName(protocolId int) string {
	switch protocolId {
	case ProtocolIdNewSwap:
		return ProtocolNameNewSwap
	case ProtocolIdUniswapV3:
		return ProtocolNameUniswapV3
	default:
		return "Unknown"
	}
}
