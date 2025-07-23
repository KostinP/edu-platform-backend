#!/bin/bash

# Скрипт автоматизирует сборку бэкенда на сервере

SERVER_PATH="/srv/edu-platform-backend"
SERVER_IP="213.139.208.67"

ssh root@$SERVER_IP << EOF
  cd $SERVER_PATH
  git pull
  docker compose build backend
  docker compose up -d backend
EOF

echo "✅ Деплой завершён"
