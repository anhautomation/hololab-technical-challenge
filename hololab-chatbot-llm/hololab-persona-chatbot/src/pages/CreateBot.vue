<template>
  <div class="card">
    <h1 class="h1">Create Persona Bot</h1>
    <div class="muted">Required: name, occupation, bio, style. Optional: allowed knowledge.</div>

    <div class="sep"></div>

    <div class="row">
      <div class="col">
        <div class="muted">Name</div>
        <input class="input" v-model.trim="form.name" placeholder="e.g., Mia" />
      </div>
      <div class="col">
        <div class="muted">Occupation</div>
        <input class="input" v-model.trim="form.job" placeholder="e.g., Product Designer" />
      </div>
    </div>

    <div style="margin-top:10px;">
      <div class="muted">Bio</div>
      <textarea class="textarea" v-model.trim="form.bio" placeholder="Background + constraints..."></textarea>
    </div>

    <div style="margin-top:10px;">
      <div class="muted">Style / Personality</div>
      <textarea class="textarea" v-model.trim="form.style" placeholder="Friendly, concise, bullet points, etc."></textarea>
    </div>

    <div style="margin-top:10px;">
      <div class="muted">Allowed Knowledge (optional)</div>
      <textarea class="textarea" v-model.trim="form.knowledge" placeholder="Only facts this bot is allowed to use..."></textarea>
    </div>

    <div v-if="error" style="margin-top:10px; color:#b00020;">{{ error }}</div>

    <div style="display:flex; gap:10px; margin-top:14px;">
      <button class="btn" @click="create" :disabled="loading">{{ loading ? "Creating..." : "Create" }}</button>
      <RouterLink class="btn secondary" to="/">Cancel</RouterLink>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { api } from "../services/api";

const router = useRouter();
const loading = ref(false);
const error = ref("");

const form = reactive({
  name: "",
  job: "",
  bio: "",
  style: "",
  knowledge: ""
});

async function create() {
  error.value = "";
  if (!form.name || !form.job || !form.bio || !form.style) {
    error.value = "Please fill: name, occupation, bio, style.";
    return;
  }
  loading.value = true;
  try {
    const bot = await api.createBot({ ...form });
    router.push(`/chat/${bot.id}`);
  } catch (e) {
    error.value = String(e?.message || e);
  } finally {
    loading.value = false;
  }
}
</script>
