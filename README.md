## Deployment


```
cf push
cf bind-service metrics-forwarder-sample metrics-service-instance
cf restage metrics-forwarder-sample
```

## monitor metrics with cf nozzle

```
~:> cf nozzle -f ValueMetric | egrep mf-sample-app
origin:"" eventType:ValueMetric timestamp:1542227248521958090 deployment:"system.domain" job:"metrics-forwarder" index:"737cc3b3-4148-4a59-b298-85ea242199e5" ip:"10.193.76.55" valueMetric:<name:"mf-sample-app" value:1 unit:"number" > 17:"\n\tsource_id\x12$a73c5bee-5dcf-4a52-8a3b-13a912bbee47" 17:"\n\vinstance_id\x12\x010"
```


## how to curl metrics forwarder

* App ID can be retrieved from `VCAP_APPLICATION` environment variable
* App Instance GUID can be retrieved from `CF_INSTANCE_GUID` environemnt variable from inside the app container
* App Index can be retrieved from `CF_INSTANCE_INDEX` environment variable from inside the app container

```
~:> curl -k https://metrics-forwarder.system.domain/v1/metrics -H 'Authorization: 0411d6ea-bb54-4fd4-4eb5-c4d2b1e56714' -H "Content-Type: application/json" -X POST -d '{"applications":[{"id":"a73c5bee-5dcf-4a52-8a3b-13a912bbee47","instances":[{"id":"2b00ab92-5d54-4124-5d40-278f","index":"0","metrics":[{"name":"mf-sample-app","type":"counter","value":3,"unit":"number","Tags":{"severity":"1"}}]}]}]}' -vvv
```