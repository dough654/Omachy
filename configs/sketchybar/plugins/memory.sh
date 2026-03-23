#!/usr/bin/env bash
# Show memory usage percentage with color-coded Nerd Font icon.

TOTAL=$(sysctl -n hw.memsize)
TOTAL_GB=$(( TOTAL / 1073741824 ))

# vm_stat reports pages; parse the page size from its header line.
VM=$(vm_stat)
PAGE_SIZE=$(echo "$VM" | awk -F'[() ]+' '/page size of/ {print $8; exit}')
if [ -z "$PAGE_SIZE" ]; then
    PAGE_SIZE=4096
fi
WIRED=$(echo "$VM"   | awk '/wired down/ {gsub(/\./,""); print $4}')
ANON=$(echo "$VM"    | awk '/Anonymous pages/ {gsub(/\./,""); print $3}')
COMPRESSED=$(echo "$VM" | awk '/occupied by compressor/ {gsub(/\./,""); print $5}')

if [ -z "$ANON" ]; then
    ANON=$(echo "$VM" | awk '/Pages active/ {gsub(/\./,""); print $3}')
fi

USED_BYTES=$(( (WIRED + ANON + COMPRESSED) * PAGE_SIZE ))
USED_GB=$(( USED_BYTES / 1073741824 ))
PCT=$(( USED_BYTES * 100 / TOTAL ))

if [ "$PCT" -ge 80 ]; then
    COLOR=0xfff38ba8   # red
elif [ "$PCT" -ge 60 ]; then
    COLOR=0xfff9e2af   # yellow
else
    COLOR=0xffa6e3a1   # green
fi

sketchybar --set "$NAME" icon="󰍛" icon.color="$COLOR" label="${USED_GB}/${TOTAL_GB}GB"
