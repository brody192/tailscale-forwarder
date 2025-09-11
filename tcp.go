package main

import (
	"context"
	"fmt"
	"main/internal/util"
	"net"
	"time"

	"github.com/northbright/iocopy"
	"golang.org/x/sync/errgroup"
)

func fwdTCP(sourceConn net.Conn, targetAddr string, targetPort int) error {
	defer sourceConn.Close()

	if tcpConn, ok := sourceConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	targetConn, err := net.Dial("tcp", net.JoinHostPort(targetAddr, fmt.Sprintf("%d", targetPort)))
	if err != nil {
		return fmt.Errorf("failed to dial target: %w", err)
	}

	defer targetConn.Close()

	if tcpConn, ok := targetConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer cancel()

		defer func() {
			if tcpConn, ok := targetConn.(*net.TCPConn); ok {
				tcpConn.CloseWrite()
			}
		}()

		if _, err := iocopy.Copy(ctx, targetConn, sourceConn); err != nil && !util.IsExpectedCopyError(err) {
			return fmt.Errorf("failed to copy data to target: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		defer cancel()

		defer func() {
			if tcpConn, ok := sourceConn.(*net.TCPConn); ok {
				tcpConn.CloseWrite()
			}
		}()

		if _, err := iocopy.Copy(ctx, sourceConn, targetConn); err != nil && !util.IsExpectedCopyError(err) {
			return fmt.Errorf("failed to copy data from source: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	return nil
}
