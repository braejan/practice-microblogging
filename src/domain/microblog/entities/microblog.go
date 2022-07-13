package entities

type MicroBlogTracking struct {
	LikesCount    int64 `json:"likes"`
	DislikesCount int64 `json:"dislikes"`
	VisitCount    int64 `json:"views"`
}

type MicroBlog struct {
	ID     string             `json:"id"`
	UserID string             `json:"user_id"`
	Text   string             `json:"text"`
	Detail *MicroBlogTracking `json:"detail,omitempty"`
}
