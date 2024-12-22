# limastart

https://github.com/tama-jp/limastart


```shell
go install github.com/tama-jp/limastart@latest
```

-------------------------------------------

```shell

go mod init github.com/tama-jp/limastart
```

```shell
limastart
```

```shell
````

HomeBrewのインストール

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

limaのインストール

```shell
brew install lima
```



limaのイメージインストール

```shell
limactl start --debug --name=debian_tools --tty=false debian_tools.yaml

```

```shell
limactl start debian_tools
```

```shell
limactl shell debian_tools 
```

```shell
limactl stop debian_tools 
limactl delete debian_tools
```

```shell
go run main.go

go build -o limastart main.go
sudo rm /usr/local/bin/limastart
sudo mv limastart /usr/local/bin/limastart
```
