FROM node:25 AS nodebuilder
RUN mkdir -p /app/ui
WORKDIR /app
COPY ui/ ./ui
WORKDIR /app/ui
RUN npm install
RUN npm run build

FROM golang:1.26 AS gobuilder
RUN mkdir -p /app /tmp/bin
WORKDIR /app
COPY . .
COPY --from=nodebuilder /app/static /app/static
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/bin/filamate .

FROM scratch AS filamate
COPY --from=gobuilder /tmp/bin/filamate /filamate
EXPOSE 9766
ENTRYPOINT ["/filamate", "daemon", "start"]
