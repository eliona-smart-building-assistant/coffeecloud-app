# Eliona app to access CoffeeCloud machines

This app collects coffee machines available in the [CoffeeCloud](https://www.scanomat.com/coffeecloud/).

## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables, the app provides an own API access.

### Registration in Eliona

To start and initialize an app in an Eliona environment, the app has to be registered in Eliona. For this, entries in database tables `public.eliona_app` and `public.eliona_store` are necessary.

This initialization can also be handled by the `reset.sql` script.

### Environment variables

The following environment variables are required:

* `CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started (e.g. `postgres://user:pass@localhost:5432/iot`).
* `INIT_CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db) for app initialization like creating schema and tables (e.g. `postgres://user:pass@localhost:5432/iot`). Default is content of `CONNECTION_STRING`.
* `API_ENDPOINT`: The endpoint of the Eliona API v2.
* `API_TOKEN`: The secret token to authenticate the app with the Eliona API.
* `API_SERVER_PORT`: (optional) The port of the API server. Defaults to 3000.
* `LOG_LEVEL`: (optional) The minimum log level. Defaults to `info`.

### Database tables

The app creates the following database tables during initialization:

* `coffecloud.configuration`: Contains the configuration of the app.
* `coffecloud.asset`: Maps machines and groups to Eliona asset IDs.

## Limitations

The app only provides the coffee machines as grouped in the CoffeeCloud environment. The deeper hierarchy of groups is not synchronized to Eliona. Machines of deeper groups are grouped under their root group.

The CoffeeCloud API limits the number of requests that can be made per unit of time. Therefore, it is important to collect data with a time interval that is long enough to avoid being banned from the API.

## References

### App API

The app provides its own API to access configuration data and other functions. The API definition is in the `openapi.yaml` file.

* [API Reference](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/coffeecloud-app/develop/openapi.yaml)

### Eliona assets

The app creates Eliona asset types and attribute sets during initialization.

### Continuous asset creation

Assets for all machines in the CoffeeCloud are created automatically when the configuration is added.

To select which assets to create, a filter can be specified in the configuration. The schema of the filter is defined in the `openapi.yaml` file. Possible filter parameters are defined in the structs marked with the `eliona:"attribute_name,filterable"` field tag.

To avoid conflicts, the Global Asset Identifier is a manufacturer's ID prefixed with asset type name as a namespace.

### Dashboard

An example dashboard is available by accessing the `/dashboard-templates` endpoint.

## Tools

### Generate API server stub

The [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) can be used to generate a server stub for the app. The easiest way to do this is to use the `generate-api-server.cmd` or `generate-api-server.sh` script.

### Generate database access

The [SQLBoiler](https://github.com/volatiletech/sqlboiler) tool can be used to generate database access code for the app. The easiest way to do this is to use the `generate-db.cmd` or `generate-db.sh` script.
