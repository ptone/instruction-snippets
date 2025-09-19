# LLM Prompt Engineering Instruction Management System

This project is a web application designed to manage and organize LLM prompt engineering instructions, specifically those intended for coding agents.

## Getting Started

### Prerequisites

*   [Go](https://golang.org/doc/install)
*   [Node.js](https://nodejs.org/en/download/)
*   Firebase project and credentials. You will need to set up a `.env` file in the root of the project with your Firebase configuration. See `.env.example` for the required variables.

### Frontend

To start the frontend development server:

```bash
cd frontend
npm install
npm run dev
```

The application will be available at `http://localhost:5173`.

### Backend

To start the backend server:

```bash
cd backend
go run main.go
```

The server will start on port `8080`.
