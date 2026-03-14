#!/bin/bash

if [ "$SELECTED" = "true" ]; then
  sketchybar --set "$NAME" \
    icon.highlight=on \
    background.drawing=on \
    background.color=0xff313244
else
  sketchybar --set "$NAME" \
    icon.highlight=off \
    background.drawing=off
fi
