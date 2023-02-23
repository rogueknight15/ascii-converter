# Ascii-Converter
Command line tool for converting PNG images to ascii art.

# Usage

```Console
$ go build ascii-converter.go
$ ./ascii-converter.exe -f <filepath>
```

# Command Line Arguments

- f - Path to a PNG image
- xscale - Scaling factor for the x axis
- yscale - Scaling factor for the y axis
- format - Output format:
    - 0 = NO_COLOR
    - 1 = FOREGROUD_COLORED
    - 2 = BACKGROUND_COLORED
    - 3 = PNG
- o - Output filepath for the converted image

# Example

```Console
$ ./ascii-converter.exe -f example.png -xscale=2 -yscale=2 -format=2 -o out.txt
```
Converts an image "example.png" to an ascii art image with ansi escaped characters with a background color and scaled by a factor of 2 on both axis as out.txt