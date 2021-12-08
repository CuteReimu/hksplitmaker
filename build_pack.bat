@echo off
go build -ldflags -H=windowsgui -o hksplitmaker.exe hksplitmaker || exit
set zipname=计时器生成器.zip
if exist %zipname% (
    del /F %zipname%
)
C:\Progra~1\WinRAR\Rar.exe a -r %zipname% hk-split-maker\src\asset\categories hk-split-maker\src\asset\icons hk-split-maker\src\asset\splits.txt LICENSE README.md hksplitmaker.exe translate.xlsx
