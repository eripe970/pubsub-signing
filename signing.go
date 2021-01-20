package signing

import (
	pubsub2 "cloud.google.com/go/pubsub"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/pubsub/v1"
)

const (
	signatureAttribute = "signature"
)

var (
	ErrInvalidSignature = errors.New("message has an invalid signature")
	ErrNotSigned        = errors.New("message does not have a signature")
)

type PushMessage struct {
	Message      pubsub.PubsubMessage `json:"message"`
	Subscription string               `json:"subscription"`
}

func ConstructMessage(payload []byte, secret string) (*PushMessage, error) {
	var m PushMessage

	if err := json.Unmarshal(payload, &m); err != nil {
		return nil, fmt.Errorf("failed to parse message json: %s", err.Error())
	}

	// Data from pub/sub is base64 encoded
	data, err := base64.StdEncoding.DecodeString(m.Message.Data)

	if err != nil {
		return nil, err
	}

	signature := m.Message.Attributes[signatureAttribute]

	if signature == "" {
		return nil, ErrNotSigned
	}

	err = validateSignature(data, secret, signature)

	if err != nil {
		return nil, err
	}

	return &m, err
}

func validateSignature(payload []byte, secret string, signature string) error {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)

	sign := mac.Sum(nil)

	decodedSign, err := hex.DecodeString(signature)

	if err != nil {
		return err
	}

	if hmac.Equal(decodedSign, sign) {
		return nil
	}

	return ErrInvalidSignature
}

func SignMessage(message *pubsub2.Message, secret string) error {
	signature := computeSignatureWithKey(message.Data, secret)

	message.Attributes[signatureAttribute] = hex.EncodeToString(signature)

	return nil
}

func SignPushMessage(message *PushMessage, secret string) error {
	data, err := base64.StdEncoding.DecodeString(message.Message.Data)

	if err != nil {
		return err
	}

	signature := computeSignatureWithKey(data, secret)

	message.Message.Attributes[signatureAttribute] = hex.EncodeToString(signature)

	return nil
}

func computeSignatureWithKey(payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return mac.Sum(nil)
}
