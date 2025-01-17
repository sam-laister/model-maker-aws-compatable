#!/bin/bash

# Paths to input USDA and output GLB files
INPUT_FILE="$1"
OUTPUT_FILE="$2"

echo "Converting $INPUT_FILE to $OUTPUT_FILE"

# Run Blender in background mode
blender --background --python ./cmd/convert_usda_to_glb.py -- "$INPUT_FILE" "$OUTPUT_FILE"