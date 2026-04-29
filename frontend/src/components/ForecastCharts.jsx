// components/ForecastCharts.jsx
// Two complementary views of the same forecast data:
//   1. Stacked bar chart — units per hour, stacked by product (operational view)
//   2. Pie chart — total units by product (prep planning view)
//
// Both use Recharts. We derive both datasets from the same entries array
// so there is one source of truth.

import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
  PieChart, Pie, Cell, Sector
} from 'recharts'
import { useState, useMemo } from 'react'
import styles from './ForecastCharts.module.css'

// A carefully chosen palette that works on dark backgrounds.
// Red-first to match the KFC brand, then warm complementary tones.
const COLORS = ['#C8102E', '#F5A623', '#E8D5B7', '#8B6355', '#D4756A', '#F0C040']

// Build the stacked bar dataset: one object per hour, with a key per product.
// e.g. { hour: '12:00', 'Zinger Burger': 24, 'Original Recipe Bucket': 18, ... }
function buildHourlyData(entries) {
  const map = {}
  entries.forEach(e => {
    const key = `${String(e.hour).padStart(2, '0')}:00`
    if (!map[key]) map[key] = { hour: key }
    map[key][e.product_name] = (map[key][e.product_name] ?? 0) + e.quantity
  })
  return Object.values(map).sort((a, b) => a.hour.localeCompare(b.hour))
}

// Build the pie dataset: total predicted units per product across the day.
function buildProductTotals(entries) {
  const map = {}
  entries.forEach(e => {
    map[e.product_name] = (map[e.product_name] ?? 0) + e.quantity
  })
  return Object.entries(map)
    .map(([name, value]) => ({ name, value: Math.round(value) }))
    .sort((a, b) => b.value - a.value)
}

// Custom tooltip for the bar chart — shows all products at that hour.
function HourTooltip({ active, payload, label }) {
  if (!active || !payload?.length) return null
  const total = payload.reduce((s, p) => s + (p.value ?? 0), 0)
  return (
    <div className={styles.tooltip}>
      <p className={styles.tooltipHour}>{label}</p>
      {payload.map(p => (
        <p key={p.dataKey} style={{ color: p.fill }}>
          {p.dataKey}: <strong>{Math.round(p.value)}</strong>
        </p>
      ))}
      <p className={styles.tooltipTotal}>Total: <strong>{Math.round(total)}</strong></p>
    </div>
  )
}

// Active pie slice — expands on hover.
function ActivePieSlice(props) {
  const { cx, cy, innerRadius, outerRadius, startAngle, endAngle, fill } = props
  return (
    <Sector
      cx={cx} cy={cy}
      innerRadius={innerRadius}
      outerRadius={outerRadius + 10}
      startAngle={startAngle}
      endAngle={endAngle}
      fill={fill}
    />
  )
}

export function ForecastCharts({ entries }) {
  const [activePieIndex, setActivePieIndex] = useState(0)

  const products = useMemo(
    () => [...new Set(entries.map(e => e.product_name))],
    [entries]
  )
  const hourlyData   = useMemo(() => buildHourlyData(entries), [entries])
  const productTotals = useMemo(() => buildProductTotals(entries), [entries])

  const totalUnits = productTotals.reduce((s, p) => s + p.value, 0)

  return (
    <div className={styles.wrapper}>
      {/* ── Bar chart: hourly breakdown ── */}
      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <h2 className={styles.chartTitle}>Hourly Forecast</h2>
          <p className={styles.chartSub}>Predicted units per hour · stacked by product</p>
        </div>
        <ResponsiveContainer width="100%" height={320}>
          <BarChart data={hourlyData} margin={{ top: 8, right: 16, left: 0, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" vertical={false} />
            <XAxis
              dataKey="hour"
              tick={{ fill: 'var(--muted)', fontSize: 11 }}
              axisLine={{ stroke: 'var(--surface-3)' }}
              tickLine={false}
            />
            <YAxis
              tick={{ fill: 'var(--muted)', fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              width={32}
            />
            <Tooltip content={<HourTooltip />} cursor={{ fill: 'rgba(255,255,255,0.04)' }} />
            <Legend
              wrapperStyle={{ fontSize: 12, color: 'var(--muted)', paddingTop: 16 }}
            />
            {products.map((p, i) => (
              <Bar key={p} dataKey={p} stackId="a" fill={COLORS[i % COLORS.length]} radius={i === products.length - 1 ? [3, 3, 0, 0] : [0,0,0,0]} />
            ))}
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* ── Pie chart: product share ── */}
      <div className={`${styles.card} ${styles.pieCard}`}>
        <div className={styles.cardHeader}>
          <h2 className={styles.chartTitle}>Product Mix</h2>
          <p className={styles.chartSub}>Total predicted units · all hours</p>
        </div>

        <div className={styles.pieLayout}>
          <ResponsiveContainer width="100%" height={260}>
            <PieChart>
              <Pie
                data={productTotals}
                cx="50%"
                cy="50%"
                innerRadius={70}
                outerRadius={100}
                dataKey="value"
                activeIndex={activePieIndex}
                activeShape={ActivePieSlice}
                onMouseEnter={(_, index) => setActivePieIndex(index)}
              >
                {productTotals.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} stroke="var(--black)" strokeWidth={2} />
                ))}
              </Pie>
              <Tooltip
                formatter={(value, name) => [`${value} units`, name]}
                contentStyle={{
                  background: 'var(--surface-2)',
                  border: '1px solid var(--surface-3)',
                  borderRadius: 8,
                  color: 'var(--bone)',
                  fontSize: 13,
                }}
              />
            </PieChart>
          </ResponsiveContainer>

          {/* Legend with percentage */}
          <div className={styles.pieLegend}>
            {productTotals.map((p, i) => (
              <div
                key={p.name}
                className={`${styles.legendRow} ${i === activePieIndex ? styles.legendActive : ''}`}
                onMouseEnter={() => setActivePieIndex(i)}
              >
                <span className={styles.legendDot} style={{ background: COLORS[i % COLORS.length] }} />
                <span className={styles.legendName}>{p.name}</span>
                <span className={styles.legendValue}>{p.value}</span>
                <span className={styles.legendPct}>
                  {totalUnits > 0 ? `${Math.round(p.value / totalUnits * 100)}%` : '—'}
                </span>
              </div>
            ))}
            <div className={styles.legendTotal}>
              <span>Total</span>
              <span>{totalUnits} units</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
