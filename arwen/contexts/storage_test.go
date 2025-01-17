package contexts

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var elrondReservedTestPrefix = []byte("RESERVED")
var epochNotifier = &worldmock.EpochNotifierStub{}

func TestNewStorageContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}
	mockBlockchain := worldmock.NewMockWorld()

	storageContext, err := NewStorageContext(host, mockBlockchain, epochNotifier, elrondReservedTestPrefix, 0)
	require.Nil(t, err)
	require.NotNil(t, storageContext)
}

func TestStorageContext_SetAddress(t *testing.T) {
	t.Parallel()

	addressA := []byte("accountA")
	addressB := []byte("accountB")
	stubOutput := &contextmock.OutputContextStub{}
	accountA := &vmcommon.OutputAccount{
		Address:        addressA,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	accountB := &vmcommon.OutputAccount{
		Address:        addressB,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	stubOutput.GetOutputAccountCalled = func(address []byte) (*vmcommon.OutputAccount, bool) {
		if bytes.Equal(address, addressA) {
			return accountA, false
		}
		if bytes.Equal(address, addressB) {
			return accountB, false
		}
		return nil, false
	}

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   stubOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, epochNotifier, elrondReservedTestPrefix, 0)

	keyA := []byte("keyA")
	valueA := []byte("valueA")

	storageContext.SetAddress(addressA)
	storageStatus, err := storageContext.SetStorage(keyA, valueA)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	foundValueA, _ := storageContext.GetStorage(keyA)
	require.Equal(t, valueA, foundValueA)
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 0)

	keyB := []byte("keyB")
	valueB := []byte("valueB")
	storageContext.SetAddress(addressB)
	storageStatus, err = storageContext.SetStorage(keyB, valueB)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	foundValueB, _ := storageContext.GetStorage(keyB)
	require.Equal(t, valueB, foundValueB)
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 1)
	foundValueA, _ = storageContext.GetStorage(keyA)
	require.Equal(t, []byte(nil), foundValueA)
}

func TestStorageContext_GetStorageUpdates(t *testing.T) {
	t.Parallel()

	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount([]byte("account"))
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	account.StorageUpdates["update"] = &vmcommon.StorageUpdate{
		Offset: []byte("update"),
		Data:   []byte("some data"),
	}

	host := &contextmock.VMHostMock{
		OutputContext: mockOutput,
	}

	mockBlockchainHook := worldmock.NewMockWorld()
	storageContext, _ := NewStorageContext(host, mockBlockchainHook, epochNotifier, elrondReservedTestPrefix, 0)

	storageUpdates := storageContext.GetStorageUpdates([]byte("account"))
	require.Equal(t, 1, len(storageUpdates))
	require.Equal(t, []byte("update"), storageUpdates["update"].Offset)
	require.Equal(t, []byte("some data"), storageUpdates["update"].Data)
}

func TestStorageContext_SetStorage(t *testing.T) {
	t.Parallel()

	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, epochNotifier, elrondReservedTestPrefix, 0)
	storageContext.SetAddress(address)

	key := []byte("key")
	value := []byte("value")

	storageStatus, err := storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	foundValue, _ := storageContext.GetStorage(key)
	require.Equal(t, value, foundValue)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageModified, storageStatus)
	foundValue, _ = storageContext.GetStorage(key)
	require.Equal(t, value, foundValue)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	foundValue, _ = storageContext.GetStorage(key)
	require.Equal(t, value, foundValue)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = nil
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageDeleted, storageStatus)
	foundValue, _ = storageContext.GetStorage(key)
	require.Equal(t, []byte{}, foundValue)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(true)
	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	foundValue, _ = storageContext.GetStorage(key)
	require.Equal(t, []byte{}, foundValue)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(false)
	key = []byte("other_key")
	value = []byte("other_value")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	foundValue, _ = storageContext.GetStorage(key)
	require.Equal(t, value, foundValue)
	require.Len(t, storageContext.GetStorageUpdates(address), 2)

	key = []byte("RESERVEDkey")
	value = []byte("doesn't matter")
	_, err = storageContext.SetStorage(key, value)
	require.Equal(t, arwen.ErrStoreElrondReservedKey, err)

	key = []byte("RESERVED")
	value = []byte("doesn't matter")
	_, err = storageContext.SetStorage(key, value)
	require.Equal(t, arwen.ErrStoreElrondReservedKey, err)
}

func TestStorageConext_SetStorage_GasUsage(t *testing.T) {
	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	storeCost := 11
	persistCost := 7
	releaseCost := 5

	gasMap := config.MakeGasMapForTests()
	gasMap["BaseOperationCost"]["StorePerByte"] = uint64(storeCost)
	gasMap["BaseOperationCost"]["PersistPerByte"] = uint64(persistCost)
	gasMap["BaseOperationCost"]["ReleasePerByte"] = uint64(releaseCost)

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(gasMap)
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, epochNotifier, elrondReservedTestPrefix, 0)
	storageContext.SetAddress(address)

	gasProvided := 100
	mockMetering.GasLeftMock = uint64(gasProvided)
	key := []byte("key")

	// Store new value
	value := []byte("value")
	storageStatus, err := storageContext.SetStorage(key, value)
	gasLeft := gasProvided - storeCost*len(value)
	storedValue, _ := storageContext.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Equal(t, gasLeft, int(mockMetering.GasLeft()))
	require.Equal(t, value, storedValue)

	// Update with longer value
	value2 := []byte("value2")
	mockMetering.GasLeftMock = uint64(gasProvided)
	storageStatus, err = storageContext.SetStorage(key, value2)
	storedValue, _ = storageContext.GetStorage(key)
	gasLeft = gasProvided - persistCost*len(value) - storeCost*(len(value2)-len(value))
	require.Nil(t, err)
	require.Equal(t, arwen.StorageModified, storageStatus)
	require.Equal(t, gasLeft, int(mockMetering.GasLeft()))
	require.Equal(t, value2, storedValue)

	// Revert to initial value
	mockMetering.GasLeftMock = uint64(gasProvided)
	storageStatus, err = storageContext.SetStorage(key, value)
	gasLeft = gasProvided - persistCost*len(value)
	gasFreed := releaseCost * (len(value2) - len(value))
	storedValue, _ = storageContext.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageModified, storageStatus)
	require.Equal(t, gasLeft, int(mockMetering.GasLeft()))
	require.Equal(t, gasFreed, int(mockMetering.GasFreedMock))
	require.Equal(t, value, storedValue)
}

func TestStorageContext_StorageProtection(t *testing.T) {
	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, epochNotifier, elrondReservedTestPrefix, 0)
	storageContext.SetAddress(address)

	key := []byte(arwen.ProtectedStoragePrefix + "something")
	value := []byte("data")

	storageStatus, err := storageContext.SetStorage(key, value)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, arwen.ErrCannotWriteProtectedKey))
	require.Len(t, storageContext.GetStorageUpdates(address), 0)

	storageContext.disableStorageProtection()
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, arwen.StorageAdded, storageStatus)
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	storageContext.enableStorageProtection()
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Equal(t, arwen.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, arwen.ErrCannotWriteProtectedKey))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)
}

func TestStorageContext_GetStorageFromAddress(t *testing.T) {
	t.Parallel()

	scAddress := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(scAddress)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &contextmock.VMHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}

	readable := []byte("readable")
	nonreadable := []byte("nonreadable")
	internalData := []byte("internalData")

	bcHook := &contextmock.BlockchainHookStub{
		GetUserAccountCalled: func(address []byte) (vmcommon.UserAccountHandler, error) {
			if bytes.Equal(readable, address) {
				return &worldmock.Account{CodeMetadata: []byte{4, 0}}, nil
			}
			if bytes.Equal(nonreadable, address) || bytes.Equal(scAddress, address) {
				return &worldmock.Account{CodeMetadata: []byte{0, 0}}, nil
			}
			return nil, nil
		},
		GetStorageDataCalled: func(accountsAddress []byte, index []byte) ([]byte, error) {
			return internalData, nil
		},
	}

	storageContext, _ := NewStorageContext(host, bcHook, epochNotifier, elrondReservedTestPrefix, 0)
	storageContext.SetAddress(scAddress)

	key := []byte("key")
	data, _ := storageContext.GetStorageFromAddress(scAddress, key)
	require.Equal(t, data, internalData)

	data, _ = storageContext.GetStorageFromAddress(readable, key)
	require.Equal(t, data, internalData)

	data, _ = storageContext.GetStorageFromAddress(nonreadable, key)
	require.Nil(t, data)
}

func TestStorageContext_LoadGasStoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_StoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	storageContext, _ := NewStorageContext(&contextmock.VMHostMock{}, &contextmock.BlockchainHookStub{}, epochNotifier, elrondReservedTestPrefix, 0)
	storageContext.PopSetActiveState()

	require.Equal(t, 0, len(storageContext.stateStack))
}

func TestStorageContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	storageContext, _ := NewStorageContext(&contextmock.VMHostMock{}, &contextmock.BlockchainHookStub{}, epochNotifier, elrondReservedTestPrefix, 0)
	storageContext.PopDiscard()

	require.Equal(t, 0, len(storageContext.stateStack))
}
