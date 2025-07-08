package main

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/sync/errgroup"
)

func fwdTCP(sourceConn net.Conn, targetAddr string, targetPort int) error {
	defer sourceConn.Close()

	targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", targetAddr, targetPort))
	if err != nil {
		return fmt.Errorf("failed to dial target: %w", err)
	}

	defer targetConn.Close()

	g := errgroup.Group{}

	g.Go(func() error {
		if _, err := io.Copy(targetConn, sourceConn); err != nil {
			return fmt.Errorf("failed to copy data to target: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		if _, err := io.Copy(sourceConn, targetConn); err != nil {
			return fmt.Errorf("failed to copy data from source: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	return nil
}
