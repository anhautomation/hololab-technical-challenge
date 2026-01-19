<template>
  <div class="wrap">
    <header class="top">
      <div class="brand">Hololab Persona Bot</div>
      <nav class="nav">
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/create">Create</RouterLink>
      </nav>
    </header>

    <main class="main">
      <RouterView />
    </main>

    <footer class="foot">
      <div>Chatbot LLM • Vue 3 + Go + SQLite</div>

      <div class="powered">
        <span>Powered by</span>

        <a
          class="icon-link"
          href="https://www.linkedin.com/in/anh-nguyen-4b400485/"
          target="_blank"
          rel="noopener"
          aria-label="LinkedIn"
          title="LinkedIn"
        >
          <svg class="icon" viewBox="0 0 24 24" aria-hidden="true">
            <path
              d="M4.98 3.5C4.98 4.88 3.87 6 2.5 6S0 4.88 0 3.5 1.12 1 2.5 1s2.48 1.12 2.48 2.5zM.5 23.5h4V7.5h-4v16zM8.5 7.5h3.8v2.2h.05c.53-1 1.83-2.2 3.77-2.2 4.03 0 4.78 2.65 4.78 6.1v9.9h-4v-8.8c0-2.1-.04-4.8-2.93-4.8-2.93 0-3.38 2.28-3.38 4.64v8.96h-4v-16z"
            />
          </svg>
        </a>
        Anh
        <span class="dot">·</span>

        <a
          class="icon-link"
          href="https://github.com/anhautomation/hololab-chatbot-llm"
          target="_blank"
          rel="noopener"
          aria-label="GitHub"
          title="GitHub"
        >
          <svg class="icon" viewBox="0 0 24 24" aria-hidden="true">
            <path
              d="M12 .5C5.73.5.5 5.74.5 12.02c0 5.1 3.29 9.43 7.86 10.96.58.1.79-.25.79-.56 0-.28-.01-1.02-.02-2-3.2.7-3.88-1.55-3.88-1.55-.52-1.33-1.27-1.69-1.27-1.69-1.04-.71.08-.7.08-.7 1.15.08 1.76 1.18 1.76 1.18 1.02 1.74 2.67 1.24 3.32.95.1-.74.4-1.24.73-1.53-2.55-.29-5.23-1.28-5.23-5.69 0-1.26.45-2.28 1.18-3.08-.12-.29-.51-1.45.11-3.02 0 0 .97-.31 3.18 1.18a11.1 11.1 0 012.9-.39c.98 0 1.97.13 2.9.39 2.2-1.49 3.17-1.18 3.17-1.18.62 1.57.23 2.73.11 3.02.73.8 1.18 1.82 1.18 3.08 0 4.42-2.69 5.39-5.25 5.68.41.35.77 1.04.77 2.1 0 1.52-.01 2.74-.01 3.11 0 .31.21.67.8.56A11.53 11.53 0 0023.5 12C23.5 5.74 18.27.5 12 .5z"
            />
          </svg>
        </a>
      </div>
    </footer>

  </div>
</template>

<script setup></script>

<style>
* { box-sizing: border-box; }
body { 
  margin:0; 
  font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial; 
  background:#f6f7fb; 
  color:#111; 
}
a { 
  color: inherit; 
  text-decoration:none; 
}
.wrap { 
  min-height:100vh; 
  display:flex; 
  flex-direction:column; 
}
.top { 
  background:#fff; 
  border-bottom:1px solid #e7e7ef; 
  padding:14px 18px; 
  display:flex; 
  justify-content:space-between; 
  align-items:center; 
}
.brand { font-weight:900; }
.nav { 
  display:flex; 
  gap:10px; 
}
.nav a { 
  padding:8px 10px; 
  border-radius:10px; 
}
.nav a.router-link-active { 
  background:#111; 
  color:#fff; 
}
.main { 
  width: min(1100px, 92vw); 
  margin: 18px auto; flex:1; 
}
.foot {
  padding: 14px 18px;
  border-top: 1px solid #e7e7ef;
  background: #fff;
  color: #666;
  font-size: 12px;
  text-align: center;
  line-height: 1.6;
}
.card { 
  background:#fff; 
  border:1px solid #e7e7ef; 
  border-radius:14px; 
  padding:16px; }
.h1 { 
  font-size:22px; 
  margin:0 0 10px; 
}
.h2 { 
  font-size:16px; 
  margin:0 0 10px; 
}
.muted { 
  color:#666; 
  font-size:13px; 
}
.sep { 
  height:1px; 
  background:#eee; 
  margin:14px 0; 
}
.row { 
  display:flex; 
  gap:12px; 
}
.col { 
  flex:1; 
}
.input, .textarea { 
  width:100%; 
  border:1px solid #ddd; 
  border-radius:12px; 
  padding:10px 12px; 
  outline:none; 
  background:#fff; 
}
.textarea { 
  min-height:110px; 
  resize:vertical; 
}
.btn { 
  border:1px solid #111; 
  background:#111; 
  color:#fff; 
  padding:10px 12px; 
  border-radius:12px; 
  cursor:pointer; 
}
.btn.secondary { 
  background:#fff; color:#111; 
  border-color:#ddd; 
}
.btn.danger { 
  background:#b00020; border-color:#b00020; 
}
.powered {
  margin-top: 6px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-wrap: nowrap;          
  white-space: nowrap;      
  line-height: 1;            
  color: #999;
  font-weight: 600;
}

.icon-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #111;
  text-decoration: none;
  border-bottom: 1px dotted #ddd;
  padding-bottom: 1px;
}
.icon-link:hover {
  border-bottom-color: #111;
}
.dot {
  color: #bbb;
}
.icon {
  width: 14px !important;
  height: 14px !important;
  display: inline-block !important;
  flex: 0 0 14px;
  fill: currentColor !important;
}
</style>
