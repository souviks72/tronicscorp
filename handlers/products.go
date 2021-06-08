package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/souviks72/tronicscorp/dbiface"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

var (
	v = validator.New()
)

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitmepty"`
	Name        string             `json:"product_name" bson:"product_name" validate:"required,max=10"`
	Price       int                `json:"price" bson:"price" validate:"required,max=2000"`
	Currency    string             `json:"currency" bson:"currency" validate:"reqired,len=3"`
	Discount    int                `json:"discount" bson:"discount"`
	Vendor      string             `json:"vendor" bson:"vendor" validate:"required"`
	Accessories []string           `json:"accessories,omitempty" bson:"accessories,omitempty"`
	IsEssential bool               `json:"is_essential" bson:"is_essential"`
}

type ProductHandler struct {
	Col dbiface.CollectionAPI
}

type ProductValidator struct {
	validator *validator.Validate
}

func (p *ProductValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

func insertProducts(ctx context.Context, products []Product, collection dbiface.CollectionAPI) ([]interface{}, error) {
	var insertIds []interface{}

	for _, product := range products {
		product.ID = primitive.NewObjectID()
		insertId, err := collection.InsertOne(ctx, product)
		if err != nil {
			log.Printf("Unable to insert into db %+v\n", err)
		}
		insertIds = append(insertIds, insertId.InsertedID)
	}

	return insertIds, nil
}

func (h *ProductHandler) CreateProducts(c echo.Context) error {
	var products []Product
	c.Echo().Validator = &ProductValidator{validator: v}

	if err := c.Bind(&products); err != nil {
		log.Fatalf("Unable to bind: %+v\n", err)
		return c.JSON(http.StatusBadRequest, "Unable to bind")
	}

	for _, product := range products {
		if err := c.Validate(product); err != nil {
			log.Printf("Unable to validate product %+v   %+v\n", product, err)
			return err
		}
	}

	IDS, err := insertProducts(context.Background(), products, h.Col)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to insert into db")
	}

	return c.JSON(http.StatusCreated, IDS)
}
