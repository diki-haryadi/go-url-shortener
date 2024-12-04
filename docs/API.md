# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Endpoints

### Create Short URL

```http
POST /urls
Content-Type: application/json

{
    "url": "https://example.com/very/long/url",
    "custom": "my-custom-path"  // optional
}
```

```curl
curl -X POST "http://localhost:8080/api/v1/shorten" \
-H "Content-Type: application/json" \
-d '{
  "url": "https://example.com",
  "short": "custom123",
  "expiry": 3600
}'

```

#### Response
```json
{
    "short_url": "http://localhost:8080/abc123",
    "original_url": "https://example.com/very/long/url",
    "created_at": "2024-12-04T12:00:00Z",
    "expires_at": null
}
```

### Resolve URL

```http
GET /{shortURL}
```

#### Response
- 307/308 Redirect to original URL
- Headers include original URL

### Get URL Statistics

```http
GET /urls/{shortURL}/stats
```

#### Response
```json
{
    "short_url": "abc123",
    "original_url": "https://example.com/very/long/url",
    "created_at": "2024-12-04T12:00:00Z",
    "visits": 42,
    "last_visit": "2024-12-04T15:30:00Z"
}
```

### Health Check

```http
GET /health
```

#### Response
```json
{
    "status": "UP",
    "redis": "UP",
    "version": "1.0.0",
    "timestamp": "2024-12-04T12:00:00Z"
}
```

### Metrics

```http
GET /metrics
```

#### Response
Prometheus formatted metrics

## Error Responses

### 400 Bad Request
```json
{
    "error": "Invalid URL format",
    "code": "INVALID_URL",
    "details": "The provided URL is not valid"
}
```

### 404 Not Found
```json
{
    "error": "URL not found",
    "code": "NOT_FOUND",
    "details": "The requested short URL does not exist"
}
```

### 429 Too Many Requests
```json
{
    "error": "Rate limit exceeded",
    "code": "RATE_LIMIT",
    "details": "Please try again later"
}
```

## Rate Limits

- 100 requests per minute per IP for URL creation
- 1000 requests per minute per IP for URL resolution
- Custom rate limits can be configured

## Headers

### Request Headers
- `Content-Type`: application/json
- `Accept`: application/json
- `X-Real-IP`: Client IP (optional)
- `User-Agent`: Client identification (optional)

### Response Headers
- `X-RateLimit-Limit`: Rate limit ceiling
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Rate limit reset time