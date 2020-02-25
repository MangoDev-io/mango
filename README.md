# algorand-asset-manager

### Setup

*Requirements*
1. `kmd` process of Algorand node must be running
2. `docker` and `docker-compose` must be installed
3. `ngrok` must be running on whatever port `kmd` is running on

*Setup*
1. Replace `psToken` in API `main.go`
2. Replace `kmdAddress` and `kmdToken` in API `main.go`

Run the following commands to start the project from the root dir:

```
docker-compose build
docker-compose up
```

The website will be running at port `:4200`.

To bring the project down, just type `docker-compose down` from the root dir.
