# Mango

Previously named Algorand Asset Manager aka AAM, which translates to Mango from Hindi.

### Setup

_Requirements_

1. `kmd` process of Algorand node must be running
2. `docker` and `docker-compose` must be installed
3. `ngrok` must be running on whatever port `kmd` is running on

_Setup_

1. Replace `psToken` in API `main.go`
2. Replace `kmdAddress` and `kmdToken` in API `main.go`

To start and stop `kmd` (from within the node directory):

```
./goal kmd start -t 3600 -d <datadir>
./goal kmd stop -d <datadir>
```

Run the following commands to start the project from the root dir:

```
docker-compose build
docker-compose up
```

The website will be running at port `:4200`.

To bring the project down, just type `docker-compose down` from the root dir.
