# netplay-lobby-server-go

Netplay lobby server written in GO. Needs Go v1.13.

## Deployment

```bash
go build
./netplay-lobby-server-go
```

## Configuration
Rename the ```config/lobby.template.yaml``` to ```config/lobby.yaml``` and place the configuration file in one of the
following directories:

 - /etc/lobby
 - $HOME/.lobby
 - ./config

## LICENSE

The server itself is licensed under AGPLv3.

This product includes GeoLite2 data created by MaxMind, available from
[https://www.maxmind.com](https://www.maxmind.com)
