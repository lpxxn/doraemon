<p align="center">
  <img src="/doraemon.png" height="100">
  <h1 align="center">
    Just Like Doraemon, Give Me All Tools I Need ~ 
  </h1>
</p>

developing...
# SSH Server management

![doraemon command](/doraemon.gif)

## install from code
```
mkdir -p $GOPATH/src/github.com/lpxxn/doraemon
cd $GOPATH/src/github.com/lpxxn/doraemon
git clone https://github.com/lpxxn/doraemon.git 
cd cmd/doraemon
export GO111MODULE=on
make install
```

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

if your publickey is encrypted you need `passphrase`

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
  passphrase = "123"
  desc = "sandbox deploy server"
```

### proxy

if you have proxy host

add a proxy

```
[[sshInfo]]
  name = "proxy"
  authMethod = "publickey"
  uri = "192.168.1.7:22"
  user = "lipeng"
  publicKeyPath = "/Users/abc/.ssh/my_test.pem"
  desc = "proxy"
```

use `proxySSHName` to specify `proxy`

```
[[sshInfo]]
  name = "my_sandbox1"
  authMethod = "publickey"
  uri = "192.2.0.8:22"
  user = "lipeng"
  publicKeyPath = "/Users/li/.ssh/test.pem"
  proxySSHName = "proxy"
  desc = "my sandbox host 1"
```