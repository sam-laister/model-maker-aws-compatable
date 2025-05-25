import cv2
import argparse
import matplotlib.pyplot as plt

def plot_image_with_bbox(image_path, bbox):
    """
    Plots an image with a bounding box.

    Parameters:
    - image_path: str, path to the image file
    - bbox: tuple, bounding box coordinates (x_min, y_min, x_max, y_max)
    """
    # Load the image
    image = cv2.imread(image_path)
    if image is None:
        raise ValueError(f"Image not found at {image_path}")

    # Convert the image from BGR to RGB
    image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)

    # Unpack the bounding box coordinates
    x_min, y_min, x_max, y_max = bbox

    # Draw the bounding box on the image
    cv2.rectangle(image, (x_min, y_min), (x_max, y_max), color=(255, 0, 0), thickness=cv2.FILLED)

    # Plot the image
    plt.imshow(image)
    plt.axis('off')  # Hide the axis
    plt.show()

if __name__ == "__main__":

    # Set up argument parser
    parser = argparse.ArgumentParser(description="Plot an image with a bounding box.")
    parser.add_argument("image_path", type=str, help="Path to the image file")
    parser.add_argument("x_min", type=int, help="Minimum x coordinate of the bounding box")
    parser.add_argument("y_min", type=int, help="Minimum y coordinate of the bounding box")
    parser.add_argument("x_max", type=int, help="Maximum x coordinate of the bounding box")
    parser.add_argument("y_max", type=int, help="Maximum y coordinate of the bounding box")

    # Parse arguments
    args = parser.parse_args()

    # Extract arguments
    image_path = args.image_path
    bbox = (args.x_min, args.y_min, args.x_max, args.y_max)

    # Plot the image with bounding box
    plot_image_with_bbox(image_path, bbox)