<p align="center">
  <img src="/doraemon.png" height="100">
  <h1 align="center">
    Just Like Doraemon, Give Me All Tools I Need ~ 
  </h1>
</p>

developing...

## config

config directory `~/.doraemon`,the program will create the config file when first run it.

### publickey

```
[[sshInfo]]
  name = "sandbox1"
  authMethod = "publickey"
  uri = "123.13.63.203:22"
  user = "abc"
  publicKeyPath = "/Users/abc/.ssh/a.pem"
  desc = "my sandbox host1"
```

if your publickey is encrypted
you need `passphrase`
```
[[sshInfo]]
  name = "gateway"
  authMethod = "publickey"
  uri = "123.13.63.203:22"
  user = "abc"
  passphrase = "dsafDFl209Dfoifqw"
  publicKeyPath = "/Users/abc/.ssh/test.pem"
  desc = "gateway jumper
```

### password

```
[[sshInfo]]
  name = "pwdservr"
  authMethod = "password"
  uri = "127.0.0.1:222"
  user = "testuser"
  startCommand = "pwd;"
  passphrase = "123"
  desc = "sandbox deploy server"
```