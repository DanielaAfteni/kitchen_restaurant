# kitchen_restaurant

This is the first lab at PR. It is related to the dining_hall_restaurant repository.

## To run the restaurant app with Docker

```bash
$ docker compose up --build
```
## To simply run the app

You need to change: `"dining_hall_url": "http://localhost:8080"` in `config/scfg.json`.

```bash
$ go run .
```