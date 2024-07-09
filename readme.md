# ShareIt!

> ⚠️ This is an overly simplistic file sharing app written in `go` meant for skill testing purpose. Definitely **NOT** meant for production!!!

## Requirements

- go compiler (1.22.5+)

## Quick start

```sh
# Server is available on port 8080
go run ./cmd
```

## Supported env-vars

- `BLOB_STORAGE_URL`: Blob storage URL. It is used to persist all files shared on the app. We support `s3://`, `file://` and `mem://` (_not for production use_). More details on those URL schemes: https://gocloud.dev/howto/blob .

- `PUBLIC_URL`: The URL on which the service is available publicly (that is, to end-users). It is useful in case the app is exposed thru a reverse proxy. If not set, the app will try to guess the public URL based on common headers. (example value: `https://shareit.example.com`)
