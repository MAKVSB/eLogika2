//go:build linux
// +build linux

package certstore

import (
	"crypto/tls"
	"fmt"
)

func LoadCert(certTP string) (*tls.Certificate, error) {
	return nil, fmt.Errorf("method not implemented for linux")
}
