FROM node:20-alpine

WORKDIR /app

COPY web/ ./

RUN npm install --omit=dev

RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]