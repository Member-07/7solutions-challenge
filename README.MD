# Note

Lottery Search System (challenge 2) open file lottery-search-system.txt

## Design decisions
- MongoDB is accessed through a repository abstraction.
- Business logic is handled by services.
- HTTP concerns are handled by Gin handlers.
- Passwords are hashed using bcrypt and are never returned in API responses.
- JWT uses HMAC SHA-256 (HS256).
- Authentication failures return the same error to avoid revealing whether an email exists.
- User email has a unique MongoDB index.

## API Testing

Assume the application is running on:

```text
http://localhost:8080
```

### 1. Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

The response will contain a JWT token:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Set the token as an environment variable:

```bash
TOKEN="YOUR_JWT_TOKEN"
```

### 3. List Users

Requires authentication.

```bash
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Get User by ID

Requires authentication.

```bash
curl http://localhost:8080/api/users/USER_ID \
  -H "Authorization: Bearer $TOKEN"
```


### 5. Update User

Requires authentication.

```bash
curl -X PUT http://localhost:8080/api/users/USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com"
  }'
```

### 6. Delete User

Requires authentication.

```bash
curl -X DELETE http://localhost:8080/api/users/USER_ID \
  -H "Authorization: Bearer $TOKEN"
```

### 6. Create User

Requires authentication.

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### API Endpoints

| Method | Endpoint             | Authentication |
| ------ | -------------------- | -------------- |
| POST   | `/api/auth/register` | No             |
| POST   | `/api/auth/login`    | No             |
| GET    | `/api/users`         | JWT            |
| GET    | `/api/users/:id`     | JWT            |
| PUT    | `/api/users/:id`     | JWT            |
| DELETE | `/api/users/:id`     | JWT            |
| POST   | `/api/users`         | JWT            |


```

ลำดับในการทดสอบคือ **Register → Login → นำ JWT ไปใส่ `$TOKEN` → ทดสอบ User APIs** ครับ.
```
