package categories

type CategoryInput struct {
	Name string
}

type CategoriesResult struct {
	Categories []Category
	Total      int64
	Page       int
	Limit      int
	TotalPages int64
}
