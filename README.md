# GCProv Service

## Goals

Implements a small service that can be used to spawn cloud instances on Google Cloud(GCE) using Compute API.
Service also configures​ ​instance(boostraps)​ ​with​ ​a​ user ​given​ ​username/password​ ​and​ ​responds​ ​with​ ​the​ ​newly created​ ​instance’s​ ID.

## Routes

gcProv ​service​ ​ ​accepts​ ​HTTP​ ​requests​ ​with​ ​the​ ​following​ ​endpoints:

● GET​ ​/healthcheck
● POST​ ​/v1/instances/create

### GET​ ​/healthcheck

`GET http://localhost:8000/healtcheck`

Implements​ ​a​ ​basic​ ​health​ ​check,​ ​returning​ ​HTTP​ ​status​ ​code​ ​200​ ​and​ ​a​ message ​page.

### POST​ ​/v1/instances/create

Receives​ ​a​ ​username​ ​and​ ​password​ ​as​ ​query​ ​parameters​ ​(ideally SSH keys ​it’s​ ​just​ ​a​ ​basic API for now)​ ​and​ ​responds​ ​with​ ​instanceID.
This​ ​endpoint​ ​creates​ ​a​ ​small​ ​cloud​ ​instance​ ​of​ ​a​ ​configurable​ ​type,​ ​configures​ ​a​ ​user​ ​with​ ​the given​ ​username​ ​and​ ​password,​ ​ensures​ ​password​ ​authentication​ ​is​ ​working,​ ​and​ ​ensures​ ​the user​ ​can​ ​be​ ​used​ ​to​ ​log​ ​in​ ​as​ ​root.
The​ ​response​ ​is​ ​sent​ ​back​ ​to​ ​the​ ​client​ ​as​ ​soon​ ​as​ ​the​ ​server​ ​is​ ​ready,​ ​so​ ​the​ ​client​ ​must​ ​expect to​ ​wait​ ​for​ ​a​ ​few​ ​seconds​ ​until​ ​the​ ​instance​ ​becomes​ ​healthy.

#### Parameters

```go
project          // project name in GCP where to launch the instance
region           // Region location under GCP
zone             // zone
username         // username to be created in the spawned instance
userpass         // userpassword so the user can log in

```

Sample parameters `project is "demo"` 
                  `region is "us-west1"`
                  `zone is "us-west1-c"`
                  `username is "demouser"`
                  `userpass is "demopass"`

### GET /v1/instances/status

`GET http://localhost:8000/status`

Receives the Google compute instance name and returns instance status

#### Status Parameters

```go
project          // project name in GCP
zone             // zone under region
instance         // instance name 

```

#### Instance customization

To customize the instance name, image and disk size or the network interfaces and service account look at attributes under provisionInstance.


## Project structure

```go
gc

-routes.go                  //routes definitions
-healthCheckHandlers.go     //Healtcheck handlers
-instanceHandlers.go        //instance prov and status

main.go                     //entrypoint for the app
```
