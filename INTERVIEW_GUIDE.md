# Trading Dashboard - Complete Interview Guide

## Table of Contents
1. [Project Overview](#project-overview)
2. [Tech Stack Deep Dive](#tech-stack-deep-dive)
3. [Architecture Overview](#architecture-overview)
4. [Backend Code Walkthrough](#backend-code-walkthrough)
5. [Frontend Code Walkthrough](#frontend-code-walkthrough)
6. [Database Design](#database-design)
7. [API Documentation](#api-documentation)
8. [Key Features Implementation](#key-features-implementation)
9. [Interview Preparation Q&A](#interview-preparation-qa)
10. [Common Interview Questions](#common-interview-questions)

---

## Project Overview

### What is This Project?
A **real-time trading dashboard** that simulates a live stock trading environment. Users can:
- View real-time stock prices (updates every 3 seconds)
- Place buy/sell orders
- View their order history
- Create accounts and authenticate securely

### Project Type
**Full-Stack Application** with:
- **Backend**: Golang (Go) - RESTful API + WebSocket server
- **Frontend**: React - Single Page Application (SPA)
- **Database**: SQLite - Persistent storage
- **Authentication**: JWT (JSON Web Tokens)

### Why This Tech Stack?
- **Golang**: High performance, excellent concurrency support (goroutines)
- **React**: Modern, component-based UI framework
- **SQLite**: Lightweight, file-based database (perfect for this project)
- **JWT**: Stateless authentication (scalable, secure)

---

## Tech Stack Deep Dive

### For MERN Developers: Understanding Golang

If you're coming from **MERN (MongoDB, Express, React, Node.js)**, here's how Golang compares:

#### Node.js vs Golang

| Aspect | Node.js (MERN) | Golang |
|--------|----------------|--------|
| **Language** | JavaScript | Go (statically typed) |
| **Concurrency** | Event loop (single-threaded) | Goroutines (multi-threaded) |
| **Performance** | Fast for I/O operations | Very fast for CPU-bound tasks |
| **Type System** | Dynamic typing | Static typing (compile-time checks) |
| **Package Manager** | npm | go mod |
| **Framework** | Express.js | Gin (similar to Express) |

#### Key Golang Concepts You Need to Know

**1. Variables and Types**
```go
// Golang (statically typed)
var username string = "admin"
var age int = 25

// JavaScript equivalent
let username = "admin";
let age = 25;
```

**2. Structs (Similar to Objects/Classes)**
```go
// Golang struct
type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
}

// JavaScript equivalent
class User {
    constructor(id, username, password) {
        this.id = id;
        this.username = username;
        this.password = password;
    }
}
```

**3. Functions**
```go
// Golang function
func createOrder(c *gin.Context) {
    // Handler logic
}

// JavaScript equivalent
function createOrder(req, res) {
    // Handler logic
}
```

**4. Goroutines (Concurrency)**
```go
// Golang - runs in background
go server.updatePrices()

// JavaScript equivalent
setInterval(() => {
    updatePrices();
}, 3000);
```

**5. Error Handling**
```go
// Golang - explicit error handling
err := db.Create(&user).Error
if err != nil {
    return err
}

// JavaScript equivalent
try {
    await db.create(user);
} catch (err) {
    return err;
}
```

### Understanding SQLite (vs MongoDB)

| Aspect | MongoDB (MERN) | SQLite |
|--------|----------------|--------|
| **Type** | NoSQL (Document-based) | SQL (Relational) |
| **Storage** | Database server | Single file |
| **Schema** | Flexible (schema-less) | Fixed schema (tables) |
| **Queries** | MongoDB queries | SQL queries |
| **Relations** | Manual references | Foreign keys |
| **Best For** | Large scale, flexible data | Small-medium apps, embedded |

#### SQLite Basics

**Tables (like MongoDB Collections)**
```sql
-- SQLite table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

-- MongoDB equivalent
db.users.insertOne({
    username: "admin",
    password: "hashed"
});
```

**GORM (Go ORM - like Mongoose)**
```go
// GORM - defines model
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}

// Mongoose equivalent
const userSchema = new mongoose.Schema({
    username: { type: String, unique: true, required: true },
    password: { type: String, required: true }
});
```

### Gin Framework (vs Express.js)

**Express.js (Node.js)**
```javascript
const express = require('express');
const app = express();

app.post('/api/orders', (req, res) => {
    // Handler
});
```

**Gin (Golang)**
```go
r := gin.Default()

r.POST("/api/orders", func(c *gin.Context) {
    // Handler
})
```

**Key Similarities:**
- Both use middleware
- Both support route groups
- Both handle JSON automatically
- Both support CORS

---

## Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    React Frontend                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │  Login   │  │  Prices │  │  Orders  │  │   Form   │ │
│  │ Component│  │  Table  │  │  Table   │  │Component │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────┘
         │                    │                    │
         │ HTTP/REST          │ WebSocket          │ HTTP/REST
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│              Golang Backend (Gin Framework)             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ REST API     │  │ WebSocket    │  │ JWT Auth     │  │
│  │ Handlers     │  │ Handler      │  │ Middleware   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│         │                    │                    │     │
│         └────────────────────┴────────────────────┘     │
│                            │                            │
│                            ▼                            │
│  ┌──────────────────────────────────────────────┐      │
│  │         GORM (ORM Layer)                     │      │
│  └──────────────────────────────────────────────┘      │
│                            │                            │
└────────────────────────────┼────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   SQLite DB     │
                    │  (trading.db)   │
                    │                 │
                    │  - users        │
                    │  - orders       │
                    └─────────────────┘
```

### Data Flow

**1. User Registration/Login Flow:**
```
User → Frontend (Login Form) 
    → POST /api/signup or /api/login
    → Backend (Validates, Hashes Password)
    → Database (Creates/Verifies User)
    → Backend (Generates JWT Token)
    → Frontend (Stores Token in localStorage)
    → User Authenticated
```

**2. Real-Time Price Updates Flow:**
```
Backend (Goroutine) 
    → Updates Stock Prices Every 3 Seconds
    → Broadcasts via WebSocket
    → All Connected Clients Receive Updates
    → Frontend Updates UI (Green/Red Colors)
```

**3. Order Placement Flow:**
```
User → Frontend (Order Form)
    → POST /api/orders (with JWT Token)
    → Backend (Validates Token)
    → Backend (Validates Order Data)
    → Database (Saves Order with User ID)
    → Backend (Returns Created Order)
    → Frontend (Refreshes Orders Table)
```

---

## Backend Code Walkthrough

### Project Structure
```
backend/
├── main.go          # Main server file (all code in one file for simplicity)
├── go.mod           # Dependencies
└── trading.db       # SQLite database (auto-created)
```

### Key Components Explained

#### 1. Data Models

```go
// User Model - represents a user in the database
type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    Username string `gorm:"unique;not null" json:"username"`
    Password string `gorm:"not null" json:"-"`  // "-" means don't return in JSON
}
```

**Explanation:**
- `gorm:"primaryKey"` - Sets ID as primary key (auto-increment)
- `gorm:"unique"` - Username must be unique
- `json:"-"` - Password never sent to client (security)

```go
// Order Model - represents a trading order
type Order struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`  // Foreign key
    Symbol    string    `gorm:"not null" json:"symbol"`
    Side      string    `gorm:"not null" json:"side"`     // "buy" or "sell"
    Quantity  int       `gorm:"not null" json:"quantity"`
    Price     float64   `gorm:"not null" json:"price"`
    Timestamp time.Time `gorm:"not null" json:"timestamp"`
}
```

**Key Point:** `UserID` links orders to users (one-to-many relationship)

#### 2. Server Initialization

```go
func NewServer() *Server {
    // 1. Connect to SQLite database
    db, err := gorm.Open(sqlite.Open("trading.db"), &gorm.Config{})
    
    // 2. Auto-migrate (creates tables if they don't exist)
    db.AutoMigrate(&User{}, &Order{})
    
    // 3. Create default admin user if database is empty
    var userCount int64
    db.Model(&User{}).Count(&userCount)
    if userCount == 0 {
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
        defaultUser := User{
            Username: "admin",
            Password: string(hashedPassword),
        }
        db.Create(&defaultUser)
    }
    
    // 4. Initialize mock stocks
    stocks := map[string]*Stock{
        "AAPL": {Symbol: "AAPL", Price: 175.50},
        "TSLA": {Symbol: "TSLA", Price: 245.30},
        // ... more stocks
    }
    
    return &Server{
        db:      db,
        stocks:  stocks,
        clients: make(map[*websocket.Conn]bool),
    }
}
```

**Key Concepts:**
- **Database Migration**: Automatically creates/updates database schema
- **bcrypt**: Password hashing algorithm (one-way encryption)
- **Map**: Go's key-value data structure (like JavaScript object)

#### 3. JWT Authentication Middleware

```go
func (s *Server) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get Authorization header
        authHeader := c.GetHeader("Authorization")
        
        // 2. Extract token from "Bearer <token>"
        tokenString := authHeader[7:]  // Skip "Bearer "
        
        // 3. Parse and validate JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        
        // 4. Extract user ID from token claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            userID := claims["user_id"].(float64)
            c.Set("user_id", uint(userID))  // Store in context
        }
        
        c.Next()  // Continue to next handler
    }
}
```

**How It Works:**
1. Client sends: `Authorization: Bearer <token>`
2. Middleware validates token signature
3. Extracts user ID from token
4. Stores user ID in request context
5. Protected routes can access user ID

#### 4. Signup Endpoint

```go
func (s *Server) signup(c *gin.Context) {
    // 1. Parse request body
    var req SignupRequest
    c.ShouldBindJSON(&req)
    
    // 2. Validate input
    if len(req.Password) < 6 {
        c.JSON(400, gin.H{"error": "Password must be at least 6 characters"})
        return
    }
    
    // 3. Check if username exists
    var existingUser User
    if s.db.Where("username = ?", req.Username).First(&existingUser).Error == nil {
        c.JSON(400, gin.H{"error": "Username already exists"})
        return
    }
    
    // 4. Hash password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    // 5. Create user
    user := User{
        Username: req.Username,
        Password: string(hashedPassword),
    }
    s.db.Create(&user)
    
    // 6. Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })
    tokenString, _ := token.SignedString(jwtSecret)
    
    // 7. Return token and user
    c.JSON(201, LoginResponse{
        Token: tokenString,
        User:  user,
    })
}
```

**Key Points:**
- **bcrypt**: One-way hashing (can't reverse password)
- **JWT**: Contains user info, expires in 24 hours
- **Status 201**: Created (for new resources)

#### 5. WebSocket Implementation

```go
func (s *Server) handleWebSocket(c *gin.Context) {
    // 1. Upgrade HTTP connection to WebSocket
    conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
    
    // 2. Register client
    s.clientsLock.Lock()
    s.clients[conn] = true
    s.clientsLock.Unlock()
    
    // 3. Send initial prices
    s.sendPricesToClient(conn)
    
    // 4. Keep connection alive
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break  // Client disconnected
        }
    }
    
    // 5. Unregister on disconnect
    s.clientsLock.Lock()
    delete(s.clients, conn)
    s.clientsLock.Unlock()
}
```

**Why WebSocket?**
- **HTTP**: Client requests → Server responds (one-way)
- **WebSocket**: Full-duplex (server can push updates anytime)

#### 6. Real-Time Price Updates (Goroutines)

```go
func (s *Server) updatePrices() {
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    ticker := time.NewTicker(3 * time.Second)  // Every 3 seconds
    
    for range ticker.C {
        // Update each stock price
        for symbol, stock := range s.stocks {
            // Random change: -2% to +2%
            changePercent := (rng.Float64()*4 - 2) / 100
            newPrice := stock.Price * (1 + changePercent)
            stock.Price = newPrice
        }
        
        // Broadcast to all connected clients
        s.broadcastPrices()
    }
}
```

**Goroutines Explained:**
- `go server.updatePrices()` - Runs in background
- Doesn't block main thread
- Similar to JavaScript `setInterval()` but more efficient

---

## Frontend Code Walkthrough

### Project Structure
```
frontend/
├── src/
│   ├── components/
│   │   ├── Login.jsx          # Authentication
│   │   ├── LivePricesTable.jsx  # Real-time prices
│   │   ├── OrderForm.jsx      # Place orders
│   │   └── OrdersTable.jsx    # Order history
│   ├── App.jsx                # Main component
│   └── main.jsx               # Entry point
```

### Key Components

#### 1. App.jsx - Main Component

```javascript
function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [token, setToken] = useState(null);
    const [user, setUser] = useState(null);
    const [orders, setOrders] = useState([]);

    // Check for existing token on mount
    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        const storedUser = localStorage.getItem('user');
        if (storedToken && storedUser) {
            setToken(storedToken);
            setUser(JSON.parse(storedUser));
            setIsAuthenticated(true);
        }
    }, []);

    // Fetch orders when authenticated
    const fetchOrders = async () => {
        const response = await fetch('http://localhost:8080/api/orders', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        const data = await response.json();
        setOrders(data);
    };

    // Show login if not authenticated
    if (!isAuthenticated) {
        return <Login onLogin={handleLogin} />;
    }

    // Show dashboard if authenticated
    return (
        <div>
            <LivePricesTable />
            <OrderForm token={token} />
            <OrdersTable orders={orders} />
        </div>
    );
}
```

**Key Concepts:**
- **useState**: Manages component state
- **useEffect**: Runs on component mount (like componentDidMount)
- **localStorage**: Browser storage (persists across refreshes)

#### 2. LivePricesTable.jsx - WebSocket Connection

```javascript
function LivePricesTable() {
    const [prices, setPrices] = useState([]);
    const [previousPrices, setPreviousPrices] = useState({});
    const wsRef = useRef(null);

    useEffect(() => {
        // Connect to WebSocket
        const ws = new WebSocket('ws://localhost:8080/ws');
        wsRef.current = ws;

        // Handle incoming messages
        ws.onmessage = (event) => {
            const newPrices = JSON.parse(event.data);
            
            // Store previous prices for comparison
            setPrices((currentPrices) => {
                const newPrev = {};
                currentPrices.forEach((stock) => {
                    newPrev[stock.symbol] = stock.price;
                });
                setPreviousPrices(newPrev);
                return newPrices;
            });
        };

        // Cleanup on unmount
        return () => {
            ws.close();
        };
    }, []);

    // Determine color based on price change
    const getPriceColor = (symbol, currentPrice) => {
        const change = getPriceChange(symbol, currentPrice);
        if (change > 0) return 'text-green-600';  // Up
        if (change < 0) return 'text-red-600';   // Down
        return 'text-gray-700';                  // No change
    };
}
```

**Key Points:**
- **WebSocket API**: Browser's built-in WebSocket client
- **useRef**: Stores WebSocket connection (persists across renders)
- **Color Coding**: Green (up), Red (down), Gray (no change)

#### 3. OrderForm.jsx - Placing Orders

```javascript
const handleSubmit = async (e) => {
    e.preventDefault();
    
    const response = await fetch('http://localhost:8080/api/orders', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,  // JWT token
        },
        body: JSON.stringify({
            symbol: formData.symbol.toUpperCase(),
            side: formData.side,
            quantity: parseInt(formData.quantity),
            price: parseFloat(formData.price),
        }),
    });

    if (response.ok) {
        // Success - refresh orders
        onOrderPlaced();
    }
};
```

**Key Points:**
- **JWT in Header**: Token sent with every request
- **Form Validation**: Client-side validation before submission
- **Error Handling**: Shows user-friendly error messages

---

## Database Design

### Schema Overview

#### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
```

**Fields:**
- `id`: Primary key (auto-increment)
- `username`: Unique identifier (indexed for fast lookup)
- `password`: Hashed with bcrypt (never stored in plain text)

#### Orders Table
```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    symbol TEXT NOT NULL,
    side TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    price REAL NOT NULL,
    timestamp DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**Relationships:**
- **One-to-Many**: One user can have many orders
- **Foreign Key**: `user_id` references `users.id`
- **Indexing**: `user_id` indexed for fast queries

### Query Examples

**Get all orders for a user:**
```go
db.Where("user_id = ?", userID).Order("timestamp DESC").Find(&orders)
```

**SQL Equivalent:**
```sql
SELECT * FROM orders 
WHERE user_id = ? 
ORDER BY timestamp DESC;
```

**Check if username exists:**
```go
db.Where("username = ?", username).First(&user)
```

**SQL Equivalent:**
```sql
SELECT * FROM users 
WHERE username = ? 
LIMIT 1;
```

---

## API Documentation

### Public Endpoints

#### POST /api/login
**Purpose:** Authenticate existing user

**Request:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin"
  }
}
```

**Error (401):**
```json
{
  "error": "Invalid credentials"
}
```

#### POST /api/signup
**Purpose:** Create new user account

**Request:**
```json
{
  "username": "newuser",
  "password": "password123"
}
```

**Validation:**
- Username must be unique
- Password must be at least 6 characters

**Response (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 2,
    "username": "newuser"
  }
}
```

#### GET /api/prices
**Purpose:** Get current stock prices (public, no auth required)

**Response (200):**
```json
[
  {"symbol": "AAPL", "price": 175.50},
  {"symbol": "TSLA", "price": 245.30},
  {"symbol": "AMZN", "price": 138.20},
  {"symbol": "INFY", "price": 18.75},
  {"symbol": "TCS", "price": 3450.00}
]
```

### Protected Endpoints (Require JWT Token)

#### POST /api/orders
**Purpose:** Place a new trading order

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request:**
```json
{
  "symbol": "AAPL",
  "side": "buy",
  "quantity": 10,
  "price": 175.50
}
```

**Validation:**
- `side` must be "buy" or "sell"
- `quantity` must be positive
- `price` must be positive

**Response (201):**
```json
{
  "id": 1,
  "user_id": 1,
  "symbol": "AAPL",
  "side": "buy",
  "quantity": 10,
  "price": 175.50,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### GET /api/orders
**Purpose:** Get all orders for authenticated user

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200):**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "symbol": "AAPL",
    "side": "buy",
    "quantity": 10,
    "price": 175.50,
    "timestamp": "2024-01-15T10:30:00Z"
  }
]
```

**Security:** Only returns orders for the authenticated user (filtered by `user_id`)

### WebSocket Endpoint

#### WS /ws
**Purpose:** Real-time price updates

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
```

**Message Format:**
```json
[
  {"symbol": "AAPL", "price": 175.50},
  {"symbol": "TSLA", "price": 245.30}
]
```

**Update Frequency:** Every 3 seconds

---

## Key Features Implementation

### 1. JWT Authentication

**Why JWT?**
- **Stateless**: Server doesn't need to store sessions
- **Scalable**: Works across multiple servers
- **Secure**: Signed with secret key (can't be tampered)

**Token Structure:**
```
Header.Payload.Signature

Header: {"alg": "HS256", "typ": "JWT"}
Payload: {"user_id": 1, "username": "admin", "exp": 1234567890}
Signature: HMACSHA256(base64(header) + "." + base64(payload), secret)
```

**Security Measures:**
- Password hashing (bcrypt)
- Token expiration (24 hours)
- HTTPS in production (not implemented here, but should be)

### 2. Real-Time Updates

**Implementation:**
1. Backend goroutine updates prices every 3 seconds
2. Broadcasts to all connected WebSocket clients
3. Frontend receives updates and re-renders

**Why WebSocket over Polling?**
- **Polling**: Client asks server every few seconds (wasteful)
- **WebSocket**: Server pushes updates when available (efficient)

### 3. User-Specific Data

**Implementation:**
- Every order includes `user_id`
- Queries filter by `user_id`
- Users can only see their own orders

**Security:**
```go
// Middleware extracts user_id from token
userID := c.Get("user_id")

// Query filters by user_id
db.Where("user_id = ?", userID).Find(&orders)
```

### 4. Price Change Visualization

**Implementation:**
```javascript
// Compare current price with previous price
const change = currentPrice - previousPrice;

// Apply color based on change
if (change > 0) return 'text-green-600';  // Up
if (change < 0) return 'text-red-600';    // Down
```

**User Experience:**
- Green = Price increased
- Red = Price decreased
- Shows change amount (↑ $2.50)

---

## Interview Preparation Q&A

### Q1: Why did you choose Golang for the backend?

**Answer:**
"I chose Golang because:
1. **Concurrency**: Goroutines make it easy to handle WebSocket connections and price updates simultaneously
2. **Performance**: Golang is compiled (not interpreted), making it faster than Node.js for CPU-intensive tasks
3. **Type Safety**: Static typing catches errors at compile-time, reducing runtime bugs
4. **Simplicity**: The language is simple and readable, making it easier to maintain

However, I'm also comfortable with Node.js/Express from my MERN experience, and I can adapt quickly to different backends."

### Q2: Explain how the WebSocket implementation works.

**Answer:**
"The WebSocket implementation has two parts:

**Backend:**
1. A goroutine runs `updatePrices()` every 3 seconds
2. It randomly changes stock prices by -2% to +2%
3. It broadcasts updated prices to all connected clients via `broadcastPrices()`
4. Clients are stored in a map: `map[*websocket.Conn]bool`

**Frontend:**
1. React component creates WebSocket connection on mount
2. `onmessage` handler receives price updates
3. Component compares new prices with previous prices
4. UI updates with color coding (green/red) based on changes

**Why WebSocket?**
- Allows server to push updates without client polling
- More efficient than HTTP polling
- Real-time updates without page refresh"

### Q3: How does JWT authentication work in your project?

**Answer:**
"JWT authentication works in three steps:

**1. Login/Signup:**
- User provides credentials
- Server validates and creates/verifies user
- Server generates JWT token with user info and expiration
- Token returned to client

**2. Token Storage:**
- Frontend stores token in `localStorage`
- Token included in `Authorization` header for protected requests

**3. Protected Routes:**
- Middleware extracts token from header
- Validates token signature and expiration
- Extracts user ID from token claims
- Stores user ID in request context
- Handler can access user ID to filter data

**Security:**
- Passwords hashed with bcrypt (one-way, can't reverse)
- Tokens signed with secret key (can't be tampered)
- Tokens expire after 24 hours"

### Q4: How do you ensure users only see their own orders?

**Answer:**
"User isolation is implemented at multiple levels:

**1. Database Level:**
- Every order has a `user_id` foreign key
- Orders are linked to users via this relationship

**2. Middleware Level:**
- JWT middleware extracts `user_id` from token
- Stores it in request context

**3. Query Level:**
```go
db.Where("user_id = ?", userID).Find(&orders)
```
- All queries filter by `user_id`
- Users can only query their own orders

**4. No Client-Side Filtering:**
- Backend enforces security (never trust client)
- Even if client tries to access other users' data, query filters prevent it"

### Q5: Explain the database schema and relationships.

**Answer:**
"The database has two main tables:

**Users Table:**
- `id` (Primary Key): Auto-incrementing integer
- `username`: Unique string (indexed for fast lookup)
- `password`: Hashed string (bcrypt)

**Orders Table:**
- `id` (Primary Key): Auto-incrementing integer
- `user_id` (Foreign Key): References `users.id`
- `symbol`: Stock symbol (e.g., "AAPL")
- `side`: "buy" or "sell"
- `quantity`: Number of shares
- `price`: Price per share
- `timestamp`: When order was placed

**Relationship:**
- One-to-Many: One user can have many orders
- Foreign key constraint ensures data integrity
- Queries use `user_id` to filter orders per user"

### Q6: How would you scale this application?

**Answer:**
"To scale this application, I would:

**1. Database:**
- Migrate from SQLite to PostgreSQL (better for production)
- Add database connection pooling
- Implement read replicas for read-heavy operations

**2. Backend:**
- Use Redis for session/token caching
- Implement rate limiting
- Add load balancing across multiple servers
- Use message queue (RabbitMQ/Kafka) for order processing

**3. WebSocket:**
- Use Redis Pub/Sub for WebSocket message broadcasting
- Allows WebSocket connections across multiple servers

**4. Frontend:**
- Implement code splitting
- Use CDN for static assets
- Add service workers for offline support

**5. Security:**
- Use HTTPS in production
- Implement refresh tokens
- Add CORS restrictions
- Rate limiting per user"

### Q7: What are the limitations of SQLite?

**Answer:**
"SQLite has some limitations:

**1. Concurrency:**
- Single writer at a time (can be slow with many writes)
- Better for read-heavy applications

**2. Scalability:**
- File-based (not suitable for distributed systems)
- No network access (must be on same machine)

**3. Features:**
- Limited compared to PostgreSQL/MySQL
- No user management
- No stored procedures

**Why I Used It:**
- Perfect for this project (single server, small scale)
- No setup required (file-based)
- Good for development and small deployments
- Easy to migrate to PostgreSQL later if needed"

### Q8: How do you handle errors in the application?

**Answer:**
"Error handling is implemented at multiple levels:

**Backend:**
```go
// Validation errors
if req.Quantity <= 0 {
    c.JSON(400, gin.H{"error": "Quantity must be positive"})
    return
}

// Database errors
if err := db.Create(&order).Error; err != nil {
    c.JSON(500, gin.H{"error": "Failed to create order"})
    return
}
```

**Frontend:**
```javascript
try {
    const response = await fetch(url);
    if (!response.ok) {
        const data = await response.json();
        setError(data.error);
    }
} catch (err) {
    setError('Error connecting to server');
}
```

**User Experience:**
- Clear error messages
- Validation before submission
- Graceful degradation (WebSocket reconnection)"

---

## Common Interview Questions

### Technical Questions

**Q: What is the difference between SQLite and PostgreSQL?**
- SQLite: File-based, embedded, single-user
- PostgreSQL: Server-based, multi-user, more features

**Q: Explain goroutines vs threads.**
- Goroutines: Lightweight, managed by Go runtime
- Threads: OS-level, heavier, more overhead

**Q: What is JWT and why use it?**
- JSON Web Token: Stateless authentication
- Contains user info, signed with secret
- No need for server-side session storage

**Q: How does WebSocket differ from HTTP?**
- HTTP: Request-response, one-way
- WebSocket: Full-duplex, persistent connection
- Server can push data anytime

**Q: What is bcrypt and why use it?**
- Password hashing algorithm
- One-way encryption (can't reverse)
- Slow by design (prevents brute force)

### Behavioral Questions

**Q: Why did you choose this tech stack?**
- Golang for performance and concurrency
- React for modern UI
- SQLite for simplicity (can migrate later)
- JWT for scalable authentication

**Q: What challenges did you face?**
- Learning Golang syntax (coming from JavaScript)
- Understanding SQLite vs MongoDB
- Implementing WebSocket correctly
- Managing state in React

**Q: How did you test the application?**
- Manual testing of all endpoints
- Tested WebSocket connections
- Verified user isolation
- Tested error scenarios

**Q: What would you improve?**
- Add unit tests
- Implement refresh tokens
- Add input validation on frontend
- Migrate to PostgreSQL for production
- Add rate limiting
- Implement order status (pending/filled)

---

## Quick Reference: MERN to Go/SQLite Translation

### Express.js → Gin

```javascript
// Express.js
app.post('/api/orders', (req, res) => {
    const order = req.body;
    res.json(order);
});
```

```go
// Gin
r.POST("/api/orders", func(c *gin.Context) {
    var order Order
    c.ShouldBindJSON(&order)
    c.JSON(200, order)
})
```

### Mongoose → GORM

```javascript
// Mongoose
const user = new User({ username: "admin" });
await user.save();
```

```go
// GORM
user := User{Username: "admin"}
db.Create(&user)
```

### MongoDB Queries → SQLite

```javascript
// MongoDB
db.orders.find({ user_id: 1 })
```

```go
// SQLite (GORM)
db.Where("user_id = ?", 1).Find(&orders)
```

### Node.js Async → Go Goroutines

```javascript
// Node.js
setInterval(() => {
    updatePrices();
}, 3000);
```

```go
// Go
go func() {
    ticker := time.NewTicker(3 * time.Second)
    for range ticker.C {
        updatePrices()
    }
}()
```

---

## Final Tips for Interview

1. **Be Honest**: Admit you're learning Golang/SQLite but show enthusiasm
2. **Show Understanding**: Explain why you made certain choices
3. **Discuss Trade-offs**: SQLite vs PostgreSQL, JWT vs sessions
4. **Talk About Improvements**: What you'd do differently next time
5. **Demonstrate Learning**: Show you can adapt to new technologies

**Key Strengths to Highlight:**
- Full-stack development experience
- Understanding of authentication/authorization
- Real-time features (WebSocket)
- Database design and relationships
- Security best practices
- Clean, modular code structure

**Areas to Acknowledge:**
- Still learning Golang (but show progress)
- SQLite limitations (but explain why it's fine for this project)
- Would add tests in production
- Would use PostgreSQL for larger scale

---

## Conclusion

This project demonstrates:
- ✅ Full-stack development skills
- ✅ Understanding of REST APIs
- ✅ Real-time features (WebSocket)
- ✅ Authentication and authorization
- ✅ Database design and relationships
- ✅ Security best practices
- ✅ Ability to learn new technologies

**Remember**: The interviewer wants to see your problem-solving ability and willingness to learn, not just perfect code. Be confident, explain your choices, and show enthusiasm for learning!

---

*Good luck with your interview! 🚀*

