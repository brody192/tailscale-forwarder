package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"main/internal/config"
	"main/internal/logger"
	"main/internal/tcp"

	"tailscale.com/tsnet"
)

func main() {
	config.Load()

	ts := &tsnet.Server{
		Hostname:     config.Cfg.TSHostname,
		AuthKey:      config.Cfg.TSAuthKey,
		ControlURL:   config.Cfg.TSControlURL,
		Dir:          config.Cfg.TSStateDir,
		RunWebClient: false,
		Ephemeral:    config.Cfg.TSEphemeral,
		UserLogf: func(format string, v ...any) {
			logger.Stdout.Info(fmt.Sprintf(format, v...))
		},
	}

	if _, err := ts.Up(context.Background()); err != nil {
		logger.Stderr.Error("failed to connect server to tailscale", logger.ErrAttr(err))
		os.Exit(1)
	}

	if err := ts.Start(); err != nil {
		logger.StderrWithSource.Error("failed to start tailscale network server", logger.ErrAttr(err))
		os.Exit(1)
	}

	defer ts.Close()

	logger.Stdout.Info("Starting tailscale_fwdr",
		slog.String("ts-hostname", config.Cfg.TSHostname),
		slog.String("ts-control-url", config.Cfg.TSControlURL),
		slog.String("ts-state-dir", config.Cfg.TSStateDir),
		slog.Bool("ts-ephemeral", config.Cfg.TSEphemeral),
		slog.Any("connection-mappings", config.Cfg.ConnectionMappings),
	)

	wg := sync.WaitGroup{}

	for _, mapping := range config.Cfg.ConnectionMappings {
		listen := ts.Listen

		if mapping.HTTPS {
			listen = ts.ListenTLS
		}

		listener, err := listen("tcp", fmt.Sprintf(":%d", mapping.SourcePort))
		if err != nil {
			logger.Stderr.Error("failed to start local listener",
				slog.Int("source_port", mapping.SourcePort),
				slog.Bool("https", mapping.HTTPS),
				logger.ErrAttr(err),
			)
			os.Exit(1)
		}

		wg.Go(func() {
			logger.Stdout.Info("listening for connections",
				slog.Int("source_port", mapping.SourcePort),
				slog.Bool("https", mapping.HTTPS),
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

					continue
				}

				go func() {
					if err := tcp.Forward(sourceConn, mapping.TargetAddr, mapping.TargetPort); err != nil {
						logger.Stderr.Error("failed to forward connection",
							slog.Int("source_port", mapping.SourcePort),
							slog.String("target_addr", mapping.TargetAddr),
							slog.Int("target_port", mapping.TargetPort),
							logger.ErrAttr(err),
						)
					}
				}()
			}
		})
	}

	wg.Wait()
}
