[![Docker Image CI](https://github.com/MihaiBlebea/trading-platform/actions/workflows/docker-image.yml/badge.svg?branch=master&event=push)](https://github.com/MihaiBlebea/trading-platform/actions/workflows/docker-image.yml)

## How to install?

- Add env variables to the .env file

	- **HTTP_PORT**

	- **POSTGRES_HOST**

	- **POSTGRES_USER**

	- **POSTGRES_PASSWORD**

	- **POSTGRES_DB**

	- **POSTGRES_PORT**

- Run docker-compose command

	- `docker-compose build && docker-compose up -d`

## How to use the API endpoints?

### Create new account:

- **POST** `/api/v1/account`

Create a new account and receive back the API Token

Success response example:

```json
{
	"success": true,
	"account": {
		"api_token": "BXhbwSF6kf", // API TOKEN
		"balance": 10000,
		"pending_balance": 0,
		"created_at": "2022-03-20T14:28:57.1174248Z"
	}
}
```

Fail response example:

```json
{
	"success": false,
	"error": "Something went wrong"
}
```


### Request an account:

- **GET** `/api/v1/account`

Retrieve information about an existing account.

Headers:

```json
{
	"Authorization": "Bearer <Api Token>"
}
```

Success response example:

```json
{
	"success": true,
	"account": {
		"api_token": "BXhbwSF6kf", // API TOKEN
		"balance": 10000,
		"pending_balance": 0,
		"created_at": "2022-03-20T14:28:57.1174248Z"
	}
}
```

Fail response example:

```json
{
	"success": false,
	"error": "Something went wrong"
}
```


### Place an order

- **POST** `/api/v1/order`

Create a new order, it can be a buy or sell order.

Headers:

```json
{
	"Authorization": "Bearer <Api Token>"
}
```

Success response example:

```json
{
    "success": true,
    "order": {
        "id": 8,
        "type": "limit",
        "status": "pending",
        "direction": "buy",
        "amount": 1000,
        "fill_price": 0,
        "amount_after_fill": 0,
        "symbol": "AAPL",
        "quantity": 0,
        "created_at": "2022-03-19T13:43:46.470965969Z"
    }
}
```

Fail response example:

```json
{
	"success": false,
	"error": "Something went wrong"
}
```


### Retrieve a list of existing orders

- **GET** `/api/v1/orders`

Get a list of existing orders for this account.

They can be pending, filled or cancelled orders.

- `amount_after_fill` key will be present in response only if the order has been filled.

- `filled_at` key is only present in response if order is filed.

Headers:

```json
{
	"Authorization": "Bearer <Api Token>"
}
```

Success response example:

```json
{
	"success": true,
	"orders": [
		{
			"id": 5,
			"type": "limit",
			"status": "pending",
			"direction": "buy",
			"amount": 1000,
			"fill_price": 164.34,
			"symbol": "AAPL",
			"quantity": 6,
			"created_at": "2022-03-19T12:35:39.280604Z"
		},
		{
			"id": 6,
			"type": "limit",
			"status": "filled",
			"direction": "sell",
			"amount": 1000,
			"fill_price": 164.23,
			"amount_after_fill": 985.38,
			"symbol": "AAPL",
			"quantity": 6,
			"filled_at": "2022-03-19T12:37:13.294119Z",
			"created_at": "2022-03-19T12:36:37.531223Z"
		}
    ]
}
```

Fail response example:

```json
{
	"success": false,
	"error": "Something went wrong"
}
```
