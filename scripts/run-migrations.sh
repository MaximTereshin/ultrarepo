#!/bin/sh

# Максимальное количество попыток
MAX_RETRIES=30
RETRY_COUNT=0

# Ожидание готовности базы данных
echo "Waiting for database..."
while ! migrate -database "${DATABASE_URL}" -path . version >/dev/null 2>&1; do
    RETRY_COUNT=$((RETRY_COUNT+1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        echo "Error: Timed out waiting for database"
        exit 1
    fi
    echo "Attempt $RETRY_COUNT of $MAX_RETRIES..."
    sleep 1
done

echo "Running migrations..."
migrate -database "${DATABASE_URL}" -path . up

# Проверка статуса выполнения миграций
if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
    exit 0
else
    echo "Migration failed"
    exit 1
fi 