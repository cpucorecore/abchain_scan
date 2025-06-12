package types

const (
	ProtocolIdUniswapV2 = iota + 1
	ProtocolIdUniswapV3
	ProtocolIdPancakeV2
	ProtocolIdPancakeV3
	ProtocolIdAerodrome
)

const (
	ProtocolNameUniswapV2 = "UniswapV2"
	ProtocolNameUniswapV3 = "UniswapV3"
	ProtocolNamePancakeV2 = "PancakeV2"
	ProtocolNamePancakeV3 = "PancakeV3"
	ProtocolNameAerodrome = "Aerodrome"
)

func GetProtocolName(protocolId int) string {
	switch protocolId {
	case ProtocolIdUniswapV2:
		return ProtocolNameUniswapV2
	case ProtocolIdUniswapV3:
		return ProtocolNameUniswapV3
	case ProtocolIdPancakeV2:
		return ProtocolNamePancakeV2
	case ProtocolIdPancakeV3:
		return ProtocolNamePancakeV3
	case ProtocolIdAerodrome:
		return ProtocolNameAerodrome
	default:
		return "Unknown"
	}
}
