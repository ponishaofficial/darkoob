TestAPI
=======

```shell
docker build -t testapi:latest .
```

```shell
docker run --name testapi -v $PWD/scenarios:/app/scenarios -d testapi:latest
```