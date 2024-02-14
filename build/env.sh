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

jobico::copy_env_files(){
	local targets=( $(jobico::all_targets) )
	for t in "${targets[@]}"; do
		pkg=$(basename "${t}")
		echo "coping ../$t/.env to ./bin/.env.$pkg" 
		cp ./$t/.env ./bin/.env.$pkg
	done
}

jobico::copy_env_files