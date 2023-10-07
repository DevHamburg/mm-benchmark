import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.util.stream.IntStream;

public class MP {

    public static void main(String[] args) {
        try {
            BufferedImage img1 = ImageIO.read(new File("../img1.jpg"));
            BufferedImage img2 = ImageIO.read(new File("../img2.jpg"));

            if (img1.getWidth() != img2.getWidth() || img1.getHeight() != img2.getHeight()) {
                System.out.println("Image dimensions do not match!");
                return;
            }

            long startTime = System.currentTimeMillis();

            BufferedImage result = new BufferedImage(img1.getWidth(), img1.getHeight(), BufferedImage.TYPE_BYTE_GRAY);

            IntStream.range(0, img1.getHeight()).parallel().forEach(y -> {
                for (int x = 0; x < img1.getWidth(); x++) {
                    int rgb1 = img1.getRGB(x, y);
                    int rgb2 = img2.getRGB(x, y);

                    int gray1 = (rgb1 >> 16) & 0xFF;  
                    int gray2 = (rgb2 >> 16) & 0xFF;  

                    int grayResult = (int) (gray1 * gray2 / 255.0);
                    int rgbResult = (grayResult << 16) | (grayResult << 8) | grayResult;

                    result.setRGB(x, y, rgbResult);
                }
            });

            ImageIO.write(result, "jpg", new File("output_image.jpg"));

            long endTime = System.currentTimeMillis();
            System.out.println("Bildverarbeitung und Matrixmultiplikation JAVA: " + (endTime - startTime) + "ms");

        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
