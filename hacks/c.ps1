mkdir .\k8s\certs
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\ctl.key -out .\k8s\certs\ctl.crt -subj "/CN=ctl/O=ctl"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\repo.key -out .\k8s\certs\repo.crt -subj "/CN=repo/O=repo"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\recorder.key -out .\k8s\certs\recorder.crt -subj "/CN=recorder/O=recorder"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\listener.key -out .\k8s\certs\listener.crt -subj "/CN=listener/O=listener"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\queue.key -out .\k8s\certs\queue.crt -subj "/CN=queue/O=queue"