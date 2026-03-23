# 🔗 url-shortener

A REST API service for shortening URLs, written in Go.

---

## 🛠 Stack

- **Go** — core language
- **chi** — HTTP router
- **MySQL** — storage
- **cleanenv** — config management
- **swagger** — API documentation
- **Docker** — containerization

---

## 📡 API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/url` | Create a short alias for a URL |
| `GET` | `/{alias}` | Redirect to the original URL |
| `DELETE` | `/url/{alias}` | Delete a short URL |

All endpoints are protected with **HTTP Basic Auth**.

Full Swagger docs available at `/swagger/index.html` after startup.

### ✂️ Create short URL

```http
POST /url
Content-Type: application/json
Authorization: Basic <credentials>

{
  "url": "https://example.com/very/long/url",
  "alias": "example"        // optional, auto-generated if omitted
}
```

```json
// 201 Created
{
  "alias": "example"
}
```

### 🔀 Redirect

```http
GET /example
// 307 Temporary Redirect → https://example.com/very/long/url
```

### 🗑 Delete

```http
DELETE /url/example
Authorization: Basic <credentials>
// 204 No Content
```

---

## ⚙️ Configuration

The service is configured via a YAML file. The path is set through the `CONFIG_PATH` environment variable.

```yaml
# config/local.yaml
env: "local"               # local | dev | prod
http_server:
  address: "0.0.0.0:5055"
  timeout: 4s
  idle_timeout: 60s
  user: "user"
service:
  alias_len: 6
```

🔐 Secrets are passed via environment variables (`.env` file for local development):

```env
CONFIG_PATH=config/local.yaml
HTTP_SERVER_PASSWORD=<bcrypt_hash>
MYSQL_ROOT_PASSWORD=<password>
```

---

## 🚀 Running locally

**🐳 With Docker (recommended):**

```bash
cp .env.example .env     # fill in your secrets
docker compose up --build
```

**Without Docker:**

```bash
# requires a running MySQL instance
go run ./cmd/url-shortener/main.go
```
---

## 🧪 Tests

```bash
go test ./...
```

Handlers are covered with unit tests using `gomock`.