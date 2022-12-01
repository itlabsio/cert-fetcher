# cert-fetcher

The goal of the project is to provide up-to-date TLS certificates for ingresses in all namespaces in a kubernetes cluster.

The program runs inside the cluster as a cronjob, takes the certificate from Vault and writes its data to the kubernetes secret in the specified namespaces.

In order not to update the data unnecessarily, a hash sum is calculated from the concatenation of the certificate and the private key and written to the itlabs.io.tls.hash.sha1 secret annotation.

## Running cert-fetcher locally

- Build binary

```bash
CGO_ENABLED=0 GOOS=linux go build -v -o refresher ./cmd/refresher/...
```

- Run the built binary

```bash
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

## Running cert-fetcher in k8s

Setup steps

- Create service account and RBAC for cert-fetcher

```bash
kubectl apply -f deploy/rbac.yaml
```

- Configure and launch cronjob

```bash
kubectl apply -n default -f deploy/cronjob.yaml
```

- Enable cert-fetcher for the Namespace

```bash
kubectl label namespace default itlabs.io/certfetcher=true
```

- Annotate the Namespace with your path to vault.

```bash
kubectl annotate namespace default "itlabs.io/certfetcher.wildcard-example-com: secret/path/to/secret/in/vault"
```

## Examples 

Manifest [**examples**](https://github.com/itlabsio/cert-fetcher/tree/main/deploy)

## License

Copyright (c) 2017-2022 ITLabs.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
