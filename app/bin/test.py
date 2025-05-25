import matplotlib.pyplot as plt
import numpy as np
from scipy.spatial import ConvexHull
from PIL import Image

# Define the data
data = [
    {"point": [626, 700], "label": "base"},
    {"point": [671, 171], "label": "base"},
    {"point": [434, 534], "label": "gun"},
    {"point": [674, 506], "label": "leg"},
    {"point": [671, 340], "label": "base"},
    {"point": [323, 506], "label": "arm"},
    {"point": [369, 706], "label": "background"},
    {"point": [623, 545], "label": "leg"},
    {"point": [447, 295], "label": "base"},
    {"point": [369, 525], "label": "arm"},
    {"point": [479, 677], "label": "background"},
    {"point": [358, 293], "label": "base"},
    {"point": [371, 619], "label": "gun"},
    {"point": [676, 171], "label": "base"},
    {"point": [558, 523], "label": "leg"},
    {"point": [314, 465], "label": "armor"},
    {"point": [623, 651], "label": "background"},
    {"point": [326, 589], "label": "gun"},
    {"point": [347, 659], "label": "background"},
    {"point": [337, 626], "label": "gun"},
]

# Load an image (replace with your image path)
image_path = './084.png'
image = np.array(Image.open(image_path))

# Create figure and plot
fig, ax = plt.subplots(figsize=(8, 8))
ax.imshow(image, extent=[0, 800, 800, 0])  # Flip y-axis for pixel coords

# Define colors for different labels
colors = {
    "base": "red",
    "gun": "blue",
    "leg": "green",
    "arm": "purple",
    "background": "gray",
    "armor": "cyan",
    "helmet": "orange",
    "shoulder": "magenta"
}

# Group points by label
grouped_points = {}
for item in data:
    label = item["label"]
    if label not in grouped_points:
        grouped_points[label] = []
    grouped_points[label].append(item["point"])

# Plot convex hulls
for label, points in grouped_points.items():
    points = np.array(points)
    if len(points) > 2:  # Convex hull requires at least 3 points
        hull = ConvexHull(points)
        hull_points = points[hull.vertices]
        ax.fill(hull_points[:, 0], hull_points[:, 1], color=colors[label], alpha=0.3, label=label)

# Plot the points
for label, points in grouped_points.items():
    points = np.array(points)
    ax.scatter(points[:, 0], points[:, 1], color=colors[label], edgecolors="black", s=50)

# Remove duplicate labels from legend
handles, labels = ax.get_legend_handles_labels()
unique_labels = dict(zip(labels, handles))
ax.legend(unique_labels.values(), unique_labels.keys())

# Show the plot
plt.show()
