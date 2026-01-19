<template>
  <div class="card">
    <div style="display:flex; justify-content:space-between; gap:12px; align-items:flex-start;">
      <div>
        <h1 class="h1">Bots</h1>
        <div class="muted">Create a persona bot and chat with it. History is stored in SQLite.</div>
        <div class="muted" style="margin-top:6px;">Hololab System: <b>{{ health }}</b></div>
      </div>
      <RouterLink class="btn" to="/create">+ Create Bot</RouterLink>
    </div>

    <div class="sep"></div>

    <div v-if="loading" class="muted">Loading...</div>
    <div v-else-if="bots.length===0" class="muted">No bots yet.</div>

    <div v-else class="grid">
      <div class="item" v-for="b in bots" :key="b.id">
        <div style="display:flex; justify-content:space-between; gap:10px;">
          <div>
            <div style="font-weight:900;">{{ b.name }}</div>
            <div class="muted">{{ b.job }}</div>
            <div class="muted" style="margin-top:6px;">id: {{ b.id }}</div>
          </div>
          <button class="btn danger" @click="remove(b.id)">Delete</button>
        </div>

        <div style="margin-top:12px;">
          <RouterLink class="btn" :to="`/chat/${b.id}`">Open Chat</RouterLink>
        </div>
      </div>
    </div>

    <div v-if="error" style="margin-top:10px; color:#b00020;">{{ error }}</div>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { api } from "../services/api";

const bots = ref([]);
const loading = ref(true);
const error = ref("");
const health = ref("checking...");

async function load() {
  loading.value = true;
  error.value = "";
  try {
    const res = await api.listBots();
    bots.value = Array.isArray(res) ? res : [];
  } catch (e) {
    error.value = String(e?.message || e);
    bots.value = [];
  } finally {
    loading.value = false;
  }
}

async function remove(id) {
  if (!confirm("Delete bot and its history?")) return;
  try {
    await api.deleteBot(id);
    await load();
  } catch (e) {
    error.value = String(e?.message || e);
  }
}

onMounted(async () => {
  try {
    const h = await api.health();
    health.value = h.ok ? "online" : "unknown";
  } catch {
    health.value = "offline";
  }
  await load();
});
</script>

<style scoped>
.grid { 
  display:grid; 
  grid-template-columns: repeat(2, minmax(0,1fr)); gap:12px; 
}
.item { 
  border:1px solid #eee; 
  border-radius:14px; padding:14px; 
}
@media (max-width: 820px) { .grid { grid-template-columns: 1fr; } }
</style>
