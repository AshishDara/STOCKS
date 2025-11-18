// API Configuration
// Uses environment variables in production, localhost in development

export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
export const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080';

// Helper function to get WebSocket URL (convert http to ws)
export const getWebSocketURL = () => {
  if (import.meta.env.VITE_WS_URL) {
    return import.meta.env.VITE_WS_URL;
  }
  
  // Auto-convert HTTP to WebSocket
  const apiUrl = API_URL;
  if (apiUrl.startsWith('https://')) {
    return apiUrl.replace('https://', 'wss://') + '/ws';
  } else if (apiUrl.startsWith('http://')) {
    return apiUrl.replace('http://', 'ws://') + '/ws';
  }
  
  return 'ws://localhost:8080/ws';
};

