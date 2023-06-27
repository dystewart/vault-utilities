package main

import (
    "fmt"
    "os"

    "github.com/hashicorp/vault/api"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Create the Vault client
    client, err := createVaultClient()
    if err != nil {
        fmt.Printf("Failed to create Vault client: %v\n", err)
        return
    }

    // Login with Kubernetes auth method
    err = loginWithKubernetes(client)
    if err != nil {
        fmt.Printf("Failed to login with Kubernetes auth method: %v\n", err)
        return
    }

    // Use the Vault client for further operations
    // ...
}

func createVaultClient() (*api.Client, error) {
    // Create a new Vault client
    client, err := api.NewClient(api.DefaultConfig())
    if err != nil {
        return nil, err
    }

    // Set the Vault server address
    client.SetAddress("http://your-vault-address:8200")

    return client, nil
}

func loginWithKubernetes(client *api.Client) error {
    // Authenticate with Vault using the Kubernetes auth method
    authPath := "auth/kubernetes/login"
    authData := map[string]interface{}{
        "jwt":   getKubernetesToken(),
        "role":  "your-kubernetes-role",
        "mount": "your-kubernetes-auth-mount",
    }

    // Perform the login request
    response, err := client.Logical().Write(authPath, authData)
    if err != nil {
        return err
    }

    // Check the login response and retrieve the Vault token if successful
    if response.Auth == nil || response.Auth.ClientToken == "" {
        return fmt.Errorf("Failed to authenticate with Kubernetes auth method")
    }

    // Set the obtained Vault token in the client for subsequent operations
    client.SetToken(response.Auth.ClientToken)

    return nil
}

func getKubernetesToken() string {
    // Create a Kubernetes client
    clientset, err := createKubernetesClient()
    if err != nil {
        fmt.Printf("Failed to create Kubernetes client: %v\n", err)
        return ""
    }

    // Retrieve the service account token
    token, err := getServiceAccountToken(clientset)
    if err != nil {
        fmt.Printf("Failed to retrieve service account token: %v\n", err)
        return ""
    }

    return token
}

func createKubernetesClient() (*kubernetes.Clientset, error) {
    // Get the path to the kubeconfig file
    kubeconfigPath := os.Getenv("KUBECONFIG")
    if kubeconfigPath == "" {
        kubeconfigPath = clientcmd.RecommendedHomeFile
    }

    // Create a Kubernetes client from the kubeconfig file
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
    }

    return clientset, nil
}

func getServiceAccountToken(clientset *kubernetes.Clientset) (string, error) {
    // Get the default service account token in the current namespace
    podName := os.Getenv("POD_NAME")
    namespace := os.Getenv("POD_NAMESPACE")

    pod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to get pod: %v", err)
    }

    // Retrieve the service account token from the pod's mounted service account token volume
    volumes := pod.Spec.Volumes
    for _, volume := range volumes {
        if volume.Secret != nil && volume.Secret.SecretName == "default-token" {
            secret, err := clientset.CoreV1().Secrets(namespace).Get(volume.Secret.SecretName, metav1.GetOptions{})
            if err != nil {
                return "", fmt.Errorf("failed to get secret: %v", err)
            }

            tokenBytes, ok := secret.Data["token"]
            if !ok {
                return "", fmt.Errorf("service account token not found in secret")
            }

            token := string(tokenBytes)
            return token, nil
        }
    }

    return "", fmt.Errorf("service account token volume not found")
}