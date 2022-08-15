package managers

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/lavash-center/getblock.io.test/src/models"
	"github.com/ybbus/jsonrpc/v2"
)

const (
	blockNumber      = "eth_blockNumber"
	getBlockByNumber = "eth_getBlockByNumber"
)

type BlocksManager interface {
	GetBlockAddress() (string, error)
}

type BlockManagerImpl struct {
	wg     *sync.WaitGroup
	mu     sync.Mutex
	rpcCli jsonrpc.RPCClient
}

func NewBlocksManagerImpl(rpcCli jsonrpc.RPCClient) *BlockManagerImpl {
	return &BlockManagerImpl{
		wg:     &sync.WaitGroup{},
		mu:     sync.Mutex{},
		rpcCli: rpcCli,
	}
}

func (g *BlockManagerImpl) getBlockNumber() (string, error) {
	resp, err := g.rpcCli.Call(blockNumber)
	if err != nil {
		return "", err
	}

	return resp.Result.(string), nil
}

func (g *BlockManagerImpl) GetBlockAddress() (string, error) {
	blockNumInHex, err := g.getBlockNumber()
	if err != nil {
		return "", err
	}

	values := make([]models.Transaction, 0)

	blockNum, err := strconv.ParseInt(blockNumInHex[2:], 16, 64)
	if err != nil {
		return "", err
	}

	for i := blockNum - 100; i <= blockNum; i++ {
		g.wg.Add(1)
		go func(block int64) {
			defer g.wg.Done()

			var txs models.BlockByNumberResponse
			_ = g.rpcCli.CallFor(&txs, getBlockByNumber, fmt.Sprintf("0x%s", strconv.FormatInt(block, 16)), true)

			if len(txs.Transactions) == 0 {
				return
			}

			g.mu.Lock()
			defer g.mu.Unlock()
			values = append(values, g.getMaxValue(txs.Transactions))
		}(i)
	}

	g.wg.Wait()

	return g.getMaxValue(values).To, nil
}

func (g *BlockManagerImpl) getMaxValue(txs []models.Transaction) models.Transaction {
	var max models.Transaction
	for i := 0; i < len(txs)-1; i++ {
		val1, _ := strconv.ParseInt(txs[i].Value, 16, 64)
		val2, _ := strconv.ParseInt(txs[i+1].Value, 16, 64)
		if val1 > val2 {
			max = txs[i]
		} else {
			max = txs[i+1]
		}
	}

	return max
}
