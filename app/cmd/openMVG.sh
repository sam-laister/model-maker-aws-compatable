#!/bin/bash

set -euo pipefail

# Define dataset and output directories
DATASET_DIR="$PWD/dataset"
OUTPUT_DIR="$PWD/output"

# Ensure output directory exists
mkdir -p "$OUTPUT_DIR"

# Run OpenMVG pipeline in Docker
sudo docker run --rm \
	-v "$DATASET_DIR":/dataset:ro \
	-v "$OUTPUT_DIR":/output \
	openmvg \
	python3 /opt/openMVG_Build/software/SfM/SfM_SequentialPipeline.py /dataset/ /output/
