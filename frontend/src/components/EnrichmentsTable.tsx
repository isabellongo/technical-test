import { useEffect, useState } from 'react'
import { fetchEnrichments } from '../hooks/useAnalytics'

type Row = {
  id_enriquecimento: string
  nome_workspace: string
  total_contatos: number
  status_processamento: string
  data_criacao: string
}

export default function EnrichmentsTable() {
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [data, setData] = useState<Row[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let mounted = true
    setLoading(true)
    fetchEnrichments(page, limit)
      .then((res) => {
        if (!mounted) return
        setData(res.data || [])
        setTotal(res.total_items || 0)
      })
      .catch(() => {})
      .finally(() => mounted && setLoading(false))
    return () => {
      mounted = false
    }
  }, [page, limit])

  const totalPages = Math.max(1, Math.ceil(total / limit))

  return (
    <div style={{ border: '1px solid #eee', padding: 12, borderRadius: 8 }}>
      <h3>Enriquecimentos</h3>
      {loading ? (
        <div>Carregando...</div>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: 6 }}>ID</th>
              <th style={{ textAlign: 'left', padding: 6 }}>Workspace</th>
              <th style={{ textAlign: 'left', padding: 6 }}>Contatos</th>
              <th style={{ textAlign: 'left', padding: 6 }}>Status</th>
              <th style={{ textAlign: 'left', padding: 6 }}>Data</th>
            </tr>
          </thead>
          <tbody>
            {data.map((r) => (
              <tr key={r.id_enriquecimento}>
                <td style={{ padding: 6, borderTop: '1px solid #f0f0f0' }}>{r.id_enriquecimento.slice(0, 8)}</td>
                <td style={{ padding: 6, borderTop: '1px solid #f0f0f0' }}>{r.nome_workspace}</td>
                <td style={{ padding: 6, borderTop: '1px solid #f0f0f0' }}>{r.total_contatos}</td>
                <td style={{ padding: 6, borderTop: '1px solid #f0f0f0' }}>{r.status_processamento}</td>
                <td style={{ padding: 6, borderTop: '1px solid #f0f0f0' }}>{new Date(r.data_criacao).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 12 }}>
        <div>
          Página {page} de {totalPages}
        </div>
        <div>
          <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page <= 1} style={{ marginRight: 8 }}>Anterior</button>
          <button onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={page >= totalPages}>Próxima</button>
        </div>
      </div>
    </div>
  )
}
