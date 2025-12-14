# Project: Online Marketplace Community (Hackathon MVP)

## Context
This project is the backend for a peer-to-peer online marketplace.
**Note:** This is a hackathon submission. Prioritize development speed and feature completion over production-grade security or scalability.

## Architecture

### High-Level Pattern
We follow a **Modular Monolith** architecture using the standard **Handler-Service-Repository** pattern. This ensures separation of concerns without the complexity of microservices.

### Layer Responsibilities & Rules

1.  **Handler (`internal/handler`)**
    * **Role:** Parse HTTP requests, validate basic inputs, call Service, send HTTP responses.
    * **Rule:** **NO** SQL queries here. **NO** complex business logic here.
    * **Pattern:** Struct-based handlers (`UserHandler` struct) holding a Service interface/struct.

2.  **Service (`internal/service`)**
    * **Role:** The "brain." Orchestrates business logic, combines data from multiple repositories if needed.
    * **Rule:** Pure Go logic. It should not know about HTTP (no `w http.ResponseWriter`).
    * **Dependency:** Accepts Repositories via interfaces (or concrete structs for speed).

3.  **Repository (`internal/repository`)**
    * **Role:** The "storage." Performs CRUD operations against the database.
    * **Rule:** The **only** place where `database/sql` is imported and SQL is written.
    * **Pattern:** Return concrete structs (`*UserRepo`), not interfaces (keep it simple for MVP).

4.  **App Wiring (`internal/app`)**
    * **Role:** Dependency Injection.
    * **Pattern:** `NewApp(db)` initializes Repos, injects them into Services, injects Services into Handlers.

### Tech Stack & Constraints
* **Language:** Go (1.25).
* **Database:** MySQL.
* **Authentication:** Firebase Auth.

## Workflow & Documentation

### Feature Tracking
- **Location:** Maintain a living record of functionality in `docs/features.md`.
- **Status:** Clearly explicitly specify whether each feature is **[Implemented]** or **[Pending]**.
- **Scope:** Track significant business capabilities only (e.g., "User Listing Creation", "Search"). Do not track trivial UI interactions (e.g., "delete confirmation dialogs", "hover states").

### Error Handling Policy
- Do not expose internal error details to clients.
- Log internal errors using the standard `log` package.
- Return generic, plain-text error messages in HTTP responses.
- Preserve appropriate HTTP status codes (e.g., 400 for bad input, 404 for not found, 500 for internal errors).

### Handler Documentation Policy
- Every HTTP handler must include a documentation comment directly above the function.
- The comment should specify, at minimum:
  - Route and method (e.g., `POST /users`).
  - Required headers (e.g., `Content-Type: application/json`).
  - Request JSON schema with field names and types, noting which fields are required or optional.
  - Success response schema (field names and types) and status code.
  - Possible error responses with status codes (use generic messages per Error Handling Policy).
- Purpose: Ensure API contracts are explicit and discoverable within the codebase.
