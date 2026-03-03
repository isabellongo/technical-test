// React import not required with automatic JSX runtime

type Props = {
  title: string
  value: string | number
  subtitle?: string
}

export default function KPICard({ title, value, subtitle }: Props) {
  return (
    <div style={{
      border: '1px solid #e6e6e6',
      padding: 12,
      borderRadius: 8,
      minWidth: 160,
      boxShadow: '0 1px 3px rgba(0,0,0,0.04)'
    }}>
      <div style={{ fontSize: 12, color: '#666' }}>{title}</div>
      <div style={{ fontSize: 24, fontWeight: 700, marginTop: 6 }}>{value}</div>
      {subtitle && <div style={{ fontSize: 12, color: '#999', marginTop: 6 }}>{subtitle}</div>}
    </div>
  )
}
