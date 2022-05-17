
# IPFS-Upload-Relay

## Environment Variables

- `MODE`: `production` or `development`
- `W3S_TOKEN`: web3.storage token

## Endpoints

- GET `/`: health check
- PUT `/upload`: upload file (multipart/form-data, file=@file)
- POST `/json`: post raw json data (application/json)
