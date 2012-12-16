#include <stdio.h>
#include <string.h>
#include <time.h>
#include <magick/api.h>

int main(int argc, char **argv)
{
	ExceptionInfo exception;
	int i = 0;
	char buffer[128];
	clock_t startTime = clock();

	// Init GraphicsMagick
	InitializeMagick(*argv);
	GetExceptionInfo(&exception);

	for(; i < 20; i++) {
		// Load image
		ImageInfo *image_info = CloneImageInfo((ImageInfo *) NULL);
		(void) strcpy(image_info->filename, "test.jpg");
		Image *image = ReadImage(image_info, &exception);

		if (image == (Image *) NULL) {
			CatchException(&exception);
		}

		// Create Thumbnail
		Image *thumbnail = ThumbnailImage(image, 200, 100, &exception);
		sprintf(buffer, "test-thumb-%d.jpg", i);
		(void) strcpy(thumbnail->filename, buffer);

		if (thumbnail == (Image *) NULL) {
			CatchException(&exception);
		}

		// Save thumbnail
		int writeResult = WriteImage(image_info, thumbnail);

		if (writeResult == 0) {
			CatchException(&thumbnail->exception);
		}

		// Cleanup
		DestroyImage(image);
		DestroyImageInfo(image_info);
	}

	printf("Time taken: %.2fs\n", (double)(clock() - startTime) / CLOCKS_PER_SEC);

	// Cleanup
	DestroyExceptionInfo(&exception);
	DestroyMagick();

	return 0;
}