# Insighta Labs+ — Profile Intelligence Platform

Insighta Labs+ is a secure, multi-interface platform for discovering and analyzing intelligence profiles. This platform was built for **HNGi14 Stage 3** to demonstrate advanced authentication, persistence, and multi-client integration.

### 🔗 Live Links
- **Live API (Azure)**: [https://stage1.doxantro.com](https://stage1.doxantro.com)
- **Web Portal (Vercel)**: [https://insighta-web-plum.vercel.app](https://insighta-web-plum.vercel.app)
- **CLI Repository**: [https://github.com/OBigVee/Insighta-cli](https://github.com/OBigVee/Insighta-cli)

## 🏛 System Architecture

The platform consists of three core components:
1. **Backend (Go)**: A high-performance API server using Chi, SQLx, and Neon PostgreSQL.
2. **Web Portal (JS/Vite)**: A premium SPA with real-time analytics and CSV exports.
3. **CLI Tool (Go/Cobra)**: A terminal-based interface with Docker and Docker Compose support.

## 🔐 Authentication & Security

We implemented a robust OAuth 2.0 flow using GitHub, designed to bypass modern browser third-party cookie restrictions:

### Unified Auth Flow
1. **Protocol**: GitHub OAuth2 with PKCE (Proof Key for Code Exchange).
2. **Token Delivery**: To ensure reliability across cross-origin deployments (Vercel to Azure), tokens are delivered via **URL Hash Fragments** and stored as **Bearer Tokens** in `localStorage`.
3. **Token Rotation**: We use short-lived Access Tokens (3m) and Refresh Tokens (5m). **Every refresh action rotates the entire token pair**, providing maximum security against replay attacks.
4. **Stateless vs. Stateful**: Access tokens are verified statelessly (JWT), while Refresh tokens are statefully tracked and hashed in the database for instant revocation support.

## 🛡 Role-Based Access Control (RBAC)

The system enforces two distinct roles via server-side middleware:

- **Analyst (Default)**: Read-only access. Can list, search, view, and export data.
- **Admin**: Full read/write access. Required for creating (`POST /api/profiles`) or deleting (`DELETE /api/profiles`) records.

*Role enforcement is handled by the `RequireRole("admin")` middleware on protected routes.*

## 🧠 Natural Language Intelligence

The search engine interprets conversational queries into structured SQL filters:
- **Keyword Inference**: Maps "men", "women", "kids", "elders" to structured categories.
- **Dynamic Age Logic**: "young" maps to `16-24`, while "above 30" uses Regex extraction.
- **Country Resolution**: Scans queries against a full ISO-3166 dictionary using word-boundary matching.

## 💻 Interfaces

### Web Portal
A premium dashboard featuring:
- Real-time statistics and data visualization.
- Secure client-side CSV generation using `apiFetch`.
- **Manual CLI Login**: An integration tool that generates a one-time login command for remote CLI users.

### CLI Tool
Built for developers and automated workflows:
- **Commands**: `login`, `whoami`, `profiles list`, `profiles search`, `profiles export`.
- **Docker Support**: Pre-configured `Dockerfile` and `docker-compose.yml` for zero-install usage.
- **Remote Mode**: Supports `auth-set` for authenticating in remote/cloud IDEs where browser redirects are blocked.

## 🛠 Setup & Environment
The backend requires the following environment variables:
- `DATABASE_URL`: Neon PostgreSQL connection string.
- `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET`: OAuth credentials.
- `JWT_SECRET`: For signing session tokens.
- `WEB_PORTAL_URL`: The production frontend URL for secure redirection.

---
