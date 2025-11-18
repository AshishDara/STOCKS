import { useState, useEffect, useRef } from 'react'
import { getWebSocketURL } from '../config'

function LivePricesTable() {
  const [prices, setPrices] = useState([])
  const [previousPrices, setPreviousPrices] = useState({})
  const wsRef = useRef(null)
  const reconnectTimeoutRef = useRef(null)

  const setupWebSocket = () => {
    const ws = new WebSocket(getWebSocketURL())
    wsRef.current = ws

    ws.onopen = () => {
      console.log('WebSocket connected')
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
        reconnectTimeoutRef.current = null
      }
    }

    ws.onmessage = (event) => {
      const newPrices = JSON.parse(event.data)
      
      // Store previous prices for comparison before updating
      setPrices((currentPrices) => {
        const newPrev = {}
        currentPrices.forEach((stock) => {
          newPrev[stock.symbol] = stock.price
        })
        setPreviousPrices(newPrev)
        return newPrices
      })
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected')
      // Attempt to reconnect after 3 seconds
      reconnectTimeoutRef.current = setTimeout(() => {
        if (wsRef.current?.readyState === WebSocket.CLOSED || !wsRef.current) {
          setupWebSocket()
        }
      }, 3000)
    }
  }

  useEffect(() => {
    setupWebSocket()

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  const getPriceChange = (symbol, currentPrice) => {
    const prevPrice = previousPrices[symbol]
    if (!prevPrice) return null
    return currentPrice - prevPrice
  }

  const getPriceColor = (symbol, currentPrice) => {
    const change = getPriceChange(symbol, currentPrice)
    if (change === null) return 'text-gray-700'
    if (change > 0) return 'text-green-600'
    if (change < 0) return 'text-red-600'
    return 'text-gray-700'
  }

  if (prices.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        Connecting to live feed...
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Symbol
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Price
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {prices.map((stock) => {
            const change = getPriceChange(stock.symbol, stock.price)
            const priceColor = getPriceColor(stock.symbol, stock.price)
            
            return (
              <tr key={stock.symbol} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">
                    {stock.symbol}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right">
                  <div className={`text-sm font-semibold ${priceColor}`}>
                    ${stock.price.toFixed(2)}
                    {change !== null && change !== 0 && (
                      <span className="ml-2 text-xs">
                        {change > 0 ? '↑' : '↓'} {Math.abs(change).toFixed(2)}
                      </span>
                    )}
                  </div>
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

export default LivePricesTable

