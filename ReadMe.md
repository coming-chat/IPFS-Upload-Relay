
# IPFS-Upload-Relay

## Environment Variables

- `MODE`: `production` or `development`

## Endpoints

- GET `/`: health check
- PUT `/upload`: upload file (multipart/form-data, file=@file)
- POST `/json`: post raw json data (application/json)
