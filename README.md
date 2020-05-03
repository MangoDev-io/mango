# Mango

Previously named Algorand Asset Manager aka AAM, which translates to Mango from Hindi.

### Setup

_Requirements_

1. `docker` and `docker-compose` must be installed

_Setup_

1. Create a `.env` file in `/api/`
2. Set the following environment variables in the file

```
API_TESTNETALGODADDRESS=<purestake node address>
API_MAINNETALGODADDRESS=<purestake node address>
API_PSTOKEN=<purestake token>

API_TOKENAUTHPASSWORD=
```

3. Edit the `docker-compose.yml` environment variables for `db` to set the Postgres database initialization configuration

4. Update the baseURL in the `web/src/app/state.service.ts` to `localhost:5000` if running on localhost, or to your hosted API address

5. Run the following commands to start the project from the root dir:

```
docker-compose build
docker-compose up
```

The website will be running at port `:4200`.

To bring the project down, just type `docker-compose down` from the root dir.
