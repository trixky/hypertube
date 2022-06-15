# Start with docker

```
docker build -t hypertube-api-scrapper .
docker run -dp 8080:8080 hypertube-api-scrapper
docker exec -it <container_id> /bin/sh
```
