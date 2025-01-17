import bpy
import sys

# Get input and output file paths from command-line arguments
input_file = sys.argv[-2]
output_file = sys.argv[-1]

# Clear all existing objects in Blender
bpy.ops.wm.read_factory_settings(use_empty=True)

# Import the USDA file
bpy.ops.wm.usd_import(filepath=input_file)

# Export the scene to GLB
bpy.ops.export_scene.gltf(filepath=output_file, export_format='GLB')