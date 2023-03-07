# Go wallpaper


### 功用
使用colly及cron套件實作[台灣觀光局桌布](https://www.taiwan.net.tw/m1.aspx?sNo=0012076)將每月的桌布下載到Windows儲存的照片

### 執行
1. 確認下載路徑及資料夾存在 (預設`./pictures`)

2. 執行cron
```=bash
go run main.go
```

3. 測試下載
```=bash
go run main.go -run
```


### Windows背景服務執行 
有兩種方法

#### Docker (推薦)

> 須安裝[Docker](https://docs.docker.com/desktop/install/windows-install/)
```
docker-compose up -d
```
#### NSSM
```=bash
go build -o wallpaper.exe
```

使用[nssm](https://nssm.cc/)將```wallpaper.exe```加入執行
