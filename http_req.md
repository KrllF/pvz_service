
## Как делать запросы?

Т.к. есть basic_auth, то всегда нужно вводить логин и пароль.

## 1. Принять заказ от курьера
```bash
curl -u "admin:adminpassword" -i -X POST "http://localhost:9000/accept/" \
-H 'Content-Type: application/json' \
-d '{"order_id": 1, "user_id": 100, "weight": 100, "price": 100, "pack_type": "film", "extra_pack": "film", "shelf_life": "2026-10-10T15:15:10Z"}'
```

## 2. Вернуть заказ курьеру
```bash
curl -u "admin:adminpassword" -i -X DELETE "http://localhost:9000/return/1000" 
```

## 3. Выдать заказы и принять возвраты клиента
```bash
curl -u "admin:adminpassword" -i -X PUT "http://localhost:9000/process" \
-H 'Content-Type: application/json' \
-d '{"user_id": 100, "operation": "issue", "order_ids": [1, 10, 100, 12, 2]}'

curl -u "admin:adminpassword" -i -X PUT "http://localhost:9000/process" \
-H 'Content-Type: application/json' \
-d '{"user_id": 100, "operation": "return", "order_ids": [100, 2, 3]}'

```

## 4. Получить список заказов для пользователя
```bash
curl -u "admin:adminpassword" -i -X GET "http://localhost:9000/history/orders/100?in_pvz=true&n=10&limit=1&page=1&pag=false&finder=true&order=1"
```

## 5. Получить список возвратов
```bash
curl -u "admin:adminpassword" -i -X GET "http://localhost:9000/history/returns?limit=2&page=3&finder=true&order=1"
```

## 6. Получить историю заказов
```bash
curl -u "admin:adminpassword" -i -X GET "http://localhost:9000/history/"
```

## 7. Принять заказы от курьера
```bash
curl -u "admin:adminpassword" -i -X GET "http://localhost:9000/accept/bulk" \
-H 'Content-Type: application/json' \
-d '{
  "path_json": "tor/por/cor/storage.json"
}'
```