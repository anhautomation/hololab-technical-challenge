# Hololab Technical Challenge – Demo Submission

This repository contains my submission for the **Hololab Technical Challenge**, consisting of **two independent tasks**:

- **Task 1 – System Design**
- **Task 2 – Programming Skills (Demo Application)**

Each task is intentionally separated to reflect different evaluation goals.

---

## Task 1 – System Design (Digital Human Platform)

This task focuses on **system architecture and technical design thinking** for a Digital Human platform.

### Scope
The design illustrates:
- A Digital Human creation and interaction platform
- Clear separation between core product logic and AI providers
- Use of third-party LLM, Voice, and Avatar/Expression AI
- Internal orchestration and control for cost, safety, and scalability
- Data, caching, monitoring, and global delivery considerations


No runnable backend or frontend code is included for this task, as it is strictly a **design exercise**.

---

## Task 2 – Programming Skills (Persona-based Chatbot Demo)

This task demonstrates the implementation of a **persona-configurable chatbot**, allowing users to:
- Create a chatbot with name, role, bio, style, and allowed knowledge
- Chat with the bot and receive responses consistent with its configured persona
- Switch between mock LLM behavior and real LLM providers (via configuration)

### Technical Stack
- **Backend:** Go
- **Frontend:** Vue.js
- **Database:** SQLite (local, lightweight)
- **LLM Integration:** Mock provider by default, optional external LLM via API key

### Purpose
The demo is designed to:
- Showcase clean architecture and system boundaries
- Demonstrate correct use of LLMs as an external capability, not core logic
- Emphasize correctness of behavior over model sophistication
- Be easy to review, run, and reason about

---

## Notes

- Task 1 and Task 2 are **intentionally decoupled**.
- The system design is illustrative and not a production blueprint.
- The demo application focuses on behavior, clarity, and architecture rather than advanced AI techniques.

Thank you for reviewing this submission.

