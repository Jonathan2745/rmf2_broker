To run all services for orion context broker and data service
```
docker compose up -d
```

All containers are running in network rmf_development_rmf-network

Exposed port is the API endpoints exposed to localhost.

| # | Container Type | HOST_NAME | INTERNAL_PORT | EXPOSED_PORT | Remarks |
| -- | --  | --- | --- | ---  | :---  | 
|1 | Postgres | postgres | 5432 | 5432* | Postgres Database used fpr Data analytics / NGSI-LD Context broker / IT Connectors Configuration (admin: ngb/ngb, )
|2 | MongoDB | mongodb | 27017 | 27016* | MongoDB used for IoT Agent and Orion Context brokder
|3 | Reverse Proxy Server | proxy | 8889 | 9999 | Bridge for services external to rmf-network to context brokers. It will be upgraded to general data gateway in the future.
|4 | Orion Context Broker | ngsi_v2 | 1026 | 1026 | Orion Context broker. only accessible within rmf-network. (for development purpose, it can be accessible via reverse proxy server (3))
|5 | NGSI-LD Broker | ngsi-ld | 9090 | 9090 | Scorpio NGSI-LD broker. only accessible within rmf-network. (for development purpose, it can be accessible via reverse proxy server (3))
|6 | REDIS | redis | 6379 | 6379* | Internal data cache (database 15 is used by cache for IT connector)
|7 | RabbitMQ | rabbitmq | 5672 / 15672 (admin) | 5672 / 15672(admin) | Service bus data broadcasting using exchange (exchange = "orion", queue = "@SYSTEM@")
|8 | IT_CONNECTOR | it_connector | 4201 | 8001 | IT data pipeline from external data source to context broker
|9 | Swagger | <not defined> | - | 8000 | Swagger Utility for NGSI-V2 and NGSI-LD


`* Port will be disabled for production build. service will only be available within the rmf-network

To run containers in the RMF service, please run the followings:
```
docker run [-p <PORT>] [-e <Environment>] --net rmf_development_rmf-network [--ip XXX.XXX.XXX.XXX] <Your Container Name>
```