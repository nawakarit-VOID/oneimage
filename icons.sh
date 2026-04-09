#!/bin/bash
set -e
export PATH=/usr/local/go/bin:$PATH

INPUT="icon.png"
OUTDIR="icons"

mkdir -p $OUTDIR

SIZES=(512 256 128 64)

for SIZE in "${SIZES[@]}"; do
  convert "$INPUT" \
    -resize ${SIZE}x${SIZE} \
    "$OUTDIR/icon-${SIZE}.png"
done

echo "✅ เสร็จแล้ว!"
