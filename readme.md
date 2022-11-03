# Go wallpaper


### 功用
使用colly及cron套件實作[台灣觀光局桌布](https://www.taiwan.net.tw/m1.aspx?sNo=0012076)將每月的桌布下載到Windows儲存的照片

### 執行
```=bash
go run main.go
```


### Windows背景服務執行
```=bash
go build -o wallpaper.exe
```

使用[nssm](https://nssm.cc/)將```wallpaper.exe```加入執行
