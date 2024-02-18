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

function jobico::copy_env_files() {
    $targets = jobico::all_targets

    foreach ($target in $targets) {
		$pkg=(Get-Item $target).Basename
		Copy-Item -Path .\$target\.env -Destination .\bin\.env.$pkg
	}
}

jobico::copy_env_files