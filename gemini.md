# Project instructions

This is to build a website as described in @system-prd.md

## Frontend

Use typescript, sveltekit (with Svelte v5), Vite, shadcn-svelte components  for the web front end

For Tailwind, always be sure to pay attention in using Tailwind version 4 conventions

To build the front end manually, don't forget to always prepend "cd frontend &&" to the npm build commands.

a dev server is used for hot reload. The dev server will be started externally, do not execute npm run dev yourself.

use 'cd frontend && npm run check' and 'cd frontend && npm run lint' to check for potential errors after making a related set of front end changes.

## Backend

For the backend use golang where needed, but prefer Firestore native web client for most web data use over an API service.

Backend services are still needed for anything more complex to support for the web app, and for back-end processing.

The backend should be built as a single golang server that will be deployed to Cloud Run

For triggers for ingestion processing, use GCP Cloud Scheduler or Firestore events as appropriate.

For any offline LLM calls, including embedding model use, use GCP Vertex and Gemini pro models.
