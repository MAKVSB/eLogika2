//go:build windows
// +build windows

package certstore

import (
	"crypto/tls"

	"github.com/google/certtostore"
)

func LoadCert(certTP string) (*tls.Certificate, error) {
	store, err := certtostore.OpenWinCertStoreCurrentUser(
		certtostore.ProviderMSSoftware,
		"",
		[]string{certTP},
		nil,
		false,
	)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	cert, ctx, err := store.CertWithContext()
	if err != nil {
		return nil, err
	}

	key, err := store.CertKey(ctx)
	if err != nil {
		return nil, err
	}

	return &tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key, // key is a crypto.Signer
	}, nil
}
