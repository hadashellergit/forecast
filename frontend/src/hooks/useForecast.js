// hooks/useForecast.js
// Encapsulates all async state for the forecast page.
// The page component stays pure UI — it just reads from this hook.

import { useState, useEffect } from 'react'
import { fetchStores, fetchForecast } from '../api'
import { format } from 'date-fns'

export function useForecast() {
  const [stores, setStores]       = useState([])
  const [storeId, setStoreId]     = useState(null)
  const [date, setDate]           = useState(format(new Date(), 'yyyy-MM-dd'))
  const [entries, setEntries]     = useState([])
  const [loading, setLoading]     = useState(false)
  const [error, setError]         = useState(null)

  // Load stores once on mount.
  useEffect(() => {
    fetchStores()
      .then(data => {
        setStores(data)
        if (data.length > 0) setStoreId(data[0].id)
      })
      .catch(err => setError(err.message))
  }, [])

  // Load forecast whenever storeId or date changes.
  useEffect(() => {
    if (!storeId || !date) return
    setLoading(true)
    setError(null)
    fetchForecast(storeId, date)
      .then(data => setEntries(data ?? []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [storeId, date])

  return { stores, storeId, setStoreId, date, setDate, entries, loading, error }
}
