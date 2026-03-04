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
    <div className="table-card">
      {loading ? (
        <div>Carregando...</div>
      ) : (
        <div className="table-scroll">
          <table className="enrichments-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Workspace</th>
                <th>Contatos</th>
                <th>Status</th>
                <th>Data</th>
              </tr>
            </thead>
            <tbody>
              {data.map((r) => (
                <tr key={r.id_enriquecimento}>
                  <td className="mono">{r.id_enriquecimento.slice(0, 8)}</td>
                  <td>{r.nome_workspace}</td>
                  <td>{r.total_contatos}</td>
                  <td>{r.status_processamento}</td>
                  <td>{new Date(r.data_criacao).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="table-footer">
        <div>
          Página {page} de {totalPages}
        </div>
        <div>
          <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page <= 1} className="btn" style={{ marginRight: 8 }}>Anterior</button>
          <button onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={page >= totalPages} className="btn">Próxima</button>
        </div>
      </div>
    </div>
  )
}
