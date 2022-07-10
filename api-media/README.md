# Start with docker

```
docker build -t hypertube-api-media .
docker run -dp 8080:8080 hypertube-api-media
docker exec -it <container_id> /bin/sh
```
