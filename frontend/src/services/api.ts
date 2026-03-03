import axios from 'axios'

const API_URL = (import.meta.env.VITE_API_URL as string) || 'http://localhost:3000'
const API_KEY = (import.meta.env.VITE_API_KEY as string) || 'driva_test_key_abc123xyz789'

const client = axios.create({
  baseURL: API_URL,
  headers: {
    Authorization: `Bearer ${API_KEY}`,
    'Content-Type': 'application/json',
  },
})

export default client
