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