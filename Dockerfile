############################
# BUILD STEP

FROM golang:latest as build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o dist/shareit ./cmd

#######################
# FINAL IMAGE

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=build /src/dist /app
ENTRYPOINT ["/app/shareit"]
