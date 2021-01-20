package signing

import (
	pubsub2 "cloud.google.com/go/pubsub"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"google.golang.org/api/pubsub/v1"
	"testing"
)

func TestConstructMessageInvalidJson(t *testing.T) {
	if _, err := ConstructMessage([]byte("not-valid-json"), "secret"); err == nil {
		t.Fatal("no error occurred for invalid json")
	}
}

func TestConstructMessageInvalidEncoding(t *testing.T) {
	message := PushMessage{
		Message: pubsub.PubsubMessage{
			Data: "not base 64 encoded",
		},
	}

	data, _ := json.Marshal(message)

	if _, err := ConstructMessage(data, "secret"); err == nil {
		t.Fatal("no error occurred for invalid message")
	}
}

func TestConstructMessageWithoutSigningAttribute(t *testing.T) {
	if _, err := ConstructMessage([]byte(`{"message:" : "QQ=="}`), "secret"); err != ErrNotSigned {
		t.Fatal("no error occurred when message missing signature")
	}
}

func TestConstructMessageWithInvalidSignature(t *testing.T) {
	signature1 := computeSignatureWithKey([]byte("some data"), "invalid")

	message := PushMessage{
		Message: pubsub.PubsubMessage{
			Data: base64.StdEncoding.EncodeToString([]byte("some data")),
			Attributes: map[string]string{
				signatureAttribute: hex.EncodeToString(signature1),
			},
		},
		Subscription: "",
	}

	data, _ := json.Marshal(message)

	if _, err := ConstructMessage(data, "secret"); err != ErrInvalidSignature {
		t.Fatal("no error occurred when message missing signature")
	}
}

func TestConstructMessageWithValidSignature(t *testing.T) {
	message := PushMessage{
		Message: pubsub.PubsubMessage{
			Data: base64.StdEncoding.EncodeToString([]byte("some data")),
		},
		Subscription: "",
	}

	_ = SignPushMessage(&message, "secret")

	data, _ := json.Marshal(message)

	if _, err := ConstructMessage(data, "secret"); err != nil {
		t.Fatal("error occurred when constructing message with valid signature")
	}
}

func TestSignPushMessageWithInvalidEncoding(t *testing.T) {
	message := PushMessage{
		Message: pubsub.PubsubMessage{
			Data: "this is not base 64 encoded",
		},
		Subscription: "",
	}

	if err := SignPushMessage(&message, "secret"); err == nil {
		t.Fatal("no error occurred when signature not base64 encoded")
	}
}

func TestSignMessage(t *testing.T) {
	message := pubsub2.Message{
		Data: []byte("foo bar"),
		Attributes: map[string]string{

		},
	}

	if err := SignMessage(&message, "secret"); err != nil {
		t.Fatal("error occurred when signing message")
	}
}

func TestValidateSignatureInvalidEncoding(t *testing.T) {
	if err := validateSignature([]byte(`some data`), "secret", "not hex decoded"); err == nil {
		t.Fatal("no error occurred when signature not hex encoded")
	}
}

func TestValidateInvalidSignature(t *testing.T) {
	signature1 := computeSignatureWithKey([]byte("test-data"), "secret")

	if err := validateSignature([]byte("some data"), "secret2", hex.EncodeToString(signature1)); err != ErrInvalidSignature {
		t.Fatal("no error occurred when signature is invalid")
	}
}

func TestValidateValidSignature(t *testing.T) {
	signature1 := computeSignatureWithKey([]byte("some data"), "secret")

	if err := validateSignature([]byte("some data"), "secret", hex.EncodeToString(signature1)); err != nil {
		t.Fatal("error occurred when signature is valid")
	}
}

func TestComputeSignatureWithKey(t *testing.T) {
	signature := computeSignatureWithKey([]byte("test-data"), "secret")

	if len(signature) == 0 {
		t.Fatal("signature is empty")
	}
}
