# Replicate HTTP requests to many hosts

Can be useful to push some state to the list of upstreams using single endpoint.

This is how it works:

```
                        ┌─────────┐
                        │SAME DATA│
                        └─────────┘
                             │
                             │
                             │                     ┌──────────┐
                 ┌───────────┴───────────┐  ┌─────▶│ UPSTREAM │
                 │                       │  │      └──────────┘
                 │                       │  │
┌─────────────┐  ▼  ┌─────────────────┐  ▼  │      ┌──────────┐
│   REQUEST   │────▶│ HTTP REPLICATOR │─────┼─────▶│ UPSTREAM │
└─────────────┘     └─────────────────┘     │      └──────────┘
                                            │
                                            │      ┌──────────┐
                                            └─────▶│ UPSTREAM │
                                                   └──────────┘
```


Replicate requests from `127.0.0.1:9999` to:

* http://one.prod:13000
* http://two.prod:13000
* http://three.prod:13000

```
docker run -rm -it -p 127.0.0.1:9999:9999 bobrik/http-replicator --listen :9999 \
    --upstreams http://one.prod:13000,http://two.prod:13000,http://three.prod:13000
```
