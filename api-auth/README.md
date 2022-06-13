# Start with docker

```
docker build -t hypertube-api .
docker run -dp 8080:8080 hypertube-api
docker exec -it <container_id> /bin/sh
```