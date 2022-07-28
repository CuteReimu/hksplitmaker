@echo off
bitsadmin /transfer n https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml %~dp0\LiveSplit.AutoSplitters.xml
bitsadmin /transfer n https://raw.githubusercontent.com/ShootMe/LiveSplit.HollowKnight/master/Components/LiveSplit.HollowKnight.dll %~dp0\LiveSplit.HollowKnight.dll
go build -ldflags "-s -w -H=windowsgui" -o hksplitmaker.exe github.com/CuteReimu/hksplitmaker
