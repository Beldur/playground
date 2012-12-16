CFLAGS=-Wall -g
EXECUTABLE=resize

all: $(EXECUTABLE)

resize:
	cc -o $(EXECUTABLE) resize.c $(CFLAGS) -O `GraphicsMagick-config --cppflags --ldflags --libs`

clean:
	rm -f resize
	rm -f test-thumb*