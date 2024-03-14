$namespace = "jobico"

function Check-Ingresses {
    $ingresses = kubectl get ing -n $namespace -o json | ConvertFrom-Json
    foreach ($ingress in $ingresses.items) {
        if (-not $ingress.status.loadBalancer.ingress) {
            return $false
        }
    }
    return $true
}

while (-not (Check-Ingresses)) {
    Write-Host "Waiting for ingresses to have an IP address..."
    Start-Sleep -Seconds 10
}

Write-Host "All ingresses have an IP address."