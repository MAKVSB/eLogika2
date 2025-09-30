package certstore

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func LoadCertByPath(certFile string, keyFile string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cert.Leaf, _ = x509.ParseCertificate(cert.Certificate[0])
	return &cert, nil
}

func LoadCertCommon(certPath *string, certTP *string) (*tls.Certificate, error) {
	if certPath != nil {
		certFile := *certPath + ".crt"
		keyFile := *certPath + ".key"

		if fileExists(certFile) && fileExists(keyFile) {
			return LoadCertByPath(certFile, keyFile)
		} else {
			return nil, fmt.Errorf("one or both certificate files are missing")
		}
	} else if certTP != nil {
		return LoadCert(*certTP)
	} else {
		return nil, fmt.Errorf("certPath or certTP must be specified")
	}
}
