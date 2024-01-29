# JDBC Source Connector

### TL;DR
To run this proof of concept/exploration, run:

```bash
make fullrun
```

You can also run setup only, with:
```bash
make setup
```

Then you can run the load tests (to generate kafka records) multiple times with:
```bash
make run
```

After setup, Kafdrop is available at `localhost:9000`.

You can clean up everything with:
```bash
make cleanup
```

### The exploration

The idea of this exploration is to set up:
- Kafka
- Postgres
- JDBC Source connector
- A simple Golang HTTP server app

Then hit the HTTP server with several requests using grafana/k6.

Kafka messages can be visualized using Kafdrop (`localhost:9000`)

### Files

- [JDBC Source Connector settings](./jdbc/register.json)
- [K6 JS file](./loadtest/loadtest.js)
- [Golang server](./app)
- [Postgres Migrations](./migrations)
- [Docker compose](./docker-compose.yml)