package helpers

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

func IsCertValid(certPath string) (bool, string) {

	certValid := PathExists(certPath)

	if certValid {
		certBytes, err := os.ReadFile(certPath)

		if err != nil {
			log.Errorf("%s - cert does not exist or cannot be read", certPath)
		}

		pub, _, _, _, err := ssh.ParseAuthorizedKey(certBytes)

		if err != nil {
			log.Errorf("Error when parsing cert: %s", err.Error())
		}

		cert, ok := pub.(*ssh.Certificate)

		if !ok {
			log.Errorf("Failed to cast to cert")
		}

		now := time.Now().UTC()
		validBefore := time.Unix(int64(cert.ValidBefore), 0).UTC()
		validAfter := time.Unix(int64(cert.ValidAfter), 0).UTC()

		validBeforeString := validBefore.Format("2006-01-01 15:04:05.5 -0700")

		if now.After(validAfter) && now.Before(validBefore) {
			return true, validBeforeString
		} else {
			return false, validBeforeString
		}
	}

	return false, "Cert not found"
}
