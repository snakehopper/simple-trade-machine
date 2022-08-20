package ftx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (a Api) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, a.Secret)
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}
