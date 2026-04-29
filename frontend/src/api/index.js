// api/index.js
// All communication with the Go backend lives here.
// Components never call fetch directly — they use these functions.

const BASE = import.meta.env.VITE_API_BASE ?? ''

/**
 * Fetches all KFC stores.
 * @returns {Promise<Array<{id: number, name: string, city: string}>>}
 */
export async function fetchStores() {
  const res = await fetch(`${BASE}/api/stores`)
  if (!res.ok) throw new Error(`stores: ${res.status}`)
  return res.json()
}

/**
 * Fetches forecast entries for a store on a date.
 * @param {number} storeId
 * @param {string} date  — "YYYY-MM-DD"
 * @returns {Promise<Array<ForecastEntry>>}
 */
export async function fetchForecast(storeId, date) {
  const res = await fetch(`${BASE}/api/forecast?store_id=${storeId}&date=${date}`)
  if (!res.ok) throw new Error(`forecast: ${res.status}`)
  return res.json()
}
