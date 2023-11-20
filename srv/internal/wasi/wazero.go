package wasi

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func InvokeModule(ctx context.Context, modname string, data string) (string, error) {
	mod := "greet"
	logger := zerolog.Ctx(ctx)
	modpath := fmt.Sprintf("worklets/%v.wasm", mod)
	var env map[string]string
	env = make(map[string]string)
	logger.Debug().Msgf("loading module %v", modpath)
	return invoke(ctx, modname, modpath, env, data)
}

func invoke(ctx context.Context, modname string, wasmPath string, env map[string]string, data string) (string, error) {
	var greetWasm []byte
	logger := zerolog.Ctx(ctx)
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	greetWasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return "", err
	}

	// exports the log function
	_, err = r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logString).Export("log").
		Instantiate(ctx)
	if err != nil {
		return "", err
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.Instantiate(ctx, greetWasm)
	if err != nil {
		return "", err
	}

	event := mod.ExportedFunction("event")
	init := mod.ExportedFunction("init")
	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// event data
	eventData := data
	eventDataSize := uint64(len(eventData))

	results, err := malloc.Call(ctx, eventDataSize)
	if err != nil {
		return "", err
	}
	eventDataPtr := results[0]
	defer free.Call(ctx, eventDataPtr)

	if !mod.Memory().Write(uint32(eventDataPtr), []byte(eventData)) {
		return "", fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d",
			eventDataPtr, eventDataSize, mod.Memory().Size())
	}

	// init, initializes the SDK.
	_, err = init.Call(ctx)
	if err != nil {
		return "", err
	}

	// calls the event function which will handle it.
	eventResultPtrSize, err := event.Call(ctx, eventDataPtr, eventDataSize)
	if err != nil {
		return "", err
	}

	eventResultPtr := uint32(eventResultPtrSize[0] >> 32)
	eventResultSize := uint32(eventResultPtrSize[0])

	if eventResultPtr != 0 {
		defer func() {
			_, err := free.Call(ctx, uint64(eventResultPtr))
			if err != nil {
				logger.Err(err)
			}
		}()
	}

	if bytes, ok := mod.Memory().Read(eventResultPtr, eventResultSize); !ok {
		return "", fmt.Errorf("Memory.Read(%d, %d) out of range of memory size %d",
			eventResultPtr, eventResultSize, mod.Memory().Size())
	} else {
		return string(bytes), nil
	}

}

func logString(ctx context.Context, m api.Module, level, offset, byteCount uint32) {
	logger := zerolog.Ctx(ctx)
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		logger.Error().Msgf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	logger.WithLevel(zerolog.Level(level)).Msg(string(buf))
}
