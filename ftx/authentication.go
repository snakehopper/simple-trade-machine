package ftx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (c *Client) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, c.Secret)
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}
