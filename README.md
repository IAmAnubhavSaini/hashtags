# hashtags

I use this to generate graph of my obsidian vaults. This is faster than their implementation.

## Usage

### Basic usage

```bash
# Create required directory and run with default path (~/notes)
mkdir -p dist && go run src/hashtags.go

# Run with custom path specified as command line argument
mkdir -p dist && go run src/hashtags.go -path=/path/to/notes

# Create a .env file with FULL_PATH=/path/to/notes for persistent configuration
echo "FULL_PATH=/path/to/notes" > .env
mkdir -p dist && go run src/hashtags.go
```

### .env file format

Create a file named `.env` in the project root with:

```
# Path to your notes directory
FULL_PATH=/path/to/your/notes
```

### Visualizing the graph

```bash
# Start a local web server
python -m http.server --bind 127.0.0.1 7979 --directory ./src
```

Then visit http://127.0.0.1:7979/ in your browser.

### Using Docker

#### Processing utility (hashtags:builder)

```bash
# Build Docker image
docker build --file Dockerfile.builder --tag hashtags:builder .

# Run with volume mounted to access your notes and output files
docker run --rm --volume /path/to/notes:/notes --volume $(pwd)/dist:/app/dist hashtags:builder -path=/notes
```

#### Visualization (hashtags:view)

```bash
# Build Docker image
docker build --file Dockerfile.view --tag hashtags:view .

# Run with the dist directory mounted to visualize the generated data
docker run --rm --volume $(pwd)/dist:/usr/share/nginx/html/dist --publish 8080:80 hashtags:view
```

Then visit http://localhost:8080/ in your browser.

#### Complete Docker Compose setup

You can also use Docker Compose to run both services together:

```bash
# Start both services
docker compose up

# Access the visualization at http://localhost:8080
```

## Features

1. Caching for faster reloads of graphs.
2. localStorage based settings, selections
3. Customizable depth control for tag-file-tag relationships
4. Ignored tags filtering to simplify complex graphs
5. Interactive visualization with zoom, drag, and click navigation
6. Search functionality for finding specific tags
7. Tag frequency tracking and bubble visualization
8. File information panel when clicking on file nodes
9. Custom color themes stored in localStorage
10. Cached tag graphs for improved performance
11. Multiple ways to specify notes directory (command line, .env file, default path)
12. Docker support for both processing and visualization components

## LICENSE

MIT License

Copyright (c) 2025 Anubhav Saini
