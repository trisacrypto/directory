## Health Check 

[Health check](https://github.com/trisacrypto/trisa/blob/e91d0c93c049573a8616ef0c608347ba42868d9d/proto/trisa/gds/api/v1beta1/api.proto#L196-L223) is an endpoint that can be hit manually or via cron to check the health of verified VASPs. It makes a GET request to the `TrisaEndpoint` (if provided) and sets the `ServiceStatus`. 

This service will send the following message: 
```
message HealthCheck {
    // The number of failed health checks that proceeded the current check.
    uint32 attempts = 1;

    // The timestamp of the last health check, successful or otherwise2.
    string last_checked_at = 2;
}
```

It expects to receive the following message: 
```
message ServiceState {
    enum Status {
        UNKNOWN = 0;
        HEALTHY = 1;
        UNHEALTHY = 2;
        DANGER = 3;
        OFFLINE = 4;
        MAINTENANCE = 5;
    }

    // Current service status as defined by the recieving system. The system is obliged
    // to respond with the closest matching status in a best-effort fashion. Alerts will
    // be triggered on service status changes if the system does not respond and the
    // previous system state was not unknown.
    Status status = 1;

    // Suggest to the directory service when to check the health status again.
    string not_before = 2;
    string not_after = 3;
}
```



### Usage 





### Options 

The VASP can return an optional `not_before` and `not_after` to suggest when the health check should be retried. These translate into the `Extra` field on the VASP model. 

`not_before` => `HealthCheckAfter` is the earliest health should be checked again. 

`not_after` => `HealthCheckBefore` is the preferred max time to check again. 

These are suggestions and the service may ignore them. 