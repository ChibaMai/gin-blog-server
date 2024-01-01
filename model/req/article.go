package req

// 新增/编辑 文章
type SaveOrUpdateArt struct {
	ID           int    `json:"id"`
	Title        string `json:"title" validate:"required" label:"文章标题"`
	Desc         string `json:"desc"`
	Content      string `json:"content" validate:"required" labe:"博客内容"`
	Img          string `json:"img"`
	Type         int8   `json:"type" validate:"required,min=1,max=3" label:"类型(1-原创 2-转载 3-翻译)"`
	Status       int8   `json:"status" validate:"required,min=1,max=3" label:"状态(1-公开 2-私密 3-评论可见)"`
	UpdateStatus int8   `json:"update_status" validate:"required,min=1,max=3" label:"更新状态(1-更新 2-断更 3-完结)"`
	IsTop        *int8  `json:"is_top" validate:"required,min=0,max=1" label:"是否置顶(0-否 1-是)"`
	OriginalUrl  string `json:"original_url"`
	Author       string `json:"author" validate:"required" label:"文章作者"`

	TagNames     []string `json:"tag_names"`
	TagIds       []uint   `json:"tag_ids"`
	CategoryName string   `json:"category_name"`
	CategoryId   int      `json:"category_id"`
}

// 条件查询文章
type GetArts struct {
	PageQuery
	Title        string `form:"title"`
	CategoryId   int    `form:"category_id"`
	TagId        int    `form:"tag_id"`
	Type         int8   `form:"type" validate:"required,min=1,max=3"`
	Status       int8   `form:"status" min=1,max=3"`
	UpdateStatus int    `form:"update_status" validate:"required,min=1,max=3"`
	IsDelete     *int8  `form:"is_delete" validate:"required,min=0,max=1"`
}

// 修改文章置顶
type UpdateArtTop struct {
	ID    int   `json:"id"`
	IsTop *int8 `json:"is_top" validate:"required,min=0,max=1" label:"是否置顶(0-否 1-是)"`
}

// 前台条件查询文章
type GetFrontArts struct {
	PageQuery
	CategoryId int `form:"category_id"`
	TagId      int `form:"tag_id"`
}

// 导入文章
type ImportArt struct {
}
