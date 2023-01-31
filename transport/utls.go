package transport

import (
	"crypto/tls"
	"github.com/Dreamacro/clash/log"
	"github.com/mroth/weightedrand/v2"
	utls "github.com/refraction-networking/utls"
)

func GetFingerprint(ClientFingerprint string) (*utls.ClientHelloID, bool) {
	if ClientFingerprint == "random" {
		chooser, _ := weightedrand.NewChooser(
			weightedrand.NewChoice("chrome", 6),
			weightedrand.NewChoice("safari", 3),
			weightedrand.NewChoice("firefox", 1),
		)
		initClient := chooser.Pick()
		log.Debugln("use random HelloID:%s", initClient)
		fingerprint, ok := Fingerprints[initClient]
		return fingerprint, ok
	}
	fingerprint, ok := Fingerprints[ClientFingerprint]
	return fingerprint, ok
}

var Fingerprints = map[string]*utls.ClientHelloID{
	"chrome":     &utls.HelloChrome_Auto,
	"firefox":    &utls.HelloFirefox_Auto,
	"safari":     &utls.HelloSafari_Auto,
	"randomized": &utls.HelloRandomized,
}

func CopyConfig(c *tls.Config) *utls.Config {
	return &utls.Config{
		RootCAs:               c.RootCAs,
		ServerName:            c.ServerName,
		InsecureSkipVerify:    c.InsecureSkipVerify,
		VerifyPeerCertificate: c.VerifyPeerCertificate,
	}
}
