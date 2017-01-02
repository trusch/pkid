pkid
====

A service for managing public key infrastructures via a REST-full interface.

# Features
* Create self-signed CA's
* Manage multiple root CA's
* Create signed sub-CA's
* Create signed server certificates
* Create signed client certificates
* RSA or ECC Keys
* Revoke Sub-CA's, clients or servers
* Automatically create CRL's

# API

## Create Certificates

These endpoints are used to create keys and issue certificates.

Options for all following endpoints are:
* `name`: string (required)
* `curve`: string (optional, default: P521)
  * valid values: P521, P384, P256, P224
* `rsaBits`: int (optional)
  * valid values: 4096, 2048, 1024
* `notBefore`: int (optional, secs since epoche, defaults to current time)
* `validFor`: string (optional, example: 12h30m, defaults to 8760h (-> 1 Year))

#### Create root CA (self signed)
* Request: `POST /ca?name=my-ca-name`
* Response: {uuid}

#### Create Sub CA
* Request: `POST /ca/{root-uuid}/ca?name=my-sub-ca`
* Response: {uuid}

#### Create Client
* Request: `POST /ca/{root-uuid}/client?name=my-client`
* Response: {uuid}

#### Create Server
* Request: `POST /ca/{root-uuid}/server?name=my-server`
* Response: {uuid}

## Get Certificates/Keys

These endpoints are used to retrieve generated certificates and keys

#### Get CA Certificate
* Request: `GET /ca/{root-uuid}/cert`
* Response: {pem certificate data}

#### Get CA Key
* Request: `GET /ca/{root-uuid}/key`
* Response: {pem key data}

#### Get Client Certificate
* Request: `GET /ca/{root-uuid}/client/{uuid}/cert`
* Response: {pem certificate data}

#### Get Client Key
* Request: `GET /ca/{root-uuid}/client/{uuid}/key`
* Response: {pem key data}

## Revoke Certificates

These endpoints can be used to revoke certificates and get the resulting CRL.

#### Revoke a CA
* Request: `POST /ca/{root-uuid}/ca/{uuid}/revoke`
* Response: "revoked"

#### Revoke a Server
* Request: `POST /ca/{root-uuid}/server/{uuid}/revoke`
* Response: "revoked"

#### Revoke a Client
* Request: `POST /ca/{root-uuid}/client/{uuid}/revoke`
* Response: "revoked"

#### Get Certificate Revocation List (CRL)
* Request: `GET /ca/{root-uuid}/crl`
* Response: {pem crl data}

## Info about CA

These endpoints can be used to gather information about a specific CA

#### Get CA info
* Request: `GET /ca/{root-uuid}`
* Response:
```json
  {
    "Entity": {
      "ID": "{uuid}",
      "Name": "my-ca",
      "IsRevoked": false,
    },
    "Revoked": [2,5,6],
    "CAs": {
      "{uuid}": "my-sub-ca"
    },
    "Clients": {
      "{uuid}": "my-client"
    },
    "Servers": {
      "{uuid}": "my-server"
    }
  }
```

#### List sub CA's
* Request: `GET /ca/{root-uuid}/ca`
* Response:
```json
  {
    "{uuid}": "my-sub-ca"
  }
```

#### List clients
* Request: `GET /ca/{root-uuid}/client`
* Response:
```json
  {
    "{uuid}": "my-client"
  }
```

#### List servers
* Request: `GET /ca/{root-uuid}/server`
* Response:
```json
  {
    "{uuid}": "my-server"
  }
```
