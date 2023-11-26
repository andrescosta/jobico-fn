package wasi

import (
	"context"
	"os"

	"github.com/andrescosta/goico/pkg/io"
	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmRuntime struct {
	module   *WasmModuleString
	cacheDir *string
	cache    wazero.CompilationCache
}

func NewWasmRuntime(ctx context.Context, tempDir string, moduleId string, mainFuncName string, wasmModule []byte) (*WasmRuntime, error) {
	wruntime := &WasmRuntime{}

	// Creates a directory to store wazero cache
	if err := io.CreateDirIfNotExist(tempDir); err != nil {
		return nil, err
	}
	cacheDir, err := os.MkdirTemp(tempDir, "cache")
	if err != nil {
		return nil, err
	}
	wruntime.cacheDir = &cacheDir

	cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
	if err != nil {
		wruntime.Close(ctx)
		return nil, err
	}
	wruntime.cache = cache

	runtimeConfig := wazero.NewRuntimeConfig().WithCompilationCache(cache)
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	// exports the log function
	_, err = runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(log).Export("log").
		Instantiate(ctx)
	if err != nil {
		wruntime.Close(ctx)
		return nil, err
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	module, err := runtime.Instantiate(ctx, wasmModule)
	if err != nil {
		wruntime.Close(ctx)
		return nil, err
	}

	w, err := newWasmModuleString(ctx, module, mainFuncName)
	if err != nil {
		wruntime.Close(ctx)
		return nil, err
	}
	wruntime.module = w
	return wruntime, nil
}

func (r *WasmRuntime) ExecuteMainFunc(ctx context.Context, data []byte) (string, error) {
	return r.module.ExecuteMainFunc(ctx, string(data))
}

func (r *WasmRuntime) Close(ctx context.Context) {
	if r.cacheDir != nil {
		os.RemoveAll(*r.cacheDir)
	}
	if r.cache != nil {
		r.cache.Close(ctx)
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
