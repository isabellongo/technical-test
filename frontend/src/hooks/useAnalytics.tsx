import { useEffect, useState } from 'react'
import api from '../services/api'

export function useOverview() {
  const [data, setData] = useState<any | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true
    setLoading(true)
    api
      .get('/analytics/overview')
      .then((res) => {
        if (mounted) setData(res.data)
      })
      .catch((err) => {
        if (mounted) setError(err?.message || 'error')
      })
      .finally(() => mounted && setLoading(false))

    return () => {
      mounted = false
    }
  }, [])

  return { data, loading, error }
}

export async function fetchEnrichments(page = 1, limit = 10, status = '', categoria = '') {
  const params: Record<string, any> = { page, limit }
  if (status) params.status_processamento = status
  if (categoria) params.categoria_tamanho_job = categoria
  const res = await api.get('/analytics/enrichments', { params })
  return res.data
}
