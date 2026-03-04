// React import not required with automatic JSX runtime

type Props = {
  title: string
  value: string | number
  subtitle?: string
}

export default function KPICard({ title, value, subtitle }: Props) {
  return (
    <div className="kpi-card">
      <div className="kpi-title">{title}</div>
      <div className="kpi-value">{value}</div>
      {subtitle && <div className="kpi-sub">{subtitle}</div>}
    </div>
  )
}
