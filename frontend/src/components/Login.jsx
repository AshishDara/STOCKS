import { useState } from 'react'
import { API_URL } from '../config'

function Login({ onLogin }) {
  const [isSignup, setIsSignup] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess('')

    const endpoint = isSignup ? '/api/signup' : '/api/login'

    try {
      const response = await fetch(`${API_URL}${endpoint}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      const data = await response.json()

      if (response.ok) {
        // Store token in localStorage
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify(data.user))
        onLogin(data.token, data.user)
      } else {
        setError(data.error || (isSignup ? 'Signup failed' : 'Login failed'))
      }
    } catch (err) {
      setError('Error connecting to server')
      console.error('Auth error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white rounded-lg shadow-md p-8 w-full max-w-md">
        <h2 className="text-3xl font-bold text-center mb-2 text-gray-800">
          Trading Dashboard
        </h2>
        <p className="text-center text-gray-600 mb-6">
          {isSignup ? 'Create a new account' : 'Login to your account'}
        </p>

        <div className="flex justify-center mb-6">
          <div className="bg-gray-100 rounded-lg p-1 inline-flex">
            <button
              type="button"
              onClick={() => {
                setIsSignup(false)
                setError('')
                setSuccess('')
              }}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                !isSignup
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              Login
            </button>
            <button
              type="button"
              onClick={() => {
                setIsSignup(true)
                setError('')
                setSuccess('')
              }}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                isSignup
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              Sign Up
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
              Username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              placeholder="Enter username"
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder={isSignup ? "At least 6 characters" : "Enter password"}
              required
              minLength={isSignup ? 6 : undefined}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            {isSignup && (
              <p className="text-xs text-gray-500 mt-1">
                Password must be at least 6 characters
              </p>
            )}
          </div>

          {error && (
            <div className="p-3 rounded-md bg-red-50 text-red-800">
              {error}
            </div>
          )}

          {success && (
            <div className="p-3 rounded-md bg-green-50 text-green-800">
              {success}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className={`w-full py-2 px-4 rounded-md font-medium ${
              loading
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700'
            } text-white transition-colors`}
          >
            {loading
              ? isSignup
                ? 'Creating account...'
                : 'Logging in...'
              : isSignup
              ? 'Sign Up'
              : 'Login'}
          </button>

          {!isSignup && (
            <div className="text-sm text-gray-500 text-center mt-4">
              <p>Default credentials:</p>
              <p className="font-mono">Username: admin</p>
              <p className="font-mono">Password: password123</p>
            </div>
          )}
        </form>
      </div>
    </div>
  )
}

export default Login
