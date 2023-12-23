## this will be integrated with the lastest version of docker images

docker run -d --name clair-db arminc/clair-db:latest
docker run -p 6060:6060 --link clair-db:postgres -d --name clair arminc/clair-local-scan:v2.0.8_fe9b059d930314b54c78f75afe265955faf4fdc1

clair-scanner --whitelist=example-nginx.yaml --clair=http://YOUR_LOCAL_IP:6060 --ip=YOUR_LOCAL_IP nginx:1.11.6-alpine