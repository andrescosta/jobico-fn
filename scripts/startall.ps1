Start-Job -Name ctl -WorkingDirectory C:\users\Andres\projects\go\jobico\bin {.\ctl.exe}
Start-Job -Name queue -WorkingDirectory C:\users\Andres\projects\go\jobico\bin {.\queue.exe}
Start-Job -Name repo -WorkingDirectory C:\users\Andres\projects\go\jobico\bin {.\repo.exe}
Start-Job -Name recorder -WorkingDirectory C:\users\Andres\projects\go\jobico\bin {.\recorder.exe}
Start-Job -Name listener -WorkingDirectory C:\users\Andres\projects\go\jobico\bin {.\listener.exe}
Start-Job -Name exec -WorkingDirectory C:\users\Andres\projects\go\jobico\bin {.\executor.exe}
