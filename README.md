# Yet Another AWS Exporter Helm Repo

To install this exporter via Helm:

```
helm repo add yaae https://cashapp.github.io/yet-another-aws-exporter
helm repo update

helm install yaae cashapp/yaae -n <NAMESPACE>
```

