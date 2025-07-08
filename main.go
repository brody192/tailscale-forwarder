package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"syscall"

	"main/internal/config"
	"main/internal/logger"

	"tailscale.com/tsnet"
)

func main() {
	ts := &tsnet.Server{
		Hostname:     config.Cfg.TSHostname,
		AuthKey:      config.Cfg.TSAuthKey,
		RunWebClient: false,
		Ephemeral:    true,
		UserLogf: func(format string, v ...any) {
			logger.Stdout.Info(fmt.Sprintf(format, v...))
		},
	}

	_, err := ts.Up(context.Background())
	if err != nil {
		logger.Stderr.Error("failed to connect server to tailscale", logger.ErrAttr(err))
		os.Exit(1)
	}

	if err := ts.Start(); err != nil {
		logger.StderrWithSource.Error("failed to start tailscale network server", logger.ErrAttr(err))
		os.Exit(1)
	}

	defer ts.Close()

	logger.Stdout.Info("🚀 Starting tailscale_fwdr",
		slog.String("ts-hostname", config.Cfg.TSHostname),
		slog.Any("connection-mappings", config.Cfg.ConnectionMappings),
	)

	wg := sync.WaitGroup{}

	for _, mapping := range config.Cfg.ConnectionMappings {
		listener, err := ts.Listen("tcp", fmt.Sprintf(":%d", mapping.SourcePort))
		if err != nil {
			logger.Stderr.Error("failed to start local listener", slog.Int("source_port", mapping.SourcePort), logger.ErrAttr(err))
			os.Exit(1)
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			logger.Stdout.Info("listening for connections",
				slog.Int("source_port", mapping.SourcePort),
				slog.String("target_addr", mapping.TargetAddr),
				slog.Int("target_port", mapping.TargetPort),
			)

			for {
				sourceConn, err := listener.Accept()
				if err != nil {
					logger.Stderr.Error("failed to accept connection",
						slog.Int("source_port", mapping.SourcePort),
						logger.ErrAttr(err),
					)

					return
				}

				go func() {
					if err := fwdTCP(sourceConn, mapping.TargetAddr, mapping.TargetPort); err != nil {
						if errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.ECONNRESET) {
							return
						}

						logger.Stderr.Error("failed to forward connection",
							slog.Int("source_port", mapping.SourcePort),
							slog.String("target_addr", mapping.TargetAddr),
							slog.Int("target_port", mapping.TargetPort),
							logger.ErrAttr(err),
						)
					}
				}()
			}
		}()
	}

	wg.Wait()
}
