export function ApiPreview({ title, resource }) {
  return (
    <section className="panel">
      <div className="panel-heading">
        <h2>{title}</h2>
        {resource.loading ? <span className="loading-dot" aria-hidden="true" /> : null}
      </div>

      {resource.error ? <p className="error-text">{resource.error}</p> : null}
      {!resource.error && !resource.loading ? (
        <pre>{JSON.stringify(resource.data, null, 2)}</pre>
      ) : null}
    </section>
  );
}
