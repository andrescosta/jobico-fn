package wasi

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type WasmFuncString struct {
	wasmfunction api.Function
	init         api.Function
	malloc       api.Function
	free         api.Function
	module       api.Module
}

func NewWasmFuncString(ctx context.Context, runtime wazero.Runtime, wasmfunction string, wasmFunc []byte) (*WasmFuncString, error) {
	mod, err := runtime.Instantiate(ctx, wasmFunc)
	if err != nil {
		return nil, err
	}
	initf := mod.ExportedFunction("init")
	wr := &WasmFuncString{
		wasmfunction: mod.ExportedFunction(wasmfunction),
		init:         initf,
		// These are undocumented, but exported. See tinygo-org/tinygo#2788
		malloc: mod.ExportedFunction("malloc"),
		free:   mod.ExportedFunction("free"),
		module: mod,
	}

	_, err = initf.Call(ctx)
	if err != nil {
		return nil, err
	}
	return wr, nil
}

func (f *WasmFuncString) Execute(ctx context.Context, data string) (string, error) {
	logger := zerolog.Ctx(ctx)

	// event data
	eventDataPtr, eventDataSize, err := f.mallocData(ctx, data)
	if err != nil {
		return "", err
	}
	defer f.free.Call(ctx, eventDataPtr)

	// calls the event function which will handle it.
	logger.Debug().Msg("calling Event method")
	eventResultPtrSize, err := f.wasmfunction.Call(ctx, eventDataPtr, eventDataSize)
	if err != nil {
		return "", err
	}

	return f.readData(ctx, eventResultPtrSize[0])
}

func (f *WasmFuncString) mallocData(ctx context.Context, eventData string) (uint64, uint64, error) {
	eventDataSize := uint64(len(eventData))
	results, err := f.malloc.Call(ctx, eventDataSize)
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

func (f *WasmFuncString) readData(ctx context.Context, eventResultPtrSize uint64) (string, error) {
	logger := zerolog.Ctx(ctx)
	eventResultPtr := uint32(eventResultPtrSize >> 32)
	eventResultSize := uint32(eventResultPtrSize)

	if eventResultPtr != 0 {
		defer func() {
			_, err := f.free.Call(ctx, uint64(eventResultPtr))
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
