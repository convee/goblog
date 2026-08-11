#!/usr/bin/env bash
set -euo pipefail

binary="${GOBLOG_BINARY:-./bin/goblog}"
config="${GOBLOG_CONFIG:-conf/config.yml}"
pid_file="${GOBLOG_PID_FILE:-.goblog.pid}"
log_file="${GOBLOG_LOG_FILE:-logs/server.log}"

running_pid() {
  if [[ ! -f "$pid_file" ]]; then
    return 1
  fi
  local pid
  pid="$(<"$pid_file")"
  if [[ "$pid" =~ ^[0-9]+$ ]] && kill -0 "$pid" 2>/dev/null; then
    printf '%s' "$pid"
    return 0
  fi
  rm -f "$pid_file"
  return 1
}

start() {
  local pid
  if pid="$(running_pid)"; then
    echo "goblog is already running (pid $pid)"
    return 0
  fi
  if [[ ! -x "$binary" ]]; then
    echo "binary not found or not executable: $binary" >&2
    exit 1
  fi
  if [[ ! -f "$config" ]]; then
    echo "config not found: $config" >&2
    exit 1
  fi
  mkdir -p "$(dirname "$log_file")"
  nohup "$binary" -c "$config" >>"$log_file" 2>&1 &
  pid=$!
  printf '%s\n' "$pid" >"$pid_file"
  echo "goblog started (pid $pid)"
}

stop() {
  local pid
  if ! pid="$(running_pid)"; then
    echo "goblog is not running"
    return 0
  fi
  kill "$pid"
  for _ in {1..20}; do
    if ! kill -0 "$pid" 2>/dev/null; then
      rm -f "$pid_file"
      echo "goblog stopped"
      return 0
    fi
    sleep 0.5
  done
  echo "goblog did not stop within 10 seconds" >&2
  exit 1
}

status() {
  local pid
  if pid="$(running_pid)"; then
    echo "goblog is running (pid $pid)"
  else
    echo "goblog is not running"
    exit 1
  fi
}

case "${1:-}" in
start) start ;;
stop) stop ;;
restart) stop; start ;;
status) status ;;
*)
  echo "usage: $0 {start|stop|restart|status}" >&2
  exit 1
  ;;
esac
