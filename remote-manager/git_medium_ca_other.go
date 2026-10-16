//go:build !windows

package main

import "crypto/x509"

func mediumHostRoots() ([]*x509.Certificate, error) { return mediumHostRootsUnix() }
