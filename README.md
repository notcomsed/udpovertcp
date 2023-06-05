# udpovertcp
tunnel udp packge over tcp connect

## How to use

#### start connect-server on remote

```bash
connect -s 192.168.1.1 8081 127.0.0.1 5353 
```

It will Listening tcp on 192.168.1.1:8081

and Decode tcp stream data to udp package & send udp to 127.0.0.1:5353

#### start bind-server on local

```bash
bind 0.0.0.0:8080 192.168.1.1:8081
```
It will bind udp on 0.0.0.0:8080 

And, it will Encode package to tcp stream & send to 192.168.1.1:8081 as tcp

---

## TEST

1.start nc in remote

```bash
nc -l -u -p 5353
```

2.start connect-server on remote

```bash
connect -s 0.0.0.0 8081 127.0.0.1 5353 
```

3.start bind-server on local

```bash
bind 127.0.0.1:8080 192.168.1.1:8081
```

4.start nc in local

```bash
nc -u 127.0.0.1 8080
```

5.enter key on nc

---

## design

```
udp -> bind-server -> tcp[tcp mux] -> connect-server ->udp
```

### package Encode

```
tcp stream
**************************************************
	|length1|Udp payload1|length2|Udp payload2|....
	|  2    |   length   |  2    |   length   |....
*************************************************
```

udp package length should't more that 1534
