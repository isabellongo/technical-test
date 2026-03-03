// React import not required with automatic JSX runtime
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts'

type Props = {
  data: Record<string, number>
}

export default function EnrichmentsChart({ data }: Props) {
  const series = Object.entries(data || {}).map(([k, v]) => ({ name: k, value: v }))

  return (
    <div style={{ width: '100%', height: 300, border: '1px solid #eee', padding: 8, borderRadius: 8 }}>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={series} margin={{ top: 20, right: 20, left: 0, bottom: 5 }}>
          <XAxis dataKey="name" />
          <YAxis />
          <Tooltip />
          <Bar dataKey="value" fill="#4f46e5" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
