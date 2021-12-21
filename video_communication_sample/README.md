# 使い方

## シグナリングサーバーの起動
EC2等の接続する全てのクライアントが接続できる場所で起動して下さい。  
`tcp:8080` ポートを使用するので必要に応じてポートを開放して下さい  
```shell
git clone git@github.com:nishina-y/p2p_cli_sample.git
cd p2p_cli_sample/signaling_server_sample
go build
./signaling_server_sample
```

## TURNサーバーの起動
EC2等の接続する全てのクライアントが接続できる場所で起動して下さい。  
`3478, 49152-65535` ポートを使用するので必要に応じてポートを開放して下さい。  
```shell
git clone https://github.com/pion/turn
cd turn/examples/turn-server/simple
go build
./simple -public-ip {GLOBAL_IP} -users username=password
```

## clientの起動
clientはmacでのみ動作確認を行っています

### 初期設定

**gstreamerのインストール**
```shell
brew install gstreamer gst-plugins-base gst-plugins-good gst-plugins-bad gst-plugins-ugly
```

**カメラのアクセス許可**  
ターミナルから以下のコマンドを実行してmacのセキュリティーとプライバシー設定でターミナルからのカメラとマイクへのアクセス許可を設定して下さい。
```shell
gst-launch-1.0 avfvideosrc ! autovideosink
gst-launch-1.0 osxaudiosrc ! autoaudiosink
```

**1つ目のクライントを起動**
```shell
git clone git@github.com:nishina-y/p2p_cli_sample.git
cd p2p_cli_sample/video_communication_sample
go build
./video_communication_sample \
    --mode answer \
    --addr {シグナリングサーバーのIP}:8080 \
    --mode answer -video-src 'autovideosrc ! videoconvert' \
    --audio-src 'osxaudiosrc ! audioresample ! audio/x-raw, rate=8000'
```

**2つ目のクライントを起動**
```shell
./video_communication_sample \
    --mode offer \
    --addr {シグナリングサーバーのIP}:8080 \
    --mode answer -video-src 'autovideosrc ! videoconvert' \
    --audio-src 'osxaudiosrc ! audioresample ! audio/x-raw, rate=8000'
```

## 動作確認
双方のPCでCLI上でのテキストチャットとカメラの映像配信が行なえます。   
※接続までに数秒時間がかかります
