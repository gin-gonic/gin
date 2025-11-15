//go:build race
// +build race

package encoder

import (
	"sync"
)

var setsMu sync.RWMutex

func CompileToGetCodeSet(ctx *RuntimeContext, typeptr uintptr) (*OpcodeSet, error) {
	if typeptr > typeAddr.MaxTypeAddr || typeptr < typeAddr.BaseTypeAddr {
		codeSet, err := compileToGetCodeSetSlowPath(typeptr)
		if err != nil {
			return nil, err
		}
		return getFilteredCodeSetIfNeeded(ctx, codeSet)
	}
	index := (typeptr - typeAddr.BaseTypeAddr) >> typeAddr.AddrShift
	setsMu.RLock()
	if codeSet := cachedOpcodeSets[index]; codeSet != nil {
		filtered, err := getFilteredCodeSetIfNeeded(ctx, codeSet)
		if err != nil {
			setsMu.RUnlock()
			return nil, err
		}
		setsMu.RUnlock()
		return filtered, nil
	}
	setsMu.RUnlock()

	codeSet, err := newCompiler().compile(typeptr)
	if err != nil {
		return nil, err
	}
	filtered, err := getFilteredCodeSetIfNeeded(ctx, codeSet)
	if err != nil {
		return nil, err
	}
	setsMu.Lock()
	cachedOpcodeSets[index] = codeSet
	setsMu.Unlock()
	return filtered, nil
}
