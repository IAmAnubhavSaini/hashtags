# hashtags

I use this to generate graph of my obsidian vaults. This is faster than their implementation.

`mkdir -p dist && go run src/hashtags.go`

`python -m http.server --bind 127.0.0.1 7979 --directory ./src` and then visit http://127.0.0.1:7979/src/viz.html

1. caching for faster reloads of graphs.
2. localStorage based settings, selections
3. Customizable depth control for tag-file-tag relationships
4. Ignored tags filtering to simplify complex graphs
5. Interactive visualization with zoom, drag, and click navigation
6. Search functionality for finding specific tags
7. Tag frequency tracking and bubble visualization
8. File information panel when clicking on file nodes
9. Custom color themes stored in localStorage
10. Cached tag graphs for improved performance

