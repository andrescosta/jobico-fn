jobico::services() {
  local targets=(ctl
		 recorder
		 queue
		 repo
		 executor
		 listener
  )
  echo "${targets[@]}"
}
ROOT=$(dirname "${BASH_SOURCE[0]}")/..

jobico::startall(){
    local files=( $(jobico::services) )
    for t in "${files[@]}"; do
       $ROOT/bin/${t} --env:basedir=$ROOT/bin --env:workdir=$ROOT/work &
    done
}

jobico::startall&
