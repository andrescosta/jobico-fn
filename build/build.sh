#!/usr/bin/env bash

readonly JOBICO_SUPPORTED_PLATFORMS=(
  linux/amd64
  windows/amd64
)

jobico::all_targets() {
  local targets=(recorder/cmd
	)
	echo "${targets[@]}"
}

#targets=(ctl/cmd recorder/cmd)

jobico::debug(){
	local targets=( $(jobico::all_targets) )
	for t in "${targets[@]}"; do
		echo $(basename "${t}")
	done
}


jobico::build(){
	local targets=( $(jobico::all_targets) )
	for t in "${targets[@]}"; do
		go build -o ./bin ../$t
	done
}
# bin=$(basename "${binary}")
# if [[ ${GOOS} == "windows" ]]; then
#    bin="${bin}.exe"
#  fi
#echo "${output_path}/${bin}"
jobico::debug
#jobico::build