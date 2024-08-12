# Simple Reverse Proxy

This is a simple reverse proxy written in golang. Useful When you don't want to setup something like nginx and just want to get up and running in a dev environment.

## How to use

Create a `config.json` file in the same folder as the program.

```
go-reverse-proxy.exe
```

Alternatively, you can pass the path to the config file directory as the first command line argument.

```
go-reverse-proxy.exe C:/Users/Alireza/
```

## Sample config file

```json
{
  "listen": ":80",
  "https": true,
  "cert": "/etc/letsencrypt/live/example.com/fullchain.pem",
  "key": "/etc/letsencrypt/live/example.com/privkey.pem",
  "proxies": [
    {
      "listen": "/",
      "connect": "http://localhost:8080"
    },
    {
      "listen": "/api/",
      "connect": "http://localhost:5050"
    }
  ]
}
```

`listen` : the address you want the reverse proxy to listen on

`https` : whether or not to use SSL

`cert` : certificate location

`key` : key location

`proxies` : array of objects containing listen address and the destination for each proxy

You can pass as many proxies as you want to `proxies` array.
