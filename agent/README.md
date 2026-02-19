# Ukoni Agent Harness

This project serves as the natural language interface for Ukoni. It is a separate service that interacts with the main Ukoni API and external LLM providers (OpenAI, Gemini).

## Features

- **Conversational Interface**: Users can interact with their household data using natural language.
- **Tool Use**: The agent can perform actions like searching products, adding items to shopping lists, and checking inventory levels.
- **Provider Agnostic**: Supports "Bring Your Own Key" for OpenAI and Gemini.
- **Secure**: Uses the same authentication as the main API, passing user credentials through.

## Architecture

The agent harness acts as a proxy:

1.  Receives request from Web Client (Prompt + Auth Token + Inventory Context).
2.  Fetches user's LLM settings (API Key) from the Main API.
3.  Constructs a prompt with available tools (based on API capabilities).
4.  Calls the LLM provider.
5.  Executes tools against the Main API using the user's Auth Token.
6.  Returns the final response to the Web Client.

## Development

### Prerequisites

- Go 1.25+
- Access to Ukoni Main API
- OpenAI or Gemini API Key

### Environment Variables

The agent server can be configured using the following environment variables:

| Variable | Description | Default |
|---|---|---|
| `PORT` | The port the server listens on | `8081` |
| `API_BASE_URL` | The base URL of the main Ukoni API | `http://localhost:8080` |
| `CORS_ALLOWED_ORIGIN` | Allowed origin for CORS | `http://localhost:5173` |
| `ENV` | The environment (development, production) | `development` |
| `DATABASE_URL` | The PostgreSQL connection string | `postgres://etin:etin@localhost:5432/ukoni?sslmode=disable` |

### Running

```bash
go run cmd/server/main.go
```
