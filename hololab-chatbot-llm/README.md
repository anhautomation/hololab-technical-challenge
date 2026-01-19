# Persona Chatbot LLM

This project is a working demo for **Task 2 – Programming Skills** of the Hololab technical challenge.

The purpose of this demo is to show the ability to:
- Build a configurable chatbot
- Enforce persona, tone, and knowledge boundaries
- Deliver a clean and testable frontend–backend flow

This is a demo project for evaluation only.

---

## What the Demo Does

The application allows users to:

1. Create a chatbot persona with:
   - Name
   - Occupation
   - Bio
   - Style / tone
   - Allowed knowledge

2. Chat with the created persona.

3. Verify that the chatbot:
   - Stays in character
   - Uses only allowed knowledge
   - Clearly refuses or defers when information is missing

---

## Project Structure

hololab-core-system  
Go backend (REST API, persona logic, LLM adapter)

hololab-persona-chatbot  
Vue 3 frontend (UI for creating bots and chatting)

---

## LLM Behavior

The demo supports two modes:

LLM_MODE = mock  
Uses an internal mock LLM to validate persona logic and scope enforcement.

LLM_MODE = openai_compat  
Uses a real OpenAI-compatible API when configured.

No application logic changes are required to switch modes.

---

## Running Locally

Backend:

go run ./cmd/api

Frontend:

npm install  
npm run dev

---

## Notes

- SQLite is used only for demo simplicity.
- Long-term persistence is not required for this challenge.
- The demo is designed to be run and evaluated live.

---

Demo submission for technical evaluation.
