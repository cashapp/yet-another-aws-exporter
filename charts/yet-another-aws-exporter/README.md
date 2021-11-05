# yaae

![Version: 0.0.2](https://img.shields.io/badge/Version-0.0.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.0.1](https://img.shields.io/badge/AppVersion-v0.0.1-informational?style=flat-square)

A Helm chart for deploying Yet Another AWS Exporter on Kubernetes

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| disabledScrapers | list | `[]` | A list of scraper IDs to disable in the exporter |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"cashapp/yet-another-aws-exporter"` |  |
| image.tag | string | `""` |  |
| imagePullSecrets | list | `[]` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| prometheus.enabled | bool | `true` |  |
| replicaCount | int | `1` |  |
| resources.limits.cpu | string | `"100m"` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"128Mi"` |  |
| scrapeInterval | string | `"300s"` | Specify the scraper interval for Prometheus  |
| securityContext | object | `{}` |  |
| serviceAccount.annotations."eks.amazonaws.com/role-arn" | string | `""` | Add the full role arn if using OIDC for pod role assumption |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| tolerations | list | `[]` |  |

