package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void getOwner(void *context, int32_t resultOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transferValue(void *context, long long gasLimit, int32_t dstOffset, int32_t sndOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t dataOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern int32_t callValue(void *context, int32_t resultOffset);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void returnData(void* context, int32_t dataOffset, int32_t length);
// extern void signalError(void* context);
// extern long long getGasLeft(void *context);
//
// extern long long getBlockTimestamp(void *context);
// extern long long getBlockNonce(void *context);
// extern long long getBlockRound(void *context);
// extern long long getBlockEpoch(void *context);
// extern void getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getStateRootHash(void *context, int32_t resultOffset);
//
// extern long long getPrevBlockTimestamp(void *context);
// extern long long getPrevBlockNonce(void *context);
// extern long long getPrevBlockRound(void *context);
// extern long long getPrevBlockEpoch(void *context);
// extern void getPrevBlockRandomSeed(void *context, int32_t resultOffset);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset);
// extern void int64finish(void* context, long long value);
import "C"

import (
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/debugging"

	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

func ElrondEImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()
	imports = imports.Namespace("env")

	imports, err := imports.Append("getOwner", getOwner, C.getOwner)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", getExternalBalance, C.getExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", blockHash, C.blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValue", transferValue, C.transferValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgument", getArgument, C.getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getFunction", getFunction, C.getFunction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", getNumArguments, C.getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", storageStore, C.storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", storageLoad, C.storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", getCaller, C.getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", callValue, C.callValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", returnData, C.returnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNonce", getBlockNonce, C.getBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRound", getBlockRound, C.getBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockEpoch", getBlockEpoch, C.getBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRandomSeed", getBlockRandomSeed, C.getBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStateRootHash", getStateRootHash, C.getStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockTimestamp", getPrevBlockTimestamp, C.getPrevBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockNonce", getPrevBlockNonce, C.getPrevBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRound", getPrevBlockRound, C.getPrevBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockEpoch", getPrevBlockEpoch, C.getPrevBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRandomSeed", getPrevBlockRandomSeed, C.getPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64getArgument", int64getArgument, C.int64getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageStore", int64storageStore, C.int64storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageLoad", int64storageLoad, C.int64storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64finish", int64finish, C.int64finish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}


//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	debugging.TraceCall("getGasLeft")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetGasLeft
	hostContext.UseGas(gasToUse)

	debugging.TraceReturnUint64(hostContext.GasLeft())
	return int64(hostContext.GasLeft())
}

//export getOwner
func getOwner(context unsafe.Pointer, resultOffset int32) {
	debugging.TraceCall("getOwner")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	owner := hostContext.GetSCAddress()
	debugging.TraceVarBytes("owner", owner)
	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, owner)
	
	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetOwner
	hostContext.UseGas(gasToUse)
}

//export signalError
func signalError(context unsafe.Pointer) {
	debugging.TraceCall("signalError")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hostContext.SignalUserError()

	gasToUse := hostContext.GasSchedule().ElrondAPICost.SignalError
	hostContext.UseGas(gasToUse)
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	debugging.TraceCall("getExternalBalance")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	balance := hostContext.GetBalance(address)
	debugging.TraceVarBytes("address", address)
	debugging.TraceVarBigIntBytes("balance", balance)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, balance)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetExternalBalance
	hostContext.UseGas(gasToUse)
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	debugging.TraceCall("blockHash")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockHash
	hostContext.UseGas(gasToUse)

	//TODO: change blockchain hook to treat actual nonce - not the offset.
	hash := hostContext.BlockHash(nonce)
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, hash)
	debugging.TraceErr("StoreBytes", err)
	if err != nil {
		return 1
	}

	debugging.TraceVarBytes("hash", hash)

	return 0
}

//export transferValue
func transferValue(context unsafe.Pointer, gasLimit int64, sndOffset int32, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	debugging.TraceCall("transferValue")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	send := arwen.LoadBytes(instCtx.Memory(), sndOffset, arwen.AddressLen)
	dest := arwen.LoadBytes(instCtx.Memory(), destOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.TransferValue
	gasToUse += hostContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
	hostContext.UseGas(gasToUse)
	debugging.TraceVarUint64("gasToUse", gasToUse)

	_, err := hostContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), data, gasLimit)
	debugging.TraceErr("Transfer", err)
	if err != nil {	
		return 1
	}

	return 0
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	debugging.TraceCall("getArgument")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetArgument
	hostContext.UseGas(gasToUse)

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		debugging.TraceErrMessage("invalid argument id")
		return -1
	}

	debugging.TraceVarBytes("argAsBytes", args[id].Bytes())
	debugging.TraceVarBigIntBytes("argAsBigInt", args[id].Bytes())
	err := arwen.StoreBytes(instCtx.Memory(), argOffset, args[id].Bytes())
	debugging.TraceErr("StoreBytes", err)
	if err != nil {
		return -1
	}

	debugging.TraceReturnInt32(int32(len(args[id].Bytes())))
	return int32(len(args[id].Bytes()))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	debugging.TraceCall("getFunction")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetFunction
	hostContext.UseGas(gasToUse)

	function := hostContext.Function()
	debugging.TraceVarString("function", function)
	err := arwen.StoreBytes(instCtx.Memory(), functionOffset, []byte(function))
	debugging.TraceErr("StoreBytes", err)
	if err != nil {
		return -1
	}

	result := int32(len(function))
	debugging.TraceReturnInt32(result)
	return result
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	debugging.TraceCall("getNumArguments")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetNumArguments
	hostContext.UseGas(gasToUse)

	result := int32(len(hostContext.Arguments()))
	debugging.TraceReturnInt32(result)
	return result
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	debugging.TraceCall("storageStore")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)
	debugging.TraceVarBytes("key", key)
	debugging.TraceVarBytes("data", data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.StorageStore
	gasToUse += hostContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(dataLength)
	hostContext.UseGas(gasToUse)
	debugging.TraceVarUint64("gasToUse", gasToUse)

	result := hostContext.SetStorage(hostContext.GetSCAddress(), key, data)
	debugging.TraceReturnInt32(result)
	return result
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, dataOffset int32) int32 {
	debugging.TraceCall("storageLoad")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)
	debugging.TraceVarBytes("key", key)
	debugging.TraceVarBytes("data", data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.StorageLoad
	gasToUse += hostContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	hostContext.UseGas(gasToUse)
	debugging.TraceVarUint64("gasToUse", gasToUse)

	err := arwen.StoreBytes(instCtx.Memory(), dataOffset, data)
	debugging.TraceErr("StoreBytes", err)
	if err != nil {
		return -1
	}

	result := int32(len(data))
	debugging.TraceReturnInt32(result)
	return result
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	debugging.TraceCall("getCaller")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	caller := hostContext.GetVMInput().CallerAddr
	debugging.TraceVarBytes("caller", caller)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetCaller
	hostContext.UseGas(gasToUse)
}

//export callValue
func callValue(context unsafe.Pointer, resultOffset int32) int32 {
	debugging.TraceCall("callValue")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	value := hostContext.GetVMInput().CallValue.Bytes()
	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	debugging.TraceVarBigIntBytes("value", value)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetCallValue
	hostContext.UseGas(gasToUse)

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, invBytes)
	debugging.TraceErr("StoreBytes", err)
	if err != nil {
		return -1
	}

	result := int32(length)
	debugging.TraceReturnInt32(result)
	return result
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	debugging.TraceCall("writeLog")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	log := arwen.LoadBytes(instCtx.Memory(), pointer, length)

	topics := make([][]byte, numTopics)
	for i := int32(0); i < numTopics; i++ {
		topics[i] = arwen.LoadBytes(instCtx.Memory(), topicPtr+i*arwen.HashLen, arwen.HashLen)
		debugging.TraceVarBytes("topic[]", topics[i])
	}

	hostContext.WriteLog(hostContext.GetSCAddress(), topics, log)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Log
	gasToUse += hostContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(numTopics*arwen.HashLen+length)
	hostContext.UseGas(gasToUse)
	debugging.TraceVarUint64("gasToUse", gasToUse)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	debugging.TraceCall("getBlockTimestamp")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().CurrentTimeStamp())
	debugging.TraceReturnInt64(result)
	return result
}

//export getBlockNonce
func getBlockNonce(context unsafe.Pointer) int64 {
	debugging.TraceCall("getBlockNonce")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockNonce
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().CurrentNonce())
	debugging.TraceReturnInt64(result)
	return result
}

//export getBlockRound
func getBlockRound(context unsafe.Pointer) int64 {
	debugging.TraceCall("getBlockRound")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRound
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().CurrentRound())
	debugging.TraceReturnInt64(result)
	return result
}

//export getBlockEpoch
func getBlockEpoch(context unsafe.Pointer) int64 {
	debugging.TraceCall("getBlockEpoch")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockEpoch
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().CurrentEpoch())
	debugging.TraceReturnInt64(result)
	return result
}

//export getBlockRandomSeed
func getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	debugging.TraceCall("getBlockRandomSeed")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	hostContext.UseGas(gasToUse)

	randomSeed := hostContext.BlockChainHook().CurrentRandomSeed()
	debugging.TraceVarBytes("randomSeed", randomSeed)
	_ = arwen.StoreBytes(instCtx.Memory(), pointer, randomSeed)
}

//export getStateRootHash
func getStateRootHash(context unsafe.Pointer, pointer int32) {
	debugging.TraceCall("getStateRootHash")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetStateRootHash
	hostContext.UseGas(gasToUse)

	stateRootHash := hostContext.BlockChainHook().GetStateRootHash()
	debugging.TraceVarBytes("randomSeed", stateRootHash)
	_ = arwen.StoreBytes(instCtx.Memory(), pointer, stateRootHash)
}

//export getPrevBlockTimestamp
func getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	debugging.TraceCall("getPrevBlockTimestamp")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().LastTimeStamp())
	debugging.TraceReturnInt64(result)
	return result
}

//export getPrevBlockNonce
func getPrevBlockNonce(context unsafe.Pointer) int64 {
	debugging.TraceCall("getPrevBlockNonce")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockNonce
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().LastNonce())
	debugging.TraceReturnInt64(result)
	return result
}

//export getPrevBlockRound
func getPrevBlockRound(context unsafe.Pointer) int64 {
	debugging.TraceCall("getPrevBlockRound")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRound
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().LastRound())
	debugging.TraceReturnInt64(result)
	return result
}

//export getPrevBlockEpoch
func getPrevBlockEpoch(context unsafe.Pointer) int64 {
	debugging.TraceCall("getPrevBlockEpoch")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockEpoch
	hostContext.UseGas(gasToUse)

	result := int64(hostContext.BlockChainHook().LastEpoch())
	debugging.TraceReturnInt64(result)
	return result
}

//export getPrevBlockRandomSeed
func getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	debugging.TraceCall("getPrevBlockRandomSeed")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	hostContext.UseGas(gasToUse)

	randomSeed := hostContext.BlockChainHook().LastRandomSeed()
	debugging.TraceVarBytes("randomSeed", randomSeed)
	_ = arwen.StoreBytes(instCtx.Memory(), pointer, randomSeed)
}

//export returnData
func returnData(context unsafe.Pointer, pointer int32, length int32) {
	debugging.TraceCall("returnData")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), pointer, length)
	hostContext.Finish(data)
	debugging.TraceVarBytes("data", data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Finish
	gasToUse += hostContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
	hostContext.UseGas(gasToUse)
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	debugging.TraceCall("int64getArgument")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64GetArgument
	hostContext.UseGas(gasToUse)

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	result := args[id].Int64()
	debugging.TraceReturnInt64(result)
	return result
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	debugging.TraceCall("int64storageStore")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := big.NewInt(value)
	debugging.TraceVarBigInt("data", data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64StorageStore
	hostContext.UseGas(gasToUse)

	result := hostContext.SetStorage(hostContext.GetSCAddress(), key, data.Bytes())
	debugging.TraceReturnInt32(result)
	return result
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32) int64 {
	debugging.TraceCall("int64storageLoad")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64StorageLoad
	hostContext.UseGas(gasToUse)

	result := bigInt.Int64()
	debugging.TraceReturnInt64(result)
	return result
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	debugging.TraceCall("int64finish")

	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hostContext.Finish(big.NewInt(0).SetInt64(value).Bytes())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64Finish
	hostContext.UseGas(gasToUse)
}