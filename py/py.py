import numpy as np
from PIL import Image
import time

def load_image_and_normalize(path):
    img = Image.open(path).convert('L')
    normalized = np.asarray(img, dtype=np.float32) / 255.0
    return normalized

def multiply_matrices(A, B):
    return A * B

def main():
    start_time = time.time()

    matrix_a = load_image_and_normalize("../img1.jpg")
    matrix_b = load_image_and_normalize("../img2.jpg")

    if matrix_a.shape != matrix_b.shape:
        print("Image dimensions do not match!")
        return

    result = multiply_matrices(matrix_a, matrix_b)

    output_image = Image.fromarray((result * 255).astype(np.uint8))
    output_image.save("output_image_python.jpg")

    elapsed_time = time.time() - start_time
    print(f"Bildverarbeitung und Matrixmultiplikation PYTHON: {elapsed_time:.3f}s")

if __name__ == "__main__":
    main()
