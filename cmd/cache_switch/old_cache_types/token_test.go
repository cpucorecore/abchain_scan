package old_cache_types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetOldToken(t *testing.T) {
	tokenString := `
{
  "Address": "0xf9ee4ce4ddbdd46bd01b55630ab98da2eddb4444",
  "Creator": "0x0000000000000000000000000000000000000000",
  "Name": "Arrion Knight",
  "Symbol": "YNCODING",
  "Decimals": 18,
  "TotalSupply": "1000000000",
  "BlockNumber": null,
  "BlockTime": "0001-01-01T00:00:00Z",
  "Filtered": false,
  "FilteredReason": 0,
  "Program": "",
  "MainPair": "0x0000000000000000000000000000000000000000"
}`
	token := &Token{}

	err := json.Unmarshal([]byte(tokenString), token)
	require.NoError(t, err)

	t.Log(token)
}
