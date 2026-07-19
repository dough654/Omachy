#!/usr/bin/env bash
# Show battery percentage with a Nerd Font icon.
PERCENTAGE=$(pmset -g batt | grep -Eo '[0-9]+%' | tr -d '%')
CHARGING=$(pmset -g batt | grep 'AC Power')

if [ -z "$PERCENTAGE" ]; then
    # No battery (Mac mini, Mac Studio, Mac Pro) — hide the item instead of
    # showing a permanent "?".
    sketchybar --set "$NAME" drawing=off
    exit 0
fi

if [ -n "$CHARGING" ]; then
    ICON="󰂄"                    # charging
    COLOR=0xffa6e3a1            # green
elif [ "$PERCENTAGE" -ge 80 ]; then
    ICON="󰁹" COLOR=0xffcdd6f4   # full
elif [ "$PERCENTAGE" -ge 60 ]; then
    ICON="󰂁" COLOR=0xffcdd6f4
elif [ "$PERCENTAGE" -ge 40 ]; then
    ICON="󰁾" COLOR=0xfff9e2af   # yellow
elif [ "$PERCENTAGE" -ge 20 ]; then
    ICON="󰁼" COLOR=0xfff38ba8   # red-ish
else
    ICON="󰁺" COLOR=0xfff38ba8   # critical
fi

sketchybar --set "$NAME" drawing=on icon="$ICON" icon.color="$COLOR" label="${PERCENTAGE}%"
