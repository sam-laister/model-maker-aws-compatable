#!/usr/bin/env python
#! -*- encoding: utf-8 -*-

import bpy
import sys
import os

def convert_ply_to_glb(input_path, output_path):
    # Clear existing mesh objects
    bpy.ops.object.select_all(action='SELECT')
    bpy.ops.object.delete()
    
    # Import PLY file
    bpy.ops.wm.ply_import(filepath=input_path)
    
    # Select all objects
    bpy.ops.object.select_all(action='SELECT')

    # Export as GLB
    bpy.ops.wm.export_scene.gltf(
        filepath=output_path,
        export_format='GLB',
        use_selection=True,
        export_draco_mesh_compression_enable=True  # Enable Draco compression
    )

if __name__ == "__main__":
    if len(sys.argv) < 5:
        print("Usage: blender -b -P convert_ply_to_glb.py -- input.ply output.glb")
        sys.exit(1)
        
    input_file = sys.argv[-2]
    output_file = sys.argv[-1]
    
    if not os.path.exists(input_file):
        print(f"Error: Input file {input_file} does not exist")
        sys.exit(1)
        
    try:
        convert_ply_to_glb(input_file, output_file)
        print(f"Successfully converted {input_file} to {output_file}")
    except Exception as e:
        print(f"Error during conversion: {str(e)}")
        sys.exit(1)