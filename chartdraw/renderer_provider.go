package chartdraw

// RendererProvider is a function that returns a renderer.
type RendererProvider func(int, int) Renderer
