#!/bin/bash

PROCESS_NAME="smart_home_hub"
OUTPUT_FILE="process_data.csv"
INTERVAL=2
DURATION=3600
PORT=5001

DISK_DEV=$(lsblk -no PKNAME $(df / | tail -1 | awk '{print $1}') | head -n 1)
if [ -z "$DISK_DEV" ]; then
    DISK_DEV=$(df / | tail -1 | awk '{print $1}' | awk -F/ '{print $NF}')
fi

PID=$(pgrep -f $PROCESS_NAME | head -n 1)

if [ -z "$PID" ]; then
    echo "ERROR: process not found"
    exit 1
fi

echo "timestamp,cpu[%],ram[MB],app_write[MB],sys_write[KBs],IO_Wait,net_rtt[ms]" > $OUTPUT_FILE

END_TIME=$((SECONDS + DURATION))

while [ $SECONDS -lt $END_TIME ]; do
    STATS=$(ps -p $PID -o %cpu,rss --no-headers)
    CPU=$(echo $STATS | awk '{print $1}' | tr -d '\r\n')
    RAM_KB=$(echo $STATS | awk '{print $2}' | tr -d '\r\n')
    RAM_MB=$(echo "scale=2; $RAM_KB / 1024" | bc | tr -d '\r\n')

    TS=$(date +"%H:%M:%S")

    APP_WRITE_BYTES=$(cat /proc/$PID/io 2>/dev/null | grep "write_bytes" | awk '{print $2}' | tr -d '\r\n')
    if [ -z "$APP_WRITE_BYTES" ]; then
        APP_WRITE_MB="0.00"
    else
        APP_WRITE_MB=$(echo "scale=2; $APP_WRITE_BYTES / 1048576" | bc | tr -d '\r\n')
    fi

    RTT=$(ss -tin "sport = :$PORT" | awk -F'rtt:' '/rtt:/ {print $2}' | cut -d'/' -f1 | awk '{sum+=$1; count++} END {if (count > 0) print sum/count; else print 0}' | tr -d '\r\n')

    SYS_WRITE_KBS=$(iostat -dk 1 2 | grep -w "$DISK_DEV" | tail -1 | awk '{print $4}' | xargs)
    if [ -z "$SYS_WRITE_KBS" ]; then SYS_WRITE_KBS="0.00"; fi

    IO_WAIT=$(top -bn1 | grep "Cpu(s)" | awk '{print $10}' | tr -d 'wa,' | xargs)
    if [ -z "$IO_WAIT" ]; then IO_WAIT="0.0"; fi

    echo "${TS},${CPU},${RAM_MB},${APP_WRITE_MB},${SYS_WRITE_KBS},${IO_WAIT},${RTT}" >> $OUTPUT_FILE

    sleep $INTERVAL
done