package utils

import (
	"blog_server/blogpb"
	"blog_server/domain"
)

func DataToBlogPb(data *domain.BlogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}
}
