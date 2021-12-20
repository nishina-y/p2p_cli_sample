# 使い方

## シグナリングサーバーの起動
EC2等の接続する全てのクライアントが接続できる場所で起動して下さい。  
`tcp:8080` ポートを使用するので必要に応じてポートを開放して下さい  
```
git clone git@github.com:nishina-y/p2p_cli_sample.git
cd p2p_cli_sample/signaling_server_sample
go build
./signaling_server_sample
```

## TURNサーバーの起動
EC2等の接続する全てのクライアントが接続できる場所で起動して下さい。  
`3478, 49152-65535` ポートを使用するので必要に応じてポートを開放して下さい。  
```
git clone https://github.com/pion/turn
cd turn/examples/turn-server/simple
go build
./simple -public-ip {GLOBAL_IP} -users username=password
```

## clientの起動
clientはmacでのみ動作確認を行っています

### 初期設定

**1つ目のクライントを起動**
```
git clone git@github.com:nishina-y/p2p_cli_sample.git
cd p2p_cli_sample/text_chat_sample

go build
./text_chat_sample --addr {シグナリングサーバーのIP}:8080 --mode answer
```

**2つ目のクライントを起動**
```
./text_chat_sample --addr {シグナリングサーバーのIP}:8080 --mode offer
```

## 動作確認
双方のPCでCLI上でのテキストチャットが行なえます。   
※接続までに数秒時間がかかります
