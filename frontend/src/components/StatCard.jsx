export function StatCard({ mark, label, value, status }) {
  return (
    <section className="stat-card">
      <div className="stat-icon" aria-hidden="true">{mark}</div>
      <div>
        <p>{label}</p>
        <strong>{value}</strong>
        <span>{status}</span>
      </div>
    </section>
  );
}
