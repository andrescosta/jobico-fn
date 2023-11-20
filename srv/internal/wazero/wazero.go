package wazero

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func InvokeModule(ctx context.Context, modname string, data string) (string, error) {
	mod := "goenv"
	logger := zerolog.Ctx(ctx)
	modpath := fmt.Sprintf("worklets/%v.wasm", mod)
	logger.Debug().Msgf("loading module %v", modpath)
	var env map[string]string
	env = make(map[string]string)
	env["data"] = data
	logger.Debug().Msgf("loading module %v", modpath)
	return invoke(ctx, modname, modpath, env, data)
}

func invoke(ctx context.Context, modname string, wasmPath string, env map[string]string, data string) (string, error) {

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	err := buildLogImports(ctx, r, modname)
	if err != nil {
		return "", err
	}

	wasmObj, err := os.ReadFile(wasmPath)
	if err != nil {
		return "", err
	}

	stdoutBuf, config := buildConfig(env)

	// It invokes the _start function by default.
	_, err = r.InstantiateWithConfig(ctx, wasmObj, config)
	if err != nil {
		return "", err
	}

	return stdoutBuf.String(), nil
}

func buildConfig(env map[string]string) (*bytes.Buffer, wazero.ModuleConfig) {
	var stdoutBuf bytes.Buffer
	config := wazero.NewModuleConfig().WithStdout(&stdoutBuf)

	for k, v := range env {
		config = config.WithEnv(k, v)
	}

	return &stdoutBuf, config
}

func buildLogImports(ctx context.Context, r wazero.Runtime, modname string) error {
	logger := zerolog.Ctx(ctx)
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(v uint32) {
			logger.Debug().Msgf("[%v]: %v", modname, v)
		}).
		Export("log_i32").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, ptr uint32, len uint32) {
			if bytes, ok := mod.Memory().Read(ptr, len); ok {
				logger.Debug().Msgf("[%v]: %v", modname, string(bytes))
			} else {
				logger.Debug().Msgf("[%v]: log_string: unable to read wasm memory", modname)
			}
		}).
		Export("log_string").
		Instantiate(ctx)
	return err
}

func logString(ctx context.Context, m api.Module, offset, byteCount uint32) {
	logger := zerolog.Ctx(ctx)
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		logger.Error().Msgf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	logger.Info().Msg(string(buf))
}

/*r := wazero.NewRuntime(ctx)
defer r.Close(ctx)
wasi_snapshot_preview1.MustInstantiate(ctx, r)

err := buildLogImports(ctx, r, modname)
if err != nil {
	return "", err
}

r.NewHostModuleBuilder("env").
NewFunctionBuilder().WithFunc(logString).Export("log").
Instantiate(ctx)

wasmObj, err := os.ReadFile(wasmPath)
if err != nil {
	return "", err
}

stdoutBuf, config := buildConfig(env)

_, err = r.InstantiateWithConfig(ctx, wasmObj, config)
if err != nil {
	return "", err
}
return stdoutBuf.String(), nil

/*exec := mod.ExportedFunction("exec")
malloc := mod.ExportedFunction("malloc")
free := mod.ExportedFunction("free")

dataSize := uint64(len(data))

results, err := malloc.Call(ctx, dataSize)
if err != nil {
	return "", err
}
dataPtr := results[0]
defer free.Call(ctx, dataPtr)

if !mod.Memory().Write(uint32(dataPtr), []byte(data)) {
	return "", fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d",
		dataPtr, dataSize, mod.Memory().Size())
}

_, err = exec.Call(ctx, dataPtr, dataSize)

if err != nil {
	return "", err
}
return stdoutBuf.String(), nil
*/
