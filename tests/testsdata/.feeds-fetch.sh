#!/bin/sh
set -u

input="${1:-.feeds.txt}"
cr=$(printf '\r')

while IFS= read -r url || [ -n "$url" ]; do
    case "$url" in ''|\#*) continue ;; esac

    # extract domain: drop scheme, then path/port/userinfo
    domain="${url#*://}"
    domain="${domain%%/*}"
    domain="${domain##*@}"
    domain="${domain%%:*}"

    case "$url" in *.json) ext=json ;; *) ext=xml ;; esac

    out="$domain.$ext"

    # Avoid collision, some feed share the same domain
    if [ "$domain" = feeds.feedburner.com ]; then
        path="${url#*://*/}"
        path="${path%/}"
        out="$domain-$(printf '%s' "$path" | tr / _).$ext"
    fi

    if [ "$domain" = www.youtube.com ]; then
        channel="${url#*channel_id=}"
        channel="${channel%%&*}"
        out="$domain-$channel.$ext"
    fi

    if wget -q --timeout=5 -O "$out" "$url"; then
        echo "ok   $url -> $out"
    else
        echo "FAIL $url" >&2
        rm -f "$out"
    fi
done < "$input"
