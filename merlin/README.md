```
docker build -t merlin .
```
```
docker run --name merlin -v /opt/1panel/apps/openresty/openresty/www/sites/bitmap.date/index/merlin:/app/merlin -it merlin
```
