package request

type CreateBlogPostRequest struct {
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Excerpt         string   `json:"excerpt"`
	Content         string   `json:"content"`
	CoverImageURL   string   `json:"cover_image_url"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Tags            []string `json:"tags"`
	Publish         bool     `json:"publish"`
}

type UpdateBlogPostRequest struct {
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Excerpt         string   `json:"excerpt"`
	Content         string   `json:"content"`
	CoverImageURL   string   `json:"cover_image_url"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Tags            []string `json:"tags"`
	Status          string   `json:"status"`
}
