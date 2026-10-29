# NASA APOD API

A REST API built with Go and Fiber that fetches NASA's Astronomy Picture of the Day (APOD), stores metadata in PostgreSQL, and stores images in MinIO.

## Technologies
- Go
- Fiber
- PostgreSQL
- MinIO
- Docker
- k8s

## API Endpoints

- `GET /api/v1/apods` - Get all APODs in storage
- `GET /api/v1/apod/today` — Get today's APOD metadata.
- `GET /api/v1/apod/today/image` — Get today's image.
- `GET /api/v1/apod/:date` — Get APOD metadata for a specific date.
- `GET /api/v1/apod/image/:date` — Get the image for a specific date.

