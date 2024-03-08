certutil -enterprise -f -v -AddStore "Root" ..\k8s\certs\ctl.crt
certutil -enterprise -f -v -AddStore "Root" ..\k8s\certs\repo.crt
certutil -enterprise -f -v -AddStore "Root" ..\k8s\certs\recorder.crt
certutil -enterprise -f -v -AddStore "Root" ..\k8s\certs\listener.crt
certutil -enterprise -f -v -AddStore "Root" ..\k8s\certs\queue.crt