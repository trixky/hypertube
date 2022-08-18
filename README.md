# hypertube

A streaming site using torrents with serverless architecture.
You can watch torrents without download them, because the platform does it for you ! __(docker-compose)__

![Recordit GIF](https://raw.githubusercontent.com/trixky/hypertube/main/.demo/demo.gif)

## Usage

### prerequisites

- docker-compose
- go
- sqlc *(./sqlc/README.md)*
- protobuf *(./proto/README.md)*


### up

```bash
source env.sh
# copy all .env.example in .env
# generate grpc endpoints (see ./protoc/README.md)
# generate sqlc methodes (see ./sqlc/README.md)
docker-compose -f docker-compose.build.yaml up
```

## Stack

- Svelte
- Go / Grpc
- Node.js / Express
- ffmpg
- nginx
