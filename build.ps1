Invoke-WebRequest -Uri https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml -OutFile LiveSplit.AutoSplitters.xml -v
Invoke-WebRequest -Uri https://raw.githubusercontent.com/ShootMe/LiveSplit.HollowKnight/master/Components/LiveSplit.HollowKnight.dll -OutFile LiveSplit.HollowKnight.dll -v
wails build
