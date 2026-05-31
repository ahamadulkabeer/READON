const defaultHeaders = {
  Accept: "application/json",
};

export async function apiRequest(path, options = {}) {
  const response = await fetch(`/api${path}`, {
    credentials: "include",
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
    ...options,
  });

  const payload = await response.json().catch(() => null);

  if (!response.ok) {
    const message = payload?.message || payload?.error || "Request failed";
    throw new Error(typeof message === "string" ? message : "Request failed");
  }

  return payload;
}

export function getBooks() {
  return apiRequest("/books");
}

export function getCategories() {
  return apiRequest("/categories");
}
