// components/StoreSelector.jsx
// Renders the list of stores as clickable cards, not a plain <select>.
// This makes the store choice feel like a real operational decision.

import styles from './StoreSelector.module.css'

export function StoreSelector({ stores, selectedId, onSelect }) {
  return (
    <div className={styles.wrapper}>
      <p className={styles.label}>Select Store</p>
      <div className={styles.grid}>
        {stores.map(store => (
          <button
            key={store.id}
            className={`${styles.card} ${store.id === selectedId ? styles.active : ''}`}
            onClick={() => onSelect(store.id)}
          >
            <span className={styles.city}>{store.city}</span>
            <span className={styles.name}>{store.name}</span>
          </button>
        ))}
      </div>
    </div>
  )
}
