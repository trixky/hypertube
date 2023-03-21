# hypertube

A streaming site using torrents with serverless architecture.  
You can watch torrents without downloading them, because the platform does it for you ! __(docker-compose)__

<img src="https://raw.githubusercontent.com/trixky/hypertube/main/.demo/demo.gif" alt="Demo gif" width="600"/>

## Usage

### Prerequisites

- docker-compose
- go
- sqlc *(./sqlc/README.md)*
- protobuf *(./proto/README.md)*

### Up

```bash
source env.sh
# copy and fill all .env from their example
# generate grpc endpoints (see ./protoc/README.md)
# generate sqlc methodes (see ./sqlc/README.md)
docker-compose -f docker-compose.build.yaml up
```

## Online

This project is online, so you can visit it by clicking [here](https://hypertube.trixky.com/)!

<img src="https://raw.githubusercontent.com/trixky/hypertube/main/.demo/login.gif" alt="Login page" width="600"/>

> The site is in __demo__ mode for legal reasons, so you can't login/register

## Stack

- Svelte
- Go / Grpc
- Node.js / Express
- ffmpeg
- nginx
- tailwind
