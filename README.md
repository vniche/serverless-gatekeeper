# serverless-gatekeeper

## Getting started

Feel free to check [my article @ medium](https://medium.com/@vcorreaniche/securing-serverless-services-in-kubernetes-with-keycloak-gatekeeper-6d07583e7382), which has the whole story behind this repository.

### Deploy function with Kubeless client
```
$ kubeless function deploy isprime \
  --from-file handler.go \
  --dependencies Gopkg.toml \
  --handler handler.IsPrime \
  --namespace poc \
  --runtime go1.10
```

## Keycloak

for this `demo`, i've used the following settings:

```
- realm: justice-league
    client: service-gatekeeper
      roles:
      - edit (will be able to access isprime function)
      - view (won't be authorized to access isprime function)
    groups:
    - privileged (is associated to service-gatekeeper's edit role)
    - unprivileged (is associated to service-gatekeeper's view role)
    users:
    - superman (belongs to privileged group)
    - clarkkent (belongs to unprivileged group)
```

## Gatekeeper

### Configure

give `gatekeeper.yaml` some attention, replace these with some real values:

```
# uncomment if client is confidential
# client-secret: <client_secret>
discovery-url: http://<keycloak_address>/auth/realms/poc
encryption_key: <random_generated_secret>
upstream-url: http://<service_address>
```

### Run in Docker
```
$ docker run -it --rm \
    -v "$(pwd)":/etc/secrets \
    --network host \
    keycloak/keycloak-gatekeeper:v2.3.0 --config=/etc/secrets/gatekeeper.yaml --enable-logging=true --enable-json-logging=true --verbose=true
```

### Run in Kubernetes

```
# first, create the secret that will store the config file
$ kubectl create secret generic gatekeeper --from-file=./gatekeeper.yaml -n poc

# then create the deployment that will use the secret
$ kubectl apply -f .k8s/deployment.yaml
```

## Use case

### non-permitted access

request token as `clarkkent`:

```
$ curl -X POST \
    ‘http://127.0.0.1:8080/auth/realms/justice-league/protocol/openid-connect/token' \
    -H “Content-Type: application/x-www-form-urlencoded” \
    -d ‘username=clarkkent&password=<clark_password>&grant_type=password&client_id=service-gatekeeper’
```

make a request to endpoint with `clarkkent` token:

```
$ curl -H 'Authorization: Bearer ...8vMTI3LjAuMC4xiJiYWNjNmNkYy05M...' --proxy http://127.0.0.1:3000 http://isprime.poc.svc:8080/isprime -d '5' -v
*   Trying 127.0.0.1...
...
* upload completely sent off: 1 out of 1 bytes
< HTTP/1.1 403 Forbidden
< Date: Sun, 11 Nov 2018 17:18:13 GMT
< Content-Length: 0
< 
* Connection #0 to host 127.0.0.1 left intact
```

### permitted access

request token as `superman`:

```
curl -X POST \
    ‘http://127.0.0.1:8080/auth/realms/justice-league/protocol/openid-connect/token' \
    -H “Content-Type: application/x-www-form-urlencoded” \
    -d ‘username=superman&password=<superman_password>&grant_type=password&client_id=service-gatekeeper’
```

make a request to endpoint with `superman` token:

```
$ curl -H ‘Authorization: Bearer …oxNTQxOTYwMzg3LCJpc…’ --proxy http://127.0.0.1:3000 http://isprime.poc.svc:8080/isprime -d ‘5’ -v
* Trying 127.0.0.1…
...
* upload completely sent off: 1 out of 1 bytes
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Content-Length: 10
< Content-Type: text/plain; charset=utf-8
< Date: Sun, 11 Nov 2018 18:30:50 GMT
< 
* Connection #0 to host 127.0.0.1 left intact
5 is prime
```