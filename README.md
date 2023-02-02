# Downloader CLI

![Preview](https://github.com/SantiagoBedoya/down/blob/main/images/image.png)

## Usage

```bash
Usage of down:
  -c int
        Concurrent workers (default 5)
  -dest string
        Destination folder (default $HOME)
  -mode int
        Download Mode (concurrent: 0 | normal: 1) (default: 0)
  -url string
        File URL to Download
```

## How to install

### Using golang

```bash
go install github.com/SantiagoBedoya/down@latest
```

### Using homebrew

```bash
brew tap SantiagoBedoya/santiagobedoya
brew install SantiagoBedoya/down
```
