package certstore

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"

	"elogika.vsb.cz/backend/initializers"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
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
		if strings.HasSuffix(*certPath, ".pfx") {
			// load PFX
			pfxData, err := os.ReadFile("elogika.pfx")
			if err != nil {
				log.Fatalf("read pfx: %v", err)
			}

			// Replace "pfx-password" with the password you used when exporting
			privateKey, cert, caCerts, err := pkcs12.DecodeChain(pfxData, *initializers.GlobalAppConfig.CERTPASS)
			if err != nil {
				log.Fatalf("pkcs12 decode: %v", err)
			}

			// Build tls.Certificate
			tlsCert := tls.Certificate{
				PrivateKey: privateKey,
				Leaf:       cert,
				Certificate: [][]byte{
					cert.Raw,
				},
			}
			// include CA chain if present
			for _, c := range caCerts {
				tlsCert.Certificate = append(tlsCert.Certificate, c.Raw)
			}
			return &tlsCert, nil
		} else {
			certFile := *certPath + ".crt"
			keyFile := *certPath + ".key"

			if fileExists(certFile) && fileExists(keyFile) {
				return LoadCertByPath(certFile, keyFile)
			} else {
				return nil, fmt.Errorf("one or both certificate files are missing")
			}
		}
	} else {
		return nil, fmt.Errorf("certPath or certTP must be specified")
	}
}
