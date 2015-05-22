// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"log"
	"os"
)

var debugLogger = log.New(os.Stdout, "[GIN-debug] ", 0)

func IsDebugging() bool {
	return ginMode == debugCode
}

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers[nuHandlers-1])
		debugPrint("%-5s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		debugLogger.Printf(format, values...)
	}
}

func debugPrintWARNING() {
	debugPrint("[WARNING] Running in DEBUG mode! Disable it before going production\n")
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\n", err)
	}
}
