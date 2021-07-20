package http

import (
	"crypto/tls"
	"net"
	"path"
)

func NewTLSListener(certificatePath string) (net.Listener, error) {
	certFile := path.Join(certificatePath, "ssl.cert")
	certKey := path.Join(certificatePath, "ssl.key")

	cer, err := tls.LoadX509KeyPair(certFile, certKey)

	if err != nil {
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	ln, err := tls.Listen("tcp", ":443", config)

	if err != nil {
		return nil, err
	}

	return ln, nil
}
