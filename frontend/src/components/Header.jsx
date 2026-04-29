// components/Header.jsx
import styles from './Header.module.css'

export function Header() {
  return (
    <header className={styles.header}>
      <div className={styles.logo}>
        {/* KFC-red stripe accent */}
        <span className={styles.stripe} />
        <div className={styles.logoText}>
          <span className={styles.kfc}>KFC</span>
          <span className={styles.tagline}>Sales Forecast</span>
        </div>
      </div>
      <p className={styles.sub}>Daily prediction engine · Next-day operations planning</p>
    </header>
  )
}
