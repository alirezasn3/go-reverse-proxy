# Simple Reverse Proxy

This is a simple reverse proxy written in golang. Useful When you don't want to setup something like nginx and just want to get up and running in dev inviroment.

## How to use

Create a `config.json` file in the same folder as the program.

```
go-reverse-proxy.exe
```

Alternatively, you can pass the path to the config file as the first command line argument.

```
go-reverse-proxy.exe C:/Users/Alireza/config.json
```

## Sample config file

```json
{
  "listen": ":80",
  "proxies": [
    {
      "listen": "localhost/",
      "connect": "http://localhost:8080"
    },
    {
      "listen": "api.localhost/",
      "connect": "http://localhost:5050"
    }
  ]
}
```

`listen` : the address you want the reverse proxy to listen on

`proxies` : array of objects containing listen address and the destination for each proxy

You can pass as many proxies as you want to `proxies` array.