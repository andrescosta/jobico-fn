package wazero

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
	logger.Debug().Msgf("loading module %v", modpath)
	var env map[string]string
	env = make(map[string]string)
	env["data"] = data
	logger.Debug().Msgf("loading module %v", modpath)
	return invoke(ctx, modname, modpath, env, data)
}

func invoke(ctx context.Context, modname string, wasmPath string, env map[string]string, data string) (string, error) {
	var greetWasm []byte
	logger := zerolog.Ctx(ctx)
	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	greetWasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return "", err
	}

	// Instantiate a Go-defined module named "env" that exports a function to
	// log to the console.
	_, err = r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logString).Export("log").
		Instantiate(ctx)
	if err != nil {
		return "", err
	}

	// Note: testdata/greet.go doesn't use WASI, but TinyGo needs it to
	// implement functions such as panic.
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Instantiate a WebAssembly module that imports the "log" function defined
	// in "env" and exports "memory" and functions we'll use in this example.
	mod, err := r.Instantiate(ctx, greetWasm)
	if err != nil {
		return "", err
	}

	// Get references to WebAssembly functions we'll use in this example.
	event := mod.ExportedFunction("event")
	init := mod.ExportedFunction("init")
	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// Let's use the argument to this main function in Wasm.
	name := data
	nameSize := uint64(len(name))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, nameSize)
	if err != nil {
		return "", err
	}
	namePtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, namePtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(namePtr), []byte(name)) {
		return "", fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size())
	}

	// Now, we can call "greet", which reads the string we wrote to memory!
	_, err = init.Call(ctx)
	if err != nil {
		return "", err
	}

	// Now, we can call "greet", which reads the string we wrote to memory!
	ptrSize, err := event.Call(ctx, namePtr, nameSize)
	if err != nil {
		return "", err
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by TinyGo.
	//ptrSize, err := greeting.Call(ctx, namePtr, nameSize)
	//if err != nil {
	//	return "", err
	//}

	eventPtr := uint32(ptrSize[0] >> 32)
	eventSize := uint32(ptrSize[0])

	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	if eventPtr != 0 {
		defer func() {
			_, err := free.Call(ctx, uint64(eventPtr))
			if err != nil {
				logger.Err(err)
			}
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := mod.Memory().Read(eventPtr, eventSize); !ok {
		return "", fmt.Errorf("Memory.Read(%d, %d) out of range of memory size %d",
			eventPtr, eventSize, mod.Memory().Size())
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

/*r := wazero.NewRuntime(ctx)
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

r := wazero.NewRuntime(ctx)
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
