function jobico::all_targets() {
    $targets = "cmd\cli",
            "cmd\ctl",
            "cmd\dashboard",
            "cmd\executor",
            "cmd\listener",
            "cmd\queue",
            "cmd\recorder",
            "cmd\repo"
    return $targets
}

function jobico::build() {
    $targets = jobico::all_targets

    foreach ($target in $targets) {
        go build -o .\bin .\$target
    }
}

jobico::build