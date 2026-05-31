# ReadOn Frontend

React + Vite frontend for the ReadOn Gin API.

## Local Development

Use Node.js 18 or newer.

```bash
cd frontend
npm install
npm run dev
```

For a more cautious first install:

```bash
npm install --ignore-scripts
npm audit
```

After `package-lock.json` exists, prefer:

```bash
npm ci --ignore-scripts
```

The frontend runs at `http://localhost:5173`.

By default, API calls to `/api/*` are proxied to `http://localhost:3000`. If the backend is running somewhere else, copy `.env.example` to `.env` and change `VITE_API_TARGET`.

```bash
cp .env.example .env
```

Example:

```env
VITE_API_TARGET=http://localhost:3001
```
