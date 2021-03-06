# dex-account
A front-end API for local account operations (register, set/recover password) using the Dex RPC.

## Installation

```bash
git clone https://github.com/norganna/dex-account
cd dex-account
make
```

## Configuration

Create a config file called `config.yaml` containing:

```yaml
## Recommended to not expose this to the internet, keep it local only,
## your backend service should have access to this but nothing else.
web-http-addr: 127.0.0.1:9088

## Optionally listen on https instead.
# web-https-addr: 127.0.0.1:9088
# web-tls-cert: domain-name.crt
# web-tls-key: domain-name.key

## Specify your backend Dex GRPC address here.
grpc-addr: 127.0.0.1:5557
## If your Dex GRPC is secured, specify a client certificate here.
# grpc-tls-client-cert: certificate.crt

## Currently there is only the memstore unless someone feels like
## adding something more permenant.
## The store only keeps short-term identity challenges.
store:
  class: memstore
```

## Running

```bash
bin/dex-account serve config.yaml
```

## Usage

### Challenge

This is the pivotal call that authorises users to access the other tow calls
in this API Server. You should call this and get a code to email, SMS or send
via some other out-of-band mechanism to the end-user to make them identify
themselves. If you don't care about such things, or are performing administrative
actions, you can just use the returned code directly.

This generates a challenge code that you can email to the client etc, you will
require the generated code to be able to create/update an account.

Generated code is valid for 24 hours.

`GET /challenge/email@address.com`

Returns:

```json
{
    "success": true,      # If successful, otherwise check message field.
    "code": "XV44T9uWH",  # The code to supply to create/update.
    "exists": false       # Whether the specified email address exists.
}
```

### Create

Creates a new account for a client.

You myst supply a valid challenge code, the account is created with
the email address from the challenge code.

The hash is a bcrypted hash of the actual password, or you can supply the
actual password in the "password" field.

There's no point supplying both hash and password, if both are supplied,
the password will be ignored.

`POST /create`

Body:
```json
{
  "code": "XV44T9uWH",       # A code generated by /challenge.
  "username": "admin",       # The username to use for the account.

  "hash": "$2y$12$qII..m5i", # The bcrypted password.
  # OR
  "password": "myPassword"   # The actual password (used if hash not supplied
                             # to generate a hash).
}
```

Returns:
```json
{
  "success": true,    # Whether the account was created.
  "error": "Message"  # If an error the error message.
}
```

### Update

Updates an existing account for a client.

You must supply a valid challenge code, the email account from the
challenge is used to identify the account to update.

The hash is a bcrypted hash of the actual password, or you can supply the
actual password in the "password" field.

There's no point supplying both hash and password, if both are supplied,
the password will be ignored.

`POST /update`

Body:

```json
{
  "code": "XV44T9uWH",       # A code generated by /challenge.
  "username": "admin",       # If supplied, updates the username of the account.

  "hash": "$2y$12$qII..m5i", # If supplied, updates the he bcrypted password.
  # OR
  "password": "myPassword"   # The actual password (used if hash not supplied
                             # to generate a hash).
}
```

Returns:
```json
{
  "success": true,    # Whether the account was created.
  "error": "Message"  # If an error the error message.
}
```
