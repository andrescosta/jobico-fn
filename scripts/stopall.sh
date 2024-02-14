jobico::killall(){
    local pids=( $(ps f | grep ../bin/ | awk '{print $1}') )
    for t in "${pids[@]}"; do
       kill ${t}
    done
}

jobico::killall
