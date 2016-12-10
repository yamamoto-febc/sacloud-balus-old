package sakura_iot_go

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type WebhookHandlerFunc func(Payload)

type WebhookHandler struct {
	Secret     string
	HandleFunc WebhookHandlerFunc
	Debug      bool
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 400
	defer func() {
		w.WriteHeader(status)
	}()

	out := func(string, ...interface{}) {}
	if h.Debug {
		out = log.Printf
	}

	out("[DEBUG] Request received\n")

	if r.Method == "POST" {
		out("[DEBUG] Request method is POST\n")

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.Bytes()

		// Secretが設定されている場合は"X-Sakura-Signature"を検証
		if h.Secret != "" {
			signature := r.Header.Get("X-Sakura-Signature")
			if !h.verifySignature([]byte(h.Secret), signature, body) {
				status = 403
				out("[DEBUG] Invalid signature:%s", signature)
				return
			}

		}

		out("[DEBUG] Request body:%s\n", string(body))

		var payload Payload
		err := json.Unmarshal(body, &payload)
		if err != nil {
			return
		}

		if h.HandleFunc == nil {
			status = 503
			return
		}

		go h.HandleFunc(payload)

		status = 200
	} else {
		out("[DEBUG] Request method is not POST\n")
	}
}

func (h *WebhookHandler) verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = ""
	const signatureLength = 40 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || (signaturePrefix != "" && !strings.HasPrefix(signature, signaturePrefix)) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature))

	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	signBody := []byte(computed.Sum(nil))

	return hmac.Equal(signBody, actual)
}
