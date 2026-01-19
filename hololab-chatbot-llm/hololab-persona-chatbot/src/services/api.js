const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:3001";

async function request(path, opts) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "content-type": "application/json" },
    ...opts,
  });
  if (!res.ok) {
    const t = await res.text().catch(() => "");
    throw new Error(`${res.status} ${t}`);
  }
  return res.json();
}

export const api = {
  health: () => request("/api/health"),
  listBots: () => request("/api/bots"),
  createBot: (payload) =>
    request("/api/bots", { method: "POST", body: JSON.stringify(payload) }),
  getBot: (id) => request(`/api/bots/${id}`),
  deleteBot: (id) => request(`/api/bots/${id}`, { method: "DELETE" }),

  listMessages: (id) => request(`/api/bots/${id}/messages`),
  sendMessage: (id, message) =>
    request(`/api/bots/${id}/messages`, {
      method: "POST",
      body: JSON.stringify({ message }),
    }),
  resetChat: (id) =>
    request(`/api/bots/${id}/messages/reset`, {
      method: "POST",
      body: JSON.stringify({}),
    }),
};
