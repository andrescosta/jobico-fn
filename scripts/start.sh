jobico::all_server_files_windows() {
  local targets=(ctl
				 recorder
				 queue
				 repo
				 executor
				 listener
  )
  echo "${targets[@]}"
}

jobico::startall(){
    local files=( $(jobico::all_server_files_windows) )
	for t in "${files[@]}"; do
		./bin/${t} --env:basedir=./bin --env:workdir=./work &
	done
}

jobico::startall &