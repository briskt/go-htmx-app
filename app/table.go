package app

type TableStructureItem[T any] struct {
	Label      string
	RenderCell func(data T) string
}
