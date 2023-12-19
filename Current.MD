# Deployment and build scripts

## Build script

## Run all

## Create docker files and push them to a registry

## Run lints (golangci)

## Run vulnerability report -> 
https://go.dev/blog/vuln 
https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck#hdr-Usage

## Run tests (future)
 
### References: 
- https://www.gnu.org/software/make/manual/html_node/index.html
- https://github.com/kubernetes/kubernetes/tree/master/hack




kube::golang::server_targets() {
  local targets=(
    cmd/kube-proxy
    cmd/kube-apiserver
    cmd/kube-controller-manager
    cmd/kubelet
    cmd/kubeadm
    cmd/kube-scheduler
    vendor/k8s.io/component-base/logs/kube-log-runner
    vendor/k8s.io/kube-aggregator
    vendor/k8s.io/apiextensions-apiserver
    cluster/gce/gci/mounter
  )