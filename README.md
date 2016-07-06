# KeyGuard

A little app to serve SSH keys over an authenticated endpoint. A helper script is used to add the key to the SSH agent with an expiry

Only [YubiKey](https://www.yubico.com/why-yubico/for-individuals/) One-time password auth at the moment.

## Usage

### 1. Create configuration

```
$ cat config.json
{
  "SSHKey": "id_rsa", # path to private key
  "LoaderScript": "loader.sh", # path to the loader script
  "PublicUrl": "https://key.yourdomain.org", # public URL where the /key endpoint can be queried
  "Auth": {
    "clientId": "12345", # yubico api credentials
    "apiKey": "apikey",
    "preferHttp": false
  }
}
```

### 2. Build

```
$ go build
```

### 3. Run

```
$ nohup ./keyguard &
```

### 4. Load key!

```
$ curl -s https://key.yourdomain.org | bash
OTP: ccccsfrhkrucdedthkkrdkkrbjdhidjkljktflhvjgcl # this is where I pressed the YubiKey button
Identity added: /tmp/tmp.2GxYjzCLaE (/tmp/tmp.2GxYjzCLaE)
Lifetime set to 32400 seconds
```

### Important

You have to create an API key at [YubiCo](https://upgrade.yubico.com/getapikey/) to use the authenticator.

## How it works

The service exposes two endpoints:
* `/`
* `/key`

`/` responds with a shell script (check `loader.sh` for an example) that makes a second call to `/keys` with the right request parameters. The successful response to the second request is the SSH key. Different authentication mechanisms may need a tailored loader script as well.

## Running it on Cloud Foundry

You can actually run KeyGuard on Cloud Foundry!

Build it, put your key and config.json in the folder and `cf push`. Don't forget to configure `PublicUrl` to the correct route beforehand.

You can use an encrypted SSH key if you are scared of pushing your key to a public cloud.

An example app manifest looks something like this:

```yaml
applications:
- name: keyguard
  memory: 32m
  buildpack: binary_buildpack
  command: ./keyguard --configPath=config.json
```

