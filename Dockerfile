FROM ubuntu:latest

COPY keyguard /app/
COPY loader.sh /app/

ENV PORT 8000
EXPOSE 8000

WORKDIR "/app"
CMD "/app/keyguard"
