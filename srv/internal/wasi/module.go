package wasi

import (
	"context"
	"fmt"

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

// Documentation: https://github.com/tetratelabs/wazero/blob/main/examples/multiple-runtimes/counter.go
func NewWasmModuleString(ctx context.Context, runtime *WasmRuntime, wasmModule []byte, mainFuncName string) (*WasmModuleString, error) {

	wazeroRuntime := wazero.NewRuntimeWithConfig(ctx, runtime.runtimeConfig)

	_, err := wazeroRuntime.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(log).Export("log").
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

func (f *WasmModuleString) ExecuteMainFunc(ctx context.Context, data string) (string, error) {
	logger := zerolog.Ctx(ctx)

	// event data
	eventDataPtr, eventDataSize, err := f.mallocForParamData(ctx, data)
	if err != nil {
		return "", err
	}
	defer f.freeFunc.Call(ctx, eventDataPtr)

	logger.Debug().Msg("calling Event method")
	eventResultPtrSize, err := f.mainFunc.Call(ctx, eventDataPtr, eventDataSize)
	if err != nil {
		return "", err
	}

	return f.readParamData(ctx, eventResultPtrSize[0])
}

func (f *WasmModuleString) mallocForParamData(ctx context.Context, eventData string) (uint64, uint64, error) {
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

func (f *WasmModuleString) readParamData(ctx context.Context, eventResultPtrSize uint64) (string, error) {
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

func log(ctx context.Context, m api.Module, level, offset, byteCount uint32) {
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
