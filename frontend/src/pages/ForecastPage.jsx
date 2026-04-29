// pages/ForecastPage.jsx
// The single page of the app. Orchestrates controls + charts.
// All data logic is in useForecast — this file is pure layout and UI state.

import { useForecast }     from '../hooks/useForecast'
import { StoreSelector }   from '../components/StoreSelector'
import { DatePicker }      from '../components/DatePicker'
import { ForecastCharts }  from '../components/ForecastCharts'
import { SummaryBar }      from '../components/SummaryBar'
import styles              from './ForecastPage.module.css'

export function ForecastPage() {
  const { stores, storeId, setStoreId, date, setDate, entries, loading, error } = useForecast()

  const selectedStore = stores.find(s => s.id === storeId)

  return (
    <main className={styles.page}>

      {/* ── Controls row ── */}
      <section className={styles.controls}>
        <StoreSelector stores={stores} selectedId={storeId} onSelect={setStoreId} />
        <DatePicker value={date} onChange={setDate} />
      </section>

      {/* ── Context heading ── */}
      {selectedStore && (
        <div className={styles.context}>
          <h1 className={styles.heading}>
            <span className={styles.headingStore}>{selectedStore.name}</span>
            <span className={styles.headingSep}> · </span>
            <span className={styles.headingDate}>{date}</span>
          </h1>
          <p className={styles.headingSub}>
            Predicted sales · generated from 7-day rolling average
          </p>
        </div>
      )}

      {/* ── States: loading / error / empty / data ── */}
      {loading && (
        <div className={styles.state}>
          <div className={styles.spinner} />
          <p>Loading forecast…</p>
        </div>
      )}

      {!loading && error && (
        <div className={`${styles.state} ${styles.stateError}`}>
          <p>⚠ {error}</p>
        </div>
      )}

      {!loading && !error && entries.length === 0 && (
        <div className={styles.state}>
          <p className={styles.emptyTitle}>No forecast for this date</p>
          <p className={styles.emptySub}>
            The daily job generates predictions for the next day.<br />
            Check back after the scheduler runs, or select a different date.
          </p>
        </div>
      )}

      {!loading && !error && entries.length > 0 && (
        <>
          <SummaryBar entries={entries} />
          <ForecastCharts entries={entries} />
        </>
      )}
    </main>
  )
}
