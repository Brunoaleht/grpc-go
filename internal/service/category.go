package service

import (
	"context"
	"grpc-example/internal/database"
	"grpc-example/internal/pb"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	categoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{categoryDB: categoryDB}
}

func (c *CategoryService) CreateCategory(ctx context.Context, input *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.categoryDB.Create(input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil

}

func (c *CategoryService) GetListCategories(ctx context.Context, input *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.categoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*pb.Category

	for _, category := range categories {

		categoryResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}

		categoriesResponse = append(categoriesResponse, categoryResponse)
	}

	return &pb.CategoryList{Categories: categoriesResponse}, nil
}

func (c *CategoryService) GetCategoryById(ctx context.Context, input *pb.GetCategoryRequest) (*pb.Category, error) {

	category, err := c.categoryDB.FindByID(input.Id)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {

	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}

		if err != nil {
			return err
		}

		createdCategory, err := c.categoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		})

	}

}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		createdCategory, err := c.categoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		categoryResponse := &pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		}

		err = stream.Send(categoryResponse)
		if err != nil {
			return err
		}

	}

}
