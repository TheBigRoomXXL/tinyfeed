---
title: Kubernetes
---

# Kubernetes

This documentation goes over a bare-bones deployment in to Kubernetes.

It's essentially the connecting pieces of running Tinyfeed in [Docker](docker.md) and using [CronJobs](cron.md)

However, this document is not designed to teach you the intricacies of Kubernetes, but give you a rough start on what YAML you will want to write.

This makes use of the `args` feature, eg: `tinyfeed --output "index.html" https://lovergne.dev/rss` which is discussed in more
details in [Basic Usage](usage.md#basic-usage)

Some assumption will be made here about your kubernetes knowledge, as running a k8s cluster requires some pre-existing knowledge. If something is not
clear, open an issue on the repo.

## Prerequisites

| Name               | Comments                                                                                                                                                       |
|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Kubernetes cluster | Any version, you can run this in minikube, k3s, gke, eks, talos etc                                                                                            |
| Storage Class      | This makes use of a storage class, however there is a hack we can make use of to get around this if your cluster is ephemeral or has no storage class          |
| Ingress            | You can use anything you like. Cloudflare tunnels, Metallb, nginx. This Documentation wont go over this however, as each cluster is very personal to the admin |

## Deploy

!!! Info 
    This documentation is not production-tested, and is aimed at home lab usage. Enough to get you over the line

The deployment is made up of 4 individual resources.

1. Persistent Volume to share the `index.html` file between the cron job and nginx
2. CronJob to generate the Static HTML from the list of feeds
3. nginx web server that serves the `index.html` file
4. Kubernetes `ClusterIP` service that we can then connect our Load Balancer to later

### How this works

The cron job, every 15 minutes (or how often you define it) will spawn and run the container with the arguments we've included.

It will then generate the `index.html` file we're used to seeing to `/data/index.html` which is then served up by the nginx container.


### Storage

In order to share the generated `index.html` file between the cronjob, where it's generated, and the nginx server where it's served from,
we need to store it somewhere.

In Kubernetes, we can use a Persistent Volume to store this file.

Create a file called `pvc.yaml` and put the below in it

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: tinyfeed
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
```

Apply this file

```shell
kubectl apply -f pvc.yaml
```

### Cron job

Create a file called `job.yaml` and in it put the below in to that file. 

``` yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: tinyfeed
spec:
  schedule: "*/15 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          securityContext:
            fsGroup: 1000
          volumes:
            - name: static
              persistentVolumeClaim:
                claimName: tinyfeed
          containers:
            - name: tinyfeed
              image: docker.io/thebigroomxxl/tinyfeed:v1.3.0
              imagePullPolicy: IfNotPresent
              args:
                - "https://lovergne.dev/rss"
                - "https://breadnet.co.uk/rss"
                - "-o"
                - "/data/index.html"
              volumeMounts:
                - mountPath: /data
                  name: static
          restartPolicy: OnFailure
```

```shell
kubectl apply -f job.yaml
```

In order to add more feeds, simply add more `args`

Eg:

```diff
              args:
                - "https://lovergne.dev/rss"
                - "https://breadnet.co.uk/rss"
+               - "https://example.com/rss"
                - "-o"
                - "/data/index.html"
```


This will use your Default storage class. 

It's important that the `AccessModes` is set to `ReadWriteMany` as the cron job will `write` to it _many_ times, and the `nginx` pod will _read from it_

### Web server

We are going to use `nginx` as our web server as it's super simple and is one of the most widely adopted web servers on the internet these days

Create a file called `deployment.yaml` and in that file, put the below

``` yaml 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tinyfeed
  labels:
    app: tinyfeed
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tinyfeed
  template:
    metadata:
      name: tinyfeed
      labels:
        app: tinyfeed
    spec:
      volumes:
        - name: static
          persistentVolumeClaim:
            claimName: tinyfeed
      containers:
        - name: tinyfeed
          image: nginx
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - mountPath: /usr/share/nginx/html
              name: static
          ports:
            - containerPort: 80
              protocol: TCP
      restartPolicy: Always
```

```shell
kubectl apply -f deployment.yaml
```

This creates an `nginx` container and then mounts the `index.html` file in to the `/usr/share/nginx/html` directory where 
the default config of this container looks for files to serve.

### Service

Depending on what your networking setup is, again this is not meant to be a conclusive document, you may decide on different method for how to serve this to your self or your users.

For this example, we will assume what ever onwards technology you use, requires a service.


Create a file called `service.yaml` and in that file put the below

```yaml 
apiVersion: v1
kind: Service
metadata:
  name: tinyfeed
spec:
  selector:
    app: tinyfeed
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: ClusterIP
```

```shell
kubectl apply -f service.yaml
```


If you are following along in MiniKube, you should [consult their documentation on how to access your service](https://minikube.sigs.k8s.io/docs/handbook/accessing/)

## Hack to avoid using a Persistent Volume

If you like to keep your clusters super ephemeral and have no persistence then the below is going to be the best solution.

The below creates a deployment of Tinyfeed with an init container that generates the `index.html` file, then once that container
exits with status of `0` (success) it spawns the _main_ container which is `nginx` to serve the site.

Every 15 minutes as defined in the Cronjob `*/15 * * * *` (different cronjob from previous example!!) it finds all the 
pods that match the label of `app=tinyfeed` and deletes them. 

This causes the cycle to start again, refreshing the feeds and then serving.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tinyfeed
  labels:
    app: tinyfeed
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tinyfeed
  template:
    metadata:
      labels:
        app: tinyfeed
    spec:
      securityContext:
        fsGroup: 1000
      volumes:
        - name: static
          emptyDir: {}
      initContainers:
        - name: tinyfeed-init
          image: docker.io/thebigroomxxl/tinyfeed:v1.3.0
          imagePullPolicy: IfNotPresent
          args:
            - "https://lovergne.dev/rss"
            - "https://breadnet.co.uk/rss"
            - "-o"
            - "/data/index.html"
          volumeMounts:
            - name: static
              mountPath: /data
      containers:
        - name: nginx
          image: nginx
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: static
              mountPath: /usr/share/nginx/html
          ports:
            - containerPort: 80
              protocol: TCP
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: tinyfeed
spec:
  selector:
    app: tinyfeed
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: ClusterIP
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: tinyfeed-roller
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: tinyfeed-roller
rules:
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - get
      - list
      - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: tinyfeed-roller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: tinyfeed-roller
subjects:
  - kind: ServiceAccount
    name: tinyfeed-roller
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: tinyfeed-roller
spec:
  concurrencyPolicy: Forbid
  schedule: "*/15 * * * *"
  jobTemplate:
    spec:
      backoffLimit: 2
      activeDeadlineSeconds: 20
      ttlSecondsAfterFinished: 10
      template:
        spec:
          serviceAccountName: tinyfeed-roller
          restartPolicy: Never
          containers:
            - name: kubectl
              image: cgr.dev/chainguard/kubectl:latest
              command:
                - 'kubectl'
                - 'delete'
                - 'pod'
                - '-l app=tinyfeed'
```

## Closing notes

This is designed to give you enough of a starting point to be able to deploy this to your own k8s cluster. It's not meant to
replace actual knowledge of Kubernetes as best practices are not being followed 100% here.

If you struggle with this, [please open an issue](https://github.com/TheBigRoomXXL/tinyfeed/issues/new?title=Kubernetes%20Docs%3A%20%3Cyour%20issue%20goes%20here%3E)
