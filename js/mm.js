const sharp = require("sharp");
const os = require("os");

const { Worker, isMainThread, parentPort, workerData } = require("worker_threads");

async function loadImageAndNormalize(path) {
  const image = await sharp(path).greyscale().raw().toBuffer();

  const normalized = new Float32Array(image.length);
  for (let i = 0; i < image.length; i++) {
    normalized[i] = image[i] / 255.0;
  }

  return normalized;
}

function multiplyMatrices(A, B, width, height, start, end) {
  const result = new Float32Array(width * (end - start));
  let index = 0;
  for (let i = start; i < end; i++) {
    for (let j = 0; j < width; j++) {
      result[index++] = A[i * width + j] * B[i * width + j];
    }
  }
  return result;
}

if (isMainThread) {
  async function main() {
    const [width, height] = [7680, 4320];
    const numWorkers = os.cpus().length;
    const rowsPerWorker = height / numWorkers;

    console.time("Bildverarbeitung und Matrixmultiplikation JAVASCRIPT");

    const matrixA = await loadImageAndNormalize("../img1.jpg");
    const matrixB = await loadImageAndNormalize("../img2.jpg");

    const workers = [];
    for (let i = 0; i < numWorkers; i++) {
      const worker = new Worker(__filename, {
        workerData: {
          A: matrixA,
          B: matrixB,
          width,
          height,
          start: i * rowsPerWorker,
          end: (i + 1) * rowsPerWorker,
        },
      });
      workers.push(worker);
    }

    const results = await Promise.all(
      workers.map((worker) => {
        return new Promise((resolve) => {
          worker.on("message", resolve);
        });
      })
    );

    const resultMatrix = new Float32Array(width * height);
    for (let i = 0; i < numWorkers; i++) {
      resultMatrix.set(results[i], i * rowsPerWorker * width);
    }

    console.timeEnd("Bildverarbeitung und Matrixmultiplikation JAVASCRIPT");

    const outputImage = new Uint8Array(resultMatrix.length);
    for (let i = 0; i < resultMatrix.length; i++) {
      outputImage[i] = resultMatrix[i] * 255;
    }
    await sharp(Buffer.from(outputImage), { raw: { width, height, channels: 1 } }).toFile("output_image.jpg");
  }

  main();
} else {
  const { A, B, width, height, start, end } = workerData;
  const result = multiplyMatrices(A, B, width, height, start, end);
  parentPort.postMessage(result);
}
