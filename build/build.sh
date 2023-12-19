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
		go build -o ./bin ../$t
	done
}

jobico::copy_env_files(){
	local targets=( $(jobico::all_targets) )
	for t in "${targets[@]}"; do
		pkg=$(basename "${t}")
		echo "coping ../$t/.env to ./bin/.env.$pkg" 
		cp ../$t/.env ./bin/.env.$pkg
	done
}

jobico::run_lints(){
	echo "Running lints, check issues file ..."
	golangci-lint run ../... > ../issues
}

jobico::deploy_locally(){
	jobico::build
	jobico::copy_env_files
}

jobico::release_locally(){
	jobico::run_lints
	jobico::deploy_locally
}

jobico::deploy_locally