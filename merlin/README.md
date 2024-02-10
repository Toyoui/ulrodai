```
docker build -t merlin .
```
```
docker run --name merlin -v /opt/1panel/apps/openresty/openresty/www/sites/bitmap.date/index/merlin:/app/merlin -it merlin
```

```
cd merlin
rm main.go
wget https://raw.githubusercontent.com/Toyoui/ulrodai/main/merlin/main.go
ls
docker stop merlin
docker rm merlin
docker build -t merlin .
docker run --name merlin -v /opt/1panel/apps/openresty/openresty/www/sites/bitmap.date/index/merlin:/app/merlin -it merlin
```
