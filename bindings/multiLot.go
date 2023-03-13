// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multiLot

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MultiLotMetaData contains all meta data concerning the MultiLot contract.
var MultiLotMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lotId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ClaimWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lotId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"}],\"name\":\"Invited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lotId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"tokenA\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"tokenBChoices\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPrivate\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isChallenge\",\"type\":\"bool\"}],\"name\":\"LotCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lotId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"token\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"size\",\"type\":\"uint256\"}],\"name\":\"LotJoined\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lotId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"size\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"winningToken\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startPriceTokenA\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startPriceTokenB\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"resolvePriceTokenA\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"resolvePriceTokenB\",\"type\":\"uint256\"}],\"name\":\"LotResolved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lotId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RefundWithdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"RATIO_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_tokenA\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_tokenBChoices\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_duration\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_isPrivate\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_isChallenge\",\"type\":\"bool\"}],\"name\":\"createLot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_priceFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_feePercentage\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lotId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"_addresses\",\"type\":\"address[]\"}],\"name\":\"invite\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lotId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_token\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"joinLot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastLotId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"priceFeed\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lotId\",\"type\":\"uint256\"}],\"name\":\"resolveLot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"}],\"name\":\"setCollateralToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_priceFeed\",\"type\":\"address\"}],\"name\":\"setPriceFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lotId\",\"type\":\"uint256\"}],\"name\":\"withdrawClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isSuccessful\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isSuccessful\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lotId\",\"type\":\"uint256\"}],\"name\":\"withdrawRefund\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isSuccessful\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// MultiLotABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiLotMetaData.ABI instead.
var MultiLotABI = MultiLotMetaData.ABI

// MultiLot is an auto generated Go binding around an Ethereum contract.
type MultiLot struct {
	MultiLotCaller     // Read-only binding to the contract
	MultiLotTransactor // Write-only binding to the contract
	MultiLotFilterer   // Log filterer for contract events
}

// MultiLotCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiLotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiLotTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiLotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiLotFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiLotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiLotSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiLotSession struct {
	Contract     *MultiLot         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultiLotCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiLotCallerSession struct {
	Contract *MultiLotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MultiLotTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiLotTransactorSession struct {
	Contract     *MultiLotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MultiLotRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiLotRaw struct {
	Contract *MultiLot // Generic contract binding to access the raw methods on
}

// MultiLotCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiLotCallerRaw struct {
	Contract *MultiLotCaller // Generic read-only contract binding to access the raw methods on
}

// MultiLotTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiLotTransactorRaw struct {
	Contract *MultiLotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiLot creates a new instance of MultiLot, bound to a specific deployed contract.
func NewMultiLot(address common.Address, backend bind.ContractBackend) (*MultiLot, error) {
	contract, err := bindMultiLot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiLot{MultiLotCaller: MultiLotCaller{contract: contract}, MultiLotTransactor: MultiLotTransactor{contract: contract}, MultiLotFilterer: MultiLotFilterer{contract: contract}}, nil
}

// NewMultiLotCaller creates a new read-only instance of MultiLot, bound to a specific deployed contract.
func NewMultiLotCaller(address common.Address, caller bind.ContractCaller) (*MultiLotCaller, error) {
	contract, err := bindMultiLot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiLotCaller{contract: contract}, nil
}

// NewMultiLotTransactor creates a new write-only instance of MultiLot, bound to a specific deployed contract.
func NewMultiLotTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiLotTransactor, error) {
	contract, err := bindMultiLot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiLotTransactor{contract: contract}, nil
}

// NewMultiLotFilterer creates a new log filterer instance of MultiLot, bound to a specific deployed contract.
func NewMultiLotFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiLotFilterer, error) {
	contract, err := bindMultiLot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiLotFilterer{contract: contract}, nil
}

// bindMultiLot binds a generic wrapper to an already deployed contract.
func bindMultiLot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiLotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiLot *MultiLotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiLot.Contract.MultiLotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiLot *MultiLotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiLot.Contract.MultiLotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiLot *MultiLotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiLot.Contract.MultiLotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiLot *MultiLotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiLot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiLot *MultiLotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiLot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiLot *MultiLotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiLot.Contract.contract.Transact(opts, method, params...)
}

// RATIOPRECISION is a free data retrieval call binding the contract method 0xd32867d0.
//
// Solidity: function RATIO_PRECISION() view returns(uint256)
func (_MultiLot *MultiLotCaller) RATIOPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "RATIO_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RATIOPRECISION is a free data retrieval call binding the contract method 0xd32867d0.
//
// Solidity: function RATIO_PRECISION() view returns(uint256)
func (_MultiLot *MultiLotSession) RATIOPRECISION() (*big.Int, error) {
	return _MultiLot.Contract.RATIOPRECISION(&_MultiLot.CallOpts)
}

// RATIOPRECISION is a free data retrieval call binding the contract method 0xd32867d0.
//
// Solidity: function RATIO_PRECISION() view returns(uint256)
func (_MultiLot *MultiLotCallerSession) RATIOPRECISION() (*big.Int, error) {
	return _MultiLot.Contract.RATIOPRECISION(&_MultiLot.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_MultiLot *MultiLotCaller) CollateralToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "collateralToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_MultiLot *MultiLotSession) CollateralToken() (common.Address, error) {
	return _MultiLot.Contract.CollateralToken(&_MultiLot.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_MultiLot *MultiLotCallerSession) CollateralToken() (common.Address, error) {
	return _MultiLot.Contract.CollateralToken(&_MultiLot.CallOpts)
}

// FeePercentage is a free data retrieval call binding the contract method 0xa001ecdd.
//
// Solidity: function feePercentage() view returns(uint256)
func (_MultiLot *MultiLotCaller) FeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "feePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeePercentage is a free data retrieval call binding the contract method 0xa001ecdd.
//
// Solidity: function feePercentage() view returns(uint256)
func (_MultiLot *MultiLotSession) FeePercentage() (*big.Int, error) {
	return _MultiLot.Contract.FeePercentage(&_MultiLot.CallOpts)
}

// FeePercentage is a free data retrieval call binding the contract method 0xa001ecdd.
//
// Solidity: function feePercentage() view returns(uint256)
func (_MultiLot *MultiLotCallerSession) FeePercentage() (*big.Int, error) {
	return _MultiLot.Contract.FeePercentage(&_MultiLot.CallOpts)
}

// LastLotId is a free data retrieval call binding the contract method 0x6a52bd45.
//
// Solidity: function lastLotId() view returns(uint256)
func (_MultiLot *MultiLotCaller) LastLotId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "lastLotId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastLotId is a free data retrieval call binding the contract method 0x6a52bd45.
//
// Solidity: function lastLotId() view returns(uint256)
func (_MultiLot *MultiLotSession) LastLotId() (*big.Int, error) {
	return _MultiLot.Contract.LastLotId(&_MultiLot.CallOpts)
}

// LastLotId is a free data retrieval call binding the contract method 0x6a52bd45.
//
// Solidity: function lastLotId() view returns(uint256)
func (_MultiLot *MultiLotCallerSession) LastLotId() (*big.Int, error) {
	return _MultiLot.Contract.LastLotId(&_MultiLot.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultiLot *MultiLotCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultiLot *MultiLotSession) Owner() (common.Address, error) {
	return _MultiLot.Contract.Owner(&_MultiLot.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultiLot *MultiLotCallerSession) Owner() (common.Address, error) {
	return _MultiLot.Contract.Owner(&_MultiLot.CallOpts)
}

// PriceFeed is a free data retrieval call binding the contract method 0x741bef1a.
//
// Solidity: function priceFeed() view returns(address)
func (_MultiLot *MultiLotCaller) PriceFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "priceFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PriceFeed is a free data retrieval call binding the contract method 0x741bef1a.
//
// Solidity: function priceFeed() view returns(address)
func (_MultiLot *MultiLotSession) PriceFeed() (common.Address, error) {
	return _MultiLot.Contract.PriceFeed(&_MultiLot.CallOpts)
}

// PriceFeed is a free data retrieval call binding the contract method 0x741bef1a.
//
// Solidity: function priceFeed() view returns(address)
func (_MultiLot *MultiLotCallerSession) PriceFeed() (common.Address, error) {
	return _MultiLot.Contract.PriceFeed(&_MultiLot.CallOpts)
}

// TotalFee is a free data retrieval call binding the contract method 0x1df4ccfc.
//
// Solidity: function totalFee() view returns(uint256)
func (_MultiLot *MultiLotCaller) TotalFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiLot.contract.Call(opts, &out, "totalFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalFee is a free data retrieval call binding the contract method 0x1df4ccfc.
//
// Solidity: function totalFee() view returns(uint256)
func (_MultiLot *MultiLotSession) TotalFee() (*big.Int, error) {
	return _MultiLot.Contract.TotalFee(&_MultiLot.CallOpts)
}

// TotalFee is a free data retrieval call binding the contract method 0x1df4ccfc.
//
// Solidity: function totalFee() view returns(uint256)
func (_MultiLot *MultiLotCallerSession) TotalFee() (*big.Int, error) {
	return _MultiLot.Contract.TotalFee(&_MultiLot.CallOpts)
}

// CreateLot is a paid mutator transaction binding the contract method 0x28335bbb.
//
// Solidity: function createLot(string _tokenA, string[] _tokenBChoices, uint256 _size, uint256 _startEpoch, uint256 _duration, bool _isPrivate, bool _isChallenge) returns()
func (_MultiLot *MultiLotTransactor) CreateLot(opts *bind.TransactOpts, _tokenA string, _tokenBChoices []string, _size *big.Int, _startEpoch *big.Int, _duration *big.Int, _isPrivate bool, _isChallenge bool) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "createLot", _tokenA, _tokenBChoices, _size, _startEpoch, _duration, _isPrivate, _isChallenge)
}

// CreateLot is a paid mutator transaction binding the contract method 0x28335bbb.
//
// Solidity: function createLot(string _tokenA, string[] _tokenBChoices, uint256 _size, uint256 _startEpoch, uint256 _duration, bool _isPrivate, bool _isChallenge) returns()
func (_MultiLot *MultiLotSession) CreateLot(_tokenA string, _tokenBChoices []string, _size *big.Int, _startEpoch *big.Int, _duration *big.Int, _isPrivate bool, _isChallenge bool) (*types.Transaction, error) {
	return _MultiLot.Contract.CreateLot(&_MultiLot.TransactOpts, _tokenA, _tokenBChoices, _size, _startEpoch, _duration, _isPrivate, _isChallenge)
}

// CreateLot is a paid mutator transaction binding the contract method 0x28335bbb.
//
// Solidity: function createLot(string _tokenA, string[] _tokenBChoices, uint256 _size, uint256 _startEpoch, uint256 _duration, bool _isPrivate, bool _isChallenge) returns()
func (_MultiLot *MultiLotTransactorSession) CreateLot(_tokenA string, _tokenBChoices []string, _size *big.Int, _startEpoch *big.Int, _duration *big.Int, _isPrivate bool, _isChallenge bool) (*types.Transaction, error) {
	return _MultiLot.Contract.CreateLot(&_MultiLot.TransactOpts, _tokenA, _tokenBChoices, _size, _startEpoch, _duration, _isPrivate, _isChallenge)
}

// Initialize is a paid mutator transaction binding the contract method 0x1794bb3c.
//
// Solidity: function initialize(address _priceFeed, address _collateralToken, uint256 _feePercentage) returns()
func (_MultiLot *MultiLotTransactor) Initialize(opts *bind.TransactOpts, _priceFeed common.Address, _collateralToken common.Address, _feePercentage *big.Int) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "initialize", _priceFeed, _collateralToken, _feePercentage)
}

// Initialize is a paid mutator transaction binding the contract method 0x1794bb3c.
//
// Solidity: function initialize(address _priceFeed, address _collateralToken, uint256 _feePercentage) returns()
func (_MultiLot *MultiLotSession) Initialize(_priceFeed common.Address, _collateralToken common.Address, _feePercentage *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.Initialize(&_MultiLot.TransactOpts, _priceFeed, _collateralToken, _feePercentage)
}

// Initialize is a paid mutator transaction binding the contract method 0x1794bb3c.
//
// Solidity: function initialize(address _priceFeed, address _collateralToken, uint256 _feePercentage) returns()
func (_MultiLot *MultiLotTransactorSession) Initialize(_priceFeed common.Address, _collateralToken common.Address, _feePercentage *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.Initialize(&_MultiLot.TransactOpts, _priceFeed, _collateralToken, _feePercentage)
}

// Invite is a paid mutator transaction binding the contract method 0x86243394.
//
// Solidity: function invite(uint256 _lotId, address[] _addresses) returns()
func (_MultiLot *MultiLotTransactor) Invite(opts *bind.TransactOpts, _lotId *big.Int, _addresses []common.Address) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "invite", _lotId, _addresses)
}

// Invite is a paid mutator transaction binding the contract method 0x86243394.
//
// Solidity: function invite(uint256 _lotId, address[] _addresses) returns()
func (_MultiLot *MultiLotSession) Invite(_lotId *big.Int, _addresses []common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.Invite(&_MultiLot.TransactOpts, _lotId, _addresses)
}

// Invite is a paid mutator transaction binding the contract method 0x86243394.
//
// Solidity: function invite(uint256 _lotId, address[] _addresses) returns()
func (_MultiLot *MultiLotTransactorSession) Invite(_lotId *big.Int, _addresses []common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.Invite(&_MultiLot.TransactOpts, _lotId, _addresses)
}

// JoinLot is a paid mutator transaction binding the contract method 0xd0a1291f.
//
// Solidity: function joinLot(uint256 _lotId, string _token, uint256 _size) returns()
func (_MultiLot *MultiLotTransactor) JoinLot(opts *bind.TransactOpts, _lotId *big.Int, _token string, _size *big.Int) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "joinLot", _lotId, _token, _size)
}

// JoinLot is a paid mutator transaction binding the contract method 0xd0a1291f.
//
// Solidity: function joinLot(uint256 _lotId, string _token, uint256 _size) returns()
func (_MultiLot *MultiLotSession) JoinLot(_lotId *big.Int, _token string, _size *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.JoinLot(&_MultiLot.TransactOpts, _lotId, _token, _size)
}

// JoinLot is a paid mutator transaction binding the contract method 0xd0a1291f.
//
// Solidity: function joinLot(uint256 _lotId, string _token, uint256 _size) returns()
func (_MultiLot *MultiLotTransactorSession) JoinLot(_lotId *big.Int, _token string, _size *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.JoinLot(&_MultiLot.TransactOpts, _lotId, _token, _size)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MultiLot *MultiLotTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MultiLot *MultiLotSession) RenounceOwnership() (*types.Transaction, error) {
	return _MultiLot.Contract.RenounceOwnership(&_MultiLot.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MultiLot *MultiLotTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MultiLot.Contract.RenounceOwnership(&_MultiLot.TransactOpts)
}

// ResolveLot is a paid mutator transaction binding the contract method 0xe30b5287.
//
// Solidity: function resolveLot(uint256 _lotId) returns()
func (_MultiLot *MultiLotTransactor) ResolveLot(opts *bind.TransactOpts, _lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "resolveLot", _lotId)
}

// ResolveLot is a paid mutator transaction binding the contract method 0xe30b5287.
//
// Solidity: function resolveLot(uint256 _lotId) returns()
func (_MultiLot *MultiLotSession) ResolveLot(_lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.ResolveLot(&_MultiLot.TransactOpts, _lotId)
}

// ResolveLot is a paid mutator transaction binding the contract method 0xe30b5287.
//
// Solidity: function resolveLot(uint256 _lotId) returns()
func (_MultiLot *MultiLotTransactorSession) ResolveLot(_lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.ResolveLot(&_MultiLot.TransactOpts, _lotId)
}

// SetCollateralToken is a paid mutator transaction binding the contract method 0x666181a9.
//
// Solidity: function setCollateralToken(address _collateralToken) returns()
func (_MultiLot *MultiLotTransactor) SetCollateralToken(opts *bind.TransactOpts, _collateralToken common.Address) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "setCollateralToken", _collateralToken)
}

// SetCollateralToken is a paid mutator transaction binding the contract method 0x666181a9.
//
// Solidity: function setCollateralToken(address _collateralToken) returns()
func (_MultiLot *MultiLotSession) SetCollateralToken(_collateralToken common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.SetCollateralToken(&_MultiLot.TransactOpts, _collateralToken)
}

// SetCollateralToken is a paid mutator transaction binding the contract method 0x666181a9.
//
// Solidity: function setCollateralToken(address _collateralToken) returns()
func (_MultiLot *MultiLotTransactorSession) SetCollateralToken(_collateralToken common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.SetCollateralToken(&_MultiLot.TransactOpts, _collateralToken)
}

// SetPriceFeed is a paid mutator transaction binding the contract method 0x724e78da.
//
// Solidity: function setPriceFeed(address _priceFeed) returns()
func (_MultiLot *MultiLotTransactor) SetPriceFeed(opts *bind.TransactOpts, _priceFeed common.Address) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "setPriceFeed", _priceFeed)
}

// SetPriceFeed is a paid mutator transaction binding the contract method 0x724e78da.
//
// Solidity: function setPriceFeed(address _priceFeed) returns()
func (_MultiLot *MultiLotSession) SetPriceFeed(_priceFeed common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.SetPriceFeed(&_MultiLot.TransactOpts, _priceFeed)
}

// SetPriceFeed is a paid mutator transaction binding the contract method 0x724e78da.
//
// Solidity: function setPriceFeed(address _priceFeed) returns()
func (_MultiLot *MultiLotTransactorSession) SetPriceFeed(_priceFeed common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.SetPriceFeed(&_MultiLot.TransactOpts, _priceFeed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MultiLot *MultiLotTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MultiLot *MultiLotSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.TransferOwnership(&_MultiLot.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MultiLot *MultiLotTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultiLot.Contract.TransferOwnership(&_MultiLot.TransactOpts, newOwner)
}

// WithdrawClaim is a paid mutator transaction binding the contract method 0xa224dcb7.
//
// Solidity: function withdrawClaim(uint256 _lotId) returns(bool isSuccessful)
func (_MultiLot *MultiLotTransactor) WithdrawClaim(opts *bind.TransactOpts, _lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "withdrawClaim", _lotId)
}

// WithdrawClaim is a paid mutator transaction binding the contract method 0xa224dcb7.
//
// Solidity: function withdrawClaim(uint256 _lotId) returns(bool isSuccessful)
func (_MultiLot *MultiLotSession) WithdrawClaim(_lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.WithdrawClaim(&_MultiLot.TransactOpts, _lotId)
}

// WithdrawClaim is a paid mutator transaction binding the contract method 0xa224dcb7.
//
// Solidity: function withdrawClaim(uint256 _lotId) returns(bool isSuccessful)
func (_MultiLot *MultiLotTransactorSession) WithdrawClaim(_lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.WithdrawClaim(&_MultiLot.TransactOpts, _lotId)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns(bool isSuccessful)
func (_MultiLot *MultiLotTransactor) WithdrawFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "withdrawFee")
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns(bool isSuccessful)
func (_MultiLot *MultiLotSession) WithdrawFee() (*types.Transaction, error) {
	return _MultiLot.Contract.WithdrawFee(&_MultiLot.TransactOpts)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns(bool isSuccessful)
func (_MultiLot *MultiLotTransactorSession) WithdrawFee() (*types.Transaction, error) {
	return _MultiLot.Contract.WithdrawFee(&_MultiLot.TransactOpts)
}

// WithdrawRefund is a paid mutator transaction binding the contract method 0x9d153495.
//
// Solidity: function withdrawRefund(uint256 _lotId) returns(bool isSuccessful)
func (_MultiLot *MultiLotTransactor) WithdrawRefund(opts *bind.TransactOpts, _lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.contract.Transact(opts, "withdrawRefund", _lotId)
}

// WithdrawRefund is a paid mutator transaction binding the contract method 0x9d153495.
//
// Solidity: function withdrawRefund(uint256 _lotId) returns(bool isSuccessful)
func (_MultiLot *MultiLotSession) WithdrawRefund(_lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.WithdrawRefund(&_MultiLot.TransactOpts, _lotId)
}

// WithdrawRefund is a paid mutator transaction binding the contract method 0x9d153495.
//
// Solidity: function withdrawRefund(uint256 _lotId) returns(bool isSuccessful)
func (_MultiLot *MultiLotTransactorSession) WithdrawRefund(_lotId *big.Int) (*types.Transaction, error) {
	return _MultiLot.Contract.WithdrawRefund(&_MultiLot.TransactOpts, _lotId)
}

// MultiLotClaimWithdrawnIterator is returned from FilterClaimWithdrawn and is used to iterate over the raw logs and unpacked data for ClaimWithdrawn events raised by the MultiLot contract.
type MultiLotClaimWithdrawnIterator struct {
	Event *MultiLotClaimWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotClaimWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotClaimWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotClaimWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotClaimWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotClaimWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotClaimWithdrawn represents a ClaimWithdrawn event raised by the MultiLot contract.
type MultiLotClaimWithdrawn struct {
	LotId  *big.Int
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClaimWithdrawn is a free log retrieval operation binding the contract event 0x98131a469190deca66117f3768bba4328c631613211a3cc4054efc6ee16cd426.
//
// Solidity: event ClaimWithdrawn(uint256 lotId, address user, uint256 amount)
func (_MultiLot *MultiLotFilterer) FilterClaimWithdrawn(opts *bind.FilterOpts) (*MultiLotClaimWithdrawnIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "ClaimWithdrawn")
	if err != nil {
		return nil, err
	}
	return &MultiLotClaimWithdrawnIterator{contract: _MultiLot.contract, event: "ClaimWithdrawn", logs: logs, sub: sub}, nil
}

// WatchClaimWithdrawn is a free log subscription operation binding the contract event 0x98131a469190deca66117f3768bba4328c631613211a3cc4054efc6ee16cd426.
//
// Solidity: event ClaimWithdrawn(uint256 lotId, address user, uint256 amount)
func (_MultiLot *MultiLotFilterer) WatchClaimWithdrawn(opts *bind.WatchOpts, sink chan<- *MultiLotClaimWithdrawn) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "ClaimWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotClaimWithdrawn)
				if err := _MultiLot.contract.UnpackLog(event, "ClaimWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimWithdrawn is a log parse operation binding the contract event 0x98131a469190deca66117f3768bba4328c631613211a3cc4054efc6ee16cd426.
//
// Solidity: event ClaimWithdrawn(uint256 lotId, address user, uint256 amount)
func (_MultiLot *MultiLotFilterer) ParseClaimWithdrawn(log types.Log) (*MultiLotClaimWithdrawn, error) {
	event := new(MultiLotClaimWithdrawn)
	if err := _MultiLot.contract.UnpackLog(event, "ClaimWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotFeeWithdrawnIterator is returned from FilterFeeWithdrawn and is used to iterate over the raw logs and unpacked data for FeeWithdrawn events raised by the MultiLot contract.
type MultiLotFeeWithdrawnIterator struct {
	Event *MultiLotFeeWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotFeeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotFeeWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotFeeWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotFeeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotFeeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotFeeWithdrawn represents a FeeWithdrawn event raised by the MultiLot contract.
type MultiLotFeeWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeWithdrawn is a free log retrieval operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_MultiLot *MultiLotFilterer) FilterFeeWithdrawn(opts *bind.FilterOpts) (*MultiLotFeeWithdrawnIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return &MultiLotFeeWithdrawnIterator{contract: _MultiLot.contract, event: "FeeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFeeWithdrawn is a free log subscription operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_MultiLot *MultiLotFilterer) WatchFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *MultiLotFeeWithdrawn) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotFeeWithdrawn)
				if err := _MultiLot.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeWithdrawn is a log parse operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_MultiLot *MultiLotFilterer) ParseFeeWithdrawn(log types.Log) (*MultiLotFeeWithdrawn, error) {
	event := new(MultiLotFeeWithdrawn)
	if err := _MultiLot.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MultiLot contract.
type MultiLotInitializedIterator struct {
	Event *MultiLotInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotInitialized represents a Initialized event raised by the MultiLot contract.
type MultiLotInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MultiLot *MultiLotFilterer) FilterInitialized(opts *bind.FilterOpts) (*MultiLotInitializedIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MultiLotInitializedIterator{contract: _MultiLot.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MultiLot *MultiLotFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MultiLotInitialized) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotInitialized)
				if err := _MultiLot.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MultiLot *MultiLotFilterer) ParseInitialized(log types.Log) (*MultiLotInitialized, error) {
	event := new(MultiLotInitialized)
	if err := _MultiLot.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotInvitedIterator is returned from FilterInvited and is used to iterate over the raw logs and unpacked data for Invited events raised by the MultiLot contract.
type MultiLotInvitedIterator struct {
	Event *MultiLotInvited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotInvitedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotInvited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotInvited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotInvitedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotInvitedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotInvited represents a Invited event raised by the MultiLot contract.
type MultiLotInvited struct {
	LotId *big.Int
	From  common.Address
	To    []common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterInvited is a free log retrieval operation binding the contract event 0x47e145642f40b14fcfba9a6525311864334f5a1e2675ef92ba4f3fe3b8e46198.
//
// Solidity: event Invited(uint256 lotId, address from, address[] to)
func (_MultiLot *MultiLotFilterer) FilterInvited(opts *bind.FilterOpts) (*MultiLotInvitedIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "Invited")
	if err != nil {
		return nil, err
	}
	return &MultiLotInvitedIterator{contract: _MultiLot.contract, event: "Invited", logs: logs, sub: sub}, nil
}

// WatchInvited is a free log subscription operation binding the contract event 0x47e145642f40b14fcfba9a6525311864334f5a1e2675ef92ba4f3fe3b8e46198.
//
// Solidity: event Invited(uint256 lotId, address from, address[] to)
func (_MultiLot *MultiLotFilterer) WatchInvited(opts *bind.WatchOpts, sink chan<- *MultiLotInvited) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "Invited")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotInvited)
				if err := _MultiLot.contract.UnpackLog(event, "Invited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInvited is a log parse operation binding the contract event 0x47e145642f40b14fcfba9a6525311864334f5a1e2675ef92ba4f3fe3b8e46198.
//
// Solidity: event Invited(uint256 lotId, address from, address[] to)
func (_MultiLot *MultiLotFilterer) ParseInvited(log types.Log) (*MultiLotInvited, error) {
	event := new(MultiLotInvited)
	if err := _MultiLot.contract.UnpackLog(event, "Invited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotLotCreatedIterator is returned from FilterLotCreated and is used to iterate over the raw logs and unpacked data for LotCreated events raised by the MultiLot contract.
type MultiLotLotCreatedIterator struct {
	Event *MultiLotLotCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotLotCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotLotCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotLotCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotLotCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotLotCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotLotCreated represents a LotCreated event raised by the MultiLot contract.
type MultiLotLotCreated struct {
	LotId         *big.Int
	TokenA        string
	TokenBChoices []string
	StartEpoch    *big.Int
	Duration      *big.Int
	Creator       common.Address
	IsPrivate     bool
	IsChallenge   bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLotCreated is a free log retrieval operation binding the contract event 0x7aa0e2508522c6e152ffa3164ff180b4d3a50b99c2d3a682a4a3176d99bcf5f5.
//
// Solidity: event LotCreated(uint256 lotId, string tokenA, string[] tokenBChoices, uint256 startEpoch, uint256 duration, address creator, bool isPrivate, bool isChallenge)
func (_MultiLot *MultiLotFilterer) FilterLotCreated(opts *bind.FilterOpts) (*MultiLotLotCreatedIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "LotCreated")
	if err != nil {
		return nil, err
	}
	return &MultiLotLotCreatedIterator{contract: _MultiLot.contract, event: "LotCreated", logs: logs, sub: sub}, nil
}

// WatchLotCreated is a free log subscription operation binding the contract event 0x7aa0e2508522c6e152ffa3164ff180b4d3a50b99c2d3a682a4a3176d99bcf5f5.
//
// Solidity: event LotCreated(uint256 lotId, string tokenA, string[] tokenBChoices, uint256 startEpoch, uint256 duration, address creator, bool isPrivate, bool isChallenge)
func (_MultiLot *MultiLotFilterer) WatchLotCreated(opts *bind.WatchOpts, sink chan<- *MultiLotLotCreated) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "LotCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotLotCreated)
				if err := _MultiLot.contract.UnpackLog(event, "LotCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLotCreated is a log parse operation binding the contract event 0x7aa0e2508522c6e152ffa3164ff180b4d3a50b99c2d3a682a4a3176d99bcf5f5.
//
// Solidity: event LotCreated(uint256 lotId, string tokenA, string[] tokenBChoices, uint256 startEpoch, uint256 duration, address creator, bool isPrivate, bool isChallenge)
func (_MultiLot *MultiLotFilterer) ParseLotCreated(log types.Log) (*MultiLotLotCreated, error) {
	event := new(MultiLotLotCreated)
	if err := _MultiLot.contract.UnpackLog(event, "LotCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotLotJoinedIterator is returned from FilterLotJoined and is used to iterate over the raw logs and unpacked data for LotJoined events raised by the MultiLot contract.
type MultiLotLotJoinedIterator struct {
	Event *MultiLotLotJoined // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotLotJoinedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotLotJoined)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotLotJoined)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotLotJoinedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotLotJoinedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotLotJoined represents a LotJoined event raised by the MultiLot contract.
type MultiLotLotJoined struct {
	LotId *big.Int
	Token string
	User  common.Address
	Size  *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLotJoined is a free log retrieval operation binding the contract event 0x2ef9c5ad696244f79a77d314d32cd61df59551c3eb7e42679ea2f2898784065c.
//
// Solidity: event LotJoined(uint256 lotId, string token, address user, uint256 size)
func (_MultiLot *MultiLotFilterer) FilterLotJoined(opts *bind.FilterOpts) (*MultiLotLotJoinedIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "LotJoined")
	if err != nil {
		return nil, err
	}
	return &MultiLotLotJoinedIterator{contract: _MultiLot.contract, event: "LotJoined", logs: logs, sub: sub}, nil
}

// WatchLotJoined is a free log subscription operation binding the contract event 0x2ef9c5ad696244f79a77d314d32cd61df59551c3eb7e42679ea2f2898784065c.
//
// Solidity: event LotJoined(uint256 lotId, string token, address user, uint256 size)
func (_MultiLot *MultiLotFilterer) WatchLotJoined(opts *bind.WatchOpts, sink chan<- *MultiLotLotJoined) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "LotJoined")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotLotJoined)
				if err := _MultiLot.contract.UnpackLog(event, "LotJoined", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLotJoined is a log parse operation binding the contract event 0x2ef9c5ad696244f79a77d314d32cd61df59551c3eb7e42679ea2f2898784065c.
//
// Solidity: event LotJoined(uint256 lotId, string token, address user, uint256 size)
func (_MultiLot *MultiLotFilterer) ParseLotJoined(log types.Log) (*MultiLotLotJoined, error) {
	event := new(MultiLotLotJoined)
	if err := _MultiLot.contract.UnpackLog(event, "LotJoined", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotLotResolvedIterator is returned from FilterLotResolved and is used to iterate over the raw logs and unpacked data for LotResolved events raised by the MultiLot contract.
type MultiLotLotResolvedIterator struct {
	Event *MultiLotLotResolved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotLotResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotLotResolved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotLotResolved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotLotResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotLotResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotLotResolved represents a LotResolved event raised by the MultiLot contract.
type MultiLotLotResolved struct {
	LotId              *big.Int
	Size               *big.Int
	WinningToken       string
	StartPriceTokenA   *big.Int
	StartPriceTokenB   *big.Int
	ResolvePriceTokenA *big.Int
	ResolvePriceTokenB *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterLotResolved is a free log retrieval operation binding the contract event 0xcb4fe1f5bd157ec52ea58b76dd52ac4ce2d6a5498e3a3fcea3b972c329964735.
//
// Solidity: event LotResolved(uint256 lotId, uint256 size, string winningToken, uint256 startPriceTokenA, uint256 startPriceTokenB, uint256 resolvePriceTokenA, uint256 resolvePriceTokenB)
func (_MultiLot *MultiLotFilterer) FilterLotResolved(opts *bind.FilterOpts) (*MultiLotLotResolvedIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "LotResolved")
	if err != nil {
		return nil, err
	}
	return &MultiLotLotResolvedIterator{contract: _MultiLot.contract, event: "LotResolved", logs: logs, sub: sub}, nil
}

// WatchLotResolved is a free log subscription operation binding the contract event 0xcb4fe1f5bd157ec52ea58b76dd52ac4ce2d6a5498e3a3fcea3b972c329964735.
//
// Solidity: event LotResolved(uint256 lotId, uint256 size, string winningToken, uint256 startPriceTokenA, uint256 startPriceTokenB, uint256 resolvePriceTokenA, uint256 resolvePriceTokenB)
func (_MultiLot *MultiLotFilterer) WatchLotResolved(opts *bind.WatchOpts, sink chan<- *MultiLotLotResolved) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "LotResolved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotLotResolved)
				if err := _MultiLot.contract.UnpackLog(event, "LotResolved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLotResolved is a log parse operation binding the contract event 0xcb4fe1f5bd157ec52ea58b76dd52ac4ce2d6a5498e3a3fcea3b972c329964735.
//
// Solidity: event LotResolved(uint256 lotId, uint256 size, string winningToken, uint256 startPriceTokenA, uint256 startPriceTokenB, uint256 resolvePriceTokenA, uint256 resolvePriceTokenB)
func (_MultiLot *MultiLotFilterer) ParseLotResolved(log types.Log) (*MultiLotLotResolved, error) {
	event := new(MultiLotLotResolved)
	if err := _MultiLot.contract.UnpackLog(event, "LotResolved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MultiLot contract.
type MultiLotOwnershipTransferredIterator struct {
	Event *MultiLotOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotOwnershipTransferred represents a OwnershipTransferred event raised by the MultiLot contract.
type MultiLotOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MultiLot *MultiLotFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MultiLotOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MultiLotOwnershipTransferredIterator{contract: _MultiLot.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MultiLot *MultiLotFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiLotOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotOwnershipTransferred)
				if err := _MultiLot.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MultiLot *MultiLotFilterer) ParseOwnershipTransferred(log types.Log) (*MultiLotOwnershipTransferred, error) {
	event := new(MultiLotOwnershipTransferred)
	if err := _MultiLot.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiLotRefundWithdrawnIterator is returned from FilterRefundWithdrawn and is used to iterate over the raw logs and unpacked data for RefundWithdrawn events raised by the MultiLot contract.
type MultiLotRefundWithdrawnIterator struct {
	Event *MultiLotRefundWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiLotRefundWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiLotRefundWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiLotRefundWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiLotRefundWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiLotRefundWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiLotRefundWithdrawn represents a RefundWithdrawn event raised by the MultiLot contract.
type MultiLotRefundWithdrawn struct {
	LotId  *big.Int
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRefundWithdrawn is a free log retrieval operation binding the contract event 0x6909eb935886ad8c734c29844350c36b0260f7006ff58559a3c286a9e7c8d878.
//
// Solidity: event RefundWithdrawn(uint256 lotId, address user, uint256 amount)
func (_MultiLot *MultiLotFilterer) FilterRefundWithdrawn(opts *bind.FilterOpts) (*MultiLotRefundWithdrawnIterator, error) {

	logs, sub, err := _MultiLot.contract.FilterLogs(opts, "RefundWithdrawn")
	if err != nil {
		return nil, err
	}
	return &MultiLotRefundWithdrawnIterator{contract: _MultiLot.contract, event: "RefundWithdrawn", logs: logs, sub: sub}, nil
}

// WatchRefundWithdrawn is a free log subscription operation binding the contract event 0x6909eb935886ad8c734c29844350c36b0260f7006ff58559a3c286a9e7c8d878.
//
// Solidity: event RefundWithdrawn(uint256 lotId, address user, uint256 amount)
func (_MultiLot *MultiLotFilterer) WatchRefundWithdrawn(opts *bind.WatchOpts, sink chan<- *MultiLotRefundWithdrawn) (event.Subscription, error) {

	logs, sub, err := _MultiLot.contract.WatchLogs(opts, "RefundWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiLotRefundWithdrawn)
				if err := _MultiLot.contract.UnpackLog(event, "RefundWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefundWithdrawn is a log parse operation binding the contract event 0x6909eb935886ad8c734c29844350c36b0260f7006ff58559a3c286a9e7c8d878.
//
// Solidity: event RefundWithdrawn(uint256 lotId, address user, uint256 amount)
func (_MultiLot *MultiLotFilterer) ParseRefundWithdrawn(log types.Log) (*MultiLotRefundWithdrawn, error) {
	event := new(MultiLotRefundWithdrawn)
	if err := _MultiLot.contract.UnpackLog(event, "RefundWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
