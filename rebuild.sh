#!/usr/bin/env bash
# Reconstruye la imagen del fork desde el fuente local y recrea el contenedor.
# Uso: ./rebuild.sh   (desde ~/mozhi)
set -e

PARENT="$(cd "$(dirname "$0")/.." && pwd)"   # directorio padre (donde viven mozhi/ y libmozhi/)

echo ">>> Construyendo mozhi-fork:local..."
cd "$PARENT"
docker build -f mozhi/Dockerfile.local -t mozhi-fork:local .

echo ">>> Recreando contenedor..."
cd "$PARENT/mozhi"
docker compose up -d --force-recreate

echo ">>> Imagen en uso:"
docker ps --filter name=mozhi --format '{{.Names}}\t{{.Image}}'
echo ">>> Listo."