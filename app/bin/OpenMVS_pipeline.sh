#!/bin/bash

set -e  # Exit on error

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <input_dir> <output_dir>"
    exit 1
fi

INPUT_DIR="$1"
OUTPUT_DIR="$2"
MVS_DIR="$OUTPUT_DIR/mvs"
OUTPUT_MODEL="$OUTPUT_DIR/final_model.ply"

mkdir -p "$MVS_DIR"

echo "0. Running OpenMVG pipeline..."
echo "Running command SfM_SequentialPipeline.py $INPUT_DIR $OUTPUT_DIR --opensfm-processes 8"
./bin/SfM_SequentialPipeline.py "$INPUT_DIR" "$OUTPUT_DIR" --opensfm-processes 8

echo "1. Converting OpenMVG SfM to OpenMVS format..."
echo "Running command openMVG_main_openMVG2openMVS -i $OUTPUT_DIR/reconstruction_sequential/sfm_data.bin -o scene.mvs"
openMVG_main_openMVG2openMVS -i "$OUTPUT_DIR/reconstruction_sequential/sfm_data.bin" -o "$MVS_DIR/scene.mvs"

echo "2. Densifying point cloud..."
echo "Running command DensifyPointCloud "scene.mvs" -o "scene_dense.mvs" -w "$MVS_DIR""
DensifyPointCloud "scene.mvs" -o "scene_dense.mvs" -w "$MVS_DIR"

echo "3. Reconstructing mesh..."
echo "Running command ReconstructMesh "scene_dense.mvs" -o "scene_mesh.mvs" -w "$MVS_DIR""
ReconstructMesh "scene_dense.mvs" -o "scene_mesh.ply" -w "$MVS_DIR"

echo "4. Refining mesh..."
echo "RefineMesh "scene.mvs" -m "scene_mesh.ply" -o "scene_dense_mesh_refine.mvs"  -w "$MVS_DIR" --scales 1 --max-face-area 16"
RefineMesh "scene.mvs" -m "scene_mesh.ply" -o "scene_dense_mesh_refine.mvs"  -w "$MVS_DIR" --scales 1 --max-face-area 16

echo "5. Texturing mesh..."
echo "Running command TextureMesh "scene_refined.mvs" -o "$OUTPUT_MODEL" -w "$MVS_DIR""
TextureMesh scene_dense.mvs -m scene_dense_mesh_refine.ply -o scene_dense_mesh_refine_texture.mvs -w "$MVS_DIR"

echo "OpenMVS processing complete!"
echo "Final model saved at: scene_dense_mesh_refine_texture.ply"
