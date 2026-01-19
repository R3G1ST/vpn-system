FROM node:18-alpine

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Используем простой сервер для development
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
