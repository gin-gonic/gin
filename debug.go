// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import "log"

func IsDebugging() bool {
	return ginMode == debugCode
}

func debugRoute(httpMethod, absolutePath string, handlers []HandlerFunc) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers[nuHandlers-1])
		debugPrint("%-5s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		log.Printf("[GIN-debug] "+format, values...)
	}
}
