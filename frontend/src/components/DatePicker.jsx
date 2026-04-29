// components/DatePicker.jsx
import styles from './DatePicker.module.css'

export function DatePicker({ value, onChange }) {
  return (
    <div className={styles.wrapper}>
      <p className={styles.label}>Forecast Date</p>
      <input
        type="date"
        className={styles.input}
        value={value}
        onChange={e => onChange(e.target.value)}
      />
    </div>
  )
}
