#!/bin/bash

namespace="jobico"

# Function to check if all ingresses have an IP address
function check_ingresses {
    local ingresses
    ingresses=$(kubectl get ing -n "$namespace" -o jsonpath='{.items[*].status.loadBalancer.ingress[*].hostname}')

    if [[ -z "$ingresses" ]]; then
        return 1
    fi

    return 0
}

# Loop until all ingresses have an IP address
while ! check_ingresses; do
    echo "Waiting for ingresses to have an IP address..."
    sleep 10
done

echo "All ingresses have an IP address."