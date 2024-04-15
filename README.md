# ExtProc Tracing

```shell
export DD_API_KEY=xxx
export DD_APP_KEY=yyy

docker-compose build && docker-compose up

../envoy/bazel-bin/source/exe/envoy-static -c envoy.yml
```

```shell
while true; do curl http://127.0.0.1:10000; sleep 0.5; done
```
