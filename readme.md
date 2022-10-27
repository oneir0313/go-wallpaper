# Go wallpaper


### 功用
台灣觀光局桌布將每月的桌布下載到Windows儲存的照片

### 執行
```=bash
go run main.go
```


### Windows背景服務執行
```=bash
go build -o wallpaper.exe
```

使用[nssm](https://nssm.cc/)將```.exe```加入執行