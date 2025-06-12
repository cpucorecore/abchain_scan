package config

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	bs, _ := json.Marshal(G)
	log.Println(string(bs))
}
