FROM node:16.17.0-alpine3.16 AS build

# Copy source
WORKDIR /app/backend
COPY . .

# Download application dependencies 
RUN npm ci

# Build application
RUN npm run build

FROM node:16.17.0-alpine3.16 AS runtime

WORKDIR /opt/app
COPY package.json package.json
COPY package-lock.json package-lock.json

RUN npm ci --production

COPY --from=build /app/backend/bin bin

USER node
CMD ["npm", "start"]