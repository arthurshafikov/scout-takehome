# Multi-stage build for Scout frontend

# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies with legacy-peer-deps flag
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Build for production
RUN npm run build

# Stage 2: Serve with nginx
FROM nginx:alpine

WORKDIR /app

# Copy nginx config
COPY <<EOF /etc/nginx/conf.d/default.conf
server {
    listen 80;
    server_name _;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # SPA routing - fallback to index.html
    root /app;
    
    location / {
        try_files $uri $uri/ /index.html;
        expires 1d;
        add_header Cache-Control "public, immutable";
    }

    # API proxy to backend
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
EOF

# Copy built assets from builder stage
COPY --from=builder /app/dist /app

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
