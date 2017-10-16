// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"runtime"

	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	log.SetFlags(0)
}

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.Release) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		s := "<<<<<<<\tHandlers Running\t>>>>>>>"
		w := 100

		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			w, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
			w -= 25
		} else {
			debugPrint("Couldn't get terminal size. Using default value...\n")
		}
		debugPrint(fmt.Sprintf("%%%ds\n", w/2), s)
		for i := 0; i < len(handlers); i += 2 {
			debugPrint("| %-50s | %-50s |\n", runtime.FuncForPC(reflect.ValueOf(handlers[i]).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(handlers[i+1]).Pointer()).Name())
		}
		nuHandlers := len(handlers)
		debugPrint("Total %d handlers found...\n\n", nuHandlers)
		handlerName := nameOfFunction(handlers.Last())
		debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func debugPrintLoadTemplate(tmpl *template.Template) {
	if IsDebugging() {
		var buf bytes.Buffer
		for _, tmpl := range tmpl.Templates() {
			buf.WriteString("\t- ")
			buf.WriteString(tmpl.Name())
			buf.WriteString("\n")
		}
		debugPrint("Loaded HTML Templates (%d): \n%s\n", len(tmpl.Templates()), buf.String())
	}
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		log.Printf("[GIN-debug] "+format, values...)
	}
}

func debugPrintWARNINGDefault() {
	debugPrint(`[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

`)
}

func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

`)
}

func debugPrintWARNINGSetHTMLTemplate() {
	debugPrint(`[WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
at initialization. ie. before any route is registered or the router is listening in a socket:

	router := gin.Default()
	router.SetHTMLTemplate(template) // << good place

`)
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\n", err)
	}
}
