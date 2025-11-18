# Trading Dashboard - Full-Stack Assignment

A real-time trading dashboard application built with Golang backend and React frontend. This application simulates live stock price updates and allows users to place buy/sell orders with secure authentication and persistent storage.

## Features

### Backend (Golang)
- ✅ REST API endpoints for prices and orders
- ✅ WebSocket endpoint for real-time price streaming
- ✅ Mock price generation with random fluctuations (0.5-2%)
- ✅ **JWT Authentication** with secure login endpoint
- ✅ **Database Persistence** using SQLite (GORM)
- ✅ Concurrent price updates using goroutines and channels
- ✅ Protected order endpoints with JWT middleware

### Frontend (React)
- ✅ **Login page** with authentication
- ✅ Live prices table with WebSocket connection
- ✅ Order form for placing buy/sell orders
- ✅ Orders history table (user-specific)
- ✅ Real-time price updates without page refresh
- ✅ Visual price change indicators (green for up, red for down)
- ✅ Modern UI with Tailwind CSS
- ✅ Token-based authentication with localStorage

## Tech Stack

- **Backend:** Golang with Gin framework
- **Frontend:** React with Vite
- **Styling:** Tailwind CSS
- **WebSocket:** Gorilla WebSocket
- **Authentication:** JWT (golang-jwt/jwt)
- **Database:** SQLite with GORM
- **Password Hashing:** bcrypt

## Project Structure

```
STOCKS/
├── backend/
│   ├── main.go          # Main server file
│   ├── go.mod           # Go dependencies
│   └── trading.db       # SQLite database (created automatically)
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Login.jsx
│   │   │   ├── LivePricesTable.jsx
│   │   │   ├── OrderForm.jsx
│   │   │   └── OrdersTable.jsx
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── index.css
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
└── README.md
```

## Setup Instructions

### Prerequisites
- Go 1.21 or higher
- Node.js 16+ and npm/yarn
- A modern web browser

### Backend Setup

#### Configure Environment Variables
1. Copy the example file and update values as needed:
   ```bash
   cd backend
   cp env.example .env
   ```
2. Update `.env` with:
   ```env
   JWT_SECRET=your-super-secret-key
   PORT=8080
   DB_PATH=trading.db
   ALLOWED_ORIGINS=http://localhost:3000
   ```
3. These values control the API port, database location, and allowed CORS origins for deployments. Leave `ALLOWED_ORIGINS` empty to allow all origins during local development.

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. (Optional) Set JWT secret as environment variable:
```bash
export JWT_SECRET="your-secret-key-here"
```
If not set, a default secret will be used (not recommended for production).

4. Run the backend server:
```bash
go run main.go
```

The backend server will start on `http://localhost:8080`

**Note:** On first run, the database will be automatically created with a default user:
- Username: `admin`
- Password: `password123`

### Frontend Setup

#### Configure Environment Variables
1. Copy the example file:
   ```bash
   cd frontend
   cp env.example .env
   ```
2. Update `.env` with your backend URLs:
   ```env
   VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
   ```
3. When deploying, replace these values with your hosted backend and WebSocket URLs (e.g., `https://api.yourdomain.com` and `wss://api.yourdomain.com/ws`).

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install npm dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## Authentication

### Default Credentials
- **Username:** `admin`
- **Password:** `password123`

### User Registration

You can create new user accounts directly from the login page:
1. Click the "Sign Up" tab on the login page
2. Enter a unique username
3. Enter a password (minimum 6 characters)
4. Click "Sign Up"
5. You'll be automatically logged in after successful registration

### How Authentication Works

1. User logs in via `/api/login` or signs up via `/api/signup` endpoint
2. Server validates credentials (or creates new user) and returns a JWT token
3. Frontend stores the token in localStorage
4. All authenticated requests include the token in the `Authorization` header: `Bearer <token>`
5. Token expires after 24 hours
6. All user data (including orders) is stored in the SQLite database

## API Endpoints

### Public Endpoints

- **POST /api/login** - User authentication
  - Request Body:
    ```json
    {
      "username": "admin",
      "password": "password123"
    }
    ```
  - Response:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": {
        "id": 1,
        "username": "admin"
      }
    }
    ```

- **POST /api/signup** - User registration (create new account)
  - Request Body:
    ```json
    {
      "username": "newuser",
      "password": "password123"
    }
    ```
  - Requirements:
    - Username must be unique
    - Password must be at least 6 characters
  - Response:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": {
        "id": 2,
        "username": "newuser"
      }
    }
    ```
  - Note: Automatically logs in the user after successful signup

- **GET /api/prices** - Get current prices for all stocks (public)
  - Response: Array of stock objects with symbol and price

- **WS /ws** - WebSocket endpoint for real-time price updates (public)
  - Connects to receive live price updates
  - Prices update every 3 seconds
  - Sends array of stock objects with updated prices

### Protected Endpoints (Require JWT Token)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

- **POST /api/orders** - Place a new order
  - Headers: `Authorization: Bearer <token>`
  - Request Body:
    ```json
    {
      "symbol": "AAPL",
      "side": "buy",
      "quantity": 10,
      "price": 175.50
    }
    ```
  - Response: Created order object with user_id

- **GET /api/orders** - Get all orders for the authenticated user
  - Headers: `Authorization: Bearer <token>`
  - Response: Array of orders (only for the logged-in user)

## Database Schema

### Users Table
- `id` (Primary Key)
- `username` (Unique, Not Null)
- `password` (Hashed with bcrypt, Not Null)

### Orders Table
- `id` (Primary Key)
- `user_id` (Foreign Key to Users, Not Null)
- `symbol` (Not Null)
- `side` (Not Null) - "buy" or "sell"
- `quantity` (Not Null)
- `price` (Not Null)
- `timestamp` (Not Null)

## Mock Stocks

The application tracks the following mock stocks:
- **AAPL** - Apple Inc.
- **TSLA** - Tesla Inc.
- **AMZN** - Amazon.com Inc.
- **INFY** - Infosys Limited
- **TCS** - Tata Consultancy Services

## How It Works

1. **Authentication:** Users must login to access order functionality. Prices and WebSocket are public.

2. **Price Updates:** The backend uses a goroutine that runs every 3 seconds to randomly update stock prices by -2% to +2%.

3. **WebSocket Streaming:** All connected clients receive real-time price updates via WebSocket connections.

4. **Order Management:** Orders are stored in SQLite database with user association. Each user can only see their own orders.

5. **Frontend Updates:** The React frontend subscribes to WebSocket updates and automatically refreshes the prices table when new data arrives.

6. **Token Management:** JWT tokens are stored in localStorage and automatically included in authenticated requests. Tokens expire after 24 hours.

## Security Features

- ✅ Password hashing with bcrypt
- ✅ JWT token-based authentication
- ✅ Protected API endpoints
- ✅ User-specific order access
- ✅ Token expiration (24 hours)
- ✅ Secure password storage (never returned in API responses)

## Development Notes

- The backend uses CORS middleware to allow cross-origin requests from the frontend
- WebSocket connections are managed with proper cleanup on disconnect
- Price changes are visually indicated with green (up) and red (down) colors
- The application uses concurrent programming patterns (goroutines, channels, mutexes) for safe concurrent access
- Database is automatically created and migrated on first run
- Default user is created if no users exist in the database

## Database File

The SQLite database file (`trading.db`) is created automatically in the `backend` directory. You can:
- Delete it to reset the database
- Use SQLite tools to inspect the data
- Backup it for data persistence

## Future Enhancements (Optional)

- [x] JWT authentication for secure access ✅
- [x] Persistent database storage ✅
- [ ] Order status tracking (pending, filled, cancelled)
- [ ] Price history charts
- [ ] Portfolio tracking
- [ ] User registration endpoint
- [ ] Password reset functionality
- [ ] Refresh token mechanism
- [ ] AWS deployment configuration

## License

This project is created for assignment purposes.
