## Initial Configuration

The API will try to load `CONFIG_PATH` from the env, with the default path as ./config.yaml

The default configuration template is
```
http:
  host_address: ':8080'
  rate_limit: 50
  metrics:
    route: '/metrics'
    host_address: ':8090'
database:
  user: %DATABASE_USER%
  password: %DATABASE_PASSWORD%
  host: '%DATABASE_HOST%'
  port: '%DATABASE_PORT%'
  database_name: %DATABASE_NAME%
  ssl: %SSL_MODE%
```

## Testing

To run basic testing, use: `make test`

To get the test coverage, run: `make coverage`

## Building

To build, use `make build`

This process will create the vendor, build the api with static dependencies on
`api\build\api`

## Folder Structure

The API is subdivided into 5 layers

* API - the main package which is the start point of the api, calling other layers
* Models - the layer which we define our business modeling and inherent logic
* PKG - where we desacoplate the external dependencies from our code
* Repository - the layer responsible for business logic
* Use Case - here lies the actual use cases of the api, which most of the time, use all layers, this is the layer which holds the interaction between the user and the other layers.

## Metrics

The API is exposing metrics for all the routes on the route and port defined inside the configuration file.

Those metrics are:

* Total Requests by Endpoint, Method, StatusCode, ResponseSize, Response Time
* Histogram of requests duration
* Histogram of requests response size