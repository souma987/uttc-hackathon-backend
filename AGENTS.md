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
