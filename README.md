# cert-fetcher

The goal of the project is to provide up-to-date TLS certificates for ingresses in all namespaces in a kubernetes cluster.

The program runs inside the cluster as a cronjob, takes the certificate from Vault and writes its data to the kubernetes secret in the specified namespaces.

In order not to update the data unnecessarily, a hash sum is calculated from the concatenation of the certificate and the private key and written to the itlabs.io.tls.hash.sha1 secret annotation.

Documentation

[**Examples**](https://github.com/itlabsio/cert-fetcher/deploy)
