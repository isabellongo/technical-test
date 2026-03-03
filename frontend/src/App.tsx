// JSX runtime automatic; não é necessário importar React explicitamente
import './App.css'
import KPICard from './components/KPICard'
import EnrichmentsChart from './components/EnrichmentsChart'
import EnrichmentsTable from './components/EnrichmentsTable'
import { useOverview } from './hooks/useAnalytics'

function App() {
  const { data, loading } = useOverview()

  const overview = data || {}
  const total = overview.total_enriquecimentos ?? 0
  const success = overview.total_sucesso ?? 0
  const percent = overview.percentual_sucesso ?? 0
  const tempoMedio = overview.tempo_medio_minutos ?? 0
  const totalPorCategoria = overview.total_por_categoria ?? {}

  return (
    <div style={{ maxWidth: 1100, margin: '24px auto', padding: 12 }}>
      <h1>Driva - Dashboard</h1>

      <div style={{ display: 'flex', gap: 12, marginBottom: 16 }}>
        <KPICard title="Total" value={total} subtitle="Enriquecimentos" />
        <KPICard title="Sucesso" value={`${success} (${percent.toFixed(1)}%)`} subtitle="Total com sucesso" />
        <KPICard title="Tempo médio (min)" value={tempoMedio.toFixed(1)} subtitle="dur. proc." />
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 420px', gap: 12, marginBottom: 16 }}>
        <div>
          <h3>Distribuição por categoria</h3>
          <EnrichmentsChart data={totalPorCategoria} />
        </div>

        <div>
          <h3>Últimos enriquecimentos</h3>
          <EnrichmentsTable />
        </div>
      </div>

      {loading && <div>Carregando dados...</div>}
    </div>
  )
}

export default App
