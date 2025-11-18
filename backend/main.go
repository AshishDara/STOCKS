package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// JWT secret key (in production, use environment variable)
var jwtSecret = []byte("your-secret-key-change-in-production")

// Stock represents a stock with its current price
type Stock struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// User represents a user in the system
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Don't return password in JSON
}

// Order represents a trading order (database model)
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Symbol    string    `gorm:"not null" json:"symbol"`
	Side      string    `gorm:"not null" json:"side"` // "buy" or "sell"
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"not null" json:"price"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
}

// OrderRequest represents an incoming order request
type OrderRequest struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupRequest represents a signup request
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Server holds the application state
type Server struct {
	db          *gorm.DB
	stocks      map[string]*Stock
	clients     map[*websocket.Conn]bool
	clientsLock sync.RWMutex
	upgrader    websocket.Upgrader
}

// NewServer creates a new server instance
func NewServer(dbPath string) *Server {
	if dbPath == "" {
		dbPath = "trading.db"
	}

	// Initialize database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database (%s): %v", dbPath, err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&User{}, &Order{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default user if it doesn't exist
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		defaultUser := User{
			Username: "admin",
			Password: string(hashedPassword),
		}
		db.Create(&defaultUser)
		log.Println("Created default user: admin / password123")
	}

	// Initialize mock stocks with starting prices
	stocks := map[string]*Stock{
		"AAPL": {Symbol: "AAPL", Price: 175.50},
		"TSLA": {Symbol: "TSLA", Price: 245.30},
		"AMZN": {Symbol: "AMZN", Price: 138.20},
		"INFY": {Symbol: "INFY", Price: 18.75},
		"TCS":  {Symbol: "TCS", Price: 3450.00},
	}

	return &Server{
		db:      db,
		stocks:  stocks,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

func main() {
	// Get JWT secret from environment or use default
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	} else {
		jwtSecret = []byte("your-secret-key-change-in-production")
	}

	dbPath := os.Getenv("DB_PATH")

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewServer(dbPath)

	// Start the price update goroutine
	go server.updatePrices()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}

	if origins := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS")); origins != "" {
		originList := strings.Split(origins, ",")
		for i := range originList {
			originList[i] = strings.TrimSpace(originList[i])
		}
		config.AllowOrigins = originList
	} else {
		config.AllowAllOrigins = true
	}

	r.Use(cors.New(config))

	// Public routes
	r.POST("/api/login", server.login)
	r.POST("/api/signup", server.signup)
	r.GET("/api/prices", server.getPrices)
	r.GET("/ws", server.handleWebSocket)

	// Protected routes (require JWT)
	api := r.Group("/api")
	api.Use(server.authMiddleware())
	{
		api.POST("/orders", server.createOrder)
		api.GET("/orders", server.getOrders)
	}

	// Start server
	log.Printf("Server starting on :%s (DB: %s)", port, dbPath)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// authMiddleware validates JWT tokens
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID from claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(userID))
			} else {
				c.JSON(401, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// login handles user authentication
func (s *Server) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Find user
	var user User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// signup handles user registration
func (s *Server) signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Validate input
	if req.Username == "" {
		c.JSON(400, gin.H{"error": "Username is required"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(400, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	// Check if username already exists
	var existingUser User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(201, LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// getPrices returns current prices for all stocks
func (s *Server) getPrices(c *gin.Context) {
	prices := make([]Stock, 0, len(s.stocks))
	for _, stock := range s.stocks {
		prices = append(prices, *stock)
	}
	c.JSON(200, prices)
}

// createOrder handles order creation
func (s *Server) createOrder(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Validate order
	if req.Side != "buy" && req.Side != "sell" {
		c.JSON(400, gin.H{"error": "Side must be 'buy' or 'sell'"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(400, gin.H{"error": "Quantity must be positive"})
		return
	}

	if req.Price <= 0 {
		c.JSON(400, gin.H{"error": "Price must be positive"})
		return
	}

	// Create order in database
	order := Order{
		UserID:    userID.(uint),
		Symbol:    req.Symbol,
		Side:      req.Side,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Timestamp: time.Now(),
	}

	if err := s.db.Create(&order).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(201, order)
}

// getOrders returns all orders for the authenticated user
func (s *Server) getOrders(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var orders []Order
	if err := s.db.Where("user_id = ?", userID).Order("timestamp DESC").Find(&orders).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(200, orders)
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	s.clientsLock.Lock()
	s.clients[conn] = true
	s.clientsLock.Unlock()

	// Send initial prices
	s.sendPricesToClient(conn)

	// Keep connection alive and handle client messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// Unregister client
	s.clientsLock.Lock()
	delete(s.clients, conn)
	s.clientsLock.Unlock()
}

// sendPricesToClient sends current prices to a specific client
func (s *Server) sendPricesToClient(conn *websocket.Conn) {
	prices := make([]Stock, 0, len(s.stocks))
	for _, stock := range s.stocks {
		prices = append(prices, *stock)
	}
	if err := conn.WriteJSON(prices); err != nil {
		log.Printf("Error sending prices: %v", err)
	}
}

// broadcastPrices sends prices to all connected clients
func (s *Server) broadcastPrices() {
	prices := make([]Stock, 0, len(s.stocks))
	for _, stock := range s.stocks {
		prices = append(prices, *stock)
	}

	s.clientsLock.RLock()
	defer s.clientsLock.RUnlock()

	for client := range s.clients {
		if err := client.WriteJSON(prices); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			client.Close()
			delete(s.clients, client)
		}
	}
}

// updatePrices simulates live price updates
func (s *Server) updatePrices() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(3 * time.Second) // Update every 3 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Update each stock price
		for symbol, stock := range s.stocks {
			// Random price change between -2% and +2%
			changePercent := (rng.Float64()*4 - 2) / 100 // -2% to +2%
			newPrice := stock.Price * (1 + changePercent)

			// Ensure price doesn't go below a minimum
			if newPrice < 1.0 {
				newPrice = 1.0
			}

			stock.Price = newPrice
			log.Printf("Updated %s price to %.2f", symbol, newPrice)
		}

		// Broadcast updated prices to all clients
		s.broadcastPrices()
	}
}
