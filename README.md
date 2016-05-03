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
