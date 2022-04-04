# Simple Load Balancer

Round-Robin load balance implementation in golang with ative and passive health check.

## Running

To simplify the process, we use a docker-compose to up all containers. To test only open a browser new tab on ``` localhost:8080``` and you will send a request and one of our servers will respond.

```
docker-compose up
Starting simple-load-balancer_api-2_1 ... done
Starting simple-load-balancer_api-3_1 ... done
Starting simple-load-balancer_api-1_1 ... done
Recreating simple-load-balancer_load-balancer_1 ... done
Attaching to simple-load-balancer_api-2_1, simple-load-balancer_api-3_1, simple-load-balancer_api-1_1, simple-load-balancer_load-balancer_1
load-balancer_1  | 2022/04/04 04:15:12 Load Balancer started at :8080
load-balancer_1  | 2022/04/04 04:15:12 http://api-1:8081 [up]
load-balancer_1  | 2022/04/04 04:15:12 http://api-2:8082 [up]
load-balancer_1  | 2022/04/04 04:15:12 http://api-3:8083 [up]
```
