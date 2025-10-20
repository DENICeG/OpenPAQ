# Configuration

There are mainly two options to get OpenPAQ running. The first option is to compile the code and run the binary directly.

This section will focus on the available docker container.
The container is available [here](ghcr.io/deniceg/openpaq:latest).


## Overview of configurable components

There are several optional configuration items available in OpenPAQ. The following diagram provides an overview about 
the available components.
For further details about each component, please have a look to the appropriate section on this page.

```mermaid
flowchart TD
    
    subgraph op[openPAQ]
        openPAQ
        subgraph optional
            TLS
            JWT
        end
    end

    subgraph ch [optional]
        Clickhouse
    end
    
    subgraph mc [optional]
        Memcached
    end

    subgraph nom [optional]
        Nominatim
    end

    op <--> ch
    op <--> nom
    op <--> mc
 
```



## Overview Environment variables

Since OpenPAQ is designed for container environments, all configuration settings are set via environment variables.
An Overview for all available settings are listed in the table below.

| Variable name            | Example Input                              | Required | Description                                                                                                             |
|--------------------------|--------------------------------------------|:--------:|-------------------------------------------------------------------------------------------------------------------------|
| CLICKHOUSE_ENABLED       | false                                      |    x     | Enables the loading of street lists from a clickhouse db server (if you want to use only nominatim set it to false)     |
| CLICKHOUSE_DB_USERNAME   | default                                    |          | User for clickhouse connection                                                                                          |
| CLICKHOUSE_DB_PASSWORD   | prettypassword                             |          | Password for clickhouse connection                                                                                      |
| CLICKHOUSE_DB_HOST       | clickhouse.server.de                       |          | Hostname for clickhouse connection without port                                                                         |
| CLICKHOUSE_DB_PORT       | 9000                                       |          | Port for clickhouse connection                                                                                          |
| CLICKHOUSE_DB_DATABASE   | address                                    |          | Database in which the street lists are stored                                                                           |
| CLICKHOUSE_DB_TABLE      | lists                                      |          | Table of street list of germany to list matching                                                                        |
| USE_TLS                  | false                                      |    x     | Enables TLS encryption to use HTTPS instead of HTTP                                                                     |  
| TLS_KEY_FILE_PATH        | /cert/key.pem                              |          | Path to the secret key file                                                                                             |
| TLS_CERT_FILE_PATH       | /cert/cert.pem                             |          | Path to the certificate file                                                                                            |
| USE_JWT                  | true                                       |    x     | To secure the usage of the api it is possible and recommended to use jwt tokens for authorization                       |
| JWT_SIGNING_KEY          | a-string-secret-at-least-256-bits-long     |          | Signing Key for jwt                                                                                                     |
| WEBSERVER_LISTEN_ADDRESS | 127.0.0.1:8001                             |          | Address OpenPAQ binds on                                                                                                |
| NOMINATIM_ADDRESS        | https://nominatim.openstreetmap.org/search |    x     | Address                                                                                                                 |
| CACHE_ENABLED            | false                                      |    x     | Use memcache                                                                                                            |
| CACHE_URL                | localhost:11211                            |          | URL to memcache                                                                                                         |
| VERSION                  | 5.0.4                                      |    x     | The name of the version is included in the response, useful during development or to distinguish different environments |
| LOG_LEVEL                | info                                       |    x     | Available log levels are: trace, debug, info, warning, error, fatal, panic                                              |
| LOG_FORMAT               | text                                       |          | An alternative log format is "text". The default is a json logging format                                               |

## Configure Clickhouse
If you want to use the list matching mechanism a clickhouse configuration is needed. Information on [Clickhouse](https://clickhouse.com/) and how to set up the database can be found on their website. 
A listmatcher for Germany is implemented and can be used out of the box. Create a table
with the schema described below. Afterwards, fill the table with data from a source of your choice (e.g. [geodata for Germany](https://www.geodaten-deutschland.de/postleitzahlen-strassen-liste-deutschland.php)).
There are several benefits by using a listmatcher:

- The results may be more accurate, as the areas and postal code regions in Nominatim are only calculated.
- The memory lookup can potentially be faster, depending on your setup.

### Schema for german listmatcher
This is an example query to create a lookup table for a german street list.

```clickhouse
CREATE TABLE default.street_list_germany
(
    `id` UInt32,
    `bundesland_name` LowCardinality(String),
    `regierungsbezirk_name` String,
    `kreis_name` String,
    `kreis_typ` LowCardinality(String),
    `gemeinde_name` String,
    `gemeinde_typ` LowCardinality(String),
    `ort_name` String,
    `ort_lat` Float64,
    `ort_lon` Float64,
    `postleitzahl` String,
    `strasse_name` String
)
ENGINE = MergeTree()
ORDER BY (bundesland_name,
 ort_name,
 postleitzahl,
 strasse_name)
SETTINGS index_granularity = 8192;
```

### Environment Variables

After setting up the database, Clickhouse must be enabled in OpenPAQ by setting following environment variables:

- CLICKHOUSE_ENABLED=true
- CLICKHOUSE_DB_HOST=<clickhouse host url\>
- CLICKHOUSE_DB_PORT=<clickhouse native port\>
- CLICKHOUSE_DB_USERNAME=<clickhouse user with read permissions\>
- CLICKHOUSE_DB_PASSWORD=<clickhouse user password\>
- CLICKHOUSE_DB_DATABASE=<clickhouse database name\>
- CLICKHOUSE_DB_TABLE=<clickhouse table name\>

## Configure SSL/TLS

By default, OpenPAQ starts an HTTP server without any encryption. This is completely fine when you run OpenPAQ behind a
reverse proxy which handles the TLS encryption. It is not recommended to run OpenPAQ without any encryption at all.

To enable the HTTPS server of OpenPAQ following environment variables must be set:

- USE_TLS=true
- TLS_KEY_FILE_PATH=<location of the keyfile in the container\>
- TLS_CERT_FILE_PATH=<location of the certfile in the container\>

### Example configuration

#### Generating a self-signed certificate

```shell 
openssl req -x509 \
-newkey rsa:4096 \
-keyout key.pem \
-out cert.pem \
-sha256 \
-days 3650 \
-nodes \
-subj "/C=DE/ST=Hessen/L=Frankfurt am Main/O=Denic e.G./OU=DataScience/CN=OpenPAQ"
```

#### Run OpenPAQ with TLS enabled

```shell 
docker run \
-e CACHE_ENABLED=false \
-e VERSION=testing \
-e CLICKHOUSE_ENABLED=false \
-e USE_TLS=true \
-e TLS_KEY_FILE_PATH=/cert/key.pem \
-e TLS_CERT_FILE_PATH=/cert/cert.pem \
-e USE_JWT=false \
-e NOMINATIM_ADDRESS=https://nominatim.openstreetmap.org/search \
-e LOG_LEVEL=debug \
-e WEBSERVER_LISTEN_ADDRESS=:8001 \
-p 8001:8001 \
-v ./key.pem:/cert/key.pem \
-v ./cert.pem:/cert/cert.pem \
ghcr.io/deniceg/openpaq:latest
```


## Configure Authentication
In some use cases it is necessary to restrict the access to the API. Therefore, OpenPAQ provides a mechanism to protect
the access with a [JWT](https://www.rfc-editor.org/rfc/rfc7519.html). To use this, following parameters must be set:

- USE_JWT=true
- JWT_SIGNING_KEY=<any string without space, tabs or linebreaks which has at least 32 chars \>

### Example

```shell
docker run \
-e CACHE_ENABLED=false \
-e VERSION=testing \
-e CLICKHOUSE_ENABLED=false \
-e USE_TLS=false \
-e USE_JWT=true \
-e JWT_SIGNING_KEY=a-string-secret-at-least-256-bits-long \
-e NOMINATIM_ADDRESS=https://nominatim.openstreetmap.org/search \
-e LOG_LEVEL=debug \
-e WEBSERVER_LISTEN_ADDRESS=:8001 \
-p 8001:8001 \
ghcr.io/deniceg/openpaq:latest
```

After setting up the service, it shouldn't be possible to access the API without a valid JWT.

```shell
curl http://127.0.0.1:8001/api/v1/check\
?street=Theodor-Stern-Kai\
&postal_code=60596\
&city=Frankfurt\
&country_code=DE\
&debug_details=false
```

The response should look like this

```json
{"message":"No token provided"}
```

### Generating an access token

There are several ways to generate JWTs. The easiest way is to use [jwt.io](https://jwt.io/).
Although it is not required, it is recommended that you set an expiration date for the token.
Now you just need to provide the generated token as an authorization header like this: 

```shell
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTc0NjI2NzUxN30.XqARRZ3xSery6SgLDjjPyR2IP2tldfNxbFI7cKOw1w0" \
http://127.0.0.1:8001/api/v1/check\
?street=Theodor-Stern-Kai\
&postal_code=60596\
&city=Frankfurt\
&country_code=DE\
&debug_details=false
```

With the provided token you should get the expected result. 


## Configure Cache

OpenPAQ is prepared to work with memcached. Activate the cache using the environment variables:

-  CACHE_ENABLED = true
-  CACHE_URL = <your server address\>

[Memcached](https://memcached.org/) is a simple object storage with a small footprint.

To start a very simple Memcached service, the following command can be used:

```shell
docker run --rm \
--name my-memcache \
-p 11211:11211 \
memcached \
memcached --memory-limit=64
```

When the memcached container is running, you can start OpenPAQ and enjoy the insane speed of the cache.

```shell
docker run \
-e CACHE_ENABLED=true \
-e CACHE_URL=localhost:11211 \
-e VERSION=testing \
-e CLICKHOUSE_ENABLED=false \
-e USE_TLS=false \
-e NOMINATIM_ADDRESS=https://nominatim.openstreetmap.org/search \
-e LOG_LEVEL=debug \
-e WEBSERVER_LISTEN_ADDRESS=:8001 \
-p 8001:8001 \
ghcr.io/deniceg/openpaq:latest
```

!!! info
    macOS and Windows docker environments are virtualized, so the hostname localhost does not work as expected.
    
    You may need to use the host IP instead of localhost as cache host.

## Advanced setup - self-hosted nominatim 

The limitation and policy of openstreetmap.org maybe a showstopper for some use cases with a high number of requests and automatically
running agents.

In these situations it may make sense to host your own nominatim server.
OpenPAQ is battle test with the docker container [mediagis/nominatim](https://github.com/mediagis/nominatim-docker).


### Example configuration

In this example a small area of Germany is used. Be aware that the import of the whole world takes a lot of time and require
a big amount of resources. See [Nominatim docs](https://nominatim.org/release-docs/5.1/admin/Installation/#hardware) for further information.

```shell
docker run -it \
  -e PBF_URL=https://download.geofabrik.de/europe/germany/hessen-latest.osm.pbf \
  -e REPLICATION_URL=https://https://download.geofabrik.de/europe/germany/hessen-updates/ \
  -p 8080:8080 \
  --name nominatim \
  mediagis/nominatim:5.1
```

!!! info
    macOS and Windows docker environments are virtualized, so the hostname localhost does not work as expected.

    You may need to use the host IP instead of localhost as cache host.


```shell
docker run \
-e CACHE_ENABLED=false \
-e VERSION=testing \
-e CLICKHOUSE_ENABLED=false \
-e USE_TLS=false \
-e NOMINATIM_ADDRESS=http://localhost:8080/search \
-e LOG_LEVEL=debug \
-e WEBSERVER_LISTEN_ADDRESS=:8001 \
-p 8001:8001 \
ghcr.io/deniceg/openpaq:latest
```


