new-Item -ItemType SymbolicLink -Path . -name stop.ps1 -Value ..\scripts\stopall.ps1 | Out-Null
New-Item -ItemType SymbolicLink -Path . -name status.ps1 -Value ..\scripts\status.ps1 | Out-Null
New-Item -ItemType SymbolicLink -Path . -name startall.ps1 -Value ..\scripts\startall.ps1 | Out-Null

@'
.\stop.ps1
rd db -Force -Recurse -ErrorAction Ignore
rd cache -Force -Recurse -ErrorAction Ignore
rd files -Force -Recurse -ErrorAction Ignore
rd log -Force -Recurse -ErrorAction Ignore
rd data -Force -Recurse -ErrorAction Ignore
del log.log
'@ | Out-File -FilePath ".\reset.ps1"

@'
.\startall.ps1
.\status.ps1
'@ | Out-File -FilePath ".\start.ps1"

@'
.\reset.ps1
.\start.ps1
'@ | Out-File -FilePath ".\restart.ps1"

Write-Output "Enviroment created."