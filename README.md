# GCP Workload Identity

## Create GCP Infrastructure and Manifests

```
./ccoctl-4.11.7 gcp create-all \
    --name=abutcher-gcp \
    --region=us-central1 \
    --project=openshift-hive-dev \
    --credentials-requests-dir=./4.11.7-credrequests \
    --output-dir=_output
```

### Create OpenShift Cluster With `ccoctl` Manifests

1. Create `install-config.yaml` for the cluster.

```
./openshift-install-4.11.7 create install-config
? SSH Public Key /home/abutcher/.ssh/id_rsa.pub
? Platform gcp
? Region us-central1
? Base Domain gcp.hive.openshift.com
? Cluster Name abutchertest
? Pull Secret [? for help] 
INFO Install-Config created in: .
```

2. Specify `credentialsMode: Manual` within the `install-config.yaml`

```
echo "credentialsMode: Manual" >> install-config.yaml
```

3. Create manifests for the cluster and copy in our `ccoctl` manifests.

```
./openshift-install-4.11.7 create manifests

cp _output/manifests/* ./manifests/

cp -a _output/tls .
```

4. Install cluster.

```
./openshift-install-4.11.7 create cluster
```

### Create Additional GCP Service Account For Our Workload

This will output a secret manifest to `manifests/open-cluster-management-observability-cloud-credentials-credentials.yaml`.

```
./ccoctl-4.11.7 gcp create-service-accounts \
    --credentials-requests-dir obs-credrequest \
    --project openshift-hive-dev \
    --workload-identity-pool abutcher-gcp \
    --workload-identity-provider abutcher-gcp \
    --name abutcher-gcp
```

### Deploy Workload

```
oc create namespace open-cluster-management-observability

# Apply credentials secret
oc apply -f manifests/open-cluster-management-observability-cloud-credentials-credentials.yaml

# Apply service account
oc apply -f serviceaccount.yaml

# Apply deployment
oc apply -f deployment.yaml
```

### Observe Workload

It lists the contents of the OIDC bucket as specified by `BUCKET_NAME` on the deployment.

```
oc logs bucket-847fb78f87-84cv4 -n open-cluster-management-observability
2023/04/07 13:28:41 [.well-known/openid-configuration keys.json]
2023/04/07 13:28:56 [.well-known/openid-configuration keys.json]
2023/04/07 13:29:11 [.well-known/openid-configuration keys.json]
2023/04/07 13:29:26 [.well-known/openid-configuration keys.json]
2023/04/07 13:29:41 [.well-known/openid-configuration keys.json]
2023/04/07 13:29:57 [.well-known/openid-configuration keys.json]
2023/04/07 13:30:12 [.well-known/openid-configuration keys.json]
2023/04/07 13:30:27 [.well-known/openid-configuration keys.json]
2023/04/07 13:30:42 [.well-known/openid-configuration keys.json]
2023/04/07 13:30:57 [.well-known/openid-configuration keys.json]
```