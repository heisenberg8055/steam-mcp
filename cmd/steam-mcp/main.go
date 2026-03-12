package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mark3labs/mcp-go/server"
)

type app struct {
	logger *slog.Logger
}

func main() {
	var transport string
	flag.StringVar(&transport, "trans", "stdio", "Transport Type(stdio, sse, streamableHTTP, inProcess)")
	flag.Parse()

	logger := slog.New(tint.NewHandler(os.Stderr, nil))

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	s := server.NewMCPServer("steam-mcp", "0.1.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithRecovery(),
		server.WithInstructions("steam mcp server"),
	)
	switch transport {
	case "stdio":
		if err := server.ServeStdio(s); err != nil {
			logger.Error("failed to start stdio transport", slog.Any("error", err), "h", "D6ED059F1B1F5D4F959FEA4B8FA8D167")
		}
	case "sse":
		httpServer := server.NewStreamableHTTPServer(s)
		if err := httpServer.Start(":8080"); err != nil {
			logger.Error("failed to start httpServer", slog.Any("error", err))
		}
	}
}
