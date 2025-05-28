package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BlockNumberResponse struct {
	JsonRpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"` // 十六进制区块高度（如 "0x4bb4e1"）
}

func GetLatestBlockNumber(rpcURL string) (uint64, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	})
	if err != nil {
		return 0, fmt.Errorf("构造请求数据失败: %v", err)
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析 JSON 响应
	var result BlockNumberResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	// 将十六进制区块高度转换为十进制
	blockNumber, err := strconv.ParseUint(result.Result[2:], 16, 64) // 去掉 "0x" 前缀
	if err != nil {
		return 0, fmt.Errorf("转换区块高度失败: %v", err)
	}

	return blockNumber, nil
}
