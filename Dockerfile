# Stage 1: Build the frontend
FROM node:20-slim AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.24.4-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go get firebase.google.com/go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Final stage: Create the production image
FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=backend /app/main .
COPY --from=frontend /app/frontend/build ./frontend/build
EXPOSE 8080
CMD ["/app/main"]
