# wm-api

Запуск из директории сервиса, чтобы `LoadConfig` подхватил `wm-api.env` так же, как остальные сервисы проекта:

```bash
cd cmd/wm-api
go run .
```

Ручка:

```bash
GET /wm_api/?feed=<uuid>&group_by=<geo|date|site>&date_start=YYYY-MM-DD&date_end=YYYY-MM-DD
```

Пример:

```bash
curl "http://localhost:8050/wm_api/?feed=0260c40b-44ad-49bf-a54d-42d6bda5c90f&group_by=geo&date_start=2026-06-01&date_end=2026-06-18"
```

Ответ:

```json
[
  {
    "group_key": "US",
    "impressions": 10,
    "clicks": 2,
    "cost": 0.15
  }
]
```
