use std::time::Instant;

fn load_image_and_normalize(path: &str) -> Vec<f32> {
    let img = image::open(path).unwrap().to_luma8();
    let (width, height) = img.dimensions();

    let mut normalized = Vec::with_capacity((width * height) as usize);
    for y in 0..height {
        for x in 0..width {
            let pixel = img.get_pixel(x, y).0[0] as f32 / 255.0;
            normalized.push(pixel);
        }
    }
    normalized
}

fn multiply_matrices(a: &[f32], b: &[f32]) -> Vec<f32> {
    a.iter()
        .zip(b.iter())
        .map(|(a_val, b_val)| a_val * b_val)
        .collect()
}

fn main() {
    let start_time = Instant::now();

    let matrix_a = load_image_and_normalize("../img1.jpg");
    let matrix_b = load_image_and_normalize("../img2.jpg");

    let _result = multiply_matrices(&matrix_a, &matrix_b);

    let elapsed_time = start_time.elapsed();
    println!(
        "Bildverarbeitung und Matrixmultiplikation RUST: {:?}",
        elapsed_time
    );
}
