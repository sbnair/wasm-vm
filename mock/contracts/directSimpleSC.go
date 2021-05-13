package contracts

import (
	"errors"
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1.3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1.3/testcommon"
	"github.com/stretchr/testify/require"
)

// WasteGasChildMock is an exposed mock contract method
func WasteGasChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("wasteGas", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByChild))
}

// FailChildMock is an exposed mock contract method
func FailChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("fail", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Runtime().FailExecution(errors.New("forced fail"))
		return instance
	})
}

// ExecOnSameCtxParentMock is an exposed mock contract method
func ExecOnSameCtxParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execOnSameCtx", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.GasUsedByParent)

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				_, err := host.ExecuteOnSameContext(input)
				require.Nil(t, err)
			}
		}

		return instance
	})
}

// ExecOnDestCtxParentMock is an exposed mock contract method
func ExecOnDestCtxParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execOnDestCtx", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.GasUsedByParent)

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				_, _, err := host.ExecuteOnDestContext(input)
				require.Nil(t, err)
			}
		}

		return instance
	})
}

// ExecOnDestCtxSingleCallParentMock is an exposed mock contract method
func ExecOnDestCtxSingleCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execOnDestCtxSingleCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 2 {
			host.Runtime().SignalUserError("need 2 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		_, _, err := host.ExecuteOnDestContext(input)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

// WasteGasParentMock is an exposed mock contract method
func WasteGasParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("wasteGas", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByParent))
}
