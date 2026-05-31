export function Topbar() {
  return (
    <nav className="topbar">
      <a className="brand" href="/">
        <span className="brand-mark" aria-hidden="true">R</span>
        ReadOn
      </a>
      <div className="nav-actions">
        <a href="/api/login">Log in</a>
        <a className="button" href="/api/signup">Sign up</a>
      </div>
    </nav>
  );
}
