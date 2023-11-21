package wasi

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type WasmFunc struct {
	event  api.Function
	init   api.Function
	malloc api.Function
	free   api.Function
	module api.Module
}

func NewWasmFunc(ctx context.Context, runtime wazero.Runtime, module string) (*WasmFunc, error) {
	wasmPath := fmt.Sprintf("worklets/%v.wasm", module)
	greetWasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, err
	}
	mod, err := runtime.Instantiate(ctx, greetWasm)
	if err != nil {
		return nil, err
	}
	initf := mod.ExportedFunction("init")
	wr := &WasmFunc{
		event: mod.ExportedFunction("event"),
		init:  initf,
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

func (f *WasmFunc) Event(ctx context.Context, data string) (string, error) {
	logger := zerolog.Ctx(ctx)

	// event data
	eventDataPtr, eventDataSize, err := f.mallocData(ctx, data)
	if err != nil {
		return "", err
	}
	defer f.free.Call(ctx, eventDataPtr)

	// calls the event function which will handle it.
	logger.Debug().Msg("calling Event method")
	eventResultPtrSize, err := f.event.Call(ctx, eventDataPtr, eventDataSize)
	if err != nil {
		return "", err
	}

	return f.readData(ctx, eventResultPtrSize[0])
}

func (f *WasmFunc) mallocData(ctx context.Context, eventData string) (uint64, uint64, error) {
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

func (f *WasmFunc) readData(ctx context.Context, eventResultPtrSize uint64) (string, error) {
	logger := zerolog.Ctx(ctx)
	eventResultPtr := uint32(eventResultPtrSize >> 32)
	eventResultSize := uint32(eventResultPtrSize)

	if eventResultPtr != 0 {
		defer func() {
			_, err := f.free.Call(ctx, uint64(eventResultPtr))
			if err != nil {
				logger.Err(err)
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
