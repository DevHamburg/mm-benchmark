#define STB_IMAGE_IMPLEMENTATION
#include "stb_image.h"
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h"
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <windows.h>
#include <time.h> 

int getNumOfCores() {
    SYSTEM_INFO sysinfo;
    GetSystemInfo(&sysinfo);
    return sysinfo.dwNumberOfProcessors;
}

int width, height, channels;
unsigned char *img1, *img2, *result;

void *multiply_matrices(void *arg) {
    int start = ((int*)arg)[0];
    int end = ((int*)arg)[1];
    for (int i = start; i < end; i++) {
        result[i] = (unsigned char)(img1[i] * img2[i] / 255.0);
    }
    return NULL;
}

int main() {
    img1 = stbi_load("../img1.jpg", &width, &height, &channels, STBI_grey);
    img2 = stbi_load("../img2.jpg", &width, &height, &channels, STBI_grey);
    result = malloc(width * height * channels);

    clock_t start_time = clock();

    int num_threads = getNumOfCores();
    pthread_t threads[num_threads];
    int rows_per_thread = height / num_threads;

    for (int i = 0; i < num_threads; i++) {
        int start = i * rows_per_thread * width;
        int end = (i + 1) * rows_per_thread * width;
        int args[2] = {start, end};
        pthread_create(&threads[i], NULL, multiply_matrices, args);
    }

    for (int i = 0; i < num_threads; i++) {
        pthread_join(threads[i], NULL);
    }

    clock_t end_time = clock();
    double elapsed_time = (double)(end_time - start_time) / CLOCKS_PER_SEC;

    printf("Bildverarbeitung und Matrixmultiplikation C: %.2f ms\n", elapsed_time * 1000.0);

    stbi_write_jpg("output_image.jpg", width, height, channels, result, 100);

    free(img1);
    free(img2);
    free(result);
    return 0;
}
