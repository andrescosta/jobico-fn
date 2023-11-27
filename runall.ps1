Start-Job -Name ctl -WorkingDirectory C:\users\Andres\projects\go\workflew\ctl\cmd\service\ {go run service.go}
Start-Job -Name queue -WorkingDirectory C:\users\Andres\projects\go\workflew\srv\cmd\queue\ {go run queue.go}
Start-Job -Name repo -WorkingDirectory C:\users\Andres\projects\go\workflew\repo\cmd\ {go run repo.go}
Start-Job -Name listener -WorkingDirectory C:\users\Andres\projects\go\workflew\srv\cmd\listener\ {go run listener.go}
Start-Job -Name recorder -WorkingDirectory C:\users\Andres\projects\go\workflew\recorder\cmd {go run recorder.go}
Start-Job -Name exec -WorkingDirectory C:\users\Andres\projects\go\workflew\srv\cmd\executor\ {go run executor.go}