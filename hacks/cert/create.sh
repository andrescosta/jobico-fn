mkdir ./k8s/certs
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/ctl.key -out ./k8s/certs/ctl.crt -subj "/CN=ctl/O=ctl" -addext "subjectAltName = DNS:ctl"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/repo.key -out ./k8s/certs/repo.crt -subj "/CN=repo/O=repo"  -addext "subjectAltName = DNS:repo"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/recorder.key -out ./k8s/certs/recorder.crt -subj "/CN=recorder/O=recorder"  -addext "subjectAltName = DNS:recorder"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/listener.key -out ./k8s/certs/listener.crt -subj "/CN=listener/O=listener"  -addext "subjectAltName = DNS:listener"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/queue.key -out ./k8s/certs/queue.crt -subj "/CN=queue/O=queue"  -addext "subjectAltName = DNS:queue"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/jaeger.key -out ./k8s/certs/jaeger.crt -subj "/CN=jaeger/O=jaeger/" -addext "subjectAltName=DNS:jaeger"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/prometheus.key -out ./k8s/certs/prometheus.crt -subj "/CN=prometheus/O=prometheus" -addext "subjectAltName = DNS:prometheus"
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout ./k8s/certs/exec.key -out ./k8s/certs/exec.crt -subj "/CN=exec/O=exec"  -addext "subjectAltName = DNS:exec"
