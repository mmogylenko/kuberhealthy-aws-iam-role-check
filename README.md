## Kuberhealthy AWS IAM Role Check

![Go](https://github.com/mmogylenko/kuberhealthy-aws-iam-role-check/workflows/Go/badge.svg) ![Gosec](https://github.com/mmogylenko/kuberhealthy-aws-iam-role-check/workflows/Gosec/badge.svg) [![GitHub tag](https://img.shields.io/github/tag/mmogylenko/kuberhealthy-aws-iam-role-check.svg)](https://github.com/mmogylenko/kuberhealthy-aws-iam-role-check/tags/)


`Kuberhealthy AWS IAM Role Check` validates if containers running within your cluster can properly make AWS service requests

#### Check Workflow

- Create AWS STS Client
- Call [Get Caller Identity](https://docs.aws.amazon.com/cli/latest/reference/sts/get-caller-identity.html) to get a role whose credentials are used to call the operation 
- Compare *TARGET_ARN* (what role we expect to be) with a role from Get Caller Identity. [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) components that are validated: Service, AccountID and ResourceID

#### Kuberhealthy AWS IAM Role Check Kube Spec Example
```yaml
apiVersion: comcast.github.io/v1
kind: KuberhealthyCheck
metadata:
  name: aws-iam-role
spec:
  runInterval: 5m
  timeout: 1m
  extraAnnotations:
    iam.amazonaws.com/role: "arn:aws:iam::000000000000:role/kubernetes-example-role" # Replace this value with your ARN
    iam.amazonaws.com/external-id: <role-external-id> # Use this if kube2iam is using external-id for roles
  podSpec:
    containers:
    - name: main
      image: ghcr.io/mmogylenko/khcheck-aws-iam-role:latest
      imagePullPolicy: IfNotPresent
      env:
        - name: TARGET_ARN
          value: "arn:aws:iam::000000000000:role/kubernetes-example-role"
        - name: DEBUG # OPTIONAL
          value: "1"
        - name: NODE_NAME # OPTIONAL. Good to know which worker is failing
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
```
, where *TARGET_ARN* is a target ARN that needs to be validated.

#### Docker Image

**Docker** is the only one requirement
```bash
➜  kuberhealthy-aws-iam-role-check git:(master) ✗ make image
docker build -f Dockerfile -t khcheck-aws-iam-role:0.0.1 /Users/mogylenk/Work/code/go/src/kuberhealthy-aws-iam-role-check
Sending build context to Docker daemon  177.2kB
Step 1/17 : FROM golang:1.15-alpine AS builder
 ---> 1a87ceb1ace5
Step 2/17 : ENV APP_NAME=khcheck-aws-iam-role
 ---> Using cache
 ---> d9cc67423f80
Step 3/17 : ENV APP_VERSION=0.0.1
 ---> Using cache
 ---> b42e694a48e3
Step 4/17 : ENV GO111MODULE=on     CGO_ENABLED=0     GOOS=linux     GOARCH=amd64
 ---> Using cache
 ---> 6643a08fac59
Step 5/17 : WORKDIR /build
 ---> Using cache
 ---> 1e38fc429b48
Step 6/17 : COPY go.mod .
 ---> Using cache
 ---> b3abebc0b899
Step 7/17 : COPY go.sum .
 ---> Using cache
 ---> 3197ca1de4b9
Step 8/17 : RUN go mod download
 ---> Using cache
 ---> 33755fb06b1d
Step 9/17 : COPY . .
 ---> 1542332f73b0
Step 10/17 : RUN date +%s > buildtime
 ---> Running in 2c0e1b4a2c17
Removing intermediate container 2c0e1b4a2c17
 ---> a077891a7357
Step 11/17 : RUN APP_BUILD_TIME=$(cat buildtime);     go build -ldflags="-X 'main.buildTime=${APP_BUILD_TIME}' -X 'main.buildVersion=${APP_VERSION}'" -o ${APP_NAME} .
 ---> Running in 1ab3c3574013
Removing intermediate container 1ab3c3574013
 ---> a23182d51dff
Step 12/17 : WORKDIR /app
 ---> Running in 318a955b424b
Removing intermediate container 318a955b424b
 ---> a3da5a415a0c
Step 13/17 : RUN cp /build/${APP_NAME} .
 ---> Running in 772f64f6126f
Removing intermediate container 772f64f6126f
 ---> 7ffd0977ad0e
Step 14/17 : FROM scratch
 --->
Step 15/17 : COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 ---> Using cache
 ---> 3a8b5b20871f
Step 16/17 : COPY --from=builder /app/${APP_NAME} /
 ---> a8919dc93f39
Step 17/17 : CMD ["/khcheck-aws-iam-role"]
 ---> Running in d4f9f919e37e
Removing intermediate container d4f9f919e37e
 ---> 063f91ea50c4
Successfully built 063f91ea50c4
Successfully tagged khcheck-aws-iam-role:0.0.1
```

#### Installation

>Make sure you are using the latest release of Kuberhealthy 2.2.0.

Run `kubectl apply` against [example spec file](example/khcheck-aws-iam-role.yaml)

```bash
kubectl apply -f khcheck-aws-iam-role.yaml -n kuberhealthy
```
##### Container Image

Image is available both from [Docker HUB](https://hub.docker.com/r/mmogylenko/khcheck-aws-iam-role) and [Github Container Registry](https://github.com/users/mmogylenko/packages/container/khcheck-aws-iam-role/)

### Licensing

This project is licensed under the Apache V2 License. See [LICENSE](LICENSE) for more information.
