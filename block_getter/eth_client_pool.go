package block_getter

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
	"sync/atomic"
	"time"
)

type EthClientPool interface {
	Get() *ethclient.Client
	Close()
}

type ethClientPool struct {
	wsUrl    string
	size     int
	rwLock   sync.RWMutex
	clients  []*ethclient.Client
	index    atomic.Int32
	stopChan chan struct{}
}

func NewEthClientPool(wsUrl string, size int) EthClientPool {
	if size <= 0 {
		size = 1
	}

	pool := &ethClientPool{
		wsUrl:    wsUrl,
		size:     size,
		clients:  make([]*ethclient.Client, size),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < size; i++ {
		if client, err := ethclient.Dial(wsUrl); err == nil {
			pool.clients[i] = client
		}
	}

	go pool.healthCheck()
	return pool
}

func (p *ethClientPool) Get() *ethclient.Client {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	idx := int(p.index.Add(1)) % p.size
	return p.clients[idx]
}

func (p *ethClientPool) healthCheck() {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopChan:
			return
		case <-ticker.C:
			p.checkAndReconnect()
		}
	}
}

func (p *ethClientPool) checkAndReconnect() {
	for i, client := range p.clients {
		if client == nil || !p.isClientHealthy(client) {
			if newClient, err := ethclient.Dial(p.wsUrl); err == nil {
				oldClient := p.clients[i]
				p.rwLock.Lock()
				p.clients[i] = newClient
				p.rwLock.Unlock()
				if oldClient != nil {
					go oldClient.Close() // 异步关闭,避免阻塞
				}
			}
		}
	}
}

func (p *ethClientPool) isClientHealthy(client *ethclient.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := client.ChainID(ctx)
	return err == nil
}

func (p *ethClientPool) Close() {
	close(p.stopChan)

	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	for i, client := range p.clients {
		if client != nil {
			client.Close()
			p.clients[i] = nil
		}
	}
}
