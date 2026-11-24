# Deploy Nominatim

OpenPAQ uses [Nominatim](https://github.com/osm-search/Nominatim) as its backend for geocoding lookups.
While the public instance at `nominatim.openstreetmap.org` works for testing, its
[usage policy](https://operations.osmfoundation.org/policies/nominatim/) limits automated and
high-volume requests. For production use, a self-hosted Nominatim instance is recommended.

OpenPAQ is battle-tested with the [mediagis/nominatim](https://github.com/mediagis/nominatim-docker)
Docker image, which bundles PostgreSQL, PostGIS and the Nominatim application into a single container.

## Prerequisites

- [Docker](https://www.docker.com/) (>= 20.10)
- Docker Compose (optional, recommended for production)
- Sufficient RAM and SSD storage (see [Hardware Requirements](#hardware-requirements))
## Quick Start

Pull the image and start Nominatim with a small dataset to verify everything works:

```bash
docker run -it \
  -e PBF_URL=https://download.geofabrik.de/europe/monaco-latest.osm.pbf \
  -p 8080:8080 \
  --name nominatim \
  mediagis/nominatim:5.1
```

After the import completes, test the API:

```bash
curl "http://localhost:8080/search?q=avenue+pasteur&format=jsonv2"
```

## Hardware Requirements

Requirements vary significantly depending on the geographic extent of the imported data.

### By Region

| Region         | RAM (Import) | RAM (Runtime) | Disk (SSD)   | Import Duration  |
|----------------|--------------|---------------|--------------|------------------|
| City / State   | 4–8 GB       | 2–4 GB        | 30–50 GB     | 30 min – 2 h     |
| Germany        | 12–16 GB     | 8–10 GB       | ~100 GB      | 8–12 h           |
| Europe         | 32–64 GB     | 16–32 GB      | 400–500 GB   | 1–3 days         |
| Full Planet    | 128 GB+      | 64 GB+        | 1 TB+        | 2–5 days         |

!!! warning
    For a full planet import, the Nominatim project strongly recommends **128 GB RAM or more**.
    Out-of-memory issues with less than 64 GB RAM are not considered bugs.
    NVMe SSDs significantly reduce import times compared to SATA SSDs.

### CPU

- **Import**: as many cores as available (8–16 recommended for a planet import)
- **Runtime**: 4+ cores for moderate query load
### Disk

- **SSD is mandatory** — HDDs are not suitable for Nominatim
- NVMe SSDs are strongly recommended for planet imports
- The import process requires significantly more disk space than the final database
### Database Size by Import Style (Planet)

| Style      | Import Time | DB Size  | After `freeze` |
|------------|-------------|----------|-----------------|
| admin      | ~4 h        | ~215 GB  | ~20 GB          |
| street     | ~22 h       | ~440 GB  | ~185 GB         |
| address    | ~36 h       | ~545 GB  | ~260 GB         |
| full       | ~54 h       | ~640 GB  | ~330 GB         |
| extratags  | ~54 h       | ~650 GB  | ~340 GB         |

!!! info
    These values grow over time as OpenStreetMap data increases. Treat them as rough estimates.

## Running Nominatim

### Single Country (e.g. Germany)

```bash
docker run -it \
  -e PBF_URL=https://download.geofabrik.de/europe/germany-latest.osm.pbf \
  -e REPLICATION_URL=https://download.geofabrik.de/europe/germany-updates/ \
  -e IMPORT_STYLE=address \
  -p 8080:8080 \
  -v nominatim-data:/var/lib/postgresql/14/main \
  --shm-size=1g \
  --name nominatim \
  mediagis/nominatim:5.1
```

### Europe

```bash
docker run -it \
  -e PBF_URL=https://download.geofabrik.de/europe-latest.osm.pbf \
  -e REPLICATION_URL=https://download.geofabrik.de/europe-updates/ \
  -e IMPORT_STYLE=address \
  -e THREADS=8 \
  -p 8080:8080 \
  -v nominatim-data:/var/lib/postgresql/14/main \
  -v nominatim-flatnode:/nominatim \
  --shm-size=4g \
  --name nominatim \
  mediagis/nominatim:5.1
```

### Full Planet

```bash
docker run -it \
  -e PBF_URL=https://planet.openstreetmap.org/pbf/planet-latest.osm.pbf \
  -e REPLICATION_URL=https://planet.openstreetmap.org/replication/day \
  -e IMPORT_STYLE=full \
  -e THREADS=16 \
  -p 8080:8080 \
  -v nominatim-data:/var/lib/postgresql/14/main \
  -v nominatim-flatnode:/nominatim \
  --shm-size=8g \
  --name nominatim \
  mediagis/nominatim:5.1
```

## Docker Compose

For production use, a `docker-compose.yml` is recommended.

### Example: Germany

```yaml
services:
  nominatim:
    image: mediagis/nominatim:5.1
    container_name: nominatim
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      PBF_URL: https://download.geofabrik.de/europe/germany-latest.osm.pbf
      REPLICATION_URL: https://download.geofabrik.de/europe/germany-updates/
      IMPORT_STYLE: address
      THREADS: 4
    volumes:
      - nominatim-data:/var/lib/postgresql/14/main
      - nominatim-flatnode:/nominatim
    shm_size: 1g
 
volumes:
  nominatim-data:
  nominatim-flatnode:
```

### Example: Europe

```yaml
services:
  nominatim:
    image: mediagis/nominatim:5.1
    container_name: nominatim
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      PBF_URL: https://download.geofabrik.de/europe-latest.osm.pbf
      REPLICATION_URL: https://download.geofabrik.de/europe-updates/
      IMPORT_STYLE: address
      THREADS: 8
    volumes:
      - nominatim-data:/var/lib/postgresql/14/main
      - nominatim-flatnode:/nominatim
    shm_size: 4g
 
volumes:
  nominatim-data:
  nominatim-flatnode:
```

### Example: Full Planet

```yaml
services:
  nominatim:
    image: mediagis/nominatim:5.1
    container_name: nominatim
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      PBF_URL: https://planet.openstreetmap.org/pbf/planet-latest.osm.pbf
      REPLICATION_URL: https://planet.openstreetmap.org/replication/day
      IMPORT_STYLE: full
      THREADS: 16
    volumes:
      - nominatim-data:/var/lib/postgresql/14/main
      - nominatim-flatnode:/nominatim
    shm_size: 8g
 
volumes:
  nominatim-data:
  nominatim-flatnode:
```

Start with:

```bash
docker compose up -d
```

## Environment Variables

| Variable             | Description                                                     | Example                                                        |
|----------------------|-----------------------------------------------------------------|----------------------------------------------------------------|
| `PBF_URL`            | URL of the OSM data file to download and import                 | `https://download.geofabrik.de/europe/germany-latest.osm.pbf` |
| `PBF_PATH`           | Local path to an already downloaded PBF file (instead of URL)   | `/data/germany.osm.pbf`                                       |
| `REPLICATION_URL`    | URL for ongoing data updates                                    | `https://download.geofabrik.de/europe/germany-updates/`       |
| `IMPORT_STYLE`       | Scope of imported data                                          | `admin`, `street`, `address`, `full`, `extratags`              |
| `THREADS`            | Number of threads for the import process                        | `4`, `8`, `16`                                                 |
| `NOMINATIM_PASSWORD` | PostgreSQL password for the Nominatim user                      | `mysecretpassword`                                             |
| `REVERSE_ONLY`       | Only enable reverse geocoding (saves ~5% disk space)            | `true`                                                         |
| `IMPORT_WIKIPEDIA`   | Import Wikipedia data for improved result ranking               | `true`                                                         |
| `FREEZE`             | Remove update structures after import (saves significant space) | `true`                                                         |

### Import Styles

| Style      | Includes                                     | Use Case                                     |
|------------|----------------------------------------------|----------------------------------------------|
| `admin`    | Administrative boundaries only               | Coarse assignment (country, state, city)      |
| `street`   | + streets                                    | Street-level geocoding                        |
| `address`  | + house numbers, postal codes                | Full address geocoding (**recommended**)      |
| `full`     | + POIs (restaurants, shops, etc.)             | Search for places and amenities               |
| `extratags`| + additional OSM tags                        | Maximum data coverage                         |

## Using a Local PBF File

If you have already downloaded the PBF file (e.g. for large imports with slow connections):

```bash
docker run -it \
  -e PBF_PATH=/data/europe-latest.osm.pbf \
  -e REPLICATION_URL=https://download.geofabrik.de/europe-updates/ \
  -e IMPORT_STYLE=address \
  -p 8080:8080 \
  -v /path/to/your/file:/data \
  -v nominatim-data:/var/lib/postgresql/14/main \
  -v nominatim-flatnode:/nominatim \
  --shm-size=4g \
  --name nominatim \
  mediagis/nominatim:5.1
```

## Keeping Data Up to Date

When `REPLICATION_URL` is set and `FREEZE` is **not** set to `true`, the container automatically fetches updates.

### Trigger Updates Manually

```bash
docker exec -it nominatim nominatim replication --once
```

### Check Update Status

```bash
docker exec -it nominatim nominatim replication --check-for-updates
```

```bash
# Freeze an already running instance
docker exec -it nominatim nominatim freeze
```

## Connecting OpenPAQ to Self-Hosted Nominatim

Point the `NOMINATIM_ADDRESS` environment variable to your local instance:

```bash
docker run \
  -e CACHE_ENABLED=false \
  -e VERSION=testing \
  -e CLICKHOUSE_ENABLED=false \
  -e USE_TLS=false \
  -e USE_JWT=false \
  -e NOMINATIM_ADDRESS=http://localhost:8080/search \
  -e LOG_LEVEL=debug \
  -e WEBSERVER_LISTEN_ADDRESS=:8001 \
  -p 8001:8001 \
  ghcr.io/deniceg/openpaq:latest
```

!!! info
    macOS and Windows Docker environments are virtualized, so the hostname `localhost` does not
    work as expected. You may need to use the host IP or `host.docker.internal` instead.

## API Endpoints

After the import completes, the Nominatim API is available at `http://localhost:8080`.

### Search (Address → Coordinates)

```bash
curl "http://localhost:8080/search?q=Brandenburger+Tor,+Berlin&format=jsonv2"
```

### Reverse Geocoding (Coordinates → Address)

```bash
curl "http://localhost:8080/reverse?lat=52.5163&lon=13.3777&format=jsonv2"
```

### Lookup (OSM ID → Details)

```bash
curl "http://localhost:8080/lookup?osm_ids=R146656,W104393803&format=jsonv2"
```

### Status

```bash
curl "http://localhost:8080/status"
```

## Monitoring

### Check Database Size

```bash
docker exec -it nominatim \
  sudo -u postgres psql -d nominatim \
  -c "SELECT pg_size_pretty(pg_database_size('nominatim'));"
```

### View Logs

```bash
docker logs -f nominatim
```

### Backup the Data Volume

```bash
# Create backup
docker run --rm \
  -v nominatim-data:/data \
  -v $(pwd):/backup \
  ubuntu tar czf /backup/nominatim-backup.tar.gz /data
 
# Restore backup
docker run --rm \
  -v nominatim-data:/data \
  -v $(pwd):/backup \
  ubuntu tar xzf /backup/nominatim-backup.tar.gz -C /
```

## Deployment Scenarios

### Single Country (e.g. Germany)

- **Server**: 16 GB RAM, 4 cores, 200 GB SSD
- **Import**: ~10 hours
- **Database size**: ~100 GB
### Europe

- **Server**: 64 GB RAM, 8+ cores, 500 GB NVMe SSD
- **Import**: 1–3 days
- **Database size**: ~250–350 GB
### Full Planet

- **Server**: 128 GB+ RAM, 16 cores, 1.5 TB+ NVMe SSD
- **Import**: 2–5 days (NVMe ~2 days, SATA SSD 4–5 days)
- **Database size**: ~640 GB (full style)
## Geofabrik Data Sources

Commonly used extracts from [download.geofabrik.de](https://download.geofabrik.de):

| Region          | PBF URL                                                                | Updates URL                                                             |
|-----------------|------------------------------------------------------------------------|-------------------------------------------------------------------------|
| Monaco (test)   | `https://download.geofabrik.de/europe/monaco-latest.osm.pbf`          | `https://download.geofabrik.de/europe/monaco-updates/`                  |
| Germany         | `https://download.geofabrik.de/europe/germany-latest.osm.pbf`         | `https://download.geofabrik.de/europe/germany-updates/`                 |
| Austria         | `https://download.geofabrik.de/europe/austria-latest.osm.pbf`         | `https://download.geofabrik.de/europe/austria-updates/`                 |
| Switzerland     | `https://download.geofabrik.de/europe/switzerland-latest.osm.pbf`     | `https://download.geofabrik.de/europe/switzerland-updates/`             |
| Europe          | `https://download.geofabrik.de/europe-latest.osm.pbf`                 | `https://download.geofabrik.de/europe-updates/`                         |
| North America   | `https://download.geofabrik.de/north-america-latest.osm.pbf`          | `https://download.geofabrik.de/north-america-updates/`                  |
| Full Planet     | `https://planet.openstreetmap.org/pbf/planet-latest.osm.pbf`          | `https://planet.openstreetmap.org/replication/day`                      |

## References

- [mediagis/nominatim-docker (GitHub)](https://github.com/mediagis/nominatim-docker)
- [mediagis/nominatim (Docker Hub)](https://hub.docker.com/r/mediagis/nominatim)
- [Nominatim Documentation](https://nominatim.org/release-docs/latest/)
- [Nominatim Hardware Requirements](https://nominatim.org/release-docs/latest/admin/Installation/#hardware)
- [OSM Data Extracts (Geofabrik)](https://download.geofabrik.de)
- [Planet Mirrors](https://wiki.openstreetmap.org/wiki/Planet.osm#Planet.osm_mirrors)
 