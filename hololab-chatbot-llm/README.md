# Persona-based LLM Chatbot System

This project presents a **persona-based LLM chatbot system** designed with clear responsibility boundaries, extensibility, and safe interaction patterns.

## Problem Scope

The system allows users to define a chatbot persona and ensures that:

- The chatbot consistently stays in character
- Responses follow the configured tone and personality
- Answers are restricted to explicitly defined knowledge
- Missing or out-of-scope information is handled safely
- Core application logic remains independent from any specific LLM provider

## System Capabilities

The application enables users to:

1. Create a chatbot persona with:
   - Name
   - Occupation
   - Background / bio
   - Communication style / tone
   - Allowed knowledge scope

2. Interact with the chatbot through a web interface

3. Verify persona consistency and knowledge boundary enforcement

## Architectural Intent

This system is structured with production-oriented thinking:

- Persona configuration is treated as structured input, not hardcoded prompts
- Prompt construction and LLM interaction are encapsulated behind a provider interface
- The backend enforces persona and knowledge constraints
- The frontend remains a thin interaction layer

These boundaries allow future scaling and provider substitution.

## High-level Architecture

The system is organized around clear responsibility and trust boundaries:

- The frontend is responsible only for persona configuration and user interaction
- The backend owns persona validation, prompt construction, and response enforcement
- LLM providers are accessed through a pluggable interface and are treated as untrusted execution components

All behavioral guarantees — persona consistency, tone adherence, and knowledge boundaries —
are enforced by the backend before and after LLM execution.

## Project Structure

hololab-core-system  
Go backend (persona logic, LLM orchestration, REST API)

hololab-persona-chatbot  
Vue 3 frontend (persona management and chat interface)

The backend is responsible for behavior enforcement, while the frontend focuses on usability and clarity.

## LLM Integration Strategy

The system supports two interchangeable execution modes.

### Internal Mock Mode

LLM_MODE = mock

- Uses an internal mock LLM implementation
- Provides deterministic and auditable behavior
- Validates persona logic and prompt boundaries
- Suitable for evaluation without external dependencies

### External LLM Mode

LLM_MODE = openai_compat

- Uses an OpenAI-compatible API
- Requires no changes in application logic
- Demonstrates clean separation between system behavior and AI providers

## Running Locally

Backend:

go run ./cmd/api

Frontend:

npm install  
npm run dev

## Live Evaluation

A running instance of the system is available for interactive evaluation.

🔗 [Open live system](https://hololab-technical-challenge.vercel.app/)

## Implementation Notes

- SQLite is used for simplicity and fast iteration
- Long-term persistence is intentionally out of scope
- The system is designed for live evaluation and walkthrough
- The architecture supports future extension such as caching layers,
  multi-provider routing, and cost-aware LLM fallback strategies
