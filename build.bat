@echo off
go build -ldflags "-s -w -H=windowsgui" -o hksplitmaker.exe github.com/CuteReimu/hksplitmaker
