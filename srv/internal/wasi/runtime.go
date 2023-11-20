package wasi

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmRuntime struct {
	funcs    map[string]*WasmFunc
	cacheDir string
	cache    wazero.CompilationCache
}

func NewWasmRuntime(ctx context.Context, path string, modules []string) (*WasmRuntime, error) {
	cacheDir, err := os.MkdirTemp(path, "cache")
	if err != nil {
		return nil, err
	}

	cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
	if err != nil {
		clean(ctx, cacheDir, nil)
		return nil, err
	}
	runtimeConfig := wazero.NewRuntimeConfig().WithCompilationCache(cache)
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	// exports the log function
	_, err = runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(log).Export("log").
		Instantiate(ctx)
	if err != nil {
		clean(ctx, cacheDir, cache)
		return nil, err
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	// init, initializes the SDK.
	wruntime := &WasmRuntime{}
	wruntime.funcs = make(map[string]*WasmFunc)
	for _, module := range modules {
		wf, err := NewWasmFunc(ctx, runtime, module)
		if err != nil {
			clean(ctx, cacheDir, cache)
			return nil, err
		}
		wruntime.funcs[module] = wf
	}
	return wruntime, nil
}

func (r *WasmRuntime) Event(ctx context.Context, wasmFunc, data string) (string, error) {
	return r.funcs[wasmFunc].Event(ctx, data)
}

func (r *WasmRuntime) Close(ctx context.Context) {
	clean(ctx, r.cacheDir, r.cache)
}

func clean(ctx context.Context, cacheDir string, cache wazero.CompilationCache) {
	os.RemoveAll(cacheDir)
	if cache != nil {
		cache.Close(ctx)
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
