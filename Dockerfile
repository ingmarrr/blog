# Use a lightweight base image
FROM node:14-alpine AS frontend

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY frontend/package*.json ./

# Install frontend dependencies
RUN npm install

# Copy the rest of the frontend files
COPY frontend/ .

# Build the frontend
RUN npm run build

# Use a minimal golang image
FROM golang:1.16-alpine AS backend

# Set the working directory
WORKDIR /app

# Copy the backend files
COPY backend/ .

# Build the backend
RUN go build -o main .

# Use a minimal alpine image to run the final executable
FROM alpine:3.13

# Set the working directory
WORKDIR /app

# Copy the frontend files from the frontend build stage
COPY --from=frontend /app/dist/ ./dist/

# Copy the backend files from the backend build stage
COPY --from=backend /app/main .

# Expose the port the backend will listen on
EXPOSE 8080

# Run the backend
CMD ["./main"]

