package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"

    "github.com/hashicorp/vault/api"
)

func main() {

    // serviceAccount token path 
    tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

    // Read the service account token from the file
    token, err := ioutil.ReadFile(tokenPath)
    if err != nil {
        log.Fatalf("Failed to read token file: %s", err)
    }

    vault_addr := "http://nerc-vault.vault.svc:8200"
    client, err := createVaultClient(vault_addr)

    fmt.Println("Token: %v", token)
    fmt.Println("tokenPath: %v", tokenPath)
    fmt.Println("vault addr: %v", client.Address())

}

func createVaultClient(vault_addr string) (*api.Client, error) {
    // Create a new Vault configuration with the specified address
    config := &api.Config{
    Address: vault_addr,
    }

    // Create a new Vault client using the configuration
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("Failed to create Vault client: %v", err)
    }

    return client, nil
}


