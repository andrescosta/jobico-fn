Start-Job -Name ctl -WorkingDirectory C:\users\Andres\projects\go\workflew\ctl\cmd\service\ {go run main.go}
Start-Job -Name queue -WorkingDirectory C:\users\Andres\projects\go\workflew\srv\cmd\queue\ {go run main.go}
Start-Job -Name repo -WorkingDirectory C:\users\Andres\projects\go\workflew\repo\cmd\ {go run main.go}
Start-Job -Name listener -WorkingDirectory C:\users\Andres\projects\go\workflew\srv\cmd\listener\ {go run main.go}
Start-Job -Name recorder -WorkingDirectory C:\users\Andres\projects\go\workflew\recorder\cmd {go run main.go}
Start-Job -Name exec -WorkingDirectory C:\users\Andres\projects\go\workflew\srv\cmd\executor\ {go run main.go}
