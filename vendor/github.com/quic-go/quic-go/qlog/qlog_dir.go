package qlog

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/qlogwriter"
)

// EventSchema is the qlog event schema for QUIC
const EventSchema = "urn:ietf:params:qlog:events:quic-12"

// DefaultConnectionTracer creates a qlog file in the qlog directory specified by the QLOGDIR environment variable.
// File names are <odcid>_<perspective>.sqlog.
// Returns nil if QLOGDIR is not set.
func DefaultConnectionTracer(_ context.Context, isClient bool, connID ConnectionID) qlogwriter.Trace {
	return defaultConnectionTracerWithSchemas(isClient, connID, []string{EventSchema})
}

func DefaultConnectionTracerWithSchemas(_ context.Context, isClient bool, connID ConnectionID, eventSchemas []string) qlogwriter.Trace {
	if !slices.Contains(eventSchemas, EventSchema) {
		eventSchemas = append([]string{EventSchema}, eventSchemas...)
	}
	return defaultConnectionTracerWithSchemas(isClient, connID, eventSchemas)
}

func defaultConnectionTracerWithSchemas(isClient bool, connID ConnectionID, eventSchemas []string) qlogwriter.Trace {
	qlogDir := os.Getenv("QLOGDIR")
	if qlogDir == "" {
		return nil
	}
	if _, err := os.Stat(qlogDir); os.IsNotExist(err) {
		if err := os.MkdirAll(qlogDir, 0o755); err != nil {
			log.Fatalf("failed to create qlog dir %s: %v", qlogDir, err)
		}
	}
	label := "server"
	if isClient {
		label = "client"
	}
	path := fmt.Sprintf("%s/%s_%s.sqlog", strings.TrimRight(qlogDir, "/"), connID, label)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create qlog file %s: %s", path, err.Error())
		return nil
	}
	fileSeq := qlogwriter.NewConnectionFileSeq(
		utils.NewBufferedWriteCloser(bufio.NewWriter(f), f),
		isClient,
		connID,
		eventSchemas,
	)
	go fileSeq.Run()
	return fileSeq
}
