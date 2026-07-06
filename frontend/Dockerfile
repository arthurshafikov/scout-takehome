# Multi-stage build for Scout frontend

# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Accept build arguments
ARG VITE_API_KEY=scout-api-key-12345

# Copy package files
COPY package*.json ./

# Install dependencies with legacy-peer-deps flag
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Build for production with environment variable
RUN VITE_API_KEY=$VITE_API_KEY npm run build

# Stage 2: Serve with nginx
FROM nginx:alpine

WORKDIR /app

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist /app

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
