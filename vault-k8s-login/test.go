package main

import (
    "fmt"
    "os"

    "github.com/hashicorp/vault/api"
)

func main() {
    // Vault address
    vaultAddr := "http://nerc-vault:8200"

     // Create a new Vault client
    client, err := createVaultClient(vaultAddr)
    if err != nil {
        log.Fatalf("Failed to create Vault client: %v", err)
    }   

    // The service-account token will be read from the path where the token's
    // Kubernetes Secret is mounted. By default, Kubernetes will mount it to
    // /var/run/secrets/kubernetes.io/serviceaccount/token, but an administrator
    // may have configured it to be mounted elsewhere.
    k8sAuth, err := auth.NewKubernetesAuth(
        "backup",
        auth.WithMountPath("/auth/kubernetes/backup"),
        auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
    )
    if err != nil {
        return "", fmt.Errorf("unable to initialize Kubernetes auth method: %w", err)
    }

    authInfo, err := client.Auth().Login(context.TODO(), k8sAuth)
    if err != nil {
        return "", fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
    }
    if authInfo == nil {
        return "", fmt.Errorf("no auth info was returned after login")
    }
    
    // get secret from Vault, from the default mount path for KV v2 in dev mode, "secret"
    secret, err := client.KVv2("nerc").Get(context.Background(), "bmc_credentials")
    if err != nil {
        return "", fmt.Errorf("unable to read secret: %w", err)
    }

    // data map can contain more than one key-value pair,
    // in this case we're just grabbing one of them
    value, ok := secret.Data["username"].(string)
    if !ok {
        return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["username"], secret.Data["username"])
    }

    return value
}

func createVaultClient(vaultAddr string) (*api.Client, error) {
    // Create a new Vault configuration with the specified address
    config := &api.Config{
        Address: vaultAddr,
    }

    // Create a new Vault client using the configuration
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("Failed to create Vault client: %v", err)
    }

    return client, nil
}

