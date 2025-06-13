package block_getter

import (
	"abchain_scan/log"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type EthClientPool interface {
	Get() *ethclient.Client
	GetWithIndex() (*ethclient.Client, int)
	Close()
}

const (
	StatusHealthy = iota
	StatusConnecting
	StatusChecking
	StatusCheckFailed
	StatusStoping
	StatusStoped
)

type ClientWrap struct {
	id     int
	url    string
	rwLock sync.RWMutex
	client *ethclient.Client
	status int
}

func NewClientWrap(id int, url string) *ClientWrap {
	w := &ClientWrap{
		id:  id,
		url: url,
	}

	w.connect()
	for {
		if !w.checkHealthy() {
			w.reconnect()
		}
		break
	}

	return w
}

func (w *ClientWrap) setStatus(status int) {
	w.rwLock.Lock()
	defer w.rwLock.Unlock()
	w.status = status
}

func (w *ClientWrap) getStatus() int {
	w.rwLock.RLock()
	defer w.rwLock.RUnlock()
	return w.status
}

func (w *ClientWrap) connect() {
	w.setStatus(StatusConnecting)

	for {
		client, err := ethclient.Dial(w.url)
		if err != nil {
			log.Logger.Info("Err: dial Ethereum err", zap.Int("id", w.id), zap.String("url", w.url), zap.Error(err))
			time.Sleep(time.Second * 2)
			continue
		}
		log.Logger.Info("dial Ethereum ok", zap.Int("id", w.id))
		w.client = client
		break
	}

	w.setStatus(StatusHealthy)
}

func (w *ClientWrap) reconnect() {
	go w.client.Close()
	w.connect()
}

func (w *ClientWrap) checkHealthy() bool {
	w.setStatus(StatusChecking)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	_, err := w.client.ChainID(ctx)
	if err == nil {
		w.setStatus(StatusHealthy)
		return true
	}

	log.Logger.Info("check health err", zap.Int("id", w.id), zap.Error(err))
	w.setStatus(StatusCheckFailed)
	return false
}

func (w *ClientWrap) isHealthy() bool {
	return w.getStatus() == StatusHealthy
}

func (w *ClientWrap) close() {
	w.setStatus(StatusStoping)
	w.client.Close()
	w.setStatus(StatusStoped)
}

type ethClientPool struct {
	url      string
	size     int
	clients  []*ClientWrap
	index    atomic.Int32
	stopChan chan struct{}
}

func NewEthClientPool(url string, size int) EthClientPool {
	if size <= 0 {
		size = 1
	}

	pool := &ethClientPool{
		url:      url,
		size:     size,
		clients:  make([]*ClientWrap, size),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < size; i++ {
		pool.clients[i] = NewClientWrap(i, url)
	}

	go pool.healthCheck()
	return pool
}

func (p *ethClientPool) Get() *ethclient.Client {
	client, _ := p.GetWithIndex()
	return client

}

func (p *ethClientPool) GetWithIndex() (*ethclient.Client, int) {
	for {
		idx := int(p.index.Add(1)) % p.size
		if p.clients[idx].isHealthy() {
			return p.clients[idx].client, idx
		}
	}
}

func (p *ethClientPool) healthCheck() {
	ticker := time.NewTicker(time.Second * 10)
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
	for _, client := range p.clients {
		if !client.checkHealthy() {
			client.reconnect()
		}
	}
}

func (p *ethClientPool) Close() {
	close(p.stopChan)
	for _, client := range p.clients {
		client.close()
	}
}
