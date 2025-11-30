# Build Stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy source code and build
COPY . .
RUN npm run build

# Production Stage
FROM nginx:stable-alpine

# Copy the Nginx configuration from the API Gateway
COPY ../api-gateway/nginx.conf /etc/nginx/conf.d/default.conf

# Remove default Nginx static content
RUN rm -rf /usr/share/nginx/html/*

# Copy the built React app to the Nginx static directory
COPY --from=builder /app/build /usr/share/nginx/html

# Expose the port
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
