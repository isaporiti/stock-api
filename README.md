# Stock Prices API

Basic Web API that serves stub stock prices for [**some supported tickers**](#supported-tickers) and users' stock portfolios.

## How to run

### With Docker
Having Docker installed, you can run
```bash
# using Makefile
make docker_build
make docker_run

# or otherwise
docker build -t stock-api .
docker run --rm -p 5001:5001 stock-api
```

### With Go
```bash
go run cmd/main.go
```

## How to try

There are two authorized usernames: `testA` and `testB`. You can try the endpoints with the following:

### User Stocks
```bash
# user testA
curl --request GET \
  --url http://localhost:5001/tickers \
  --header 'Authorization: Basic dGVzdEE6'

# user testB
curl --request GET \
  --url http://localhost:5001/tickers \
  --header 'Authorization: Basic dGVzdEI6'
```

### Stocks history
``` bash
curl --request GET \
  --url http://localhost:5001/tickers/FB/history \
  --header 'Authorization: Basic dGVzdEI6'
```

### Unauthorized
```bash
curl --request GET \
  --url http://localhost:5001/tickers \
  --header 'Authorization: Basic dW5rbm93bjo='
```

## Supported Tickers
You can query the history price for these tickers. Querying for anything else will result in a `404 Not found` response.
- `AAPL`
- `MSFT`
- `GOOG`
- `AMZN`
- `FB`
- `TSLA`
- `NVDA`
- `JPM`
- `BABA`
- `JNJ`
- `WMT`
- `PG`
- `PYPL`
- `DIS`
- `ADBE`
- `PFE`
- `V`
- `MA`
- `CRM`
- `NFLX`
