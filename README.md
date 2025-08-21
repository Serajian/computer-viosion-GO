# 🖼️ GoCV Face Detection

A lightweight Go application using **GoCV** and **Haar Cascade** for real-time webcam-based face detection. Features a live preview with a heads-up display (HUD) showing FPS and hotkeys, plus eye validation to reduce false positives.

## Features

- **Real-time Preview**: Webcam feed with HUD showing FPS and hotkeys.
- **Face Detection**: Haar Cascade with eye validation for accuracy.
- **Hotkeys**:
    - `ESC`: Exit the application.
    - `S`: Save a snapshot.
    - `F`: Flip (mirror) the feed.
    - `G`: Toggle grayscale mode.
    - `E`: Apply Canny edge detection.
    - `N`: Revert to normal mode.

## Prerequisites

- **OS**: macOS (tested on Apple Silicon)
- **Go**: ≥ 1.21
- **OpenCV**: 4.x (via Homebrew)
- **Tools**: `pkg-config`

## Installation

### 1. Install Dependencies

```bash
brew update
brew install opencv pkg-config
```

### 2. Configure Environment

Add OpenCV to `PKG_CONFIG_PATH`:

```bash
echo 'export PKG_CONFIG_PATH="$(brew --prefix opencv)/lib/pkgconfig:$PKG_CONFIG_PATH"' >> ~/.zshrc
source ~/.zshrc
```

Verify OpenCV installation:

```bash
pkg-config --modversion opencv4
```

### 3. Set Up Project

```bash
git clone <repository-url>
cd <repository-directory>
go mod tidy
```

## Usage

Run the application:

```bash
go run .
```

### Hotkeys

- `ESC`: Quit
- `S`: Save snapshot
- `F`: Flip feed
- `G`: Grayscale
- `E`: Canny edges
- `N`: Normal mode

## Contributing

1. Fork the repository.
2. Create a branch (`git checkout -b feature/your-feature`).
3. Commit changes (`git commit -m 'Add feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

Ensure code follows the project’s style and includes tests.

## License

[MIT License](LICENSE.txt)

## Acknowledgments

- [GoCV](https://gocv.io/)
- [OpenCV](https://opencv.org/)
- Haar Cascade classifiers
