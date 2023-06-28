package main

import (
    "fmt"
    "io/ioutil"
    "log"

    "github.com/hashicorp/vault/api"
)

func main() {
    // Vault address
    vaultAddr := "http://nerc-vault.vault.svc:8200"

    // Create a new Vault client
    client, err := createVaultClient(vaultAddr)
    if err != nil {
        log.Fatalf("Failed to create Vault client: %v", err)
    }

    // Retrieve the service account token from the file
    tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
    token, err := ioutil.ReadFile(tokenPath)
    if err != nil {
        log.Fatalf("Failed to read token file: %s", err)
    }

    // Authenticate with Vault using the service account token
    secret, err := authenticateWithVault(client, string(token))
    if err != nil {
      log.Fatalf("Failed to authenticate with Vault: %v", err)
    }

    // Obtain the Vault token from the authentication response
    vaultToken, ok := secret.Auth.ClientToken()
    if !ok {
      log.Fatal("Failed to retrieve Vault token")
    }

    // Set the Vault token in the Vault client
    client.SetToken(vaultToken)

    // Perform a sample operation to validate the login
    secretData, err := readSecret(client, "nerc/nerc-ocp-test/postgres")
    if err != nil {
        log.Fatalf("Failed to read secret: %v", err)
    }

    fmt.Println("Secret: ", secretData)
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

func authenticateWithVault(client *api.Client, token string) (*api.Secret, error) {
    // Authenticate with Vault using the service account token
    authPath := "auth/kubernetes/backup/"
    authData := map[string]interface{}{
        "role":       "backup",
        "jwt":        token,
        "kubernetes": "true",
    }

    secret, err := client.Logical().Write(authPath, authData)
    if err != nil {
        return nil, fmt.Errorf("Failed to authenticate with Vault: %v", err)
    }

    return secret, nil
}

func readSecret(client *api.Client, path string) (map[string]interface{}, error) {
    // Read a secret from Vault
    secret, err := client.Logical().Read(path)
    if err != nil {
        return nil, fmt.Errorf("Failed to read secret: %v", err)
    }

    if secret == nil || secret.Data == nil {
        return nil, fmt.Errorf("Secret not found")
    }

    return secret.Data, nil
}