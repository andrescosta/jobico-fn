jobico::all_server_files_windows() {
  local targets=(ctl
				 executor
				 listener
				 queue
				 recorder
				 repo
  )
	echo "${targets[@]}"
}

jobico::killall(){
    local pids=( $(ps | grep ../bin/ | awk '{print $1}') )
	for t in "${pids[@]}"; do
		kill ${t}
	done
}

jobico::killall