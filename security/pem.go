// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package security

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/fs"
)

func ReadPrivateKeyFromPEMFile[KeyType any](file fs.File) (*KeyType, error) {
	pemBytes, _ := io.ReadAll(file)
	pemBlock, _ := pem.Decode([]byte(pemBytes))
	if pemBlock == nil {
		return nil, fmt.Errorf("no PEM data found")
	}
	key, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*KeyType)
	if !ok {
		return nil, fmt.Errorf("unexpected key type")
	}
	return rsaKey, nil
}
