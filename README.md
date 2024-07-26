# quick-share

Simple secrets sharing app with E2EE. The server only ever sees an encrypted
blob of data. The server does not have the key to decrypt the data. The key is
only shared with the recipient of the secret.

## Usage

Deploy the server:

```bash
$ docker run -d -p 8888:80 --name quick-share ghcr.io/simonschneider/quick-share
```

or use the docker-compose file:

```bash
$ docker-compose up -d
```

or you can build the binary yourself and just use that (the binary embeds all static files to be selfcontained and does not need any other files):

```bash
$ go build -o quick-share *.go
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
