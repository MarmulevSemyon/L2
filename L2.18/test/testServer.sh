#!/usr/bin/env bash
set -u

BASE_URL="${BASE_URL:-http://localhost:8080}"

pass_count=0
fail_count=0

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
  echo
  echo "=================================================="
  echo "$1"
  echo "=================================================="
}

ok() {
  echo -e "${GREEN}[OK]${NC}   $1"
  pass_count=$((pass_count + 1))
}

fail() {
  echo -e "${RED}[FAIL]${NC} $1"
  fail_count=$((fail_count + 1))
}

info() {
  echo -e "${BLUE}$1${NC}"
}

check_status() {
  local actual="$1"
  local expected="$2"
  local name="$3"

  if [[ "$actual" == "$expected" ]]; then
    ok "$name -> status $actual"
  else
    fail "$name -> expected status $expected, got $actual"
  fi
}

check_body_contains() {
  local body="$1"
  local needle="$2"
  local name="$3"

  if [[ "$body" == *"$needle"* ]]; then
    ok "$name -> body contains: $needle"
  else
    fail "$name -> body does not contain: $needle"
    echo "       body: $body"
  fi
}

request() {
  local method="$1"
  local path="$2"
  local data="${3:-}"

  local tmp_body
  tmp_body="$(mktemp)"

  local status
  if [[ -n "$data" ]]; then
    status=$(curl -sS -o "$tmp_body" -w "%{http_code}" \
      -X "$method" \
      -H "Content-Type: application/json" \
      -d "$data" \
      "$BASE_URL$path")
  else
    status=$(curl -sS -o "$tmp_body" -w "%{http_code}" \
      -X "$method" \
      "$BASE_URL$path")
  fi

  RESPONSE_BODY="$(cat "$tmp_body")"
  RESPONSE_STATUS="$status"
  rm -f "$tmp_body"
}

info "Перед запуском тестов сервер должен быть запущен на пустых данных на 8080 порте."
info "Если сервер уже был запущен ранее, останови его и запусти заново."
info "Например:"
echo "./bin/httpserver --port 8080"
info "Или:"
echo "go run ./cmd/httpserver --port 8080"
echo

read -r -p "Нажми Enter после запуска или перезапуска сервера..."

if ! curl -sS "$BASE_URL/events_for_day?user_id=1&date=2026-03-27" >/dev/null 2>&1; then
  fail "Сервер недоступен по адресу $BASE_URL"
  echo "Убедись, что сервер запущен и слушает нужный порт."
  exit 1
fi

ok "Сервер отвечает по адресу $BASE_URL"

ok "Сервер отвечает по адресу $BASE_URL"

print_header "1. Создание события"

request "POST" "/create_event" '{"user_id":1,"date":"2026-03-27","event":"meeting"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "create_event valid"
check_body_contains "$RESPONSE_BODY" '"user_id":1' "create_event valid"
check_body_contains "$RESPONSE_BODY" '"event":"meeting"' "create_event valid"

print_header "2. Получение событий за день"

request "GET" "/events_for_day?user_id=1&date=2026-03-27"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "events_for_day valid"
check_body_contains "$RESPONSE_BODY" '"event":"meeting"' "events_for_day valid"

print_header "3. Обновление события"

request "POST" "/update_event" '{"id":1,"user_id":1,"date":"2026-03-28","event":"updated meeting"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "update_event valid"
check_body_contains "$RESPONSE_BODY" '"id":1' "update_event valid"
check_body_contains "$RESPONSE_BODY" '"event":"updated meeting"' "update_event valid"

print_header "4. Проверка старой даты после update"

request "GET" "/events_for_day?user_id=1&date=2026-03-27"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "events_for_day old date after update"
check_body_contains "$RESPONSE_BODY" '"result":null' "events_for_day old date after update"

print_header "5. Проверка новой даты после update"

request "GET" "/events_for_day?user_id=1&date=2026-03-28"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "events_for_day new date after update"
check_body_contains "$RESPONSE_BODY" '"event":"updated meeting"' "events_for_day new date after update"

print_header "6. Создание событий для недели и месяца"

request "POST" "/create_event" '{"user_id":1,"date":"2026-03-23","event":"monday"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "create monday event"

request "POST" "/create_event" '{"user_id":1,"date":"2026-03-29","event":"sunday"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "create sunday event"

request "POST" "/create_event" '{"user_id":1,"date":"2026-03-30","event":"next week"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "create next week event"

request "POST" "/create_event" '{"user_id":1,"date":"2026-04-01","event":"april event"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "create april event"

request "POST" "/create_event" '{"user_id":2,"date":"2026-03-28","event":"other user event"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "create other user event"

print_header "7. Получение событий за неделю"

request "GET" "/events_for_week?user_id=1&date=2026-03-27"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "events_for_week valid"
check_body_contains "$RESPONSE_BODY" '"event":"monday"' "events_for_week valid"
check_body_contains "$RESPONSE_BODY" '"event":"updated meeting"' "events_for_week valid"
check_body_contains "$RESPONSE_BODY" '"event":"sunday"' "events_for_week valid"

print_header "8. Получение событий за месяц"

request "GET" "/events_for_month?user_id=1&date=2026-03-27"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "events_for_month valid"
check_body_contains "$RESPONSE_BODY" '"event":"monday"' "events_for_month valid"
check_body_contains "$RESPONSE_BODY" '"event":"updated meeting"' "events_for_month valid"
check_body_contains "$RESPONSE_BODY" '"event":"sunday"' "events_for_month valid"

print_header "9. Удаление события"

request "POST" "/delete_event" '{"id":1,"user_id":1}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "delete_event valid"
check_body_contains "$RESPONSE_BODY" '"result":"deleted"' "delete_event valid"

print_header "10. Проверка, что удалённое событие исчезло"

request "GET" "/events_for_day?user_id=1&date=2026-03-28"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "200" "events_for_day after delete"

if [[ "$RESPONSE_BODY" == *'"event":"updated meeting"'* ]]; then
  fail "deleted event is still present"
else
  ok "deleted event is absent"
fi

print_header "11. Негативные сценарии"

request "POST" "/create_event" '{"user_id":1,"date":"2026-99-99","event":"bad date"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "400" "create_event invalid date"
check_body_contains "$RESPONSE_BODY" '"error":"invalid date"' "create_event invalid date"

request "GET" "/events_for_day?user_id=abc&date=2026-03-27"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "400" "events_for_day invalid user_id"
check_body_contains "$RESPONSE_BODY" '"error":"invalid user_id"' "events_for_day invalid user_id"

request "GET" "/events_for_day?user_id=1&date=bad-date"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "400" "events_for_day invalid date"
check_body_contains "$RESPONSE_BODY" '"error":"invalid date"' "events_for_day invalid date"

request "POST" "/update_event" '{"id":999,"user_id":1,"date":"2026-03-27","event":"missing"}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "503" "update_event not found"
check_body_contains "$RESPONSE_BODY" '"error":"event not found"' "update_event not found"

request "POST" "/delete_event" '{"id":999,"user_id":1}'
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "503" "delete_event not found"
check_body_contains "$RESPONSE_BODY" '"error":"event not found"' "delete_event not found"

request "GET" "/events_for_day"
info "BODY: $RESPONSE_BODY"
check_status "$RESPONSE_STATUS" "400" "events_for_day missing params"
check_body_contains "$RESPONSE_BODY" '"error":"missing params' "events_for_day missing params"

print_header "Итог"

echo -e "${GREEN}Passed: $pass_count${NC}"
echo -e "${RED}Failed: $fail_count${NC}"

if [[ "$fail_count" -ne 0 ]]; then
  exit 1
fi

echo -e "${GREEN}Все проверки прошли успешно.${NC}" 