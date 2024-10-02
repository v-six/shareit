# ShareIt!

> ⚠️ This is an overly simplistic file sharing app written in `go` meant for skill testing purpose. Definitely **NOT** meant for production!!!

## Requirements

- go compiler (1.22.5+)
- optional: [Taskfile](https://taskfile.dev/)
- optional: [ko](https://ko.build/)
- optional: docker

## Quick start

```sh
# Server is available on port 8080
go run ./cmd
```

## Supported env-vars

- `BLOB_STORAGE_URL`: Blob storage URL. It is used to persist all files shared on the app. We support `s3://`, `file://` and `mem://` (_not for production use_). [More details on those URL schemes here](https://gocloud.dev/howto/blob). More specifically for [S3 compatible storage](https://gocloud.dev/howto/blob/#s3-compatible).

- `PUBLIC_URL`: The URL on which the service is available publicly (that is, to end-users). It is useful in case the app is exposed thru a reverse proxy. If not set, the app will try to guess the public URL based on common headers. (example value: `https://shareit.example.com`)

## About building docker image

- Can be built from the provided `Dockerfile`
- Can be built with [ko](https://ko.build/)

## About the `/healthz` endpoint

An health check endpoint is available at `/healthz`. It is advised to use it as it checks the availability of the underlying storage. A status code of 200 OK means everything is fine.

## About CI / CD

### CI

CI is running on every push:
- testing *(disabled for now)*
- linting
- building 
 
### Develop CD

Develop CD is running on every push on `develop` branch.
- URL: https://shareit-dev.do-interview2.cw.substance3d.io
- S3: https://shareit-dev-storage.ams3.digitaloceanspaces.com
- Docker tags: `dev` `${github.short_sha}-dev`

### Production CD

Production CD is running on every push on `main` branch.
- URL: https://shareit.do-interview2.cw.substance3d.io
- S3: https://shareit-storage.ams3.digitaloceanspaces.com
- Docker tags: `latest` `${github.short_sha}`

The build is optimized with `-ldflags="-s -w"`

## About branches

- `main`: production protected branch. Merge changes with PR only.
- `develop`: develoment branch. Use it daily to add new features.