# Quickstart

TOC:
- [Prerequisities](#prerequisities)
- [As Platform Engineer, set up Kind cluster with Radius](#as-platform-engineer-set-up-kind-cluster-with-radius)
- [As Developer, describe your Workload with a Score file](#as-developer-describe-your-workload-with-a-score-file)
- [In CI/CD, generate Radius's app.bicep from Score](#in-cicd-generate-radiuss-appbicep-from-score)
- [In CI/CD, deploy the generated Radius's app.bicep to Radius](#in-cicd-deploy-the-generated-radiuss-appbicep-to-radius)
- [See the application in the Radius dashboard](#see-the-application-in-the-radius-dashboard)
- [As Developer, deploy my Workload locally with Docker Compose](#as-developer-deploy-my-workload-locally-with-docker-compose)

## Prerequisities

- Have Docker installed locally
- Have `kind` CLi installed locally
- Have `rad` CLi installed locally
- Have `score-radius` CLi installed locally

## As Platform Engineer, set up Kind cluster with Radius

```bash
./scripts/setup-kind-cluster.sh
```

```bash
rad workspace create kubernetes default
rad group create default --workspace default
rad env create default --group default
```

## As Developer, describe your Workload with a Score file

First, create a dedicated folder where you will have all the files needed to successfully deploy your application to Radius:
```bash
mkdir podinfo
cd podinfo
```

Create a `score.yaml` with this content:
```yaml
apiVersion: score.dev/v1b1
metadata:
  name: podinfo
service:
  ports:
    http:
      port: 9898
      targetPort: 9898
      protocol: TCP
containers:
  podinfo:
    image: .
    command:
      - ./podinfo
    args:
      - --port=9898
      - --cache-server=tcp://$(CONNECTION_REDIS_HOST):$(CONNECTION_REDIS_PORT)
    variables:
      PODINFO_UI_COLOR: "#34577c"
      PODINFO_UI_MESSAGE: ""
      CONNECTION_REDIS_HOST: ${resources.cache.host}
      CONNECTION_REDIS_PORT: ${resources.cache.port}
    livenessProbe:
      httpGet:
        port: 9898
        path: /healthz
    readinessProbe:
      httpGet:
        port: 9898
        path: /readyz
resources:
  cache:
    type: redis
    params:
      disableDefaultEnvVars: true
```

At this stage, here are the files you have locally:
```none
.
└──score.yaml
```

## In CI/CD, generate Radius's app.bicep from Score

Initialize the current `score-radius` workspace:
```bash
../score-radius init --no-sample

cp ../examples/provisioners/redis.provisioners.yaml .score-radius/
```

See the available resource types:
```bash
../score-radius provisioners list
```

```none
+-------+---------+-----------------------+--------------------------------------------------+----------------------------------------------------------------+
| TYPE  |  CLASS  |        PARAMS         |                     OUTPUTS                      |                          DESCRIPTION                           |
+-------+---------+-----------------------+--------------------------------------------------+----------------------------------------------------------------+
| redis | default | disableDefaultEnvVars | connectionString, host, password, port, username | Generates a Applications.Datastores/redisCaches bicep resource |
+-------+---------+-----------------------+--------------------------------------------------+----------------------------------------------------------------+
```

Generate the Radius's `app.bicep` file from the Score file:
```bash
../score-radius generate score.yaml \
    -i ghcr.io/stefanprodan/podinfo:latest \
    -o app.bicep
```

At this stage, here are the files you have locally:
```none
.
├── app.bicep
├── .score-radius
│   ├── redis.provisioners.yaml
│   └── state.yaml
└── score.yaml
```

## In CI/CD, deploy the generated Radius's app.bicep to Radius

Create a `bicepconfig.json`:
```json
{
	"experimentalFeaturesEnabled": {
		"extensibility": true
	},
	"extensions": {
		"radius": "br:biceptypes.azurecr.io/radius:latest"
	}
}
```

At this stage, here are the files you have locally:
```none
.
├── app.bicep
├── bicepconfig.json
├── .score-radius
│   ├── redis.provisioners.yaml
│   └── state.yaml
└── score.yaml
```

```bash
rad run app.bicep --group default --application podinfo --environment default
```

Or:
```bash
rad deploy app.bicep --group default --application podinfo --environment default

kubectl port-forward FIXME
```

## See the application in the Radius dashboard

![](./images/demo-with-redis-graph.png)

## As Developer, deploy my Workload locally with Docker Compose

```bash
score-compose init --no-sample

score-compose generate score.yaml \
    --image ghcr.io/stefanprodan/podinfo:latest \
    --publish 9898:podinfo:9898
```

At this stage, here are the files you have locally (files related to Radius removed):
```none
.
├── compose.yaml
├── .score-compose
│   ├── mounts
│   │   └── redis-igglYg
│   │       └── redis.conf
│   ├── state.yaml
│   └── zz-default.provisioners.yaml
└── score.yaml
```

Deploy the generated `compose.yaml` file with Docker Compose:
```bash
docker compose up --build -d --wait
```