<p align="center">
  <img src="/doraemon.png" height="100">
  <h2 align="center">
    Just Like Doraemon, Give Me All Tools I Need ~ 
  </h2>
</p>

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

## Features

- [x] SSH server management
- [x] Simple http service for sharing files
- [x] Run custom command

# SSH Server management

run `doraemon`    
![doraemon ssh](/resources/doraemon_ssh.gif)

## config

The configuration file is in directory `~/.doraemon`,the configuration file will be created when the program is run for
the first time.    
you can run `doraemon`, type `openConfigDir` open the configuration directory

### password

connection to the server by username and password.

```
[[sshInfo]]
  name = "pwdservr"
  authMethod = "password"
  uri = "127.0.0.1:222"
  user = "testuser"
  passphrase = "123"
  desc = "sandbox deploy server"
```

### publickey

if your ssh server connection uses publickey

```
[[sshInfo]]
  name = "sandbox1"
  authMethod = "publickey"
  uri = "123.13.63.203:22"
  user = "abc"
  publicKeyPath = "/Users/abc/.ssh/a.pem"
  desc = "my sandbox host1"
```

if your publickey is encrypted, use `passphrase` to specify the ciphertext

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

### proxy

if there is a proxy service, you must first configure the proxy server, other configurations use `proxySSHName` to
specify the proxy server

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

you can use `startCommand` to run custom command after login service

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

# simple http service for sharing files

Share folder via http service    
```
doraemon srv .
```    
![doraemon srv](/resources/doraemon_srv.gif)



# custom command
run `doraemon cmd`    
![doraemon cmd](/resources/doraemon_cmd.gif)

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
