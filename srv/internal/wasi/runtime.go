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
	funcs    map[string]*WasmFuncString
	cacheDir string
	cache    wazero.CompilationCache
}

type Func struct {
	ModuleId   string
	WasmModule []byte
	FuncName   string
}

func NewWasmRuntime(ctx context.Context, tempDir string, funcs []*Func) (*WasmRuntime, error) {
	cacheDir, err := os.MkdirTemp(tempDir, "cache")
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
	wruntime.funcs = make(map[string]*WasmFuncString)
	for _, f := range funcs {
		wf, err := NewWasmFuncString(ctx, runtime, f.FuncName, f.WasmModule)
		if err != nil {
			clean(ctx, cacheDir, cache)
			return nil, err
		}
		wruntime.funcs[f.ModuleId] = wf
	}
	return wruntime, nil
}

func (r *WasmRuntime) Execute(ctx context.Context, wasmRuntime string, data []byte) (string, error) {
	return r.funcs[wasmRuntime].Execute(ctx, string(data))
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
