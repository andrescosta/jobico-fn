jobico::all_server_files_windows() {
  local targets=(ctl.exe
				 executor.exe
				 listener.exe
				 queue.exe
				 recorder.exe
				 repo.exe
  )
	echo "${targets[@]}"
}

jobico::killall(){
    local pids=( $(ps | grep jobico/bin/ | awk '{print $1}') )
	for t in "${pids[@]}"; do
		kill ${t}
	done
}

jobico::killall