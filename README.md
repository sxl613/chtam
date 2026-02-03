# chtam

A small PDF pagination utility written in Go.

## What it does

Adds page numbers to PDF files. That's it.

I wrote this because I kept needing to add page numbers to documents before printing, and I wanted a laser-focused tool for exactly that. Also, I got to play around with the PDF format (I don't like it).

## Installation
Build from source:

```bash
git clone https://github.com/sxl613/chtam.git
cd chtam
go build
```

## Usage

```bash
# Basic usage - adds page numbers to output.pdf
./stamper -input input.pdf -output output.pdf

```

Run `./stamper -h` for all options.

## Dependencies

Uses [pdfcpu](https://github.com/pdfcpu/pdfcpu) for PDF manipulation.
