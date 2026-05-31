import React from "react";
import ReactDOM from "react-dom/client";
import { getBooks, getCategories } from "./api/client";
import { ApiPreview } from "./components/ApiPreview";
import { StatCard } from "./components/StatCard";
import { Topbar } from "./components/Topbar";
import { useApiResource } from "./hooks/useApiResource";
import "./styles.css";

function App() {
  const books = useApiResource(getBooks);
  const categories = useApiResource(getCategories);

  const bookCount = Array.isArray(books.data?.data) ? books.data.data.length : "API";
  const categoryCount = Array.isArray(categories.data?.data) ? categories.data.data.length : "API";

  return (
    <main className="app-shell">
      <Topbar />

      <section className="hero">
        <div>
          <p className="eyebrow">React frontend starter</p>
          <h1>Storefront workbench for your ReadOn API.</h1>
          <p className="lede">
            This app runs separately from Gin in development and talks to your existing backend through a Vite proxy.
          </p>
        </div>
      </section>

      <section className="stats-grid" aria-label="API status">
        <StatCard
          mark="B"
          label="Books endpoint"
          value={bookCount}
          status={books.loading ? "Checking..." : books.error ? "Needs backend" : "Connected"}
        />
        <StatCard
          mark="C"
          label="Categories endpoint"
          value={categoryCount}
          status={categories.loading ? "Checking..." : categories.error ? "Needs backend" : "Connected"}
        />
      </section>

      <section className="preview-grid">
        <ApiPreview title="Books response" resource={books} />
        <ApiPreview title="Categories response" resource={categories} />
      </section>
    </main>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
