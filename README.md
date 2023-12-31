# barbar

Output barcodes to the terminal using the [kitty graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/).

## Installation

You need to have a working go toolchain.

To install: 

```
go build -o barbar .
```

## Usage

To generate a Code 128 barcode, use the `barcode` verb followed by the string you want to encode.

![an example of running barbar to generate a barcode](https://github.com/go-go-golems/barbar/assets/128441/0f497d33-3c29-4079-a427-f231b16f242f)

To generate a QR code, use the `qrcode` verb followed by the string you want to encode.

![an example of running bar to generate a qrcode](https://github.com/go-go-golems/barbar/assets/128441/f35500e8-3d83-45d8-8e5a-5cc69f801a7c)

