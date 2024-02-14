#!/usr/bin/env bash

readonly JOBICO_SUPPORTED_PLATFORMS=(
  linux/amd64
  windows/amd64
)

jobico::all_targets() {
  local targets=(cmd/cli
				 cmd/ctl
				 cmd/dashboard
				 cmd/executor
				 cmd/listener
				 cmd/queue
				 cmd/recorder
				 cmd/repo
  )
	echo "${targets[@]}"
}

jobico::debug(){
	local targets=( $(jobico::all_targets) )
	for t in "${targets[@]}"; do
		echo $(basename "${t}")
	done
}


jobico::build(){
	local targets=( $(jobico::all_targets) )
	for t in "${targets[@]}"; do
		echo "Building $t ..."
		go build -o ./bin ./$t
	done
}

jobico::build