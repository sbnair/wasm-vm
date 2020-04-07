package host

import (
	"encoding/hex"
	"fmt"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var parentKeyA = []byte("parentKeyA......................")
var parentKeyB = []byte("parentKeyB......................")
var childKey = []byte("childKey........................")
var parentDataA = []byte("parentDataA")
var parentDataB = []byte("parentDataB")
var childData = []byte("childData")
var parentFinishA = []byte("parentFinishA")
var parentFinishB = []byte("parentFinishB")
var childFinish = []byte("childFinish")
var parentTransferReceiver = []byte("parentTransferReceiver..........")
var childTransferReceiver = []byte("childTransferReceiver...........")
var parentTransferValue = int64(42)
var parentTransferData = []byte("parentTransferData")

var recursiveIterationCounterKey = []byte("recursiveIterationCounter.......")
var recursiveIterationBigCounterKey = []byte("recursiveIterationBigCounter....")

var gasProvided = uint64(1000000)

var parentCompilationCost_SameCtx = uint64(3578)
var childCompilationCost_SameCtx = uint64(3269)

var parentCompilationCost_DestCtx = uint64(3268)
var childCompilationCost_DestCtx = uint64(1810)

var vaultAddress = []byte("vaultAddress....................")
var thirdPartyAddress = []byte("thirdPartyAddress...............")

func expectedVMOutput_SameCtx_Prepare() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)
	AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(135)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_WrongContractCalled() *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare()

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(180)
	executeAPICost := uint64(39)
	gasLostOnFailure := uint64(50000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_OutOfGas() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	AddFinishData(vmOutput, parentFinishA)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(124)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(36)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_SuccessfulChildCall() *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare()

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		3,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		96,
		[]byte("qwerty"),
	)

	// The child SC has stored data on the parent's storage
	SetStorageUpdate(parentAccount, childKey, childData)

	// The called child SC will output some arbitrary data, but also data that it
	// has read from the storage of the parent.
	AddFinishData(vmOutput, childFinish)
	AddFinishData(vmOutput, parentDataA)
	for _, c := range parentDataA {
		AddFinishData(vmOutput, []byte{c})
	}
	AddFinishData(vmOutput, parentDataB)
	for _, c := range parentDataB {
		AddFinishData(vmOutput, []byte{c})
	}
	AddFinishData(vmOutput, []byte("child ok"))
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(188)
	executeAPICost := uint64(39)
	childExecutionCost := uint64(431)
	finalCost := uint64(139)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_SameCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_SuccessfulChildCall_BigInts() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	// parentAccount.BalanceDelta = big.NewInt(-99)

	_ = AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)

	// The child SC will output "child ok" if it could read some expected Big
	// Ints directly from the parent's context.
	AddFinishData(vmOutput, []byte("child ok"))
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(143)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(107)
	finalCost := uint64(67)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_SameCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_Direct(recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		SetStorageUpdateStrings(account, key, value)
	}

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_MutualMethods(recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	AddFinishData(vmOutput, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		SetStorageUpdateStrings(account, key, value)
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_MutualSCs(recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	if recursiveCalls%2 == 1 {
		parentAccount.BalanceDelta = big.NewInt(-5)
		childAccount.BalanceDelta = big.NewInt(5)
	} else {
		parentAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
		childAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	}

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
		}
		SetStorageUpdateStrings(parentAccount, key, value)
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutput_DestCtx_Prepare() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)
	AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(135)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_DestCtx_WrongContractCalled() *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare()

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-42)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(180)
	executeAPICost := uint64(42)
	gasLostOnFailure := uint64(10000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_DestCtx_OutOfGas() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	AddFinishData(vmOutput, parentFinishA)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(124)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(36)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_DestCtx_SuccessfulChildCall() *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare()

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		99-12,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		12,
		[]byte("Second sentence."),
	)

	SetStorageUpdate(childAccount, childKey, childData)

	AddFinishData(vmOutput, childFinish)
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(188)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	finalCost := uint64(65)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_DestCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutput_DestCtx_SuccessfulChildCall_BigInts() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)

	// The child SC will output "child ok" if it could NOT read the Big Ints from
	// the parent's context.
	AddFinishData(vmOutput, []byte("child ok"))
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(143)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(101)
	finalCost := uint64(68)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_DestCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_Direct(recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		SetStorageUpdateStrings(account, key, value)
	}

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_MutualMethods(recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	AddFinishData(vmOutput, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		SetStorageUpdateStrings(account, key, value)
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_MutualSCs(recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentIterations := (recursiveCalls / 2) + (recursiveCalls % 2)
	childIterations := recursiveCalls - parentIterations

	balanceDelta := int64(5*parentIterations - 3*childIterations)

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.BalanceDelta = big.NewInt(-balanceDelta)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	childAccount.BalanceDelta = big.NewInt(balanceDelta)

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
			SetStorageUpdateStrings(parentAccount, key, value)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
			SetStorageUpdateStrings(childAccount, key, value)
		}
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	counterValue := (recursiveCalls + recursiveCalls%2) / 2
	SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(counterValue + 1)})
	SetStorageUpdate(childAccount, recursiveIterationCounterKey, []byte{byte(counterValue)})
	if recursiveCalls%2 == 0 {
		SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	} else {
		SetStorageUpdate(childAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	}

	return vmOutput
}

func expectedVMOutput_AsyncCall() *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)

	_ = AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		6,
		[]byte("hello there"),
	)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	SetStorageUpdate(childAccount, childKey, childData)

	_ = AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		nil,
	)

	AddFinishData(vmOutput, []byte("thirdparty"))
	AddFinishData(vmOutput, []byte("vault"))

	return vmOutput
}

func createDataBytes(returnCode vmcommon.ReturnCode, returnData [][]byte) []byte {
	data := []byte("@" + hex.EncodeToString([]byte(returnCode.String())))

	for _, returnDataPiece := range returnData {
		returnDataPieceHex := []byte("@" + hex.EncodeToString(returnDataPiece))
		data = append(data, returnDataPieceHex...)
	}

	return data
}
