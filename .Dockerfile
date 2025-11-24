FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY ai-api /bin/ai-api

CMD ["ai-api"]