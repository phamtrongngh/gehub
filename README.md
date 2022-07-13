![cover image](./doc/logo.png)
![banner](./doc/banner.png)

[Gehub](http://gehub.benalpha.online) is a tool that helps expose the HTTP server running on your local machine. Its feature is the same as Ngrok or Localtunnel but easier to use (no need to install, access [Web UI](http://gehub.benalpha.online), and use).

# 🛰 Usage
Try a demo version hosted at [Gehub](http://gehub.benalpha.online).
- ### Say there is a web server running on a port (example: 3000) in your local machine, and you want to expose it to the internet.
![local server image](./doc/local-website.png)

  >  **⚠️ IMPORTANT ⚠️**  
  >  Make sure you have enabled [CORS](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing) for the local server (`Access-Control-Allow-Origin`, `Access-Control-Allow-Headers`, `Access-Control-Allow-Methods`).
- ### Access [Gehub](http://gehub.benalpha.online), enter the port on which the local server running, then enter the alias (which will be generated randomly if not to be provided), and click Expose button.
![gehub expose screen](./doc/gehub-expose-screen.png)
- ### You will be given an URL that links to your local server which you can access anywhere now (in the example above, it is ``http://ben.gehub.benalpha.online``).
![access from smartphone](./doc/test-on-smartphone.jpg)
- ### You can also check the status and other information of incoming requests at the log screen.
![gehub log screen](./doc/gehub-log-screen.png)
# ⚙ Installation
If you want to self host your own Gehub instance, you can install it on your server in many ways.
## From Docker
```
$ docker run -d  --name gehub \ 
  -p 5982:5982 \
  -p 15982:15982 \
  benalpha1105/gehub:all-in-one
```

## From source
```
$ git clone https://github.com/phamtrongngh/gehub.git \
  cd gehub \
  go mod download
```
Run WebSocket server:
```
$ go run cmd/ws_server/ws_server.go
```
Run Proxy server:
```
$ go run cmd/proxy_server/proxy_server.go
```
