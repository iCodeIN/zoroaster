package trigger

import (
	"fmt"
	"github.com/HAL-xyz/ethrpc"
	"github.com/HAL-xyz/zoroaster/rpc"
	"github.com/HAL-xyz/zoroaster/tokenapi"
	"github.com/HAL-xyz/zoroaster/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func getLogsForBlock(client rpc.IEthRpc, blockNo int, addresses []string) ([]ethrpc.Log, error) {
	filter := ethrpc.FilterParams{
		FromBlock: fmt.Sprintf("0x%x", blockNo),
		ToBlock:   fmt.Sprintf("0x%x", blockNo),
		Address:   addresses,
	}
	return client.EthGetLogs(filter)
}

type mockETHCli struct {
	rpc.IEthRpc
}

func (cli mockETHCli) EthCall(transaction ethrpc.T, tag string) (string, error) {
	return "0x0000000000000000000000000000000000000000000000212321d502f10fbc4a", nil
}

func (z mockETHCli) EthGetTransactionByHash(hash string) (*ethrpc.Transaction, error) {
	return &ethrpc.Transaction{
		From: "0x000",
		To:   "0x111",
	}, nil
}

var mockCli mockETHCli

var mockTokenApi = tokenapi.New(mockCli)

type mockTApiCurrency struct {
	tokenapi.ITokenAPI
}

func (t mockTApiCurrency) GetExchangeRate(tokenAddress, fiatCurrency string) (float32, error) {
	return float32(202.37), nil
}

func (t mockTApiCurrency) Decimals(address string) string {
	return "6"
}

func (t mockTApiCurrency) FromWei(wei interface{}, units interface{}) string {
	return mockTokenApi.FromWei(wei, units)
}

func (t mockTApiCurrency) GetRPCCli() rpc.IEthRpc {
	return TokenApiMainnet.GetRPCCli()
}

func (t mockTApiCurrency) EncodeMethod(methodName, cntABI string, inputs []tokenapi.Input) (string, error) {
	return "0x", nil
}

func (t mockTApiCurrency) MakeEthRpcCall(cntAddress, data string, blockNumber int) (string, error) {
	return "0x0000000000000000000000000000000000000000000000212321d502f10fbc4a", nil
}

var mockTApiCurr mockTApiCurrency

var logs550, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5690550, []string{"0x494b4a86212fee251aa9019fe3cdb92a54d9efa1"})
var logs551, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5690551, []string{"0x494b4a86212fee251aa9019fe3cdb92a54d9efa1"})
var logs552, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5690552, []string{"0x494b4a86212fee251aa9019fe3cdb92a54d9efa1"})

func TestValidateFilterLog(t *testing.T) {

	logs, err := GetLogsFromFile("../resources/events/logs1.json")
	assert.NoError(t, err)

	tg, err := GetTriggerFromFile("../resources/triggers/ev1.json")
	assert.NoError(t, err)

	abiObj, err := abi.JSON(strings.NewReader(tg.ContractABI))

	res, err := validateFilterLog(&logs[0], &tg.Filters[0], &abiObj, tg.Filters[0].EventName, mockTokenApi)
	assert.NoError(t, err)
	assert.True(t, res)

	res2, err := validateFilterLog(&logs[0], &tg.Filters[1], &abiObj, tg.Filters[0].EventName, mockTokenApi)
	assert.NoError(t, err)
	assert.True(t, res2)

	res3 := validateTriggerLog(&logs[0], tg, &abiObj, "Transfer", mockTokenApi)
	assert.True(t, res3)

	res4 := validateTriggerLog(&logs[1], tg, &abiObj, "Transfer", mockTokenApi)
	assert.False(t, res4)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressFixedArrayEqAtPosition0(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0x165d3cc520d1B29718E6B7C34d81f627a927E22e",
                    "Predicate": "Eq"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v5",
                "ParameterType": "address[3]",
                "Index": 0
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressFixedArrayEqAtPosition0",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressFixedArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0x165d3cc520d1B29718E6B7C34d81f627a927E22e",
                    "Predicate": "IsIn"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v5",
                "ParameterType": "address[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressFixedArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressFixedArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v5",
                "ParameterType": "address[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressFixedArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressFixedArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "4",
                    "Predicate": "SmallerThan"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v5",
                "ParameterType": "address[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressFixedArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressFixedArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v5",
                "ParameterType": "address[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressFixedArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressFixedArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "Eq"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v5",
                "ParameterType": "address[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressFixedArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressDynamicArrayEqAtPosition0(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0xE38eCea1316b05C5E49B8b440742633E345AaE02",
                    "Predicate": "Eq"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v4",
                "ParameterType": "address[]",
                "Index": 0
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressDynamicArrayEqAtPosition0",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressDynamicArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0xE38eCea1316b05C5E49B8b440742633E345AaE02",
                    "Predicate": "IsIn"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v4",
                "ParameterType": "address[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressDynamicArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressDynamicArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v4",
                "ParameterType": "address[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressDynamicArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressDynamicArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "4",
                    "Predicate": "SmallerThan"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v4",
                "ParameterType": "address[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressDynamicArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressDynamicArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v4",
                "ParameterType": "address[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressDynamicArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x1e5aebb232ae66459d6c6144e6bfe8269362db01ae47a7a8f89b4df6feff8271
func TestAddressDynamicArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "2",
                    "Predicate": "Eq"
                },
                "EventName": "address_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v4",
                "ParameterType": "address[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestAddressDynamicArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolFixedArrayEqAtPosition1(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "false",
                    "Predicate": "Eq"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v10",
                "ParameterType": "bool[3]",
                "Index": 1
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolFixedArrayEqAtPosition1",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolFixedArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "False",
                    "Predicate": "IsIn"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v10",
                "ParameterType": "bool[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolFixedArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolFixedArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v10",
                "ParameterType": "bool[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolFixedArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// // https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolFixedArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "4",
                    "Predicate": "SmallerThan"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v10",
                "ParameterType": "bool[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolFixedArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolFixedArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v10",
                "ParameterType": "bool[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolFixedArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolFixedArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "Eq"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v10",
                "ParameterType": "bool[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolFixedArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolDynamicArrayEqAtPosition1(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "true",
                    "Predicate": "Eq"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v9",
                "ParameterType": "bool[]",
                "Index": 1
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolDynamicArrayEqAtPosition1",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolDynamicArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "False",
                    "Predicate": "IsIn"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v9",
                "ParameterType": "bool[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolDynamicArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolDynamicArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v9",
                "ParameterType": "bool[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolDynamicArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolDynamicArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "SmallerThan"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v9",
                "ParameterType": "bool[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolDynamicArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolDynamicArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v9",
                "ParameterType": "bool[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolDynamicArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

// https://rinkeby.etherscan.io/tx/0x56bce35c702186f21f4a102a116bc9c822f879a81f6d29d75502547945721d5e
func TestBoolDynamicArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "2",
                    "Predicate": "Eq"
                },
                "EventName": "bool_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v9",
                "ParameterType": "bool[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBoolDynamicArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs550, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690550, matches[0].Log.BlockNumber)
}

func TestBytes16EqWithOX(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0x20010db8000000000000000000000001",
                    "Predicate": "Eq"
                },
                "EventName": "bytes_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v11",
                "ParameterType": "bytes16",
                "Index": 1
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBytes16EqWithOX",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestBytes16Eq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "20010db8000000000000000000000001",
                    "Predicate": "Eq"
                },
                "EventName": "bytes_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v11",
                "ParameterType": "bytes16",
                "Index": 1
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestBytes16Eq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256FixedArrayEqAtPosition1(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "101112",
                    "Predicate": "Eq"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]",
                "Index": 1
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayEqAtPosition1",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256FixedArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "789",
                    "Predicate": "IsIn"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256FixedArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256FixedArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "4",
                    "Predicate": "SmallerThan"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256FixedArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256FixedArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "Eq"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256DinamicArrayEqAtPosition0(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "123",
                    "Predicate": "Eq"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v6",
                "ParameterType": "int256[]",
                "Index": 0
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256DinamicArrayEqAtPosition0",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256DinamicArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "123",
                    "Predicate": "IsIn"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v6",
                "ParameterType": "int256[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256DinamicArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256DinamicArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v6",
                "ParameterType": "int256[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256DinamicArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256DinamicArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "SmallerThan"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v6",
                "ParameterType": "int256[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256DinamicArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256DinamicArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v6",
                "ParameterType": "int256[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256DinamicArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestInt256DinamicArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "2",
                    "Predicate": "Eq"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v6",
                "ParameterType": "int256[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256DinamicArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690551, matches[0].Log.BlockNumber)
}

func TestStringFixedArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "againindex0",
                    "Predicate": "IsIn"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v3",
                "ParameterType": "string[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringFixedArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringFixedArrayEqAtPosition0(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "againindex0",
                    "Predicate": "Eq"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v3",
                "ParameterType": "string[3]",
                "Index": 0
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringFixedArrayEqAtPosition0",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringFixedArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v3",
                "ParameterType": "string[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringFixedArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringFixedArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v3",
                "ParameterType": "string[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringFixedArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringFixedArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "5",
                    "Predicate": "SmallerThan"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v3",
                "ParameterType": "string[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringFixedArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringFixedArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "Eq"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v3",
                "ParameterType": "string[3]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringFixedArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringDinamicArrayLengthInBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "SmallerThan"
                },
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v2",
                "ParameterType": "string[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringDinamicArrayLengthInBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringDinamicArrayLengthSmallerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "3",
                    "Predicate": "SmallerThan"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v2",
                "ParameterType": "string[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringDinamicArrayLengthSmallerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringDinamicArrayLengthBiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1",
                    "Predicate": "BiggerThan"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v2",
                "ParameterType": "string[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringDinamicArrayLengthBiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringDinamicArrayLengthEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "2",
                    "Predicate": "Eq"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v2",
                "ParameterType": "string[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringDinamicArrayLengthEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringDinamicArrayEqAtPosition0(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "index0",
                    "Predicate": "Eq"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v2",
                "ParameterType": "string[]",
                "Index": 0
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringDinamicArrayEqAtPosition0",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringDinamicArrayIsIn(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "index0",
                    "Predicate": "IsIn"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v2",
                "ParameterType": "string[]"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringDinamicArrayIsIn",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestStringEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "test",
                    "Predicate": "Eq"
                },
                "EventName": "string_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v1",
                "ParameterType": "string"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestStringEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs552, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5690552, matches[0].Log.BlockNumber)
}

func TestUint8Eq(t *testing.T) {
	js := `{
       "Filters": [
           {
               "Condition": {
                   "Attribute": "0000000000000000000000000000000000000000000000000000000000000001",
                   "Predicate": "Eq"
               },
               "EventName": "OrderApprovedPartOne",
               "FilterType": "CheckEventParameter",
               "ParameterName": "side",
               "ParameterType": "uint8"
           }
       ],
       "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenTransferProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"target\",\"type\":\"address\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"extradata\",\"type\":\"bytes\"}],\"name\":\"staticCall\",\"outputs\":[{\"name\":\"result\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumMakerProtocolFee\",\"type\":\"uint256\"}],\"name\":\"changeMinimumMakerProtocolFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumTakerProtocolFee\",\"type\":\"uint256\"}],\"name\":\"changeMinimumTakerProtocolFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"array\",\"type\":\"bytes\"},{\"name\":\"desired\",\"type\":\"bytes\"},{\"name\":\"mask\",\"type\":\"bytes\"}],\"name\":\"guardedArrayReplace\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumTakerProtocolFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"codename\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"testCopyAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"arrToCopy\",\"type\":\"bytes\"}],\"name\":\"testCopy\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"}],\"name\":\"calculateCurrentPrice_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newProtocolFeeRecipient\",\"type\":\"address\"}],\"name\":\"changeProtocolFeeRecipient\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"buyCalldata\",\"type\":\"bytes\"},{\"name\":\"buyReplacementPattern\",\"type\":\"bytes\"},{\"name\":\"sellCalldata\",\"type\":\"bytes\"},{\"name\":\"sellReplacementPattern\",\"type\":\"bytes\"}],\"name\":\"orderCalldataCanMatch\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"validateOrder_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"basePrice\",\"type\":\"uint256\"},{\"name\":\"extra\",\"type\":\"uint256\"},{\"name\":\"listingTime\",\"type\":\"uint256\"},{\"name\":\"expirationTime\",\"type\":\"uint256\"}],\"name\":\"calculateFinalPrice\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeRecipient\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"}],\"name\":\"hashOrder_\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[14]\"},{\"name\":\"uints\",\"type\":\"uint256[18]\"},{\"name\":\"feeMethodsSidesKindsHowToCalls\",\"type\":\"uint8[8]\"},{\"name\":\"calldataBuy\",\"type\":\"bytes\"},{\"name\":\"calldataSell\",\"type\":\"bytes\"},{\"name\":\"replacementPatternBuy\",\"type\":\"bytes\"},{\"name\":\"replacementPatternSell\",\"type\":\"bytes\"},{\"name\":\"staticExtradataBuy\",\"type\":\"bytes\"},{\"name\":\"staticExtradataSell\",\"type\":\"bytes\"}],\"name\":\"ordersCanMatch_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"},{\"name\":\"orderbookInclusionDesired\",\"type\":\"bool\"}],\"name\":\"approveOrder_\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumMakerProtocolFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"}],\"name\":\"hashToSign_\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelledOrFinalized\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"exchangeToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder_\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[14]\"},{\"name\":\"uints\",\"type\":\"uint256[18]\"},{\"name\":\"feeMethodsSidesKindsHowToCalls\",\"type\":\"uint8[8]\"},{\"name\":\"calldataBuy\",\"type\":\"bytes\"},{\"name\":\"calldataSell\",\"type\":\"bytes\"},{\"name\":\"replacementPatternBuy\",\"type\":\"bytes\"},{\"name\":\"replacementPatternSell\",\"type\":\"bytes\"},{\"name\":\"staticExtradataBuy\",\"type\":\"bytes\"},{\"name\":\"staticExtradataSell\",\"type\":\"bytes\"},{\"name\":\"vs\",\"type\":\"uint8[2]\"},{\"name\":\"rssMetadata\",\"type\":\"bytes32[5]\"}],\"name\":\"atomicMatch_\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[7]\"},{\"name\":\"uints\",\"type\":\"uint256[9]\"},{\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"name\":\"side\",\"type\":\"uint8\"},{\"name\":\"saleKind\",\"type\":\"uint8\"},{\"name\":\"howToCall\",\"type\":\"uint8\"},{\"name\":\"calldata\",\"type\":\"bytes\"},{\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"name\":\"staticExtradata\",\"type\":\"bytes\"}],\"name\":\"validateOrderParameters_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INVERSE_BASIS_POINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[14]\"},{\"name\":\"uints\",\"type\":\"uint256[18]\"},{\"name\":\"feeMethodsSidesKindsHowToCalls\",\"type\":\"uint8[8]\"},{\"name\":\"calldataBuy\",\"type\":\"bytes\"},{\"name\":\"calldataSell\",\"type\":\"bytes\"},{\"name\":\"replacementPatternBuy\",\"type\":\"bytes\"},{\"name\":\"replacementPatternSell\",\"type\":\"bytes\"},{\"name\":\"staticExtradataBuy\",\"type\":\"bytes\"},{\"name\":\"staticExtradataSell\",\"type\":\"bytes\"}],\"name\":\"calculateMatchPrice_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"approvedOrders\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"registryAddress\",\"type\":\"address\"},{\"name\":\"tokenTransferProxyAddress\",\"type\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"protocolFeeAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerRelayerFee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerRelayerFee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerProtocolFee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerProtocolFee\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"feeRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"feeMethod\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"side\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"saleKind\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"target\",\"type\":\"address\"}],\"name\":\"OrderApprovedPartOne\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"howToCall\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"calldata\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"replacementPattern\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"staticTarget\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"staticExtradata\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"paymentToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"basePrice\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"extra\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"listingTime\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expirationTime\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"salt\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"orderbookInclusionDesired\",\"type\":\"bool\"}],\"name\":\"OrderApprovedPartTwo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"OrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"buyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"sellHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"metadata\",\"type\":\"bytes32\"}],\"name\":\"OrdersMatched\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]",
       "ContractAdd": "0x7be8076f4ea4a4ad08075c2508e481d6c946d12b",
       "TriggerName": "WAE - TestUint8Eq",
       "TriggerType": "WatchEvents"
   }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	logs, err := getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252401, []string{"0x7be8076f4ea4a4ad08075c2508e481d6c946d12b"})
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 2, len(matches))
	assert.Equal(t, 9252401, matches[0].Log.BlockNumber)
}

func TestBytes32Eq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0xc9d56e63913927c79cc33d1593c8e155e9517bf309bb63ade5d90259484d8dd1",
                    "Predicate": "Eq"
                },
                "EventName": "Fill",
                "FilterType": "CheckEventParameter",
                "ParameterName": "orderHash",
                "ParameterType": "bytes32"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"toAddress\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderToDepositAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"preSign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_GAS_LIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"name\":\"assetProxies\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"KECCAK256_ETH_ASSET_DATA\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"targetOrderEpoch\",\"type\":\"uint256\"}],\"name\":\"cancelOrdersUpTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"depositAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"name\":\"getAssetProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"target\",\"type\":\"address\"}],\"name\":\"addWithdrawOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"validatorAddress\",\"type\":\"address\"},{\"name\":\"approval\",\"type\":\"bool\"}],\"name\":\"setSignatureValidatorApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"getOrdersInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"\",\"type\":\"tuple[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"preSigned\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ETH_ASSET_DATA\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"target\",\"type\":\"address\"}],\"name\":\"removeWithdrawOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrder\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"executeTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isWithdrawOperator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"registerAssetProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_DOMAIN_HASH\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentContextAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOrderHash\",\"type\":\"bytes32\"},{\"name\":\"newOfferAmount\",\"type\":\"uint256\"},{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orderToBeCanceled\",\"type\":\"tuple\"}],\"name\":\"updateOrder\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"newOrderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"newAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"oldOrderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"oldAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"}],\"name\":\"DepositChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"target\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"WithdrawOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"target\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"WithdrawOperatorRemoved\",\"type\":\"event\"}]",
        "ContractAdd": "0x7a6425c9b3f5521bfa5d71df710a2fb80508319b",
        "TriggerName": "WAE - TestBytes32Eq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252045, []string{"0x7a6425c9b3f5521bfa5d71df710a2fb80508319b"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9252045, matches[0].Log.BlockNumber)
}

// https://etherscan.io/tx/0x55ae08e51da4e787b7589ba9342a81091ae76f29a86b723c2e96eb32be7303d0
func TestBytesEqStartingWith0x(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0x0000000000000000000000000000000000000000000000000000000000000000",
                    "Predicate": "Eq"
                },
                "EventName": "Sent",
                "FilterType": "CheckEventParameter",
                "ParameterName": "data",
                "ParameterType": "bytes"
            }
        ],
        "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"defaultOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"holder\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"granularity\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorSend\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"authorizeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"isOperatorFor\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"holder\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"revokeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorBurn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"defaultOperators\",\"type\":\"address[]\"},{\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"name\":\"feeReceiver\",\"type\":\"address\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Sent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"AuthorizedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"RevokedOperator\",\"type\":\"event\"}]",
        "ContractAdd": "0xc2058f5d9736e8df8ba03ca3582b7cd6ac613658",
        "TriggerName": "WAE - TestBytesEqStartingWith0x",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9243327, []string{"0xc2058f5d9736e8df8ba03ca3582b7cd6ac613658"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9243327, matches[0].Log.BlockNumber)
}

// https://etherscan.io/tx/0x55ae08e51da4e787b7589ba9342a81091ae76f29a86b723c2e96eb32be7303d0
func TestBytesEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0000000000000000000000000000000000000000000000000000000000000000",
                    "Predicate": "Eq"
                },
                "EventName": "Sent",
                "FilterType": "CheckEventParameter",
                "ParameterName": "data",
                "ParameterType": "bytes"
            }
        ],
        "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"defaultOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"holder\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"granularity\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorSend\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"authorizeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"isOperatorFor\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"holder\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"revokeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorBurn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"defaultOperators\",\"type\":\"address[]\"},{\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"name\":\"feeReceiver\",\"type\":\"address\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Sent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"AuthorizedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"RevokedOperator\",\"type\":\"event\"}]",
        "ContractAdd": "0xc2058f5d9736e8df8ba03ca3582b7cd6ac613658",
        "TriggerName": "WAE - TestBytesEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9243327, []string{"0xc2058f5d9736e8df8ba03ca3582b7cd6ac613658"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9243327, matches[0].Log.BlockNumber)
}

func TestBoolEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "true",
                    "Predicate": "Eq"
                },
                "EventName": "TokenUpdated",
                "FilterType": "CheckEventParameter",
                "ParameterName": "emergencyUnlock",
                "ParameterType": "bool"
            }
        ],
        "ContractABI": "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address payable\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"minAmount\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"airdropDate\",\"type\":\"uint256\"}],\"name\":\"AirdropAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"AssetClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endDate\",\"type\":\"uint256\"}],\"name\":\"AssetLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"FeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenInactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"emergencyUnlock\",\"type\":\"bool\"}],\"name\":\"TokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensAirdropped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"WalletChanged\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"activateToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"}],\"name\":\"addToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"claim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"claimable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getAirdrops\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"destTokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"numerators\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"denominators\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"dates\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getAssetIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getLockedAsset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"enum Lock.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getTokenInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"emergencyUnlock\",\"type\":\"bool\"},{\"internalType\":\"enum Lock.TokenStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"getTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"minAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bool[]\",\"name\":\"emergencyUnlocks\",\"type\":\"bool[]\"},{\"internalType\":\"enum Lock.TokenStatus[]\",\"name\":\"statuses\",\"type\":\"uint8[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getWallet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"inactivateToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"address payable\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"numerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"date\",\"type\":\"uint256\"}],\"name\":\"setAirdrop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address payable\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"setWallet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"baseToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"numerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"date\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"updateAirdrop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"emergencyUnlock\",\"type\":\"bool\"}],\"name\":\"updateToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
        "ContractAdd": "0x73866e69c6f6f74fc48539dd541a6df8c8059e04",
        "TriggerName": "WAE - TestBoolEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9133542, []string{"0x73866e69c6f6f74fc48539dd541a6df8c8059e04"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9133542, matches[0].Log.BlockNumber)
}

func TestUint64Eq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "1578650455",
                    "Predicate": "Eq"
                },
                "EventName": "LogKill",
                "FilterType": "CheckEventParameter",
                "ParameterName": "timestamp",
                "ParameterType": "uint64"
            }
        ],
        "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"matchingEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sell_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"getBestOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"min_fill_amount\",\"type\":\"uint256\"}],\"name\":\"sellAllAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"stop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"name\":\"buy_amt\",\"type\":\"uint128\"}],\"name\":\"make\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint256\"}],\"name\":\"getBuyAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"insert\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"last_offer_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"matchingEnabled_\",\"type\":\"bool\"}],\"name\":\"setMatchingEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"del_rank\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"maxTakeAmount\",\"type\":\"uint128\"}],\"name\":\"take\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"}],\"name\":\"getMinSell\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dustId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNextUnsortedOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"close_time\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_span\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_best\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stopped\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id_\",\"type\":\"bytes32\"}],\"name\":\"bump\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sell_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"getOfferCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"max_fill_amount\",\"type\":\"uint256\"}],\"name\":\"buyAllAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isActive\",\"outputs\":[{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"offers\",\"outputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFirstUnsortedOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBetterOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_dust\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getWorseOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_near\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"dust\",\"type\":\"uint256\"}],\"name\":\"setMinSell\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isClosed\",\"outputs\":[{\"name\":\"closed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_rank\",\"outputs\":[{\"name\":\"next\",\"type\":\"uint256\"},{\"name\":\"prev\",\"type\":\"uint256\"},{\"name\":\"delb\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isOfferSorted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buyEnabled_\",\"type\":\"bool\"}],\"name\":\"setBuyEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"buy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pos\",\"type\":\"uint256\"},{\"name\":\"rounding\",\"type\":\"bool\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"buyEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"}],\"name\":\"getPayAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"close_time\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":true,\"inputs\":[{\"indexed\":true,\"name\":\"sig\",\"type\":\"bytes4\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"foo\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"bar\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"fax\",\"type\":\"bytes\"}],\"name\":\"LogNote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogItemUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogMake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogBump\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"take_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"give_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogTake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogKill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"LogSetAuthority\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"LogSetOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"LogBuyEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"min_amount\",\"type\":\"uint256\"}],\"name\":\"LogMinSell\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"LogMatchingEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogUnsortedOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogSortedOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogInsert\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogDelete\",\"type\":\"event\"}]",
        "ContractAdd": "0x39755357759ce0d7f32dc8dc45414cca409ae24e",
        "TriggerName": "WAE - TestUint64Eq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252369, []string{"0x39755357759ce0d7f32dc8dc45414cca409ae24e"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9252369, matches[0].Log.BlockNumber)
}

func TestUint128Eq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "13012580915",
                    "Predicate": "Eq"
                },
                "EventName": "LogKill",
                "FilterType": "CheckEventParameter",
                "ParameterName": "buy_amt",
                "ParameterType": "uint128"
            }
        ],
        "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"matchingEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sell_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"getBestOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"min_fill_amount\",\"type\":\"uint256\"}],\"name\":\"sellAllAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"stop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"name\":\"buy_amt\",\"type\":\"uint128\"}],\"name\":\"make\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint256\"}],\"name\":\"getBuyAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"insert\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"last_offer_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"matchingEnabled_\",\"type\":\"bool\"}],\"name\":\"setMatchingEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"del_rank\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"maxTakeAmount\",\"type\":\"uint128\"}],\"name\":\"take\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"}],\"name\":\"getMinSell\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dustId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNextUnsortedOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"close_time\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_span\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_best\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stopped\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id_\",\"type\":\"bytes32\"}],\"name\":\"bump\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sell_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"getOfferCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"max_fill_amount\",\"type\":\"uint256\"}],\"name\":\"buyAllAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isActive\",\"outputs\":[{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"offers\",\"outputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFirstUnsortedOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBetterOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_dust\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getWorseOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_near\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"dust\",\"type\":\"uint256\"}],\"name\":\"setMinSell\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isClosed\",\"outputs\":[{\"name\":\"closed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_rank\",\"outputs\":[{\"name\":\"next\",\"type\":\"uint256\"},{\"name\":\"prev\",\"type\":\"uint256\"},{\"name\":\"delb\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isOfferSorted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buyEnabled_\",\"type\":\"bool\"}],\"name\":\"setBuyEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"buy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pos\",\"type\":\"uint256\"},{\"name\":\"rounding\",\"type\":\"bool\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"buyEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"}],\"name\":\"getPayAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"close_time\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":true,\"inputs\":[{\"indexed\":true,\"name\":\"sig\",\"type\":\"bytes4\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"foo\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"bar\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"fax\",\"type\":\"bytes\"}],\"name\":\"LogNote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogItemUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogMake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogBump\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"take_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"give_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogTake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogKill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"LogSetAuthority\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"LogSetOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"LogBuyEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"min_amount\",\"type\":\"uint256\"}],\"name\":\"LogMinSell\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"LogMatchingEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogUnsortedOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogSortedOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogInsert\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogDelete\",\"type\":\"event\"}]",
        "ContractAdd": "0x39755357759ce0d7f32dc8dc45414cca409ae24e",
        "TriggerName": "WAE - TestUint128Eq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252369, []string{"0x39755357759ce0d7f32dc8dc45414cca409ae24e"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9252369, matches[0].Log.BlockNumber)
}

func XXXTestUint128EqBis(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "0000000000000000000000000000000000000000000001bd3a4571b054e80800",
                    "Predicate": "Eq"
                },
                "EventName": "LogKill",
                "FilterType": "CheckEventParameter",
                "ParameterName": "buy_amt",
                "ParameterType": "uint128"
            }
        ],
        "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"matchingEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sell_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"getBestOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"min_fill_amount\",\"type\":\"uint256\"}],\"name\":\"sellAllAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"stop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"name\":\"buy_amt\",\"type\":\"uint128\"}],\"name\":\"make\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"pay_amt\",\"type\":\"uint256\"}],\"name\":\"getBuyAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"pos\",\"type\":\"uint256\"}],\"name\":\"insert\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"last_offer_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"matchingEnabled_\",\"type\":\"bool\"}],\"name\":\"setMatchingEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"del_rank\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"maxTakeAmount\",\"type\":\"uint128\"}],\"name\":\"take\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"}],\"name\":\"getMinSell\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dustId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNextUnsortedOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"close_time\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_span\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_best\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stopped\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id_\",\"type\":\"bytes32\"}],\"name\":\"bump\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sell_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"getOfferCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"max_fill_amount\",\"type\":\"uint256\"}],\"name\":\"buyAllAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isActive\",\"outputs\":[{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"offers\",\"outputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFirstUnsortedOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBetterOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"_dust\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getWorseOffer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_near\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"dust\",\"type\":\"uint256\"}],\"name\":\"setMinSell\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isClosed\",\"outputs\":[{\"name\":\"closed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_rank\",\"outputs\":[{\"name\":\"next\",\"type\":\"uint256\"},{\"name\":\"prev\",\"type\":\"uint256\"},{\"name\":\"delb\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isOfferSorted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buyEnabled_\",\"type\":\"bool\"}],\"name\":\"setBuyEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"buy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"pos\",\"type\":\"uint256\"},{\"name\":\"rounding\",\"type\":\"bool\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"offer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"buyEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pay_gem\",\"type\":\"address\"},{\"name\":\"buy_gem\",\"type\":\"address\"},{\"name\":\"buy_amt\",\"type\":\"uint256\"}],\"name\":\"getPayAmount\",\"outputs\":[{\"name\":\"fill_amt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"close_time\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":true,\"inputs\":[{\"indexed\":true,\"name\":\"sig\",\"type\":\"bytes4\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"foo\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"bar\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"fax\",\"type\":\"bytes\"}],\"name\":\"LogNote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogItemUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"buy_gem\",\"type\":\"address\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogMake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogBump\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"take_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"give_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogTake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"pair\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buy_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"pay_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"buy_amt\",\"type\":\"uint128\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"LogKill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"LogSetAuthority\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"LogSetOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"LogBuyEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pay_gem\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"min_amount\",\"type\":\"uint256\"}],\"name\":\"LogMinSell\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"LogMatchingEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogUnsortedOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogSortedOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogInsert\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LogDelete\",\"type\":\"event\"}]",
        "ContractAdd": "0x39755357759ce0d7f32dc8dc45414cca409ae24e",
        "TriggerName": "WAE - TestUint128EqBis",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252460, []string{"0x39755357759ce0d7f32dc8dc45414cca409ae24e"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 2, len(matches))
	assert.Equal(t, 9252460, matches[0].Log.BlockNumber)
}

func TestAddressEqNotDecoded(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0x000000000000000000000000be65b13f63203c9af771e6836e9636f0026982a6",
                "Predicate": "Eq"
            },
            "EventName": "Burn",
            "FilterType": "CheckEventParameter",
            "ParameterName": "burner",
            "ParameterType": "address"
        }
    ],
    "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"burntTokenReserved\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initialPrice\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"baseRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalAssetBorrow\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"loanOrderData\",\"outputs\":[{\"name\":\"loanOrderHash\",\"type\":\"bytes32\"},{\"name\":\"leverageAmount\",\"type\":\"uint256\"},{\"name\":\"initialMarginAmount\",\"type\":\"uint256\"},{\"name\":\"maintenanceMarginAmount\",\"type\":\"uint256\"},{\"name\":\"maxDurationUnixTimestampSec\",\"type\":\"uint256\"},{\"name\":\"index\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rateMultiplier\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wethContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenizedRegistry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newTarget\",\"type\":\"address\"}],\"name\":\"setTarget\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"burntTokenReserveList\",\"outputs\":[{\"name\":\"lender\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"loanTokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bZxVault\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bZxOracle\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bZxContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"leverageList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"spreadMultiplier\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"burntTokenReserveListIndex\",\"outputs\":[{\"name\":\"index\",\"type\":\"uint256\"},{\"name\":\"isSet\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"loanOrderHashes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_newTarget\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"minter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"borrower\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"borrowAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"interestRate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"collateralTokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tradeTokenToFillAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"withdrawOnOpen\",\"type\":\"bool\"}],\"name\":\"Borrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"claimant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"assetAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"remainingTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"Claim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]",
    "ContractAdd": "0x14094949152eddbfcd073717200da82fed8dc960",
    "TriggerName": "WAE - TestAddressEqNotDecoded",
    "TriggerType": "WatchEvents"
}`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252175, []string{"0x14094949152eddbfcd073717200da82fed8dc960"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9252175, matches[0].Log.BlockNumber)
}

func TestUint256Eq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "664245",
                    "Predicate": "Eq"
                },
                "EventName": "AuctionCreated",
                "FilterType": "CheckEventParameter",
                "ParameterName": "tokenId",
                "ParameterType": "uint256"
            },
            {
                "Condition": {
                    "Attribute": "22000000000000000",
                    "Predicate": "Eq"
                },
                "EventName": "AuctionCreated",
                "FilterType": "CheckEventParameter",
                "ParameterName": "startingPrice",
                "ParameterType": "uint256"
            },
            {
                "Condition": {
                    "Attribute": "11000000000000000",
                    "Predicate": "Eq"
                },
                "EventName": "AuctionCreated",
                "FilterType": "CheckEventParameter",
                "ParameterName": "endingPrice",
                "ParameterType": "uint256"
            },
            {
                "Condition": {
                    "Attribute": "31536000",
                    "Predicate": "Eq"
                },
                "EventName": "AuctionCreated",
                "FilterType": "CheckEventParameter",
                "ParameterName": "duration",
                "ParameterType": "uint256"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_startingPrice\",\"type\":\"uint256\"},{\"name\":\"_endingPrice\",\"type\":\"uint256\"},{\"name\":\"_duration\",\"type\":\"uint256\"},{\"name\":\"_seller\",\"type\":\"address\"}],\"name\":\"createAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"bid\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isSiringClockAuction\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"getAuction\",\"outputs\":[{\"name\":\"seller\",\"type\":\"address\"},{\"name\":\"startingPrice\",\"type\":\"uint256\"},{\"name\":\"endingPrice\",\"type\":\"uint256\"},{\"name\":\"duration\",\"type\":\"uint256\"},{\"name\":\"startedAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ownerCut\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"cancelAuctionWhenPaused\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"cancelAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"getCurrentPrice\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nonFungibleContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_nftAddr\",\"type\":\"address\"},{\"name\":\"_cut\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"startingPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"endingPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"AuctionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"winner\",\"type\":\"address\"}],\"name\":\"AuctionSuccessful\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"AuctionCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]",
        "ContractAdd": "0xc7af99fe5513eb6710e6d5f44f9989da40f27f26",
        "TriggerName": "WAE - TestUint256Eq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9252357, []string{"0xc7af99fe5513eb6710e6d5f44f9989da40f27f26"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9252357, matches[0].Log.BlockNumber)
}

func TestUint256InBetween(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "339587",
                    "Predicate": "BiggerThan"
                },
                "EventName": "LogResult",
                "FilterType": "CheckEventParameter",
                "ParameterName": "ResultSerialNumber",
                "ParameterType": "uint256"
            },
            {
                "Condition": {
                    "Attribute": "339589",
                    "Predicate": "SmallerThan"
                },
                "EventName": "LogResult",
                "FilterType": "CheckEventParameter",
                "ParameterName": "ResultSerialNumber",
                "ParameterType": "uint256"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"newCallbackGasPrice\",\"type\":\"uint256\"}],\"name\":\"ownerSetCallbackGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWon\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitAsPercentOfHouse\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newHouseEdge\",\"type\":\"uint256\"}],\"name\":\"ownerSetHouseEdge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"payoutsPaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\"}],\"name\":\"ownerSetTreasury\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressToCheck\",\"type\":\"address\"}],\"name\":\"playerGetPendingTxByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newContractBalanceInWei\",\"type\":\"uint256\"}],\"name\":\"ownerUpdateContractBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newPayoutStatus\",\"type\":\"bool\"}],\"name\":\"ownerPausePayouts\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ownerChangeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMaxProfitAsPercent\",\"type\":\"uint256\"}],\"name\":\"ownerSetMaxProfitAsPercentOfHouse\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWagered\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumBet\",\"type\":\"uint256\"}],\"name\":\"ownerSetMinBet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newStatus\",\"type\":\"bool\"}],\"name\":\"ownerPauseGame\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasForOraclize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerTransferEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contractBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minBet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"playerWithdrawPendingTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalBets\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"randomQueryID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gamePaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"originalPlayerBetId\",\"type\":\"bytes32\"},{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"originalPlayerProfit\",\"type\":\"uint256\"},{\"name\":\"originalPlayerBetValue\",\"type\":\"uint256\"}],\"name\":\"ownerRefundPlayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newSafeGasToOraclize\",\"type\":\"uint32\"}],\"name\":\"ownerSetOraclizeSafeGas\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"ownerkill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdge\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rollUnder\",\"type\":\"uint256\"}],\"name\":\"playerRollDice\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdgeDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxPendingPayouts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RewardValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ProfitValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BetValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"RandomQueryID\",\"type\":\"uint256\"}],\"name\":\"LogBet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ResultSerialNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"DiceResult\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Status\",\"type\":\"int256\"},{\"indexed\":false,\"name\":\"Proof\",\"type\":\"bytes\"}],\"name\":\"LogResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RefundValue\",\"type\":\"uint256\"}],\"name\":\"LogRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"SentToAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"AmountTransferred\",\"type\":\"uint256\"}],\"name\":\"LogOwnerTransfer\",\"type\":\"event\"}]",
        "ContractAdd": "0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d",
        "TriggerName": "WAE - TestUint256InBetween",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9130794, []string{"0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9130794, matches[0].Log.BlockNumber)
}

func TestUint256BiggerThan(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "339587",
                    "Predicate": "BiggerThan"
                },
                "EventName": "LogResult",
                "FilterType": "CheckEventParameter",
                "ParameterName": "ResultSerialNumber",
                "ParameterType": "uint256"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"newCallbackGasPrice\",\"type\":\"uint256\"}],\"name\":\"ownerSetCallbackGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWon\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitAsPercentOfHouse\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newHouseEdge\",\"type\":\"uint256\"}],\"name\":\"ownerSetHouseEdge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"payoutsPaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\"}],\"name\":\"ownerSetTreasury\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressToCheck\",\"type\":\"address\"}],\"name\":\"playerGetPendingTxByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newContractBalanceInWei\",\"type\":\"uint256\"}],\"name\":\"ownerUpdateContractBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newPayoutStatus\",\"type\":\"bool\"}],\"name\":\"ownerPausePayouts\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ownerChangeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMaxProfitAsPercent\",\"type\":\"uint256\"}],\"name\":\"ownerSetMaxProfitAsPercentOfHouse\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWagered\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumBet\",\"type\":\"uint256\"}],\"name\":\"ownerSetMinBet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newStatus\",\"type\":\"bool\"}],\"name\":\"ownerPauseGame\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasForOraclize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerTransferEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contractBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minBet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"playerWithdrawPendingTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalBets\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"randomQueryID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gamePaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"originalPlayerBetId\",\"type\":\"bytes32\"},{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"originalPlayerProfit\",\"type\":\"uint256\"},{\"name\":\"originalPlayerBetValue\",\"type\":\"uint256\"}],\"name\":\"ownerRefundPlayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newSafeGasToOraclize\",\"type\":\"uint32\"}],\"name\":\"ownerSetOraclizeSafeGas\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"ownerkill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdge\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rollUnder\",\"type\":\"uint256\"}],\"name\":\"playerRollDice\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdgeDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxPendingPayouts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RewardValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ProfitValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BetValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"RandomQueryID\",\"type\":\"uint256\"}],\"name\":\"LogBet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ResultSerialNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"DiceResult\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Status\",\"type\":\"int256\"},{\"indexed\":false,\"name\":\"Proof\",\"type\":\"bytes\"}],\"name\":\"LogResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RefundValue\",\"type\":\"uint256\"}],\"name\":\"LogRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"SentToAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"AmountTransferred\",\"type\":\"uint256\"}],\"name\":\"LogOwnerTransfer\",\"type\":\"event\"}]",
        "ContractAdd": "0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d",
        "TriggerName": "WAE - TestUint256BiggerThan",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9130794, []string{"0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9130794, matches[0].Log.BlockNumber)
}

func TestUint256EqBytes32EqAddressEq(t *testing.T) {
	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "339588",
                    "Predicate": "Eq"
                },
                "EventName": "LogResult",
                "FilterType": "CheckEventParameter",
                "ParameterName": "ResultSerialNumber",
                "ParameterType": "uint256"
            },
            {
                "Condition": {
                    "Attribute": "0x20772cb5ef5a19914c0a66bb89a98ab945fa025436738e2eb9daa1e9a695569b",
                    "Predicate": "Eq"
                },
                "EventName": "LogResult",
                "FilterType": "CheckEventParameter",
                "ParameterName": "BetID",
                "ParameterType": "bytes32"
            },
            {
                "Condition": {
                    "Attribute": "0x01da8b0481e8f8bd6af2352a778970dd61b3a067",
                    "Predicate": "Eq"
                },
                "EventName": "LogResult",
                "FilterType": "CheckEventParameter",
                "ParameterName": "PlayerAddress",
                "ParameterType": "address"
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"newCallbackGasPrice\",\"type\":\"uint256\"}],\"name\":\"ownerSetCallbackGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWon\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitAsPercentOfHouse\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newHouseEdge\",\"type\":\"uint256\"}],\"name\":\"ownerSetHouseEdge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"payoutsPaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\"}],\"name\":\"ownerSetTreasury\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressToCheck\",\"type\":\"address\"}],\"name\":\"playerGetPendingTxByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newContractBalanceInWei\",\"type\":\"uint256\"}],\"name\":\"ownerUpdateContractBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newPayoutStatus\",\"type\":\"bool\"}],\"name\":\"ownerPausePayouts\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ownerChangeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMaxProfitAsPercent\",\"type\":\"uint256\"}],\"name\":\"ownerSetMaxProfitAsPercentOfHouse\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWagered\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumBet\",\"type\":\"uint256\"}],\"name\":\"ownerSetMinBet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newStatus\",\"type\":\"bool\"}],\"name\":\"ownerPauseGame\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasForOraclize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerTransferEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contractBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minBet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"playerWithdrawPendingTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalBets\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"randomQueryID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gamePaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"originalPlayerBetId\",\"type\":\"bytes32\"},{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"originalPlayerProfit\",\"type\":\"uint256\"},{\"name\":\"originalPlayerBetValue\",\"type\":\"uint256\"}],\"name\":\"ownerRefundPlayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newSafeGasToOraclize\",\"type\":\"uint32\"}],\"name\":\"ownerSetOraclizeSafeGas\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"ownerkill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdge\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rollUnder\",\"type\":\"uint256\"}],\"name\":\"playerRollDice\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdgeDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxPendingPayouts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RewardValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ProfitValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BetValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"RandomQueryID\",\"type\":\"uint256\"}],\"name\":\"LogBet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ResultSerialNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"DiceResult\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Status\",\"type\":\"int256\"},{\"indexed\":false,\"name\":\"Proof\",\"type\":\"bytes\"}],\"name\":\"LogResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RefundValue\",\"type\":\"uint256\"}],\"name\":\"LogRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"SentToAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"AmountTransferred\",\"type\":\"uint256\"}],\"name\":\"LogOwnerTransfer\",\"type\":\"event\"}]",
        "ContractAdd": "0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d",
        "TriggerName": "WAE - TestUint256EqBytes32EqAddressEq",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9130794, []string{"0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9130794, matches[0].Log.BlockNumber)
}

func TestMatchEvent9(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "50",
                "Predicate": "BiggerThan"
            },
            "EventName": "Deposit",
            "FilterType": "CheckEventParameter",
            "ParameterName": "amount",
            "ParameterType": "uint256"
        }
    ],
    "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"assertion\",\"type\":\"bool\"}],\"name\":\"assert\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"},{\"name\":\"feeWithdrawal\",\"type\":\"uint256\"}],\"name\":\"adminWithdraw\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastActiveTransaction\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"withdrawn\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"},{\"name\":\"isAdmin\",\"type\":\"bool\"}],\"name\":\"setAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"invalidOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"out\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a\",\"type\":\"uint256\"},{\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"safeSub\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"invalidateOrdersBefore\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a\",\"type\":\"uint256\"},{\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"safeMul\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"traded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"expiry\",\"type\":\"uint256\"}],\"name\":\"setInactivityReleasePeriod\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a\",\"type\":\"uint256\"},{\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"safeAdd\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tradeValues\",\"type\":\"uint256[8]\"},{\"name\":\"tradeAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8[2]\"},{\"name\":\"rs\",\"type\":\"bytes32[4]\"}],\"name\":\"trade\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"inactivityReleasePeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orderFills\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"feeAccount_\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"Order\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"get\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"give\",\"type\":\"address\"}],\"name\":\"Trade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"}]",
    "ContractAdd": "0x2a0c0dbecc7e4d658f48e01e3fa353f44050c208",
    "TriggerName": "test WAE depositi > 50",
    "TriggerType": "WatchEvents"
}`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	logs, _ := GetLogsFromFile("../resources/events/logs3.json")
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 2, len(matches))
	assert.Equal(t, "0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7", matches[0].Log.Topics[0])
	assert.Equal(t, "0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7", matches[1].Log.Topics[0])
}

func TestMatchEvent8(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0x2b13d1463b3821dd8a625e8935ab079251f1376d",
                "Predicate": "Eq"
            },
            "EventName": "Withdrawal",
            "FilterType": "CheckEventParameter",
            "ParameterName": "src",
            "ParameterType": "address"
        }
    ],
    "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"guy\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]",
    "ContractAdd": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
    "TriggerName": "WAE - Method",
    "TriggerType": "WatchEvents"
}`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9222611, []string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9222611, matches[0].Log.BlockNumber)
}

func TestMatchEvent7(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0x875c04fcadcd0ae4b369679d6d8eefaf3080de016142bddfabdfe543b430ffac",
                "Predicate": "Eq"
            },
            "EventName": "address_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v4",
            "ParameterType": "address[]"
        }
    ],
    "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":true,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":true,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":true,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":true,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
    "ContractAdd": "0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1",
    "TriggerName": "test event",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5693736, []string{"0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5693736, matches[0].Log.BlockNumber)
}

func TestMatchEvent6(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "true",
                "Predicate": "Eq"
            },
            "EventName": "bool_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v8",
            "ParameterType": "bool"
        },
        {
            "Condition": {
                "Attribute": "0xa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb49",
                "Predicate": "Eq"
            },
            "EventName": "int_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v9",
            "ParameterType": "bool[]"
        }
    ],
    "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":true,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":true,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":true,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":true,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
    "ContractAdd": "0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1",
    "TriggerName": "test event",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5693736, []string{"0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5693736, matches[0].Log.BlockNumber)
}

func TestMatchEvent5(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0x7250aa92a65150fcccca5852ea1a09d977f749a84405ea8fcbfbc6727ec4f515",
                "Predicate": "Eq"
            },
            "EventName": "int_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v6",
            "ParameterType": "int[]"
        }
    ],
    "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":true,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":true,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":true,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":true,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
    "ContractAdd": "0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1",
    "TriggerName": "test event",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5693738, []string{"0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5693738, matches[0].Log.BlockNumber)
}

func TestMatchEvent4(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658",
                "Predicate": "Eq"
            },
            "EventName": "string_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v1",
            "ParameterType": "string"
        },
        {
            "Condition": {
                "Attribute": "0x4c314d3ec7fe572b5f0ca8d4231464e89602faea100b49f287ef5148bfb5b776",
                "Predicate": "Eq"
            },
            "EventName": "string_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v2",
            "ParameterType": "string[]"
        },
        {
            "Condition": {
                "Attribute": "0xe7cc7564e647aae3b7253c8ab67ae03afc76f838e01ce364433dba9960e50afd",
                "Predicate": "Eq"
            },
            "EventName": "string_event",
            "FilterType": "CheckEventParameter",
            "ParameterName": "v3",
            "ParameterType": "string[3]"
        }
    ],
    "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":true,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":true,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":true,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":true,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
    "ContractAdd": "0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1",
    "TriggerName": "test event",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiRinkeby.GetRPCCli(), 5693738, []string{"0x63cbf20c5e2a2a6599627fdce8b9f0cc3b782be1"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 5693738, matches[0].Log.BlockNumber)
}

func TestMatchEvent3(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d",
                "Predicate": "Eq"
            },
            "FilterType": "BasicFilter",
            "ParameterName": "To"
        },
        {
            "Condition": {
                "Attribute": "335632",
                "Predicate": "Eq"
            },
            "EventName": "LogResult",
            "FilterType": "CheckEventParameter",
            "ParameterName": "ResultSerialNumber",
            "ParameterType": "uint256"
        },
        {
            "Condition": {
                "Attribute": "0x49c2381f46efd87fbc3e6662593bf4992a6e027ca569bddd35b3dce2c2f9ec23",
                "Predicate": "Eq"
            },
            "EventName": "LogResult",
            "FilterType": "CheckEventParameter",
            "ParameterName": "BetID",
            "ParameterType": "bytes32"
        },
        {
            "Condition": {
                "Attribute": "0x4236daa27a262fe6baf9bb43ade5e41f8f7498f9",
                "Predicate": "Eq"
            },
            "EventName": "LogResult",
            "FilterType": "CheckEventParameter",
            "ParameterName": "PlayerAddress",
            "ParameterType": "address"
        }
    ],
    "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"newCallbackGasPrice\",\"type\":\"uint256\"}],\"name\":\"ownerSetCallbackGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWon\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitAsPercentOfHouse\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newHouseEdge\",\"type\":\"uint256\"}],\"name\":\"ownerSetHouseEdge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"payoutsPaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\"}],\"name\":\"ownerSetTreasury\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressToCheck\",\"type\":\"address\"}],\"name\":\"playerGetPendingTxByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newContractBalanceInWei\",\"type\":\"uint256\"}],\"name\":\"ownerUpdateContractBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newPayoutStatus\",\"type\":\"bool\"}],\"name\":\"ownerPausePayouts\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ownerChangeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMaxProfitAsPercent\",\"type\":\"uint256\"}],\"name\":\"ownerSetMaxProfitAsPercentOfHouse\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWagered\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumBet\",\"type\":\"uint256\"}],\"name\":\"ownerSetMinBet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newStatus\",\"type\":\"bool\"}],\"name\":\"ownerPauseGame\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasForOraclize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerTransferEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contractBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minBet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"playerWithdrawPendingTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalBets\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"randomQueryID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gamePaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"originalPlayerBetId\",\"type\":\"bytes32\"},{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"originalPlayerProfit\",\"type\":\"uint256\"},{\"name\":\"originalPlayerBetValue\",\"type\":\"uint256\"}],\"name\":\"ownerRefundPlayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newSafeGasToOraclize\",\"type\":\"uint32\"}],\"name\":\"ownerSetOraclizeSafeGas\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"ownerkill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdge\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rollUnder\",\"type\":\"uint256\"}],\"name\":\"playerRollDice\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdgeDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxPendingPayouts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RewardValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ProfitValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BetValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"RandomQueryID\",\"type\":\"uint256\"}],\"name\":\"LogBet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ResultSerialNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"DiceResult\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Status\",\"type\":\"int256\"},{\"indexed\":false,\"name\":\"Proof\",\"type\":\"bytes\"}],\"name\":\"LogResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RefundValue\",\"type\":\"uint256\"}],\"name\":\"LogRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"SentToAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"AmountTransferred\",\"type\":\"uint256\"}],\"name\":\"LogOwnerTransfer\",\"type\":\"event\"}]",
    "ContractAdd": "0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d",
    "TriggerName": "WAE",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	logs, _ := GetLogsFromFile("../resources/events/logs2.json")
	matches := MatchEvent(tg, logs, TokenApiMainnet)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9098826, matches[0].Log.BlockNumber)
}

func TestMatchEvent2(t *testing.T) {
	js := `{
    "Filters": [
        {
            "Condition": {
                "Attribute": "0xd3a6fdd4408f5fd15623abbae9041025a337314d",
                "Predicate": "Eq"
            },
            "EventName": "Fill",
            "FilterType": "CheckEventParameter",
            "ParameterName": "makerAddress",
            "ParameterType": "address"
        },
        {
            "Condition": {
                "Attribute": "0x0000000000000000000000000d056bb17ad4df5593b93a1efc29cb35ba4aa38d",
                "Predicate": "Eq"
            },
            "EventName": "Fill",
            "FilterType": "CheckEventParameter",
            "ParameterName": "feeRecipientAddress",
            "ParameterType": "address"
        }
    ],
    "ContractABI": "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"preSign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"right\",\"type\":\"tuple\"},{\"name\":\"leftMakerAssetSpreadAmount\",\"type\":\"uint256\"}],\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrderNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"name\":\"assetProxies\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrKillOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"targetOrderEpoch\",\"type\":\"uint256\"}],\"name\":\"cancelOrdersUpTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrdersNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"name\":\"getAssetProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrKillOrder\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"validatorAddress\",\"type\":\"address\"},{\"name\":\"approval\",\"type\":\"bool\"}],\"name\":\"setSignatureValidatorApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"getOrdersInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"\",\"type\":\"tuple[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"preSigned\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrder\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"executeTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"registerAssetProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ZRX_ASSET_DATA\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_DOMAIN_HASH\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentContextAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_zrxAssetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"}]",
    "ContractAdd": "0x080bf510fcbf18b91105470639e9561022937712",
    "TriggerName": "WAE",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 9099675, []string{"0x080bf510fcbf18b91105470639e9561022937712"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 9099675, matches[0].Log.BlockNumber)
}

func TestMatchEvent1(t *testing.T) {

	logs, _ := GetLogsFromFile("../resources/events/logs1.json")

	tg1, err := GetTriggerFromFile("../resources/triggers/ev1.json")
	assert.NoError(t, err)
	matches1 := MatchEvent(tg1, logs, mockTokenApi)

	assert.Equal(t, 1, len(matches1))
	assert.Equal(t, "677420000", matches1[0].EventParams["value"])

	tg2, err := GetTriggerFromFile("../resources/triggers/ev2.json")
	assert.NoError(t, err)
	matches2 := MatchEvent(tg2, logs, mockTokenApi)

	assert.Equal(t, 3, len(matches2))
	assert.Equal(t, "677420000", matches2[0].EventParams["value"])
	assert.Equal(t, "771470000", matches2[1].EventParams["value"])
	assert.Equal(t, "607760000", matches2[2].EventParams["value"])

	assert.Equal(t, 3, len(matches2[0].EventParams))
	assert.Equal(t, "0xf750f050e5596eb9480523eef7260b1535a689bd", matches2[0].EventParams["from"])
	assert.Equal(t, "0xcd95b32c98423172e04b1c76841e5a73f4532a7f", matches2[0].EventParams["to"])
	assert.Equal(t, "677420000", matches2[0].EventParams["value"])

	// testing ToPersistent()
	persistentJson, err := utils.GimmePrettyJson(matches1[0].ToPersistent())
	expectedJsn := `{
  "ContractAdd": "0xdac17f958d2ee523a2206206994597c13d831ec7",
  "EventName": "Transfer",
  "EventData": {
    "EventParameters": {
      "from": "0xf750f050e5596eb9480523eef7260b1535a689bd",
      "to": "0xcd95b32c98423172e04b1c76841e5a73f4532a7f",
      "value": "677420000"
    },
    "Data": "0x0000000000000000000000000000000000000000000000000000000028609be0",
    "Topics": [
      "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
      "0x000000000000000000000000f750f050e5596eb9480523eef7260b1535a689bd",
      "0x000000000000000000000000cd95b32c98423172e04b1c76841e5a73f4532a7f"
    ]
  },
  "Transaction": {
    "BlockHash": "0xf3d70d822816015f26843d378b8c1d5d5da62f5d346f3e86d91a0c2463d30543",
    "BlockNumber": 8496661,
    "BlockTimestamp": 0,
    "Hash": "0xf44984a4b533ac0e7b608c881a856eff44ee8c17b9f4dcf8b4ee74e9c10c0455"
  }
}`
	ok, err := utils.AreEqualJSON(persistentJson, expectedJsn)

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestMatchEventEmitted(t *testing.T) {

	js := `{
  "TriggerName":"Watch an Event",
  "TriggerType":"WatchEvents",
  "ContractAdd":"0xdac17f958d2ee523a2206206994597c13d831ec7",
  "ContractABI":"[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_upgradedAddress\",\"type\":\"address\"}],\"name\":\"deprecate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_evilUser\",\"type\":\"address\"}],\"name\":\"addBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradedAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_maker\",\"type\":\"address\"}],\"name\":\"getBlackListStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBasisPoints\",\"type\":\"uint256\"},{\"name\":\"newMaxFee\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlackListed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_clearedUser\",\"type\":\"address\"}],\"name\":\"removeBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blackListedUser\",\"type\":\"address\"}],\"name\":\"destroyBlackFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialSupply\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"Deprecate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"feeBasisPoints\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maxFee\",\"type\":\"uint256\"}],\"name\":\"Params\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_blackListedUser\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlackFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"AddedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"RemovedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]",
  "Filters":[
    {
      "FilterType":"CheckEventEmitted",
      "EventName": "Transfer"
    }
  ]
}`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	logs, _ := GetLogsFromFile("../resources/events/logs1.json")

	matches := MatchEvent(tg, logs, mockTokenApi)
	assert.Equal(t, 7, len(matches))
}

func TestHandleNullTerminatedStrings(t *testing.T) {

	js := `
{
  "Filters": [
    {
      "EventName": "PositionTokensMinted",
      "FilterType": "CheckEventParameter",
      "ParameterName":"contractName",
      "ParameterType":"string",
      "Condition":{
        "Predicate":"Eq",
        "Attribute":"MRI-BTC-28D-20200817"
      }
    }
  ],
  "ContractABI": "[{\"constant\":true,\"inputs\":[],\"name\":\"getLatestMarketCollateralPool\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_DURATION_DAYS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINTER_BRIDGE_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_DURATION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"qtyToMint\",\"type\":\"uint256\"},{\"name\":\"longTokenRecipient\",\"type\":\"address\"},{\"name\":\"shortTokenRecipient\",\"type\":\"address\"}],\"name\":\"mintPositionTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"inputAddress\",\"type\":\"address\"}],\"name\":\"getUserAddressOrDSProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"COLLATERAL_TOKEN_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateRequiredCollateral\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"createDSProxyWallet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_honeyLemonOracleAddress\",\"type\":\"address\"}],\"name\":\"setOracleAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"market\",\"type\":\"address\"}],\"name\":\"getCollateralPool\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"currentMRI\",\"type\":\"uint256\"},{\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"generateContractSpecs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[7]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_minterBridgeAddress\",\"type\":\"address\"}],\"name\":\"setMinterBridgeAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"mri\",\"type\":\"uint256\"},{\"name\":\"marketContractAddress\",\"type\":\"address\"}],\"name\":\"settleMarketContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenAddresses\",\"type\":\"address[]\"},{\"name\":\"tokensToRedeem\",\"type\":\"uint256[]\"}],\"name\":\"batchRedeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"lookbackIndexValue\",\"type\":\"uint256\"},{\"name\":\"currentIndexValue\",\"type\":\"uint256\"},{\"name\":\"marketAndsTokenNames\",\"type\":\"bytes32[3]\"},{\"name\":\"newMarketExpiration\",\"type\":\"uint256\"}],\"name\":\"dailySettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"marketContracts\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLatestMri\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getExpiringMarketContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isDailyContractDeployed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLatestMarketContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"marketContractSpecs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"makerAddress\",\"type\":\"address\"}],\"name\":\"getFillableAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"dSProxyToAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAllMarketContracts\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"makerAddresses\",\"type\":\"address[]\"}],\"name\":\"getFillableAmounts\",\"outputs\":[{\"name\":\"fillableAmounts\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_COLLATERAL_RATIO\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"marketContractFactoryMPX\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"HONEY_LEMON_ORACLE_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"addressToDSProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_marketContractFactoryMPX\",\"type\":\"address\"},{\"name\":\"_honeyLemonOracle\",\"type\":\"address\"},{\"name\":\"_minterBridge\",\"type\":\"address\"},{\"name\":\"_wBTCTokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"qtyToMint\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"marketId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"contractName\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"longTokenRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"longTokenDSProxy\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"shortTokenRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"shortTokenDSProxy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"latestMarketContract\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"longTokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"shortTokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"PositionTokensMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"revenuePerUnit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"MarketContractSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"currentMRI\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"contractName\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"expiration\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"collateralPerUnit\",\"type\":\"uint256\"}],\"name\":\"MarketContractDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"DSProxy\",\"type\":\"address\"}],\"name\":\"dSProxyCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]",
  "ContractAdd": "0x9ceb5486eD0F3F2DBCaE906E4192472e88657983",
  "TriggerName": "Honeylemon 28day Position Tokens Minted",
  "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 10679595, []string{"0x9ceb5486eD0F3F2DBCaE906E4192472e88657983"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Len(t, matches, 1)
}

func TestIntConversion(t *testing.T) {

	js := `
{
  "Filters": [
    {
      "EventName": "LogTaskSubmitted",
      "FilterType": "CheckEventParameter",
      "ParameterName":"taskReceiptId",
      "ParameterType":"uint256",
      "Condition":{
        "Predicate":"Eq",
        "Attribute":"68"
      }
    }
  ],
  "ContractABI": "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"gelatoGasPriceOracle\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"oracleRequestData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"gelatoMaxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"internalGasRequirement\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minExecutorStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"executorSuccessShare\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sysAdminSuccessShare\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSuccessShare\",\"type\":\"uint256\"}],\"internalType\":\"struct IGelatoSysAdmin.GelatoSysAdminInitialState\",\"name\":\"_\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskReceiptId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"LogCanExecFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskReceiptId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"executorRefund\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"LogExecReverted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskReceiptId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"executorSuccessFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sysAdminSuccessFee\",\"type\":\"uint256\"}],\"name\":\"LogExecSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldExecutor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newExecutor\",\"type\":\"address\"}],\"name\":\"LogExecutorAssignedExecutor\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawAmount\",\"type\":\"uint256\"}],\"name\":\"LogExecutorBalanceWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldStake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newStake\",\"type\":\"uint256\"}],\"name\":\"LogExecutorStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldShare\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newShare\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"LogExecutorSuccessShareSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"}],\"name\":\"LogExecutorUnstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newProviderFunds\",\"type\":\"uint256\"}],\"name\":\"LogFundsProvided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"realWithdrawAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newProviderFunds\",\"type\":\"uint256\"}],\"name\":\"LogFundsUnprovided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOracle\",\"type\":\"address\"}],\"name\":\"LogGelatoGasPriceOracleSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMaxGas\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMaxGas\",\"type\":\"uint256\"}],\"name\":\"LogGelatoMaxGasSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldRequirment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRequirment\",\"type\":\"uint256\"}],\"name\":\"LogInternalGasRequirementSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMin\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMin\",\"type\":\"uint256\"}],\"name\":\"LogMinExecutorStakeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newData\",\"type\":\"bytes\"}],\"name\":\"LogOracleRequestDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldExecutor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newExecutor\",\"type\":\"address\"}],\"name\":\"LogProviderAssignedExecutor\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"LogProviderModuleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"LogProviderModuleRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"LogSysAdminFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldShare\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newShare\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"LogSysAdminSuccessShareSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskReceiptId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"cancellor\",\"type\":\"address\"}],\"name\":\"LogTaskCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"taskSpecHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldTaskSpecGasPriceCeil\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newTaskSpecGasPriceCeil\",\"type\":\"uint256\"}],\"name\":\"LogTaskSpecGasPriceCeilSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"taskSpecHash\",\"type\":\"bytes32\"}],\"name\":\"LogTaskSpecProvided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"taskSpecHash\",\"type\":\"bytes32\"}],\"name\":\"LogTaskSpecUnprovided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskReceiptId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"taskReceiptHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"struct TaskReceipt\",\"name\":\"taskReceipt\",\"type\":\"tuple\"}],\"name\":\"LogTaskSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXEC_TX_OVERHEAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NO_CEIL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IGelatoProviderModule[]\",\"name\":\"_modules\",\"type\":\"address[]\"}],\"name\":\"addProviderModules\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskReceipt\",\"name\":\"_TR\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gelatoGasPrice\",\"type\":\"uint256\"}],\"name\":\"canExec\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task\",\"name\":\"_task\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_expiryDate\",\"type\":\"uint256\"}],\"name\":\"canSubmitTask\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskReceipt\",\"name\":\"_TR\",\"type\":\"tuple\"}],\"name\":\"cancelTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentTaskReceiptId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskReceipt\",\"name\":\"_TR\",\"type\":\"tuple\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskReceipt\",\"name\":\"taskReceipt\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gelatoGasPrice\",\"type\":\"uint256\"}],\"name\":\"executionWrapper\",\"outputs\":[{\"internalType\":\"enum GelatoCore.ExecutionResult\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newExecutor\",\"type\":\"address\"}],\"name\":\"executorAssignsExecutor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"executorByProvider\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"executorProvidersCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"executorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_gas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPrice\",\"type\":\"uint256\"}],\"name\":\"executorSuccessFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executorSuccessShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gelatoGasPriceOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gelatoMaxGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskReceipt\",\"name\":\"_TR\",\"type\":\"tuple\"}],\"name\":\"hashTaskReceipt\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition[]\",\"name\":\"conditions\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskSpec\",\"name\":\"_taskSpec\",\"type\":\"tuple\"}],\"name\":\"hashTaskSpec\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"internalGasRequirement\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"}],\"name\":\"isExecutorAssigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"}],\"name\":\"isExecutorMinStaked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"isModuleProvided\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_gelatoMaxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gelatoGasPrice\",\"type\":\"uint256\"}],\"name\":\"isProviderLiquid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task\",\"name\":\"_task\",\"type\":\"tuple\"}],\"name\":\"isTaskProvided\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"res\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"contract IGelatoCondition[]\",\"name\":\"conditions\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskSpec\",\"name\":\"_taskSpec\",\"type\":\"tuple\"}],\"name\":\"isTaskSpecProvided\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_gelatoMaxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gelatoGasPrice\",\"type\":\"uint256\"}],\"name\":\"minExecProviderFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minExecutorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"provider\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cycleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"submissionsLeft\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskReceipt[]\",\"name\":\"_taskReceipts\",\"type\":\"tuple[]\"}],\"name\":\"multiCancelTasks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"contract IGelatoCondition[]\",\"name\":\"conditions\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskSpec[]\",\"name\":\"_taskSpecs\",\"type\":\"tuple[]\"},{\"internalType\":\"contract IGelatoProviderModule[]\",\"name\":\"_modules\",\"type\":\"address[]\"}],\"name\":\"multiProvide\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_providers\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"_newExecutor\",\"type\":\"address\"}],\"name\":\"multiReassignProviders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"contract IGelatoCondition[]\",\"name\":\"conditions\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskSpec[]\",\"name\":\"_taskSpecs\",\"type\":\"tuple[]\"},{\"internalType\":\"contract IGelatoProviderModule[]\",\"name\":\"_modules\",\"type\":\"address[]\"}],\"name\":\"multiUnprovide\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleRequestData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"provideFunds\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition[]\",\"name\":\"conditions\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskSpec[]\",\"name\":\"_taskSpecs\",\"type\":\"tuple[]\"}],\"name\":\"provideTaskSpecs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newExecutor\",\"type\":\"address\"}],\"name\":\"providerAssignsExecutor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task\",\"name\":\"_task\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_gelatoGasPrice\",\"type\":\"uint256\"}],\"name\":\"providerCanExec\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"providerFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userProxy\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task\",\"name\":\"_task\",\"type\":\"tuple\"}],\"name\":\"providerModuleChecks\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"providerModules\",\"outputs\":[{\"internalType\":\"contract IGelatoProviderModule[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IGelatoProviderModule[]\",\"name\":\"_modules\",\"type\":\"address[]\"}],\"name\":\"removeProviderModules\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_percentage\",\"type\":\"uint256\"}],\"name\":\"setExecutorSuccessShare\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOracle\",\"type\":\"address\"}],\"name\":\"setGelatoGasPriceOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newMaxGas\",\"type\":\"uint256\"}],\"name\":\"setGelatoMaxGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newRequirement\",\"type\":\"uint256\"}],\"name\":\"setInternalGasRequirement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newMin\",\"type\":\"uint256\"}],\"name\":\"setMinExecutorStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_requestData\",\"type\":\"bytes\"}],\"name\":\"setOracleRequestData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_percentage\",\"type\":\"uint256\"}],\"name\":\"setSysAdminSuccessShare\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_taskSpecHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceCeil\",\"type\":\"uint256\"}],\"name\":\"setTaskSpecGasPriceCeil\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeExecutor\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task\",\"name\":\"_task\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_expiryDate\",\"type\":\"uint256\"}],\"name\":\"submitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"_tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_sumOfRequestedTaskSubmits\",\"type\":\"uint256\"}],\"name\":\"submitTaskChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"contract IGelatoProviderModule\",\"name\":\"module\",\"type\":\"address\"}],\"internalType\":\"struct Provider\",\"name\":\"_provider\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition\",\"name\":\"inst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"struct Condition[]\",\"name\":\"conditions\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfProviderGasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct Task[]\",\"name\":\"_tasks\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_cycles\",\"type\":\"uint256\"}],\"name\":\"submitTaskCycle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sysAdminFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_gas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPrice\",\"type\":\"uint256\"}],\"name\":\"sysAdminSuccessFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sysAdminSuccessShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"taskReceiptHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"taskSpecGasPriceCeil\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSuccessShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawAmount\",\"type\":\"uint256\"}],\"name\":\"unprovideFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"realWithdrawAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contract IGelatoCondition[]\",\"name\":\"conditions\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enum Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"enum DataFlow\",\"name\":\"dataFlow\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"termsOkCheck\",\"type\":\"bool\"}],\"internalType\":\"struct Action[]\",\"name\":\"actions\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeil\",\"type\":\"uint256\"}],\"internalType\":\"struct TaskSpec[]\",\"name\":\"_taskSpecs\",\"type\":\"tuple[]\"}],\"name\":\"unprovideTaskSpecs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstakeExecutor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawAmount\",\"type\":\"uint256\"}],\"name\":\"withdrawExcessExecutorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"realWithdrawAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address payable\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"withdrawSysAdminFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"realWithdrawAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
  "ContractAdd": "0x1d681d76ce96E4d70a88A00EBbcfc1E47808d0b8",
  "TriggerName": "LogTaskSubmission",
  "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 10696118, []string{"0x1d681d76ce96E4d70a88A00EBbcfc1E47808d0b8"})
	matches := MatchEvent(tg, logs, mockTokenApi)

	assert.Len(t, matches, 1)
	assert.Equal(t, "68", matches[0].EventParams["taskReceiptId"])
	assert.Equal(t, "0x1c48f17fd863b1c440a148fb04e33a37d4a817dae10233ad315f1e68ab800600", matches[0].EventParams["taskReceiptHash"])
}

func TestBasicFilters(t *testing.T) {

	js := `
{
  "Filters": [
    {
      "EventName": "CreateStream",
      "FilterType": "CheckEventEmitted"
    },
    {
      "FilterType":"BasicFilter",
      "ParameterName":"From",
      "Condition":{
        "Predicate":"Eq",
        "Attribute":"0xf3ad7a80c7debe37db5cee1e3ed45f31a5629e5e"
      }
    },
    {
      "FilterType":"BasicFilter",
      "ParameterName":"To",
      "Condition":{
        "Predicate":"Eq",
        "Attribute":"0xbd6a40bb904aea5a49c59050b5395f7484a4203d"
      }
    }
  ],
  "ContractABI": "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getEarnings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextStreamId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getCompoundingStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exchangeRateInitial\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"cancelStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFromStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"createCompoundingStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"interestOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feePercentage\",\"type\":\"uint256\"}],\"name\":\"updateFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"isCompoundingStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"deltaOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"takeEarnings\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"createStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cTokenManager\",\"outputs\":[{\"internalType\":\"contract ICTokenManager\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"mantissa\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cTokenManagerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"exchangeRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"CreateCompoundingStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"name\":\"PayInterest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TakeEarnings\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"CreateStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFromStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientBalance\",\"type\":\"uint256\"}],\"name\":\"CancelStream\",\"type\":\"event\"}]",
  "ContractAdd": "0xa4fc358455febe425536fd1878be67ffdbdec59a",
  "TriggerName": "Sablier All New Streams",
  "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 11020360, []string{"0xa4fc358455febe425536fd1878be67ffdbdec59a"})

	matches := MatchEvent(tg, logs, TokenApiMainnet)

	assert.Len(t, matches, 1)
	assert.Equal(t, "0xf3ad7a80c7debe37db5cee1e3ed45f31a5629e5e", matches[0].TxFrom)
}

func TestBasicFilters2(t *testing.T) {

	js := `
{
  "Filters": [
    {
      "EventName": "CreateStream",
      "FilterType": "CheckEventEmitted"
    },
    {
      "FilterType":"BasicFilter",
      "ParameterName":"From",
      "Condition":{
        "Predicate":"Eq",
        "Attribute":"0x0000000000000000000000000000000000000000"
      }
    },
    {
      "FilterType":"BasicFilter",
      "ParameterName":"To",
      "Condition":{
        "Predicate":"Eq",
        "Attribute":"0xbd6a40bb904aea5a49c59050b5395f7484a4203d"
      }
    }
  ],
  "ContractABI": "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getEarnings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextStreamId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getCompoundingStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exchangeRateInitial\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"cancelStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFromStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"createCompoundingStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"interestOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feePercentage\",\"type\":\"uint256\"}],\"name\":\"updateFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"isCompoundingStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"deltaOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"takeEarnings\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"createStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cTokenManager\",\"outputs\":[{\"internalType\":\"contract ICTokenManager\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"mantissa\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cTokenManagerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"exchangeRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"CreateCompoundingStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"name\":\"PayInterest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TakeEarnings\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"CreateStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFromStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientBalance\",\"type\":\"uint256\"}],\"name\":\"CancelStream\",\"type\":\"event\"}]",
  "ContractAdd": "0xa4fc358455febe425536fd1878be67ffdbdec59a",
  "TriggerName": "Sablier All New Streams",
  "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 11020360, []string{"0xa4fc358455febe425536fd1878be67ffdbdec59a"})

	matches := MatchEvent(tg, logs, TokenApiMainnet)

	assert.Len(t, matches, 0)
}

func TestBasicFilters3(t *testing.T) {

	js := `
{
 "Filters": [
   {
     "FilterType":"BasicFilter",
     "ParameterName":"From",
     "Condition":{
       "Predicate":"Eq",
       "Attribute":"0xf3ad7a80c7debe37db5cee1e3ed45f31a5629e5e"
     }
   },
   {
     "EventName": "CreateStream",
     "FilterType": "CheckEventEmitted"
   },
   {
     "FilterType":"BasicFilter",
     "ParameterName":"To",
     "Condition":{
       "Predicate":"Eq",
       "Attribute":"0xbd6a40bb904aea5a49c59050b5395f7484a4203d"
     }
   }
 ],
 "ContractABI": "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getEarnings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextStreamId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getCompoundingStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exchangeRateInitial\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"cancelStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFromStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"createCompoundingStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"interestOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feePercentage\",\"type\":\"uint256\"}],\"name\":\"updateFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"isCompoundingStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"deltaOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"takeEarnings\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"createStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cTokenManager\",\"outputs\":[{\"internalType\":\"contract ICTokenManager\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"mantissa\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cTokenManagerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"exchangeRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"CreateCompoundingStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"name\":\"PayInterest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TakeEarnings\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"CreateStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFromStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientBalance\",\"type\":\"uint256\"}],\"name\":\"CancelStream\",\"type\":\"event\"}]",
 "ContractAdd": "0xa4fc358455febe425536fd1878be67ffdbdec59a",
 "TriggerName": "Sablier All New Streams",
 "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 11020360, []string{"0xa4fc358455febe425536fd1878be67ffdbdec59a"})

	matches := MatchEvent(tg, logs, TokenApiMainnet)

	assert.Len(t, matches, 1)
	assert.Equal(t, "0xf3ad7a80c7debe37db5cee1e3ed45f31a5629e5e", matches[0].TxFrom)
}

func TestTxFromAndToWithoutBasicFilters(t *testing.T) {

	js := `
{
 "Filters": [
   {
     "EventName": "CreateStream",
     "FilterType": "CheckEventEmitted"
   }
 ],
 "ContractABI": "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getEarnings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextStreamId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getCompoundingStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exchangeRateInitial\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"cancelStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFromStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"createCompoundingStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"getStream\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ratePerSecond\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"interestOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feePercentage\",\"type\":\"uint256\"}],\"name\":\"updateFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"isCompoundingStream\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"}],\"name\":\"deltaOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"takeEarnings\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"createStream\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cTokenManager\",\"outputs\":[{\"internalType\":\"contract ICTokenManager\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"mantissa\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cTokenManagerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"exchangeRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderSharePercentage\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientSharePercentage\",\"type\":\"uint256\"}],\"name\":\"CreateCompoundingStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientInterest\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sablierInterest\",\"type\":\"uint256\"}],\"name\":\"PayInterest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TakeEarnings\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stopTime\",\"type\":\"uint256\"}],\"name\":\"CreateStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFromStream\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"streamId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"senderBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recipientBalance\",\"type\":\"uint256\"}],\"name\":\"CancelStream\",\"type\":\"event\"}]",
 "ContractAdd": "0xa4fc358455febe425536fd1878be67ffdbdec59a",
 "TriggerName": "Sablier All New Streams",
 "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 11020360, []string{"0xa4fc358455febe425536fd1878be67ffdbdec59a"})

	matches := MatchEvent(tg, logs, TokenApiMainnet)

	assert.Len(t, matches, 1)
	assert.Equal(t, "0xf3ad7a80c7debe37db5cee1e3ed45f31a5629e5e", matches[0].TxFrom)
}

func TestCurrencyWithExplicitCurrenciesData(t *testing.T) {
	// where ParameterCurrency is a token address (AAVE), and AttributeCurrency is a fiat symbol (USD)
	// ParameterName is in the data field as a single integer
	js := `
{
  "Filters": [
	{
	  "FilterType":"CheckEventParameter",
	  "EventName": "ProtectionAdded",
	  "ParameterName":"_reserveAmount",
	  "ParameterType":"uint256",
	  "ParameterCurrency": "0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9",
	  "Condition":{
		"Predicate":"Eq",
		"Attribute":"5622040970000000",
		"AttributeCurrency": "usd"
	  }
	}
  ],
  "ContractABI": "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_expirationTime\",\"type\":\"uint256\"}],\"name\":\"BalanceLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"BalanceUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_added\",\"type\":\"bool\"}],\"name\":\"PoolWhitelistUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_poolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_reserveAmount\",\"type\":\"uint256\"}],\"name\":\"ProtectionAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_poolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_reserveAmount\",\"type\":\"uint256\"}],\"name\":\"ProtectionRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevPoolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevReserveAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newPoolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newReserveAmount\",\"type\":\"uint256\"}],\"name\":\"ProtectionUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newAmount\",\"type\":\"uint256\"}],\"name\":\"SystemBalanceUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_expirationTime\",\"type\":\"uint256\"}],\"name\":\"addLockedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"}],\"name\":\"addPoolToWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_poolAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserveRateN\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserveRateD\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"addProtectedLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"decSystemBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"incSystemBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"}],\"name\":\"isPoolWhitelisted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"lockedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"lockedBalanceCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_endIndex\",\"type\":\"uint256\"}],\"name\":\"lockedBalanceRange\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"protectedLiquidity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"contract IDSToken\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"contract IERC20Token\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"protectedLiquidityCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"protectedLiquidityId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"protectedLiquidityIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"removeLockedBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"}],\"name\":\"removePoolFromWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"removeProtectedLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"systemBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"}],\"name\":\"totalProtectedPoolAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"}],\"name\":\"totalProtectedReserveAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newPoolAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newReserveAmount\",\"type\":\"uint256\"}],\"name\":\"updateProtectedLiquidityAmounts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"whitelistedPool\",\"outputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"whitelistedPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"whitelistedPools\",\"outputs\":[{\"internalType\":\"contract IConverterAnchor[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
  "ContractAdd": "0xf5fab5dbd2f3bf675de4cb76517d4767013cfb55",
  "TriggerName": "NewAdd1",
  "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 11646700, []string{"0xf5fab5dbd2f3bf675de4cb76517d4767013cfb55"})

	matches := MatchEvent(tg, logs, mockTApiCurr)
	assert.Len(t, matches, 1)
}

func TestCurrencyWithImplicitCurrenciesData(t *testing.T) {
	// where ParameterCurrency is an event field (_reserveToken), and AttributeCurrency is a fiat symbol (USD)
	// ParamaterName is in the data field
	js := `
{
  "Filters": [
	{
	  "FilterType":"CheckEventParameter",
	  "EventName": "ProtectionAdded",
	  "ParameterName":"_reserveAmount",
	  "ParameterType":"uint256",
	  "ParameterCurrency": "_reserveToken",
	  "Condition":{
		"Predicate":"Eq",
		"Attribute":"5622040970000000",
		"AttributeCurrency": "usd"
	  }
	}
  ],
  "ContractABI": "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_expirationTime\",\"type\":\"uint256\"}],\"name\":\"BalanceLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"BalanceUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_added\",\"type\":\"bool\"}],\"name\":\"PoolWhitelistUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_poolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_reserveAmount\",\"type\":\"uint256\"}],\"name\":\"ProtectionAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_poolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_reserveAmount\",\"type\":\"uint256\"}],\"name\":\"ProtectionRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevPoolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevReserveAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newPoolAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newReserveAmount\",\"type\":\"uint256\"}],\"name\":\"ProtectionUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newAmount\",\"type\":\"uint256\"}],\"name\":\"SystemBalanceUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_expirationTime\",\"type\":\"uint256\"}],\"name\":\"addLockedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"}],\"name\":\"addPoolToWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_poolAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserveRateN\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserveRateD\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"addProtectedLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"decSystemBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"incSystemBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"}],\"name\":\"isPoolWhitelisted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"lockedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"lockedBalanceCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_endIndex\",\"type\":\"uint256\"}],\"name\":\"lockedBalanceRange\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"protectedLiquidity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"contract IDSToken\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"contract IERC20Token\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"protectedLiquidityCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"protectedLiquidityId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"}],\"name\":\"protectedLiquidityIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"removeLockedBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"_poolAnchor\",\"type\":\"address\"}],\"name\":\"removePoolFromWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"removeProtectedLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"systemBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"}],\"name\":\"totalProtectedPoolAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IDSToken\",\"name\":\"_poolToken\",\"type\":\"address\"},{\"internalType\":\"contract IERC20Token\",\"name\":\"_reserveToken\",\"type\":\"address\"}],\"name\":\"totalProtectedReserveAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newPoolAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newReserveAmount\",\"type\":\"uint256\"}],\"name\":\"updateProtectedLiquidityAmounts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"whitelistedPool\",\"outputs\":[{\"internalType\":\"contract IConverterAnchor\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"whitelistedPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"whitelistedPools\",\"outputs\":[{\"internalType\":\"contract IConverterAnchor[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20Token\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
  "ContractAdd": "0xf5fab5dbd2f3bf675de4cb76517d4767013cfb55",
  "TriggerName": "NewAdd1",
  "TriggerType": "WatchEvents"
}
`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	var logs, _ = getLogsForBlock(TokenApiMainnet.GetRPCCli(), 11646700, []string{"0xf5fab5dbd2f3bf675de4cb76517d4767013cfb55"})

	matches := MatchEvent(tg, logs, mockTApiCurr)
	assert.Len(t, matches, 1)
}

func TestCurrencyWithExplicitCurrenciesDataArray(t *testing.T) {
	// where ParameterCurrency is a token address (AAVE), and AttributeCurrency is a fiat symbol
	// ParameterName is in the data field as an []int

	js := `{
        "Filters": [
            {
                "Condition": {
                    "Attribute": "20.479844",
                    "Predicate": "Eq",
					"AttributeCurrency": "usd"
                },
                "EventName": "int_event",
                "FilterType": "CheckEventParameter",
                "ParameterName": "v7",
                "ParameterType": "int256[3]",
			    "ParameterCurrency": "0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9",
                "Index": 1
            }
        ],
        "ContractABI": "[{\"constant\":false,\"inputs\":[],\"name\":\"bool_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"address_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"string_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bytes_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"int_event_f\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"v1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"v2\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"string[3]\",\"name\":\"v3\",\"type\":\"string[3]\"}],\"name\":\"string_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"v4\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[3]\",\"name\":\"v5\",\"type\":\"address[3]\"}],\"name\":\"address_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"v6\",\"type\":\"int256[]\"},{\"indexed\":false,\"internalType\":\"int256[3]\",\"name\":\"v7\",\"type\":\"int256[3]\"}],\"name\":\"int_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"v8\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool[]\",\"name\":\"v9\",\"type\":\"bool[]\"},{\"indexed\":false,\"internalType\":\"bool[3]\",\"name\":\"v10\",\"type\":\"bool[3]\"}],\"name\":\"bool_event\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"v11\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v12\",\"type\":\"bytes\"}],\"name\":\"bytes_event\",\"type\":\"event\"}]",
        "ContractAdd": "0x494b4a86212fee251aa9019fe3cdb92a54d9efa1",
        "TriggerName": "WAE - TestInt256FixedArrayEqAtPosition1",
        "TriggerType": "WatchEvents"
    }`

	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	matches := MatchEvent(tg, logs551, mockTApiCurr)

	assert.Equal(t, 1, len(matches))
}

func TestCurrencyWithExplicitCurrenciesTopicInt(t *testing.T) {
	// where ParameterCurrency is a token address (AAVE), and AttributeCurrency is a fiat symbol
	// ParameterName is in the topic field as an uint256
	js := `{
	"Filters": [
		{
			"Condition": {
				"Attribute": "67.935608999999999996",
				"Predicate": "Eq",
				"AttributeCurrency": "usd"
			},
			"EventName": "LogResult",
			"FilterType": "CheckEventParameter",
			"ParameterName": "ResultSerialNumber",
			"ParameterType": "uint256",
			"ParameterCurrency": "0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9"
		}
	],
    "ContractABI": "[{\"constant\":false,\"inputs\":[{\"name\":\"newCallbackGasPrice\",\"type\":\"uint256\"}],\"name\":\"ownerSetCallbackGasPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWon\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitAsPercentOfHouse\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newHouseEdge\",\"type\":\"uint256\"}],\"name\":\"ownerSetHouseEdge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"payoutsPaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\"}],\"name\":\"ownerSetTreasury\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"myid\",\"type\":\"bytes32\"},{\"name\":\"result\",\"type\":\"string\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"__callback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressToCheck\",\"type\":\"address\"}],\"name\":\"playerGetPendingTxByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newContractBalanceInWei\",\"type\":\"uint256\"}],\"name\":\"ownerUpdateContractBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfitDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newPayoutStatus\",\"type\":\"bool\"}],\"name\":\"ownerPausePayouts\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"ownerChangeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMaxProfitAsPercent\",\"type\":\"uint256\"}],\"name\":\"ownerSetMaxProfitAsPercentOfHouse\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWeiWagered\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newMinimumBet\",\"type\":\"uint256\"}],\"name\":\"ownerSetMinBet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newStatus\",\"type\":\"bool\"}],\"name\":\"ownerPauseGame\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasForOraclize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerTransferEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contractBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minBet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"playerWithdrawPendingTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxProfit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalBets\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"randomQueryID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gamePaused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"originalPlayerBetId\",\"type\":\"bytes32\"},{\"name\":\"sendTo\",\"type\":\"address\"},{\"name\":\"originalPlayerProfit\",\"type\":\"uint256\"},{\"name\":\"originalPlayerBetValue\",\"type\":\"uint256\"}],\"name\":\"ownerRefundPlayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newSafeGasToOraclize\",\"type\":\"uint32\"}],\"name\":\"ownerSetOraclizeSafeGas\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"ownerkill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdge\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rollUnder\",\"type\":\"uint256\"}],\"name\":\"playerRollDice\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"houseEdgeDivisor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxPendingPayouts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RewardValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ProfitValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BetValue\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"RandomQueryID\",\"type\":\"uint256\"}],\"name\":\"LogBet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ResultSerialNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"PlayerNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"DiceResult\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Status\",\"type\":\"int256\"},{\"indexed\":false,\"name\":\"Proof\",\"type\":\"bytes\"}],\"name\":\"LogResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"BetID\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"PlayerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"RefundValue\",\"type\":\"uint256\"}],\"name\":\"LogRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"SentToAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"AmountTransferred\",\"type\":\"uint256\"}],\"name\":\"LogOwnerTransfer\",\"type\":\"event\"}]",
    "ContractAdd": "0xa52e014b3f5cc48287c2d483a3e026c32cc76e6d",
    "TriggerName": "WAE",
    "TriggerType": "WatchEvents"
}`
	tg, err := NewTriggerFromJson(js)
	assert.NoError(t, err)

	logs, _ := GetLogsFromFile("../resources/events/logs2.json")
	matches := MatchEvent(tg, logs, mockTApiCurr)

	assert.Equal(t, 1, len(matches))
}
