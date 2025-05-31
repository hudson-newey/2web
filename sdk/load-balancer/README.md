# 2Web Load Balancer Config

By default, 2Web will use [Caddy](https://caddyserver.com/) for the load
balancer.

To run the load balancer, you can run the following command with the `caddyfile`
in the same directory as the `docker-compose.yml` file.

```sh
$ docker compose up load-balancer
>
```

Note that Caddy has `https://` by default, meaning that you shouldn't have to
configure HTTPS if you modify the `Caddyfile` with your domain.
"HTTPS by default" also includes serving through localhost and development
servers.
