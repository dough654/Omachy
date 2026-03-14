#!/bin/bash

PERCENTAGE="$(pmset -g batt | grep -Eo "\d+%" | cut -d% -f1)"
CHARGING="$(pmset -g batt | grep 'AC Power')"

if [ -z "$PERCENTAGE" ]; then
  exit 0
fi

if [ -n "$CHARGING" ]; then
  ICON="󰂄"
  COLOR=0xffa6e3a1
elif [ "$PERCENTAGE" -gt 80 ]; then
  ICON="󰁹"
  COLOR=0xffa6e3a1
elif [ "$PERCENTAGE" -gt 60 ]; then
  ICON="󰂀"
  COLOR=0xffcdd6f4
elif [ "$PERCENTAGE" -gt 40 ]; then
  ICON="󰁾"
  COLOR=0xfff9e2af
elif [ "$PERCENTAGE" -gt 20 ]; then
  ICON="󰁼"
  COLOR=0xfffab387
else
  ICON="󰁺"
  COLOR=0xfff38ba8
fi

sketchybar --set "$NAME" icon="$ICON" icon.color="$COLOR" label="${PERCENTAGE}%"
