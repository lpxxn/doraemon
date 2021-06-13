<p align="center">
  <img src="/doraemon.png" height="100">
  <h1 align="center">
    Just Like Doraemon, Give Me All Tools I Need ~ 
  </h1>
</p>

developing...

# install 

## install from code
```
mkdir -p $GOPATH/src/github.com/lpxxn/doraemon
cd $GOPATH/src/github.com/lpxxn/doraemon
git clone https://github.com/lpxxn/doraemon.git 
cd cmd/doraemon
export GO111MODULE=on
make install
```

## install from go get

```
GO111MODULE=on go get -u github.com/lpxxn/doraemon/cmd/doraemon
```

### help
```
doraemon -h
```
```
type exit or :q or \q to exit app

ssh manager and .....

Usage:
  doraemon [flags]
  doraemon [command]

Available Commands:
  cmd         custom cmd
  help        Help about any command

Flags:
  -h, --help      help for doraemon
  -l, --loopRun   not exist until type :q or \q

Use "doraemon [command] --help" for more information about a command.
bye ~ 👋👋👋
```

## Features

- [x] SSH server management
- [x] Run custom command

# SSH Server management
run `doraemon`
![doraemon ssh](/doraemon_ssh.gif)

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

if there is a proxy service, you must first configure the proxy server, other configurations use `proxySSHName` to specify the proxy server

```
[[sshInfo]]
  name = "proxy"
  authMethod = "publickey"
  uri = "192.168.1.7:22"
  user = "lipeng"
  publicKeyPath = "/Users/abc/.ssh/my_test.pem"
  desc = "proxy"
```

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
## run command after login service
you can use `startCommand` run command after login service
```
[[sshInfo]]
  name = "sandbox1"
  authMethod = "publickey"
  uri = "123.13.63.203:22"
  user = "abc"
  publicKeyPath = "/Users/abc/.ssh/a.pem"
  desc = "my sandbox host1"
  startCommand = "cd /home/abc/app/"

[[sshInfo]]
  name = "sandbox1"
  authMethod = "publickey"
  uri = "123.13.63.203:22"
  user = "abc"
  publicKeyPath = "/Users/abc/.ssh/a.pem"
  desc = "my sandbox host1"
  startCommand = "tmux attach -t abc"
  
```

# custom command
run `doraemon cmd`
![doraemon cmd](/doraemon_cmd.gif)

manage commonly used command
```
[[cmdInfo]]
  name = "cd_test"
  cmd = """ 
  cd /; pwd; ls -al;
  """
  desc = "test command"
```
you can run `doraemon cmd`
