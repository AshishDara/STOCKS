# Quick Reference Cheat Sheet

## Project Overview
- **Type**: Full-Stack Trading Dashboard
- **Backend**: Golang (Gin) + SQLite
- **Frontend**: React + Tailwind CSS
- **Auth**: JWT
- **Real-time**: WebSocket

## Tech Stack Comparison (MERN → This Project)

| MERN | This Project | Purpose |
|------|--------------|---------|
| Node.js | Golang | Backend runtime |
| Express.js | Gin | Web framework |
| MongoDB | SQLite | Database |
| Mongoose | GORM | ORM |
| JWT (same) | JWT | Authentication |

## Key Golang Concepts

### Variables
```go
var name string = "admin"
var age int = 25
count := 10  // Short declaration
```

### Structs (like Classes)
```go
type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
}
```

### Functions
```go
func createUser(c *gin.Context) {
    // Handler code
}
```

### Error Handling
```go
err := db.Create(&user).Error
if err != nil {
    return err
}
```

### Goroutines (Concurrency)
```go
go server.updatePrices()  // Runs in background
```

## SQLite vs MongoDB

| Feature | MongoDB | SQLite |
|---------|---------|--------|
| Type | NoSQL | SQL |
| Schema | Flexible | Fixed |
| Storage | Server | File |
| Queries | MongoDB syntax | SQL |

## API Endpoints

### Public
- `POST /api/login` - Login
- `POST /api/signup` - Register
- `GET /api/prices` - Get prices
- `WS /ws` - WebSocket

### Protected (Need JWT)
- `POST /api/orders` - Place order
- `GET /api/orders` - Get orders

## Database Schema

### Users
- `id` (PK)
- `username` (unique)
- `password` (hashed)

### Orders
- `id` (PK)
- `user_id` (FK → users)
- `symbol`
- `side` (buy/sell)
- `quantity`
- `price`
- `timestamp`

## JWT Flow

1. User logs in → Server generates JWT
2. Token stored in localStorage
3. Token sent in `Authorization: Bearer <token>` header
4. Middleware validates token
5. User ID extracted from token

## WebSocket Flow

1. Client connects: `new WebSocket('ws://localhost:8080/ws')`
2. Backend goroutine updates prices every 3 seconds
3. Server broadcasts to all clients
4. Frontend receives and updates UI

## Security Features

- ✅ Password hashing (bcrypt)
- ✅ JWT tokens
- ✅ User-specific data (filtered by user_id)
- ✅ Token expiration (24 hours)
- ✅ Input validation

## Common Interview Answers

**Q: Why Golang?**
- Great concurrency (goroutines)
- High performance
- Type safety

**Q: Why SQLite?**
- Simple, no setup
- Perfect for this scale
- Easy to migrate later

**Q: How does JWT work?**
- Token contains user info
- Signed with secret key
- Stateless (no server storage)

**Q: How do you ensure user isolation?**
- Every order has user_id
- Queries filter by user_id
- Backend enforces (not client)

## Key Code Snippets

### Signup Handler
```go
// Hash password
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

// Create user
user := User{Username: req.Username, Password: string(hashedPassword)}
db.Create(&user)

// Generate JWT
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": user.ID,
    "exp": time.Now().Add(time.Hour * 24).Unix(),
})
```

### WebSocket Broadcast
```go
func (s *Server) broadcastPrices() {
    prices := make([]Stock, 0, len(s.stocks))
    for _, stock := range s.stocks {
        prices = append(prices, *stock)
    }
    for client := range s.clients {
        client.WriteJSON(prices)
    }
}
```

### Price Updates (Goroutine)
```go
func (s *Server) updatePrices() {
    ticker := time.NewTicker(3 * time.Second)
    for range ticker.C {
        // Update prices
        // Broadcast to clients
    }
}
```

## Frontend Key Concepts

### WebSocket Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
    const prices = JSON.parse(event.data);
    setPrices(prices);
};
```

### Authenticated Request
```javascript
fetch('http://localhost:8080/api/orders', {
    headers: {
        'Authorization': `Bearer ${token}`
    }
});
```

### State Management
```javascript
const [token, setToken] = useState(localStorage.getItem('token'));
```

## Project Structure

```
STOCKS/
├── backend/
│   ├── main.go       # All backend code
│   ├── go.mod        # Dependencies
│   └── trading.db    # SQLite database
└── frontend/
    └── src/
        ├── App.jsx
        └── components/
            ├── Login.jsx
            ├── LivePricesTable.jsx
            ├── OrderForm.jsx
            └── OrdersTable.jsx
```

## Running the Project

### Backend
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Default Credentials
- Username: `admin`
- Password: `password123`

## Key Features

1. ✅ User registration and login
2. ✅ JWT authentication
3. ✅ Real-time price updates (WebSocket)
4. ✅ Place buy/sell orders
5. ✅ View order history
6. ✅ User-specific data isolation
7. ✅ Password hashing
8. ✅ Database persistence

## Interview Tips

1. **Explain your choices**: Why Golang? Why SQLite?
2. **Show understanding**: Explain JWT, WebSocket, database relationships
3. **Discuss trade-offs**: SQLite limitations, when to use PostgreSQL
4. **Talk improvements**: Tests, refresh tokens, rate limiting
5. **Be honest**: Admit learning Golang, show enthusiasm

## Common Questions

**Q: How does authentication work?**
- User logs in → Server validates → Returns JWT → Client stores token → Token sent with requests

**Q: How do you handle real-time updates?**
- Backend goroutine updates prices → Broadcasts via WebSocket → Frontend receives and updates UI

**Q: How do you ensure security?**
- Password hashing, JWT tokens, user-specific queries, input validation

**Q: What would you improve?**
- Add tests, refresh tokens, PostgreSQL for scale, rate limiting

