# KeyGuard

A little app to serve SSH keys over an authenticated endpoint. A helper script
is used to add the key to the SSH agent with an expiry

Only [YubiKey](https://www.yubico.com/why-yubico/for-individuals/) One-time password auth at the moment.

## Usage

The app is configured via environment variables:

```
KG_PUBLIC_URL: the public URL where the /key can be queried, required
KG_PRIVATE_KEY: path to the ssh private key, default: id_rsa
KG_LOADER_SCRIPT: path to the loader script, default: loader.sh
KG_AUTH_MODULE: name of the authentication module; default: yubikey
KG_PORT: the http listen port, default: 8000

# yubikey options
KG_YUBI_CLIENT_ID: the yubico client id, required
KG_YUBI_API_KEY: the yubico api key, required
KG_YUBI_API_HOST: the yubi auth server, default: api.yubico.com/wsapi/2.0/verify
KG_YUBI_USE_HTTPS: protocol for contacting the auth server, default: true
```

### Load key!

Deploy the app to your favourite application platform and:

```
$ curl -s https://key.yourdomain.org | bash
OTP: ccccsfrhkrucdedthkkrdkkrbjdhidjkljktflhvjgcl # this is where I pressed the YubiKey button
Identity added: /tmp/tmp.2GxYjzCLaE (/tmp/tmp.2GxYjzCLaE)
Lifetime set to 32400 seconds
```

#### Retrieve the public key
Sometimes it's rather handy to get the public key when you want to add it to
certain services such as GitHub.

```
curl -s https://key.yourdomain.org/pubkey
```

### Important

You have to create an API key at [YubiCo](https://upgrade.yubico.com/getapikey/)
to use the authenticator.

## Building

```
$ go build
```

## Testing

```
bin/test
```

## How it works

The service exposes three endpoints:
* `/:expiry`
* `/key`
* `/pubkey`

`/` responds with a shell script (check `loader.sh` for an example) that makes a
second call to `/keys` with the right request parameters. The successful
response to the second request is the SSH key. Different authentication
mechanisms may need a tailored loader script as well.

Epiry in hours can be specified with a single integer parameter to the route.
eg: `/3`

`/pubkey` just responds with the public key without authentication.

## K8s deployment

Check the [ci/k8s](https://github.com/cromega/keyguard/tree/master/ci/k8s) folder for an example

## Running it on Cloud Foundry

Build it, put your key in the folder and `cf push`.

You can use an encrypted SSH key if you are not comfortable with pushing your private key somewhere

An example app manifest looks something like this:

```yaml
applications:
- name: keyguard
  memory: 32m
  buildpack: binary_buildpack
  command: ./keyguard
  env:
    KG_PUBLIC_URL: https://key.yourdomain.org
    KG_YUBI_CLIENT_ID: 1234
    KG_YUBI_API_KEY: foobar
```

