# Pub/sub utilty for message signing
Utilities for signing messages

## To sign a message 
To sign a message before sending it to GCP pub/sub. 

````
message := pubsub.Message{
		Data: data,
	}

signing.SignMessage(&message, "signing-secret")

````

## Verify the signature of a message
To constrct a message and to verify the signature. 

````
signing.ConstructMessage(payload, secret)
````

The signature is stored as a hex encoded sting as an attribute on the message, like.

````
{
  "message": {
    "attributes": {
      "signature": "cc6d0457d8b4ec0e994d793302cd6962a0d12101abbc79e561f532a826eca4ee"
    },
    "data": "eyJOYW1lIjoibmV3LW5hbWUifQ=="
  },
  "subscription": "cloud-run-events"
}
````

See examples for a solution for publishing messages on Cloud Run. 
