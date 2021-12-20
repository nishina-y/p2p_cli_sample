# pion webrtcを使った動画配信のサンプル実装
pion webrtcを使ったCLIプログラム同士の双方向通信とテキストチャット・動画配信を行うサンプル実装です。

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

**gstreamerのインストール**
```
brew install gstreamer gst-plugins-base gst-plugins-good gst-plugins-bad gst-plugins-ugly
```

**カメラのアクセス許可**  
ターミナルで以下のコマンドを実行してgstreamerでカメラが起動する事を確認及びアクセス許可設定を行って下さい
```
gst-launch-1.0 avfvideosrc ! autovideosink
```

**1つ目のクライントを起動**
```
git clone git@github.com:nishina-y/p2p_cli_sample.git
cd p2p_cli_sample/video_communication_sample
go build
./video_communication_sample --addr {シグナリングサーバーのIP}:8080 --mode answer -video-src 'autovideosrc ! videoconvert'
```

**2つ目のクライントを起動**
```
./video_communication_sample --addr {シグナリングサーバーのIP}:8080 --mode answer -video-src 'autovideosrc ! videoconvert'
```

## 動作確認
双方のPCでCLI上でのテキストチャットとカメラの映像配信が行なえます。   
※接続までに数秒時間がかかります
