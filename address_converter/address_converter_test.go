package address_converter

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestToABChainAddressString(t *testing.T) {
	tests := []struct {
		ethAddressHex        string
		expectABChainAddress string
	}{
		{
			ethAddressHex:        "0x54db4d39cedeaf78af1d15aa540570ca168fe303",
			expectABChainAddress: "NEW182Ln6vMMXGRwPbESbyGCxc3U3sPU1PrcKw8",
		},
		{
			ethAddressHex:        "0x46FCae97CE42645c10aF9794bbf6b9B46Dfc2985",
			expectABChainAddress: "NEW182KWmSyNtMbXkKRGpyK7t5iHStkVxSpTaBf",
		},
	}

	for _, test := range tests {
		ethAddress := common.HexToAddress(test.ethAddressHex)
		abChainAddress := EthAddr2ABChainAddrStr(ethAddress)
		if abChainAddress != test.expectABChainAddress {
			t.Errorf("expected %s, got %s", test.expectABChainAddress, abChainAddress)
		}
		t.Logf("ETH Address: %s, ABChain Address: %s", ethAddress.Hex(), abChainAddress)
	}
}
