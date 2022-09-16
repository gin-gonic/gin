/*
 * MIT License
 *
 * Copyright (c) 2019-present Fenny and Contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 * Code Sources In Here Come From: https://github.com/gofiber/fiber/blob/master/prefork.go
 */
package gin

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	reuseport "github.com/libp2p/go-reuseport"
)

const (
	envPreforkChildKey = "GIN_PREFORK_CHILD"
	envPreforkChildVal = "1"
)

// Holds process of childs
var children = map[int]*exec.Cmd{}

// IsChild determines if the current process is a child of Prefork
func IsChild() bool {
	return os.Getenv(envPreforkChildKey) == envPreforkChildVal
}

// watchMaster watches child procs
func watchMaster() {
	if runtime.GOOS == "windows" {
		// finds parent process,
		// and waits for it to exit
		p, err := os.FindProcess(os.Getppid())
		if err == nil {
			_, _ = p.Wait()
		}
		os.Exit(1)
	}
	// if it is equal to 1 (init process ID),
	// it indicates that the master process has exited
	for range time.NewTicker(time.Millisecond * 500).C {
		if os.Getppid() == 1 {
			os.Exit(1)
		}
	}
}

// prefork manages child processes to make use of the OS REUSEPORT or REUSEADDR feature
func prefork(addr string, engine *Engine) (err error) {
	// 👶 child process 👶
	if IsChild() {
		// use 1 cpu core per child process
		runtime.GOMAXPROCS(1)

		// kill current child proc when master exits
		go watchMaster()

		// Run child
		listener, err := reuseport.Listen("tcp", addr)
		if err != nil {
			return err
		}
		defer listener.Close()

		return http.Serve(listener, engine.Handler())
	}

	// child structure to be used in error returning
	type child struct {
		pid int
		err error
	}
	// create variables
	max := runtime.GOMAXPROCS(0)
	channel := make(chan child, max)

	// kill child procs when master exits
	defer func() {
		for _, proc := range children {
			_ = proc.Process.Kill()
		}
	}()

	// launch child procs
	for i := 0; i < max; i++ {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Add gin prefork child flag into child proc env
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("%s=%s", envPreforkChildKey, envPreforkChildVal),
		)

		if err = cmd.Start(); err != nil {
			return fmt.Errorf("failed to start a child prefork process, error: %v", err)
		}

		// Store child process ids
		pid := cmd.Process.Pid
		children[pid] = cmd

		// notify master if child crashes
		go func() {
			channel <- child{pid, cmd.Wait()}
		}()
	}

	// return error if child crashes
	return (<-channel).err
}
