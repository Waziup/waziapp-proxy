# WaziApp Proxy

WaziApp proxy is a simple http proxy that is intended to listen on the WaziApp unix socket `/var/lib/waziapp/proxy.sock` and forwards to a local webserver. Easy spoken, this tool converts a _normal_ ip:port socket into a unix domain socket file.

More general/similar software includes [`socat`](https://linux.die.net/man/1/socat).

Most webservers can be configured to listen on a unix socket already, but some software can only listen on a `ip:port` socket - in that cas this tool can help you.

Here is a WaziApp proxy command with the `socat` equivalent:

```sh
socat TCP:localhost:12346 UNIX-LISTEN:/var/lib/waziapp/proxy.sock

waziapp-proxy localhost:12346
```

Both commands will create a `proxy.sock` file that forwards to `localhost:12346`.

```sh
curl --unix-socket proxy.sock localhost/index.html

# ... is now the same as ...

curl localhost:12346/index.html
```

In difference to `socat`, this software

- does not stop after the first connection is served,
- can add path prefixes like `/prefix`,
- has a clean log output to see the forwarded requests,
- can only forward from a unix socket to a ip:port socket,
- can't be configured that much...

## WaziApp's UI

WaziApps can have a UI by serving some HTML files with a simple (maybe static) webserver, which they must serve on the unix domain socket named above. This socket is made available to the WaziGate by mapping the `/var/lib/waziapp` folder between the app (docker container) and the WaziGate (docker host). It is not possible to serve a UI on a port (e.g. 80 or 8080) of the container, as the port is not bridged to host and is therefore not available from the outside.

## Usage

```txt
Usage: ./waziapp-proxy {addr}
Create a unix domain socket (for use inside WaziApps) forwarding to local address.
Example:
  ./waziapp-proxy http://localhost:8080/test
  unix:/var/lib/waziapp/proxy.sock --> http://localhost:8080/test
Use env WAZIAPP_ADDR to override the default waziapp proxy socket address.
```


