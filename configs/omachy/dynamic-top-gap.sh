#!/usr/bin/env bash
set -eu

HOME_DIR="${HOME}"
AEROSPACE_CONFIG="${HOME_DIR}/.config/aerospace/aerospace.toml"
SKETCHYBAR_CONFIG="${HOME_DIR}/.config/sketchybar/sketchybarrc"
STATE_DIR="${HOME_DIR}/.omachy"
MIN_GAP=16
MAX_GAP=64

if [ ! -f "$AEROSPACE_CONFIG" ] || [ ! -f "$SKETCHYBAR_CONFIG" ]; then
    exit 0
fi

if ! command -v swift >/dev/null 2>&1; then
    exit 0
fi

if ! command -v aerospace >/dev/null 2>&1; then
    exit 0
fi

bar_val() {
    key="$1"
    default="$2"
    value=$(awk -F'=' -v key="$key" '
        $0 ~ "^[[:space:]]*" key "[[:space:]]*=" {
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2)
            print $2
            exit
        }
    ' "$SKETCHYBAR_CONFIG")
    if [ -z "$value" ]; then
        value="$default"
    fi
    echo "$value"
}

HEIGHT=$(bar_val "height" "36")
MARGIN=$(bar_val "margin" "8")
Y_OFFSET=$(bar_val "y_offset" "4")

case "$HEIGHT" in ''|*[!0-9]*) HEIGHT=36 ;; esac
case "$MARGIN" in ''|*[!0-9]*) MARGIN=8 ;; esac
case "$Y_OFFSET" in ''|*[!0-9]*) Y_OFFSET=4 ;; esac

FOOTPRINT=$((HEIGHT + MARGIN + Y_OFFSET))

SWIFT_OUT=$(swift -e '
import AppKit

for s in NSScreen.screens {
    let frame = s.frame
    let visible = s.visibleFrame
    let topInset = Int(round(frame.maxY - visible.maxY))
    let name = (s.localizedName as NSString).lowercased
    print("\(name)|\(topInset)")
}
' 2>/dev/null || true)

if [ -z "$SWIFT_OUT" ]; then
    exit 0
fi

built_in_inset=""
external_max_gap=0

while IFS='|' read -r name top_inset; do
    [ -z "$name" ] && continue
    case "$top_inset" in ''|*[!0-9]*) continue ;; esac

    gap=$((FOOTPRINT - top_inset))
    if [ "$gap" -lt "$MIN_GAP" ]; then
        gap=$MIN_GAP
    fi
    if [ "$gap" -gt "$MAX_GAP" ]; then
        gap=$MAX_GAP
    fi

    case "$name" in
        *built-in*)
            built_in_inset="$gap"
            ;;
        *)
            if [ "$gap" -gt "$external_max_gap" ]; then
                external_max_gap="$gap"
            fi
            ;;
    esac
done <<EOF
$SWIFT_OUT
EOF

if [ -z "$built_in_inset" ]; then
    built_in_inset=$MIN_GAP
fi

if [ "$external_max_gap" -eq 0 ]; then
    external_max_gap=$((FOOTPRINT))
    if [ "$external_max_gap" -lt "$MIN_GAP" ]; then
        external_max_gap=$MIN_GAP
    fi
    if [ "$external_max_gap" -gt "$MAX_GAP" ]; then
        external_max_gap=$MAX_GAP
    fi
fi

NEW_OUTER_TOP="    outer.top        = [{monitor.'built-in' = ${built_in_inset}}, ${external_max_gap}]"

CURRENT_OUTER_TOP=$(awk '/^[[:space:]]*outer\.top[[:space:]]*=/{print; exit}' "$AEROSPACE_CONFIG" || true)

if [ "$CURRENT_OUTER_TOP" = "$NEW_OUTER_TOP" ]; then
    exit 0
fi

tmp_file="${AEROSPACE_CONFIG}.tmp"
awk -v replacement="$NEW_OUTER_TOP" '
    /^[[:space:]]*outer\.top[[:space:]]*=/ {
        print replacement
        next
    }
    { print }
' "$AEROSPACE_CONFIG" > "$tmp_file"
mv "$tmp_file" "$AEROSPACE_CONFIG"

aerospace reload-config >/dev/null 2>&1 || true
