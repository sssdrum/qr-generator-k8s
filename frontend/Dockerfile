FROM oven/bun:alpine

WORKDIR /app
COPY package.json .
RUN bun install
RUN bun add -g serve
COPY . .
RUN bun run build
EXPOSE 3000
CMD ["serve", "-s", "dist"]
