#!/bin/bash

# Тестовый скрипт для Lottery API
# Убедитесь, что сервер запущен на localhost:8080

API_URL="http://localhost:8080"
LAST_RESPONSE=""  # Глобальная переменная для хранения последнего ответа

echo "=== Тестирование Lottery API ==="
echo

# Функция для красивого вывода JSON
pretty_json() {
    if command -v jq &> /dev/null; then
        echo "$1" | jq .
    else
        echo "$1"
    fi
}

# Функция для выполнения запроса и отображения результата
make_request() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    
    echo "--- $description ---"
    echo "Запрос: $method $url"
    if [ ! -z "$data" ]; then
        echo "Данные: $data"
    fi
    echo
    
    if [ -z "$data" ]; then
        response=$(curl -s -X $method "$API_URL$url")
    else
        response=$(curl -s -X $method -H "Content-Type: application/json" -d "$data" "$API_URL$url")
    fi
    
    # Сохраняем ответ в глобальную переменную
    LAST_RESPONSE="$response"
    
    echo "Ответ:"
    pretty_json "$response"
    echo
    echo "----------------------------------------"
    echo
}

# Тест 1: Создание тиража
echo "🎯 Тест 1: Создание тиража"
make_request "POST" "/draws" "" "Создание нового тиража"

# Получим ID созданного тиража из последнего ответа
DRAW_ID=$(echo "$LAST_RESPONSE" | jq -r '.id')
if [ -z "$DRAW_ID" ] || [ "$DRAW_ID" == "null" ]; then
    echo "Ошибка: не удалось получить ID тиража из ответа"
    exit 1
fi
echo "Полученный ID тиража: $DRAW_ID"
echo

# Тест 2: Попытка создать второй тираж (должна вернуть ошибку)
echo "🎯 Тест 2: Попытка создать второй активный тираж"
make_request "POST" "/draws" "" "Создание второго тиража (должна быть ошибка)"

# Тест 3: Покупка билетов
echo "🎯 Тест 3: Покупка билетов"
make_request "POST" "/tickets" "{\"numbers\": [1, 5, 12, 23, 36], \"draw_id\": $DRAW_ID}" "Покупка билета #1"
make_request "POST" "/tickets" "{\"numbers\": [2, 8, 15, 24, 35], \"draw_id\": $DRAW_ID}" "Покупка билета #2"
make_request "POST" "/tickets" "{\"numbers\": [3, 9, 16, 25, 34], \"draw_id\": $DRAW_ID}" "Покупка билета #3"

# Тест 4: Валидация билетов
echo "🎯 Тест 4: Тестирование валидации"
make_request "POST" "/tickets" "{\"numbers\": [1, 2, 3, 4], \"draw_id\": $DRAW_ID}" "Билет с 4 числами (ошибка)"
make_request "POST" "/tickets" "{\"numbers\": [1, 2, 3, 4, 5, 6], \"draw_id\": $DRAW_ID}" "Билет с 6 числами (ошибка)"
make_request "POST" "/tickets" "{\"numbers\": [0, 2, 3, 4, 5], \"draw_id\": $DRAW_ID}" "Билет с числом 0 (ошибка)"
make_request "POST" "/tickets" "{\"numbers\": [1, 2, 3, 4, 37], \"draw_id\": $DRAW_ID}" "Билет с числом 37 (ошибка)"
make_request "POST" "/tickets" "{\"numbers\": [1, 1, 3, 4, 5], \"draw_id\": $DRAW_ID}" "Билет с повторяющимися числами (ошибка)"

# Тест 5: Получение информации о тираже
echo "🎯 Тест 5: Получение информации о тираже"
make_request "GET" "/draws/$DRAW_ID" "" "Получение информации о тираже"

# Тест 6: Получение всех тиражей
echo "🎯 Тест 6: Получение всех тиражей"
make_request "GET" "/draws" "" "Получение списка всех тиражей"

# Тест 7: Закрытие тиража
echo "🎯 Тест 7: Закрытие тиража"
make_request "POST" "/draws/$DRAW_ID/close" "" "Закрытие тиража и определение победителей"

# Тест 8: Получение результатов
echo "🎯 Тест 8: Получение результатов"
make_request "GET" "/draws/$DRAW_ID/results" "" "Получение результатов тиража"

# Тест 9: Попытка купить билет в закрытом тираже
echo "🎯 Тест 9: Попытка купить билет в закрытом тираже"
make_request "POST" "/tickets" "{\"numbers\": [10, 11, 12, 13, 14], \"draw_id\": $DRAW_ID}" "Покупка билета в закрытом тираже (ошибка)"

# Тест 10: Попытка закрыть уже закрытый тираж
echo "🎯 Тест 10: Попытка закрыть уже закрытый тираж"
make_request "POST" "/draws/$DRAW_ID/close" "" "Закрытие уже закрытого тиража (ошибка)"

# Тест 11: Создание нового тиража после закрытия предыдущего
echo "🎯 Тест 11: Создание нового тиража после закрытия предыдущего"
make_request "POST" "/draws" "" "Создание нового тиража после закрытия предыдущего"

echo "=== Тестирование завершено ==="
echo
echo "Для более детальной проверки результатов запустите:"
echo "curl -s $API_URL/draws | jq ."
echo "curl -s $API_URL/draws/$DRAW_ID/results | jq ."