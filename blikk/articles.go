package blikk

import (
	"context"
	"iter"
	"strconv"
)

const (
	articlesPath          = "v1/admin/articles"
	additionArticlesKind  = "addition"
	allowanceArticlesKind = "allowance"
	equipmentArticlesKind = "equipment"
	materialArticlesKind  = "material"
	timeArticlesKind      = "time"
)

type ArticlesService struct {
	client *Client
}

type ArticleType string

const (
	ArticleTypeAll       ArticleType = "All"
	ArticleTypeAllowance ArticleType = "Allowance"
	ArticleTypeMaterial  ArticleType = "Material"
	ArticleTypeEquipment ArticleType = "Equipment"
	ArticleTypeAddition  ArticleType = "Addition"
	ArticleTypeTime      ArticleType = "Time"
)

func (t *ArticleType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, ArticleTypeAll, ArticleTypeAllowance, ArticleTypeMaterial, ArticleTypeEquipment, ArticleTypeAddition, ArticleTypeTime)
}

type ArticleAdditionType string

const (
	ArticleAdditionTypeNumber  ArticleAdditionType = "Number"
	ArticleAdditionTypeBoolean ArticleAdditionType = "Boolean"
)

func (t *ArticleAdditionType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, ArticleAdditionTypeNumber, ArticleAdditionTypeBoolean)
}

const (
	ArticleSortByName          = "name"
	ArticleSortByArticleNumber = "articleNumber"
	ArticleSortByCreatedDate   = "createdDate"
	ArticleSortByUpdatedDate   = "updatedDate"
)

// ArticleListParams.Type defaults to ArticleTypeMaterial on the API side when left empty.
type ArticleListParams struct {
	ListOptions
	Type       ArticleType `query:"type"`
	Query      string      `query:"filter.query"`
	ExternalID string      `query:"filter.externalId"`
}

type ArticleListItem struct {
	ID            int       `json:"id"`
	ArticleNumber string    `json:"articleNumber"`
	ExternalID    string    `json:"externalId"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Cost          float64   `json:"cost"`
	Unit          *NameID   `json:"unit"`
	Description   string    `json:"description"`
	IsActive      bool      `json:"isActive"`
	CreatedBy     *NameID   `json:"createdBy"`
	UpdatedBy     *NameID   `json:"updatedBy"`
	Created       *DateTime `json:"created"`
	Updated       *DateTime `json:"updated"`
}

type Article struct {
	ID                        int                 `json:"id"`
	Name                      string              `json:"name"`
	Description               string              `json:"description"`
	Type                      ArticleType         `json:"type"`
	ArticleNumber             string              `json:"articleNumber"`
	Price                     float64             `json:"price"`
	Cost                      float64             `json:"cost"`
	Unit                      *NameID             `json:"unit"`
	IsActive                  bool                `json:"isActive"`
	ExternalID                string              `json:"externalId"`
	VatCode                   VatCode             `json:"vatCode"`
	SalaryCode                string              `json:"salaryCode"`
	SalaryCodeForHour         string              `json:"salaryCodeForHour"`
	ArticleGroup              string              `json:"articleGroup"`
	EAN                       string              `json:"ean"`
	IsFavorite                bool                `json:"isFavorite"`
	Manufacturer              string              `json:"manufacturer"`
	ManufacturerArticleNumber string              `json:"manufacturerArticleNumber"`
	AffectsCompTimeBalance    bool                `json:"affectsCompTimeBalance"`
	Percentage                float64             `json:"percentage"`
	SalaryCompensation        float64             `json:"salaryCompensation"`
	Start1                    string              `json:"start1"`
	End1                      string              `json:"end1"`
	Start2                    string              `json:"start2"`
	End2                      string              `json:"end2"`
	NormalWorkDays            bool                `json:"normalWorkDays"`
	WeekDays                  bool                `json:"weekDays"`
	MajorPublic               bool                `json:"majorPublic"`
	Supplier                  *NameID             `json:"supplier"`
	AdditionType              ArticleAdditionType `json:"additionType"`
	CostCenter                *CostCenterRef      `json:"costCenter"`
	Model                     string              `json:"model"`
	SerialNumber              string              `json:"serialNumber"`
	YearOfPurchase            int                 `json:"yearOfPurchase"`
	EquipmentType             string              `json:"equipmentType"`
	CreatedBy                 *NameID             `json:"createdBy"`
	UpdatedBy                 *NameID             `json:"updatedBy"`
	Created                   *DateTime           `json:"created"`
	Updated                   *DateTime           `json:"updated"`
}

type ArticleCommon struct {
	ID            int       `json:"id"`
	ArticleNumber string    `json:"articleNumber"`
	Cost          float64   `json:"cost"`
	Description   string    `json:"description"`
	ExternalID    string    `json:"externalId"`
	IsActive      bool      `json:"isActive"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Unit          NameID    `json:"unit"`
	VatCode       VatCode   `json:"vatCode"`
	CreatedBy     NameID    `json:"createdBy"`
	UpdatedBy     *NameID   `json:"updatedBy"`
	Created       DateTime  `json:"created"`
	Updated       *DateTime `json:"updated"`
}

type AdditionArticle struct {
	ArticleCommon
	AdditionType      ArticleAdditionType `json:"additionType"`
	SalaryCode        string              `json:"salaryCode"`
	SalaryCodeForHour string              `json:"salaryCodeForHour"`
}

type AllowanceArticle struct {
	ArticleCommon
	SalaryCode        string `json:"salaryCode"`
	SalaryCodeForHour string `json:"salaryCodeForHour"`
}

type EquipmentArticle struct {
	ArticleCommon
	EquipmentTypeID int    `json:"equipmentTypeId"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serialNumber"`
	YearOfPurchase  int    `json:"yearOfPurchase"`
}

type MaterialArticle struct {
	ArticleCommon
	ArticleGroupID            int    `json:"articleGroupId"`
	EAN                       string `json:"ean"`
	IsFavorite                bool   `json:"isFavorite"`
	ManufacturerArticleNumber string `json:"manufacturerArticleNumber"`
	Manufacturer              string `json:"manufacturer"`
	SupplierID                int    `json:"supplierId"`
}

type TimeArticle struct {
	ArticleCommon
	SalaryCode        string `json:"salaryCode"`
	SalaryCodeForHour string `json:"salaryCodeForHour"`
}

// ArticleInput.IsActive is a pointer because the API keeps the current value when it is omitted.
type ArticleInput struct {
	Name          string  `json:"name"`
	UnitID        int     `json:"unitId"`
	VatCodeID     int     `json:"vatCodeId"`
	ArticleNumber string  `json:"articleNumber,omitempty"`
	Cost          float64 `json:"cost,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Description   string  `json:"description,omitempty"`
	ExternalID    string  `json:"externalId,omitempty"`
	IsActive      *bool   `json:"isActive,omitempty"`
}

// Requests take the addition type by number while responses report it by name as ArticleAdditionType.
type AdditionTypeID int

const (
	AdditionTypeIDNumber  AdditionTypeID = 1
	AdditionTypeIDBoolean AdditionTypeID = 2
)

type AdditionArticleInput struct {
	ArticleInput
	AdditionType      AdditionTypeID `json:"additionType,omitempty"`
	SalaryCode        string         `json:"salaryCode,omitempty"`
	SalaryCodeForHour string         `json:"salaryCodeForHour,omitempty"`
}

type AllowanceArticleInput struct {
	ArticleInput
	SalaryCode        string `json:"salaryCode,omitempty"`
	SalaryCodeForHour string `json:"salaryCodeForHour,omitempty"`
}

type EquipmentArticleInput struct {
	ArticleInput
	EquipmentTypeID int    `json:"equipmentTypeId,omitempty"`
	Model           string `json:"model,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	YearOfPurchase  int    `json:"yearOfPurchase,omitempty"`
}

type MaterialArticleInput struct {
	ArticleInput
	ArticleGroupID            int    `json:"articleGroupId,omitempty"`
	EAN                       string `json:"ean,omitempty"`
	IsFavorite                bool   `json:"isFavorite"`
	Manufacturer              string `json:"manufacturer,omitempty"`
	ManufacturerArticleNumber string `json:"manufacturerArticleNumber,omitempty"`
	SupplierID                int    `json:"supplierId,omitempty"`
}

type TimeArticleInput struct {
	ArticleInput
	SalaryCode        string `json:"salaryCode,omitempty"`
	SalaryCodeForHour string `json:"salaryCodeForHour,omitempty"`
}

func (s *ArticlesService) List(ctx context.Context, params ArticleListParams) (*Page[ArticleListItem], error) {
	return listPage[ArticleListItem](ctx, s.client, articlesPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ArticlesService) ListAll(ctx context.Context, params ArticleListParams) iter.Seq2[ArticleListItem, error] {
	return listAll[ArticleListItem](ctx, s.client, articlesPath, params)
}

func (s *ArticlesService) Get(ctx context.Context, id int) (*Article, error) {
	return getResource[Article](ctx, s.client, articlesPath, strconv.Itoa(id))
}

func (s *ArticlesService) CreateAddition(ctx context.Context, input AdditionArticleInput) (*AdditionArticle, error) {
	return createResource[AdditionArticle](ctx, s.client, input, articlesPath, additionArticlesKind)
}

func (s *ArticlesService) CreateAllowance(ctx context.Context, input AllowanceArticleInput) (*AllowanceArticle, error) {
	return createResource[AllowanceArticle](ctx, s.client, input, articlesPath, allowanceArticlesKind)
}

func (s *ArticlesService) CreateEquipment(ctx context.Context, input EquipmentArticleInput) (*EquipmentArticle, error) {
	return createResource[EquipmentArticle](ctx, s.client, input, articlesPath, equipmentArticlesKind)
}

func (s *ArticlesService) CreateMaterial(ctx context.Context, input MaterialArticleInput) (*MaterialArticle, error) {
	return createResource[MaterialArticle](ctx, s.client, input, articlesPath, materialArticlesKind)
}

func (s *ArticlesService) CreateTime(ctx context.Context, input TimeArticleInput) (*TimeArticle, error) {
	return createResource[TimeArticle](ctx, s.client, input, articlesPath, timeArticlesKind)
}

func (s *ArticlesService) UpdateAddition(ctx context.Context, id int, input AdditionArticleInput) (*AdditionArticle, error) {
	return updateResourceReturning[AdditionArticle](ctx, s.client, input, articlesPath, additionArticlesKind, strconv.Itoa(id))
}

func (s *ArticlesService) UpdateAllowance(ctx context.Context, id int, input AllowanceArticleInput) (*AllowanceArticle, error) {
	return updateResourceReturning[AllowanceArticle](ctx, s.client, input, articlesPath, allowanceArticlesKind, strconv.Itoa(id))
}

func (s *ArticlesService) UpdateEquipment(ctx context.Context, id int, input EquipmentArticleInput) (*EquipmentArticle, error) {
	return updateResourceReturning[EquipmentArticle](ctx, s.client, input, articlesPath, equipmentArticlesKind, strconv.Itoa(id))
}

func (s *ArticlesService) UpdateMaterial(ctx context.Context, id int, input MaterialArticleInput) (*MaterialArticle, error) {
	return updateResourceReturning[MaterialArticle](ctx, s.client, input, articlesPath, materialArticlesKind, strconv.Itoa(id))
}

func (s *ArticlesService) UpdateTime(ctx context.Context, id int, input TimeArticleInput) (*TimeArticle, error) {
	return updateResourceReturning[TimeArticle](ctx, s.client, input, articlesPath, timeArticlesKind, strconv.Itoa(id))
}
