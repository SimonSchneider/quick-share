# quick-share

Simple secrets sharing app. End-to-end encrypted. The server only ever sees an encrypted
blob of data. The server does not have the key to decrypt the data.

## Usage

Deploy the server:

```bash
$ docker run -d -p 8080:8080 --name quick-share-server ghcr.io/quick-share/quick-share-server
```

Navigate to the page in your browser and start sharing secrets.

## Details

```mermaid
sequenceDiagram
    actor a as Alice
    participant ab as Alice Browser
    participant srv as Server
    participant bb as Bob Browser
    actor b as Bob
    a->>ab: Input secret
    a->>ab: Create secret
    activate ab
    ab ->> ab: Generate encryption key
    ab ->> ab: Encrypt secret with key

    ab ->> srv: encrypted blob
    activate srv
    srv ->> ab: secretId
    deactivate srv

    ab->>a: link = server.com/secrets/{id}#35;{key}
    deactivate ab
    
    a ->> b: send link to Bob or let him use the QR code

    b ->> bb: Navigate to link
    activate bb
    bb ->> srv: GET server.com/secrets/{id}
    activate srv
    Note right of srv: url fragments (anything after #35;)<br/>are never sent to the server
    srv ->> bb: encrypted blob
    deactivate srv

    bb ->> bb: decrypt blob using key from url fragment
    deactivate bb
    bb -> b: show secret

```

## More Features

- [ ] Allow uploading files
- [ ] Allow setting a password for the secret in addition to the encryption key. This would allow the secret to be shared even more securely by sharing the password over a different channel (ie. voice) and only allowing the link to be a one time use.