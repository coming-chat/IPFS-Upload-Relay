
# IPFS-Upload-Relay

## Environment Variables

Work mode:

- `MODE`: `production` or `development`, default mode is `production`

System related:

- `REDIS_CONNECTION_STRING`

S3 related settings:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `FOREVERLAND_BUCKET`

## Endpoints

- GET `/`: health check
- PUT `/upload`: upload file (multipart/form-data, file=@file)
- POST `/json`: post raw json data (application/json)
- POST `/video`: post a new video (?url=...)

## Other platforms support

Fork and add as U wish :P
