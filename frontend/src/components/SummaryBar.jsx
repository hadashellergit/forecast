// components/SummaryBar.jsx
// Three quick-read KPI tiles above the charts.
// Derived from entries — no extra API call needed.

import styles from './SummaryBar.module.css'

function peakHour(entries) {
  const byHour = {}
  entries.forEach(e => {
    byHour[e.hour] = (byHour[e.hour] ?? 0) + e.quantity
  })
  const [hour] = Object.entries(byHour).sort((a, b) => b[1] - a[1])[0] ?? []
  return hour != null ? `${String(hour).padStart(2, '0')}:00` : '—'
}

function topProduct(entries) {
  const byProduct = {}
  entries.forEach(e => {
    byProduct[e.product_name] = (byProduct[e.product_name] ?? 0) + e.quantity
  })
  const [name] = Object.entries(byProduct).sort((a, b) => b[1] - a[1])[0] ?? []
  return name ?? '—'
}

export function SummaryBar({ entries }) {
  const total   = Math.round(entries.reduce((s, e) => s + e.quantity, 0))
  const peak    = peakHour(entries)
  const top     = topProduct(entries)

  const kpis = [
    { label: 'Total Units', value: total.toLocaleString(), accent: false },
    { label: 'Peak Hour',   value: peak,                   accent: false },
    { label: 'Top Product', value: top,                    accent: true  },
  ]

  return (
    <div className={styles.bar}>
      {kpis.map(k => (
        <div key={k.label} className={styles.tile}>
          <p className={styles.tileLabel}>{k.label}</p>
          <p className={`${styles.tileValue} ${k.accent ? styles.accent : ''}`}>{k.value}</p>
        </div>
      ))}
    </div>
  )
}
