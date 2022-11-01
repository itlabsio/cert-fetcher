# cert-fetcher

The goal of the project is to provide up-to-date TLS certificates for ingresses in all namespaces in a kubernetes cluster.

The program runs inside the cluster as a cronjob, takes the certificate from Vault and writes its data to the kubernetes secret in the specified namespaces.

In order not to update the data unnecessarily, a hash sum is calculated from the concatenation of the certificate and the private key and written to the itlabs.io.tls.hash.sha1 secret annotation.

## How to build
```
CGO_ENABLED=0 GOOS=linux go build -v -o refresher ./cmd/refresher/...
```

## How to launch
```
refresher -help
Usage of ./refresher:
  -kubeconfig string
        path to kubeconfig (default "$HOME/.kube/config")
  -vault-addr string
        Vault server address (default "http://127.0.0.1:8200")
  -vault-jwt string
        Vault jwt token
  -vault-path string
        Vault authmethod path if not token (default "kubernetes")
  -vault-role string
        vault role (default "default")
  -vault-token string
        Vault client token
```


Documentation

[**Examples**](https://github.com/itlabsio/cert-fetcher/tree/main/deploy)
