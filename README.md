# faasflow-redis-datastore
A **[faasflow](https://github.com/s8sg/faasflow)** datastore implementation that uses redis  to store data  

## Getting Stated

### Deploy redis

#### Deploy in Kubernets


```bash
kubectl apply -f resource/redis-k8s-standalone.yml
```
#### Deploy in Swarm

TODO

### Use redis dataStore in `faasflow`
* Set the `stack.yml` with the necessary environments
```yaml
      redis_url: "redis.default.svc.cluster.local:6379"
      redis_master: ""
```
* Use the `faasflowMinioDataStore` as a DataStore on `handler.go`
```go
redisDataStore "github.com/chennqqi/faas-flow-redis-datastore"

func DefineDataStore() (faasflow.DataStore, error) {

       // initialize redis DataStore
       ds, err := redisDataStore.InitFromEnv()
       if err != nil {
               return nil, err
       }

       return ds, nil
}
```
