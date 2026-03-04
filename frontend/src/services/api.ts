import axios from 'axios'

// Prefer a relative path (e.g. '/api') so the frontend can be served from the
// same origin and nginx will proxy requests to the internal API service.
const API_URL = (import.meta.env.VITE_API_URL as string) || ''
const API_KEY = (import.meta.env.VITE_API_KEY as string) || 'driva_test_key_abc123xyz789'

const client = axios.create({
  baseURL: API_URL || undefined,
  headers: {
    Authorization: `Bearer ${API_KEY}`,
    'Content-Type': 'application/json',
  },
})

export default client
