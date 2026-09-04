package abi

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

const erc20ABIJSON = `[
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "from", "type": "address"},
			{"indexed": true, "name": "to", "type": "address"},
			{"indexed": false, "name": "value", "type": "uint256"}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "owner", "type": "address"},
			{"indexed": true, "name": "spender", "type": "address"},
			{"indexed": false, "name": "value", "type": "uint256"}
		],
		"name": "Approval",
		"type": "event"
	}
]`

func TestABIDecoder(t *testing.T) {
	reg := NewRegistry()
	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") // USDC

	err := reg.RegisterContract(contractAddr, erc20ABIJSON, []string{"Transfer"})
	if err != nil {
		t.Fatalf("failed to register ABI: %v", err)
	}

	parsed, _ := abi.JSON(strings.NewReader(erc20ABIJSON))
	transferEvent := parsed.Events["Transfer"]

	fromAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	toAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	val := big.NewInt(1000000)

	valBytes := common.LeftPadBytes(val.Bytes(), 32)
	fromTopic := common.BytesToHash(fromAddr.Bytes())
	toTopic := common.BytesToHash(toAddr.Bytes())

	log := ethTypes.Log{
		Address:     contractAddr,
		Topics:      []common.Hash{transferEvent.ID, fromTopic, toTopic},
		Data:        valBytes,
		BlockNumber: 12345,
		TxHash:      common.HexToHash("0xabc"),
		Index:       2,
	}

	decoded, err := reg.Decode(log, 1)
	if err != nil {
		t.Fatalf("failed to decode log: %v", err)
	}
	if decoded == nil {
		t.Fatal("expected decoded event, got nil")
	}

	if decoded.EventName != "Transfer" {
		t.Fatalf("expected event 'Transfer', got '%s'", decoded.EventName)
	}

	if decoded.Data["value"].(*big.Int).Cmp(val) != 0 {
		t.Fatalf("expected value %v, got %v", val, decoded.Data["value"])
	}

	if decoded.Data["from"].(common.Address) != fromAddr {
		t.Fatalf("expected from %v, got %v", fromAddr, decoded.Data["from"])
	}

	if decoded.Data["to"].(common.Address) != toAddr {
		t.Fatalf("expected to %v, got %v", toAddr, decoded.Data["to"])
	}
}

func BenchmarkABIDecoder(b *testing.B) {
	reg := NewRegistry()
	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	_ = reg.RegisterContract(contractAddr, erc20ABIJSON, nil)

	parsed, _ := abi.JSON(strings.NewReader(erc20ABIJSON))
	transferEvent := parsed.Events["Transfer"]

	fromAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	toAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	val := big.NewInt(50000000)

	log := ethTypes.Log{
		Address:     contractAddr,
		Topics:      []common.Hash{transferEvent.ID, common.BytesToHash(fromAddr.Bytes()), common.BytesToHash(toAddr.Bytes())},
		Data:        common.LeftPadBytes(val.Bytes(), 32),
		BlockNumber: 15000000,
		TxHash:      common.HexToHash("0x999"),
		Index:       0,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ev, err := reg.Decode(log, 1)
		if err != nil || ev == nil {
			b.Fatal(err)
		}
	}
}
