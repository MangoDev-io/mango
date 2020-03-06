# Mango

Previously named Algorand Asset Manager aka AAM, which translates to Mango from Hindi.

### Setup

_Requirements_

1. `docker` and `docker-compose` must be installed

_Setup_

1. Replace `psToken` in API `main.go`

Run the following commands to start the project from the root dir:

```
docker-compose down
docker-compose build
docker-compose up
```

The website will be running at port `:4200`.

To bring the project down, just type `docker-compose down` from the root dir.
