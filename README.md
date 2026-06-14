# 空洞骑士计时器生成器（中文版）

![](https://img.shields.io/github/go-mod/go-version/CuteReimu/hksplitmaker "Language")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/hksplitmaker/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/hksplitmaker/actions/workflows/golangci-lint.yml "Analysis")
[![](https://img.shields.io/github/license/CuteReimu/hksplitmaker)](https://github.com/CuteReimu/hksplitmaker/blob/master/LICENSE "LICENSE")

## 如何使用

https://cutereimu.cn/daily/hollowknight/hksplitmaker-faq.html

## 编译说明

首先需要 Go 和 Nodejs，然后安装 wails：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

然后使用以下命令就可以调试或打包了：

```bash
# 提前下载LiveSplit计时器插件
curl -O https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml
curl -O https://raw.githubusercontent.com/ShootMe/LiveSplit.HollowKnight/master/Components/LiveSplit.HollowKnight.dll

# 本地调试
wails dev

# 打包
wails build -platform=windows/amd64 -webview2 embed
```

## 特别鸣谢

空洞骑士LiveSplit计时器插件：https://github.com/ShootMe/LiveSplit.HollowKnight

本程序的所有**非代码部分**（图标和模板）全部都来自：https://github.com/slaurent22/hk-split-maker
