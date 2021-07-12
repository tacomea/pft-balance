package main

import (
	"blog_client/blogpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Creating blog request")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not conntect: %v\n", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create blog
	blogId := createBlog(c)
	_ = createBlog(c)
	_ = createBlog(c)

	// read blog
	readBlog(c, blogId)

	// update blog
	updateBlog(c, blogId)

	// delete blog
	deleteBlog(c, blogId)

	// list blogs
	listBlog(c)

}

func createBlog(c blogpb.BlogServiceClient) string {
	blog := blogpb.Blog{
		AuthorId: "Takumi",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: &blog})
	if err != nil {
		log.Fatalf("unexpected error : %v\n", err)
	}
	blogId := createBlogRes.GetBlog().GetId()
	fmt.Println("blog has been created: ", createBlogRes)
	return blogId
}

func readBlog(c blogpb.BlogServiceClient, blogId string) {
	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "jeafed"})
	if err != nil {
		fmt.Printf("correct - error happened while reading: %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	readBlogRes, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Printf("wrong - error happened while reading: %v\n", err)
	}
	fmt.Printf("blog has been read: %v\n", readBlogRes)
}

func updateBlog(c blogpb.BlogServiceClient, blogId string) {
	newBlog := blogpb.Blog{
		Id:       blogId,
		AuthorId: "Takumi",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: &newBlog})
	if err != nil {
		fmt.Println("wrong - error while updating")
	}
	fmt.Printf("blog has been updated: %v\n", updateRes)
}

func deleteBlog(c blogpb.BlogServiceClient, blogId string) {
	deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})
	if err != nil {
		fmt.Printf("wrong - error happened while deleting: %v\n", err)
	}
	fmt.Printf("blog has been updated: %v\n", deleteRes)
}

func listBlog(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error in ListBlog(): %v \n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something happened: %v \n", err)
		}
		fmt.Println("list of blogs: ", res.GetBlog())
	}
}
