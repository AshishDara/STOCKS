import { useState, useEffect } from 'react'
import { API_URL } from './config'
import Login from './components/Login'
import LivePricesTable from './components/LivePricesTable'
import OrderForm from './components/OrderForm'
import OrdersTable from './components/OrdersTable'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [token, setToken] = useState(null)
  const [user, setUser] = useState(null)
  const [orders, setOrders] = useState([])

  // Check for existing token on mount
  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    if (storedToken && storedUser) {
      setToken(storedToken)
      setUser(JSON.parse(storedUser))
      setIsAuthenticated(true)
    }
  }, [])

  // Fetch orders when authenticated
  const fetchOrders = async () => {
    if (!token) return

    try {
      const response = await fetch(`${API_URL}/api/orders`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (response.status === 401) {
        // Token expired or invalid
        handleLogout()
        return
      }

      const data = await response.json()
      setOrders(data)
    } catch (error) {
      console.error('Error fetching orders:', error)
    }
  }

  useEffect(() => {
    if (isAuthenticated && token) {
      fetchOrders()
      // Refresh orders every 5 seconds
      const interval = setInterval(fetchOrders, 5000)
      return () => clearInterval(interval)
    }
  }, [isAuthenticated, token])

  const handleLogin = (newToken, newUser) => {
    setToken(newToken)
    setUser(newUser)
    setIsAuthenticated(true)
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
    setIsAuthenticated(false)
    setOrders([])
  }

  const handleOrderPlaced = () => {
    fetchOrders()
  }

  if (!isAuthenticated) {
    return <Login onLogin={handleLogin} />
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800">
            Trading Dashboard
          </h1>
          <div className="flex items-center gap-4">
            <span className="text-gray-600">Welcome, {user?.username}</span>
            <button
              onClick={handleLogout}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-md transition-colors"
            >
              Logout
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-semibold mb-4 text-gray-700">
              Live Stock Prices
            </h2>
            <LivePricesTable />
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-semibold mb-4 text-gray-700">
              Place Order
            </h2>
            <OrderForm onOrderPlaced={handleOrderPlaced} token={token} />
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-4 text-gray-700">
            Order History
          </h2>
          <OrdersTable orders={orders} />
        </div>
      </div>
    </div>
  )
}

export default App
