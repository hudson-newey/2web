# TODO: Because 2web applications optionally use the database, load balancer,
# and ssr, we should support not having these installed in the project.
# If the load balancer does not have an explicit config, we should still use
# caddy to serve the files, but just use the default sdk config.
# TODO: I should also support non-docker alternatives
FROM node:22-slim

WORKDIR /app

# Copy server files
COPY dist/server/ssr.js ./dist/server/
COPY server/package.json .
COPY server/package-lock.json .

# Install dependencies
RUN npm install

# Copy database and load balancer docker compose scripts
COPY database/docker-compose.yml ./database/
COPY load-balancer/docker-compose.yml ./load-balancer/

# Install Docker CLI (needed to run docker-compose)
RUN apt-get update && apt-get install -y docker.io

# Create startup script
RUN echo '#!/bin/bash\n\
docker-compose -f database/docker-compose.yml up -d\n\
docker-compose -f load-balancer/docker-compose.yml up -d\n\
node dist/server/ssr.js' > /app/start.sh && \
chmod +x /app/start.sh

EXPOSE 3000

CMD ["/app/start.sh"]
