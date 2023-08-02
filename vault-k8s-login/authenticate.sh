#!/bin/bash

export VAULT_ADDR="http://nerc-vault:8200"
export VAULT_TOKEN=$(vault write -field=token auth/kubernetes/backup/login role=backup jwt=$(cat /run/secrets/kubernetes.io/serviceaccount/token))
