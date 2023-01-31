package vmess

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/transport"
	utls "github.com/refraction-networking/utls"

	tlsC "github.com/Dreamacro/clash/component/tls"
	C "github.com/Dreamacro/clash/constant"
)

type TLSConfig struct {
	Host              string
	SkipCertVerify    bool
	FingerPrint       string
	ClientFingerprint string
	NextProtos        []string
}

func StreamTLSConn(conn net.Conn, cfg *TLSConfig) (net.Conn, error) {
	tlsConfig := &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: cfg.SkipCertVerify,
		NextProtos:         cfg.NextProtos,
	}

	if len(cfg.FingerPrint) == 0 {
		tlsConfig = tlsC.GetGlobalTLSConfig(tlsConfig)
	} else {
		var err error
		if tlsConfig, err = tlsC.GetSpecifiedFingerprintTLSConfig(tlsConfig, cfg.FingerPrint); err != nil {
			return nil, err
		}
	}

	if len(cfg.ClientFingerprint) != 0 {
		if fingerprint, exists := transport.GetFingerprint(cfg.ClientFingerprint); exists {
			log.Debugln("using HelloID:%s", fingerprint)

			utlsConn := utls.UClient(conn, transport.CopyConfig(tlsConfig), utls.ClientHelloID{
				Client:  fingerprint.Client,
				Version: fingerprint.Version,
				Seed:    nil,
			})

			ctx, cancel := context.WithTimeout(context.Background(), C.DefaultTLSTimeout)
			defer cancel()

			err := utlsConn.HandshakeContext(ctx)
			return utlsConn, err
		}
	}

	tlsConn := tls.Client(conn, tlsConfig)

	ctx, cancel := context.WithTimeout(context.Background(), C.DefaultTLSTimeout)
	defer cancel()

	err := tlsConn.HandshakeContext(ctx)
	return tlsConn, err
}
