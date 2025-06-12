package abi

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	b, e := json.Marshal(Topic2ProtocolIds)
	require.Nil(t, e)
	t.Log(string(b))

	b, e = json.Marshal(FactoryAddress2ProtocolId)
	require.Nil(t, e)
	t.Log(string(b))

	b, e = json.Marshal(Topic2FactoryAddresses)
	require.Nil(t, e)
	t.Log(string(b))
}
