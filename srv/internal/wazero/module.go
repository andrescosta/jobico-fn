package wazero

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmModuleString struct {
	mainFunc   api.Function
	initFunc   api.Function
	mallocFunc api.Function
	freeFunc   api.Function
	module     api.Module
}

type EventFuncResult struct {
	ResponseCode  uint64
	ResultPtrSize uint64
}

func NewWasmModuleString(ctx context.Context, runtime *WasmRuntime, wasmModule []byte, mainFuncName string) (*WasmModuleString, error) {

	wazeroRuntime := wazero.NewRuntimeWithConfig(ctx, runtime.runtimeConfig)

	_, err := wazeroRuntime.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logForExport).Export("log").
		Instantiate(ctx)
	if err != nil {
		return nil, err
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, wazeroRuntime)

	module, err := wazeroRuntime.Instantiate(ctx, wasmModule)
	if err != nil {
		return nil, err
	}

	initf := module.ExportedFunction("init")
	wr := &WasmModuleString{
		mainFunc: module.ExportedFunction(mainFuncName),
		initFunc: initf,
		// These are undocumented, but exported. See tinygo-org/tinygo#2788
		mallocFunc: module.ExportedFunction("malloc"),
		freeFunc:   module.ExportedFunction("free"),
		module:     module,
	}

	// Call the init function to initialize the module
	_, err = initf.Call(ctx)
	if err != nil {
		return nil, err
	}

	return wr, nil
}

func (f *WasmModuleString) ExecuteMainFunc(ctx context.Context, data string) (uint64, string, error) {
	logger := zerolog.Ctx(ctx)

	// reserve memory for the string parameter
	funcParameterPtr, funcParameterSize, err := f.writeParameterToMemory(ctx, data)
	if err != nil {
		return 0, "", err
	}
	defer func() {
		_, err := f.freeFunc.Call(ctx, uint64(funcParameterPtr))
		if err != nil {
			logger.Warn().AnErr("err", err)
		}
	}()

	resultFuncPtr, resultFuncSize, err := f.reserveMemoryForResult(ctx)
	if err != nil {
		return 0, "", err
	}
	defer func() {
		_, err := f.freeFunc.Call(ctx, uint64(resultFuncPtr))
		if err != nil {
			logger.Warn().AnErr("err", err)
		}
	}()

	logger.Debug().Msg("calling Event method")
	// The result of the call will be stored in struct pointed by resultFuncPtr
	_, err = f.mainFunc.Call(ctx, resultFuncPtr, funcParameterPtr, funcParameterSize)
	if err != nil {
		return 0, "", err
	}

	code, res, err := f.readResultFromMemory(ctx, resultFuncPtr, resultFuncSize)
	if err != nil {
		return 0, "", err
	}
	return code, res, nil
}

func (f *WasmModuleString) reserveMemoryForResult(ctx context.Context) (uint64, uint64, error) {
	eventDataSize := uint64(unsafe.Sizeof(EventFuncResult{}))
	results, err := f.mallocFunc.Call(ctx, eventDataSize)
	if err != nil {
		return 0, 0, err
	}
	eventDataPtr := results[0]
	return eventDataPtr, eventDataSize, nil
}

func (f *WasmModuleString) writeParameterToMemory(ctx context.Context, eventData string) (uint64, uint64, error) {
	eventDataSize := uint64(len(eventData))
	results, err := f.mallocFunc.Call(ctx, eventDataSize)
	if err != nil {
		return 0, 0, err
	}
	eventDataPtr := results[0]

	if !f.module.Memory().Write(uint32(eventDataPtr), []byte(eventData)) {
		return 0, 0, fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d",
			eventDataPtr, eventDataSize, f.module.Memory().Size())
	}
	return eventDataPtr, eventDataSize, nil
}

func (f *WasmModuleString) readResultFromMemory(ctx context.Context, eventResultPtr uint64, eventResultSize uint64) (uint64, string, error) {

	if data, ok := f.module.Memory().Read(uint32(eventResultPtr), uint32(eventResultSize)); ok {
		var result EventFuncResult
		err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &result)
		if err != nil {
			return 0, "", err
		}
		responseString, err := f.readDataFromMemory(ctx, result.ResultPtrSize)
		if err != nil {
			return 0, "", err
		}
		return result.ResponseCode, responseString, nil
	} else {
		return 0, "", fmt.Errorf("Memory.Read(%d, %d) out of range of memory size %d",
			eventResultPtr, eventResultSize, f.module.Memory().Size())
	}
}

func (f *WasmModuleString) readDataFromMemory(ctx context.Context, eventResultPtrSize uint64) (string, error) {
	logger := zerolog.Ctx(ctx)
	eventResultPtr := uint32(eventResultPtrSize >> 32)
	eventResultSize := uint32(eventResultPtrSize)

	if eventResultPtr != 0 {
		defer func() {
			_, err := f.freeFunc.Call(ctx, uint64(eventResultPtr))
			if err != nil {
				logger.Err(err).Msg("error freeing memory")
			}
		}()
	}

	if bytes, ok := f.module.Memory().Read(eventResultPtr, eventResultSize); !ok {
		return "", fmt.Errorf("Memory.Read(%d, %d) out of range of memory size %d",
			eventResultPtr, eventResultSize, f.module.Memory().Size())
	} else {
		return string(bytes), nil
	}
}

func logForExport(ctx context.Context, m api.Module, level, offset, byteCount uint32) {
	logger := zerolog.Ctx(ctx)
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		logger.Error().Msgf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	logger.WithLevel(zerolog.Level(level)).Msg(string(buf))
}

func (f *WasmModuleString) Close(ctx context.Context) {
	f.module.Close(ctx)
}
