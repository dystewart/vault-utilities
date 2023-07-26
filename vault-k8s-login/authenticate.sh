#!/bin/bash

# Obtain the service account token from the pod
SERVICE_ACCOUNT_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)

# Set the Vault URL and Kubernetes Auth Method endpoint
VAULT_URL="https://nerc-vault.vault.svc.cluster.local"
K8S_AUTH_ENDPOINT="/v1/auth/kubernetes/backup"

# Authenticate to Vault using the Kubernetes Auth Method
RESPONSE=$(curl -s --request POST --data "{\"jwt\": \"$SERVICE_ACCOUNT_TOKEN\", \"role\": \"backup\"}" "$VAULT_URL$K8S_AUTH_ENDPOINT")
VAULT_TOKEN=$(echo "$RESPONSE" | jq -r '.auth.client_token')

# Now you can use the VAULT_TOKEN to interact with Vault API as an authenticated user
# For example: curl -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_URL/v1/secret/data/<secret-name>"

# Your additional commands using the authenticated VAULT_TOKEN go here...
