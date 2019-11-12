package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func CreateGasConfig(gasMap map[string]uint64) (*GasCost, error) {
	baseOps := &BaseOperationCost{}
	err := mapstructure.Decode(gasMap, baseOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(*baseOps)
	if err != nil {
		return nil, err
	}

	elrondOps := &ElrondAPICost{}
	err = mapstructure.Decode(gasMap, elrondOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(*elrondOps)
	if err != nil {
		return nil, err
	}

	bigIntOps := &BigIntAPICost{}
	err = mapstructure.Decode(gasMap, bigIntOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(*bigIntOps)
	if err != nil {
		return nil, err
	}

	ethOps := &EthAPICost{}
	err = mapstructure.Decode(gasMap, ethOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(*ethOps)
	if err != nil {
		return nil, err
	}

	cryptOps := &CryptoAPICost{}
	err = mapstructure.Decode(gasMap, cryptOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(*cryptOps)
	if err != nil {
		return nil, err
	}

	opcodeCosts := &WASMOpcodeCost{}
	err = mapstructure.Decode(gasMap, opcodeCosts)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(*opcodeCosts)
	if err != nil {
		return nil, err
	}

	gasCost := &GasCost{
		BaseOperationCost: *baseOps,
		BigIntAPICost:     *bigIntOps,
		EthAPICost:        *ethOps,
		ElrondAPICost:     *elrondOps,
		CryptoAPICost:     *cryptOps,
		WASMOpcodeCost:    *opcodeCosts,
	}

	return gasCost, nil
}

func checkForZeroUint64Fields(arg interface{}) error {
	v := reflect.ValueOf(arg)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() != reflect.Uint64 && field.Kind() != reflect.Uint32 {
			continue
		}
		if field.Uint() == 0 {
			name := v.Type().Field(i).Name
			return errors.New(fmt.Sprintf("field %s has the value 0", name))
		}
	}

	return nil
}

func MakeGasMap(value uint64) map[string]uint64 {
	gasMap := make(map[string]uint64)
	gasMap = FillGasMap(gasMap, value)
	return gasMap
}

func FillGasMap(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap = FillGasMap_BaseOperationCosts(gasMap, value)
	gasMap = FillGasMap_ElrondAPICosts(gasMap, value)
	gasMap = FillGasMap_EthereumAPICosts(gasMap, value)
	gasMap = FillGasMap_BigIntAPICosts(gasMap, value)
	gasMap = FillGasMap_CryptoAPICosts(gasMap, value)
	gasMap = FillGasMap_WASMOpcodeValues(gasMap, value)

	return gasMap
}

func FillGasMap_BaseOperationCosts(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap["StorePerByte"] = value
	gasMap["DataCopyPerByte"] = value

	return gasMap
}

func FillGasMap_ElrondAPICosts(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap["GetOwner"] = value
	gasMap["GetExternalBalance"] = value
	gasMap["GetBlockHash"] = value
	gasMap["TransferValue"] = value
	gasMap["GetArgument"] = value
	gasMap["GetFunction"] = value
	gasMap["GetNumArguments"] = value
	gasMap["StorageStore"] = value
	gasMap["StorageLoad"] = value
	gasMap["GetCaller"] = value
	gasMap["GetCallValue"] = value
	gasMap["Log"] = value
	gasMap["Finish"] = value
	gasMap["SignalError"] = value
	gasMap["GetBlockTimeStamp"] = value
	gasMap["GetGasLeft"] = value
	gasMap["Int64GetArgument"] = value
	gasMap["Int64StorageStore"] = value
	gasMap["Int64StorageLoad"] = value
	gasMap["Int64Finish"] = value
	gasMap["GetStateRootHash"] = value
	gasMap["GetBlockNonce"] = value
	gasMap["GetBlockEpoch"] = value
	gasMap["GetBlockRound"] = value
	gasMap["GetBlockRandomSeed"] = value

	return gasMap
}

func FillGasMap_EthereumAPICosts(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap["UseGas"] = value
	gasMap["GetAddress"] = value
	gasMap["GetExternalBalance"] = value
	gasMap["GetBlockHash"] = value
	gasMap["Call"] = value
	gasMap["CallDataCopy"] = value
	gasMap["GetCallDataSize"] = value
	gasMap["CallCode"] = value
	gasMap["CallDelegate"] = value
	gasMap["CallStatic"] = value
	gasMap["StorageStore"] = value
	gasMap["StorageLoad"] = value
	gasMap["GetCaller"] = value
	gasMap["GetCallValue"] = value
	gasMap["CodeCopy"] = value
	gasMap["GetCodeSize"] = value
	gasMap["GetBlockCoinbase"] = value
	gasMap["Create"] = value
	gasMap["GetBlockDifficulty"] = value
	gasMap["ExternalCodeCopy"] = value
	gasMap["GetExternalCodeSize"] = value
	gasMap["GetGasLeft"] = value
	gasMap["GetBlockGasLimit"] = value
	gasMap["GetTxGasPrice"] = value
	gasMap["Log"] = value
	gasMap["GetBlockNumber"] = value
	gasMap["GetTxOrigin"] = value
	gasMap["Finish"] = value
	gasMap["Revert"] = value
	gasMap["GetReturnDataSize"] = value
	gasMap["ReturnDataCopy"] = value
	gasMap["SelfDestruct"] = value
	gasMap["GetBlockTimeStamp"] = value

	return gasMap
}

func FillGasMap_BigIntAPICosts(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap["BigIntNew"] = value
	gasMap["BigIntByteLength"] = value
	gasMap["BigIntGetBytes"] = value
	gasMap["BigIntSetBytes"] = value
	gasMap["BigIntIsInt64"] = value
	gasMap["BigIntGetInt64"] = value
	gasMap["BigIntSetInt64"] = value
	gasMap["BigIntAdd"] = value
	gasMap["BigIntSub"] = value
	gasMap["BigIntMul"] = value
	gasMap["BigIntCmp"] = value
	gasMap["BigIntFinish"] = value
	gasMap["BigIntStorageLoad"] = value
	gasMap["BigIntStorageStore"] = value
	gasMap["BigIntGetArgument"] = value
	gasMap["BigIntGetCallValue"] = value
	gasMap["BigIntGetExternalBalance"] = value

	return gasMap
}

func FillGasMap_CryptoAPICosts(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap["SHA256"] = value
	gasMap["Keccak256"] = value

	return gasMap
}

func FillGasMap_WASMOpcodeValues(gasMap map[string]uint64, value uint64) map[string]uint64 {
	gasMap["Unreachable"] = value
	gasMap["Nop"] = value
	gasMap["Block"] = value
	gasMap["Loop"] = value
	gasMap["If"] = value
	gasMap["Else"] = value
	gasMap["End"] = value
	gasMap["Br"] = value
	gasMap["BrIf"] = value
	gasMap["BrTable"] = value
	gasMap["Return"] = value
	gasMap["Call"] = value
	gasMap["CallIndirect"] = value
	gasMap["Drop"] = value
	gasMap["Select"] = value
	gasMap["GetLocal"] = value
	gasMap["SetLocal"] = value
	gasMap["TeeLocal"] = value
	gasMap["GetGlobal"] = value
	gasMap["SetGlobal"] = value
	gasMap["I32Load"] = value
	gasMap["I64Load"] = value
	gasMap["F32Load"] = value
	gasMap["F64Load"] = value
	gasMap["I32Load8S"] = value
	gasMap["I32Load8U"] = value
	gasMap["I32Load16S"] = value
	gasMap["I32Load16U"] = value
	gasMap["I64Load8S"] = value
	gasMap["I64Load8U"] = value
	gasMap["I64Load16S"] = value
	gasMap["I64Load16U"] = value
	gasMap["I64Load32S"] = value
	gasMap["I64Load32U"] = value
	gasMap["I32Store"] = value
	gasMap["I64Store"] = value
	gasMap["F32Store"] = value
	gasMap["F64Store"] = value
	gasMap["I32Store8"] = value
	gasMap["I32Store16"] = value
	gasMap["I64Store8"] = value
	gasMap["I64Store16"] = value
	gasMap["I64Store32"] = value
	gasMap["MemorySize"] = value
	gasMap["MemoryGrow"] = value
	gasMap["I32Const"] = value
	gasMap["I64Const"] = value
	gasMap["F32Const"] = value
	gasMap["F64Const"] = value
	gasMap["RefNull"] = value
	gasMap["RefIsNull"] = value
	gasMap["I32Eqz"] = value
	gasMap["I32Eq"] = value
	gasMap["I32Ne"] = value
	gasMap["I32LtS"] = value
	gasMap["I32LtU"] = value
	gasMap["I32GtS"] = value
	gasMap["I32GtU"] = value
	gasMap["I32LeS"] = value
	gasMap["I32LeU"] = value
	gasMap["I32GeS"] = value
	gasMap["I32GeU"] = value
	gasMap["I64Eqz"] = value
	gasMap["I64Eq"] = value
	gasMap["I64Ne"] = value
	gasMap["I64LtS"] = value
	gasMap["I64LtU"] = value
	gasMap["I64GtS"] = value
	gasMap["I64GtU"] = value
	gasMap["I64LeS"] = value
	gasMap["I64LeU"] = value
	gasMap["I64GeS"] = value
	gasMap["I64GeU"] = value
	gasMap["F32Eq"] = value
	gasMap["F32Ne"] = value
	gasMap["F32Lt"] = value
	gasMap["F32Gt"] = value
	gasMap["F32Le"] = value
	gasMap["F32Ge"] = value
	gasMap["F64Eq"] = value
	gasMap["F64Ne"] = value
	gasMap["F64Lt"] = value
	gasMap["F64Gt"] = value
	gasMap["F64Le"] = value
	gasMap["F64Ge"] = value
	gasMap["I32Clz"] = value
	gasMap["I32Ctz"] = value
	gasMap["I32Popcnt"] = value
	gasMap["I32Add"] = value
	gasMap["I32Sub"] = value
	gasMap["I32Mul"] = value
	gasMap["I32DivS"] = value
	gasMap["I32DivU"] = value
	gasMap["I32RemS"] = value
	gasMap["I32RemU"] = value
	gasMap["I32And"] = value
	gasMap["I32Or"] = value
	gasMap["I32Xor"] = value
	gasMap["I32Shl"] = value
	gasMap["I32ShrS"] = value
	gasMap["I32ShrU"] = value
	gasMap["I32Rotl"] = value
	gasMap["I32Rotr"] = value
	gasMap["I64Clz"] = value
	gasMap["I64Ctz"] = value
	gasMap["I64Popcnt"] = value
	gasMap["I64Add"] = value
	gasMap["I64Sub"] = value
	gasMap["I64Mul"] = value
	gasMap["I64DivS"] = value
	gasMap["I64DivU"] = value
	gasMap["I64RemS"] = value
	gasMap["I64RemU"] = value
	gasMap["I64And"] = value
	gasMap["I64Or"] = value
	gasMap["I64Xor"] = value
	gasMap["I64Shl"] = value
	gasMap["I64ShrS"] = value
	gasMap["I64ShrU"] = value
	gasMap["I64Rotl"] = value
	gasMap["I64Rotr"] = value
	gasMap["F32Abs"] = value
	gasMap["F32Neg"] = value
	gasMap["F32Ceil"] = value
	gasMap["F32Floor"] = value
	gasMap["F32Trunc"] = value
	gasMap["F32Nearest"] = value
	gasMap["F32Sqrt"] = value
	gasMap["F32Add"] = value
	gasMap["F32Sub"] = value
	gasMap["F32Mul"] = value
	gasMap["F32Div"] = value
	gasMap["F32Min"] = value
	gasMap["F32Max"] = value
	gasMap["F32Copysign"] = value
	gasMap["F64Abs"] = value
	gasMap["F64Neg"] = value
	gasMap["F64Ceil"] = value
	gasMap["F64Floor"] = value
	gasMap["F64Trunc"] = value
	gasMap["F64Nearest"] = value
	gasMap["F64Sqrt"] = value
	gasMap["F64Add"] = value
	gasMap["F64Sub"] = value
	gasMap["F64Mul"] = value
	gasMap["F64Div"] = value
	gasMap["F64Min"] = value
	gasMap["F64Max"] = value
	gasMap["F64Copysign"] = value
	gasMap["I32WrapI64"] = value
	gasMap["I32TruncSF32"] = value
	gasMap["I32TruncUF32"] = value
	gasMap["I32TruncSF64"] = value
	gasMap["I32TruncUF64"] = value
	gasMap["I64ExtendSI32"] = value
	gasMap["I64ExtendUI32"] = value
	gasMap["I64TruncSF32"] = value
	gasMap["I64TruncUF32"] = value
	gasMap["I64TruncSF64"] = value
	gasMap["I64TruncUF64"] = value
	gasMap["F32ConvertSI32"] = value
	gasMap["F32ConvertUI32"] = value
	gasMap["F32ConvertSI64"] = value
	gasMap["F32ConvertUI64"] = value
	gasMap["F32DemoteF64"] = value
	gasMap["F64ConvertSI32"] = value
	gasMap["F64ConvertUI32"] = value
	gasMap["F64ConvertSI64"] = value
	gasMap["F64ConvertUI64"] = value
	gasMap["F64PromoteF32"] = value
	gasMap["I32ReinterpretF32"] = value
	gasMap["I64ReinterpretF64"] = value
	gasMap["F32ReinterpretI32"] = value
	gasMap["F64ReinterpretI64"] = value
	gasMap["I32Extend8S"] = value
	gasMap["I32Extend16S"] = value
	gasMap["I64Extend8S"] = value
	gasMap["I64Extend16S"] = value
	gasMap["I64Extend32S"] = value
	gasMap["I32TruncSSatF32"] = value
	gasMap["I32TruncUSatF32"] = value
	gasMap["I32TruncSSatF64"] = value
	gasMap["I32TruncUSatF64"] = value
	gasMap["I64TruncSSatF32"] = value
	gasMap["I64TruncUSatF32"] = value
	gasMap["I64TruncSSatF64"] = value
	gasMap["I64TruncUSatF64"] = value
	gasMap["MemoryInit"] = value
	gasMap["DataDrop"] = value
	gasMap["MemoryCopy"] = value
	gasMap["MemoryFill"] = value
	gasMap["TableInit"] = value
	gasMap["ElemDrop"] = value
	gasMap["TableCopy"] = value
	gasMap["TableGet"] = value
	gasMap["TableSet"] = value
	gasMap["TableGrow"] = value
	gasMap["TableSize"] = value
	gasMap["Wake"] = value
	gasMap["I32Wait"] = value
	gasMap["I64Wait"] = value
	gasMap["Fence"] = value
	gasMap["I32AtomicLoad"] = value
	gasMap["I64AtomicLoad"] = value
	gasMap["I32AtomicLoad8U"] = value
	gasMap["I32AtomicLoad16U"] = value
	gasMap["I64AtomicLoad8U"] = value
	gasMap["I64AtomicLoad16U"] = value
	gasMap["I64AtomicLoad32U"] = value
	gasMap["I32AtomicStore"] = value
	gasMap["I64AtomicStore"] = value
	gasMap["I32AtomicStore8"] = value
	gasMap["I32AtomicStore16"] = value
	gasMap["I64AtomicStore8"] = value
	gasMap["I64AtomicStore16"] = value
	gasMap["I64AtomicStore32"] = value
	gasMap["I32AtomicRmwAdd"] = value
	gasMap["I64AtomicRmwAdd"] = value
	gasMap["I32AtomicRmw8UAdd"] = value
	gasMap["I32AtomicRmw16UAdd"] = value
	gasMap["I64AtomicRmw8UAdd"] = value
	gasMap["I64AtomicRmw16UAdd"] = value
	gasMap["I64AtomicRmw32UAdd"] = value
	gasMap["I32AtomicRmwSub"] = value
	gasMap["I64AtomicRmwSub"] = value
	gasMap["I32AtomicRmw8USub"] = value
	gasMap["I32AtomicRmw16USub"] = value
	gasMap["I64AtomicRmw8USub"] = value
	gasMap["I64AtomicRmw16USub"] = value
	gasMap["I64AtomicRmw32USub"] = value
	gasMap["I32AtomicRmwAnd"] = value
	gasMap["I64AtomicRmwAnd"] = value
	gasMap["I32AtomicRmw8UAnd"] = value
	gasMap["I32AtomicRmw16UAnd"] = value
	gasMap["I64AtomicRmw8UAnd"] = value
	gasMap["I64AtomicRmw16UAnd"] = value
	gasMap["I64AtomicRmw32UAnd"] = value
	gasMap["I32AtomicRmwOr"] = value
	gasMap["I64AtomicRmwOr"] = value
	gasMap["I32AtomicRmw8UOr"] = value
	gasMap["I32AtomicRmw16UOr"] = value
	gasMap["I64AtomicRmw8UOr"] = value
	gasMap["I64AtomicRmw16UOr"] = value
	gasMap["I64AtomicRmw32UOr"] = value
	gasMap["I32AtomicRmwXor"] = value
	gasMap["I64AtomicRmwXor"] = value
	gasMap["I32AtomicRmw8UXor"] = value
	gasMap["I32AtomicRmw16UXor"] = value
	gasMap["I64AtomicRmw8UXor"] = value
	gasMap["I64AtomicRmw16UXor"] = value
	gasMap["I64AtomicRmw32UXor"] = value
	gasMap["I32AtomicRmwXchg"] = value
	gasMap["I64AtomicRmwXchg"] = value
	gasMap["I32AtomicRmw8UXchg"] = value
	gasMap["I32AtomicRmw16UXchg"] = value
	gasMap["I64AtomicRmw8UXchg"] = value
	gasMap["I64AtomicRmw16UXchg"] = value
	gasMap["I64AtomicRmw32UXchg"] = value
	gasMap["I32AtomicRmwCmpxchg"] = value
	gasMap["I64AtomicRmwCmpxchg"] = value
	gasMap["I32AtomicRmw8UCmpxchg"] = value
	gasMap["I32AtomicRmw16UCmpxchg"] = value
	gasMap["I64AtomicRmw8UCmpxchg"] = value
	gasMap["I64AtomicRmw16UCmpxchg"] = value
	gasMap["I64AtomicRmw32UCmpxchg"] = value
	gasMap["V128Load"] = value
	gasMap["V128Store"] = value
	gasMap["V128Const"] = value
	gasMap["I8x16Splat"] = value
	gasMap["I8x16ExtractLaneS"] = value
	gasMap["I8x16ExtractLaneU"] = value
	gasMap["I8x16ReplaceLane"] = value
	gasMap["I16x8Splat"] = value
	gasMap["I16x8ExtractLaneS"] = value
	gasMap["I16x8ExtractLaneU"] = value
	gasMap["I16x8ReplaceLane"] = value
	gasMap["I32x4Splat"] = value
	gasMap["I32x4ExtractLane"] = value
	gasMap["I32x4ReplaceLane"] = value
	gasMap["I64x2Splat"] = value
	gasMap["I64x2ExtractLane"] = value
	gasMap["I64x2ReplaceLane"] = value
	gasMap["F32x4Splat"] = value
	gasMap["F32x4ExtractLane"] = value
	gasMap["F32x4ReplaceLane"] = value
	gasMap["F64x2Splat"] = value
	gasMap["F64x2ExtractLane"] = value
	gasMap["F64x2ReplaceLane"] = value
	gasMap["I8x16Eq"] = value
	gasMap["I8x16Ne"] = value
	gasMap["I8x16LtS"] = value
	gasMap["I8x16LtU"] = value
	gasMap["I8x16GtS"] = value
	gasMap["I8x16GtU"] = value
	gasMap["I8x16LeS"] = value
	gasMap["I8x16LeU"] = value
	gasMap["I8x16GeS"] = value
	gasMap["I8x16GeU"] = value
	gasMap["I16x8Eq"] = value
	gasMap["I16x8Ne"] = value
	gasMap["I16x8LtS"] = value
	gasMap["I16x8LtU"] = value
	gasMap["I16x8GtS"] = value
	gasMap["I16x8GtU"] = value
	gasMap["I16x8LeS"] = value
	gasMap["I16x8LeU"] = value
	gasMap["I16x8GeS"] = value
	gasMap["I16x8GeU"] = value
	gasMap["I32x4Eq"] = value
	gasMap["I32x4Ne"] = value
	gasMap["I32x4LtS"] = value
	gasMap["I32x4LtU"] = value
	gasMap["I32x4GtS"] = value
	gasMap["I32x4GtU"] = value
	gasMap["I32x4LeS"] = value
	gasMap["I32x4LeU"] = value
	gasMap["I32x4GeS"] = value
	gasMap["I32x4GeU"] = value
	gasMap["F32x4Eq"] = value
	gasMap["F32x4Ne"] = value
	gasMap["F32x4Lt"] = value
	gasMap["F32x4Gt"] = value
	gasMap["F32x4Le"] = value
	gasMap["F32x4Ge"] = value
	gasMap["F64x2Eq"] = value
	gasMap["F64x2Ne"] = value
	gasMap["F64x2Lt"] = value
	gasMap["F64x2Gt"] = value
	gasMap["F64x2Le"] = value
	gasMap["F64x2Ge"] = value
	gasMap["V128Not"] = value
	gasMap["V128And"] = value
	gasMap["V128Or"] = value
	gasMap["V128Xor"] = value
	gasMap["V128Bitselect"] = value
	gasMap["I8x16Neg"] = value
	gasMap["I8x16AnyTrue"] = value
	gasMap["I8x16AllTrue"] = value
	gasMap["I8x16Shl"] = value
	gasMap["I8x16ShrS"] = value
	gasMap["I8x16ShrU"] = value
	gasMap["I8x16Add"] = value
	gasMap["I8x16AddSaturateS"] = value
	gasMap["I8x16AddSaturateU"] = value
	gasMap["I8x16Sub"] = value
	gasMap["I8x16SubSaturateS"] = value
	gasMap["I8x16SubSaturateU"] = value
	gasMap["I8x16Mul"] = value
	gasMap["I16x8Neg"] = value
	gasMap["I16x8AnyTrue"] = value
	gasMap["I16x8AllTrue"] = value
	gasMap["I16x8Shl"] = value
	gasMap["I16x8ShrS"] = value
	gasMap["I16x8ShrU"] = value
	gasMap["I16x8Add"] = value
	gasMap["I16x8AddSaturateS"] = value
	gasMap["I16x8AddSaturateU"] = value
	gasMap["I16x8Sub"] = value
	gasMap["I16x8SubSaturateS"] = value
	gasMap["I16x8SubSaturateU"] = value
	gasMap["I16x8Mul"] = value
	gasMap["I32x4Neg"] = value
	gasMap["I32x4AnyTrue"] = value
	gasMap["I32x4AllTrue"] = value
	gasMap["I32x4Shl"] = value
	gasMap["I32x4ShrS"] = value
	gasMap["I32x4ShrU"] = value
	gasMap["I32x4Add"] = value
	gasMap["I32x4Sub"] = value
	gasMap["I32x4Mul"] = value
	gasMap["I64x2Neg"] = value
	gasMap["I64x2AnyTrue"] = value
	gasMap["I64x2AllTrue"] = value
	gasMap["I64x2Shl"] = value
	gasMap["I64x2ShrS"] = value
	gasMap["I64x2ShrU"] = value
	gasMap["I64x2Add"] = value
	gasMap["I64x2Sub"] = value
	gasMap["F32x4Abs"] = value
	gasMap["F32x4Neg"] = value
	gasMap["F32x4Sqrt"] = value
	gasMap["F32x4Add"] = value
	gasMap["F32x4Sub"] = value
	gasMap["F32x4Mul"] = value
	gasMap["F32x4Div"] = value
	gasMap["F32x4Min"] = value
	gasMap["F32x4Max"] = value
	gasMap["F64x2Abs"] = value
	gasMap["F64x2Neg"] = value
	gasMap["F64x2Sqrt"] = value
	gasMap["F64x2Add"] = value
	gasMap["F64x2Sub"] = value
	gasMap["F64x2Mul"] = value
	gasMap["F64x2Div"] = value
	gasMap["F64x2Min"] = value
	gasMap["F64x2Max"] = value
	gasMap["I32x4TruncSF32x4Sat"] = value
	gasMap["I32x4TruncUF32x4Sat"] = value
	gasMap["I64x2TruncSF64x2Sat"] = value
	gasMap["I64x2TruncUF64x2Sat"] = value
	gasMap["F32x4ConvertSI32x4"] = value
	gasMap["F32x4ConvertUI32x4"] = value
	gasMap["F64x2ConvertSI64x2"] = value
	gasMap["F64x2ConvertUI64x2"] = value
	gasMap["V8x16Swizzle"] = value
	gasMap["V8x16Shuffle"] = value
	gasMap["I8x16LoadSplat"] = value
	gasMap["I16x8LoadSplat"] = value
	gasMap["I32x4LoadSplat"] = value
	gasMap["I64x2LoadSplat"] = value

	return gasMap
}