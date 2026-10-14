package blikk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticlesListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "articles_list.json"))

	_, err := client.Articles.List(context.Background(), ArticleListParams{
		ListOptions: ListOptions{SortBy: ArticleSortByArticleNumber, SortOrder: SortAscending},
		Type:        ArticleTypeAll,
		Query:       "lor",
		ExternalID:  "EXT-001",
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/admin/articles", request.Path)
	assert.Equal(t, map[string][]string{
		"page":              {"1"},
		"sortBy":            {"articleNumber"},
		"sortOrder":         {"ascending"},
		"type":              {"All"},
		"filter.query":      {"lor"},
		"filter.externalId": {"EXT-001"},
	}, request.Query)
}

func TestArticlesListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/admin/articles", "articles_list.json"))

	page, err := client.Articles.List(context.Background(), ArticleListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 2)
	article := page.Items[0]
	assert.Equal(t, 1, article.ID)
	assert.Equal(t, "A024", article.ArticleNumber)
	assert.Equal(t, "EXT-001", article.ExternalID)
	assert.Equal(t, "Lorem ipsum", article.Name)
	assert.Equal(t, 500.0, article.Price)
	assert.Equal(t, 345.5, article.Cost)
	assert.Equal(t, &NameID{ID: 1, Name: "st"}, article.Unit)
	assert.Equal(t, "Lorem ipsum dolar sit amet", article.Description)
	assert.True(t, article.IsActive)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, article.CreatedBy)
	assert.Nil(t, article.UpdatedBy)
	assert.Equal(t, "2026-10-03 09:40:57", article.Created.String())
	assert.Nil(t, article.Updated)
	assert.Empty(t, page.Items[1].ExternalID)
	assert.Zero(t, page.Items[1].Cost)
}

func TestArticlesListAllWalksPages(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "articles_list.json"))

	var ids []int
	for article, err := range client.Articles.ListAll(context.Background(), ArticleListParams{Type: ArticleTypeTime}) {
		require.NoError(t, err)
		ids = append(ids, article.ID)
	}

	assert.Equal(t, []int{1, 2}, ids)
}

func TestArticlesGetDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/admin/articles/42", "articles_detail.json"))

	article, err := client.Articles.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, article.ID)
	assert.Equal(t, "Article 1", article.Name)
	assert.Equal(t, "Description of article 1", article.Description)
	assert.Equal(t, ArticleTypeMaterial, article.Type)
	assert.Equal(t, "ABC-123", article.ArticleNumber)
	assert.Equal(t, 200.0, article.Price)
	assert.Equal(t, 100.0, article.Cost)
	assert.Equal(t, &NameID{ID: 1, Name: "Unit 1"}, article.Unit)
	assert.True(t, article.IsActive)
	assert.Equal(t, "EXT-123", article.ExternalID)
	assert.Equal(t, VatCode{Name: "VAT 1", Code: "1", VAT: 25}, article.VatCode)
	assert.Equal(t, "123", article.SalaryCode)
	assert.Equal(t, "555", article.SalaryCodeForHour)
	assert.Equal(t, "123", article.ArticleGroup)
	assert.Equal(t, "1234567890123", article.EAN)
	assert.True(t, article.IsFavorite)
	assert.Equal(t, "Toyota", article.Manufacturer)
	assert.Equal(t, "123-123", article.ManufacturerArticleNumber)
	assert.False(t, article.AffectsCompTimeBalance)
	assert.Equal(t, 5.0, article.Percentage)
	assert.Equal(t, 100.0, article.SalaryCompensation)
	assert.Empty(t, article.Start1)
	assert.Empty(t, article.End2)
	assert.True(t, article.NormalWorkDays)
	assert.True(t, article.WeekDays)
	assert.False(t, article.MajorPublic)
	assert.Equal(t, &NameID{ID: 796, Name: "A & Management AB"}, article.Supplier)
	assert.Empty(t, article.AdditionType)
	assert.Equal(t, &CostCenterRef{ID: 1, Name: "Cost center 1", Code: "CC1"}, article.CostCenter)
	assert.Equal(t, "Model 1", article.Model)
	assert.Equal(t, "1234567890123", article.SerialNumber)
	assert.Equal(t, 2020, article.YearOfPurchase)
	assert.Empty(t, article.EquipmentType)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, article.CreatedBy)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, article.UpdatedBy)
	assert.Equal(t, "2026-10-09 09:40:57", article.Created.String())
	assert.Equal(t, "2026-10-11 09:40:57", article.Updated.String())
}

func TestArticleEnumsDecodeCaseInsensitively(t *testing.T) {
	var types []ArticleType
	require.NoError(t, json.Unmarshal([]byte(`["material","Time","ADDITION","other"]`), &types))
	var additionTypes []ArticleAdditionType
	require.NoError(t, json.Unmarshal([]byte(`["number","Boolean",""]`), &additionTypes))

	assert.Equal(t, []ArticleType{ArticleTypeMaterial, ArticleTypeTime, ArticleTypeAddition, "other"}, types)
	assert.Equal(t, []ArticleAdditionType{ArticleAdditionTypeNumber, ArticleAdditionTypeBoolean, ""}, additionTypes)
}

func TestArticlesCreateAdditionSendsFullInputAndDecodesResponse(t *testing.T) {
	client, handler := newRecordingClient(t, respondWithArticleFixture(t, http.StatusCreated, "articles_addition.json"))

	created, err := client.Articles.CreateAddition(context.Background(), AdditionArticleInput{
		ArticleInput: ArticleInput{
			Name:          "New Addition 1",
			UnitID:        817,
			VatCodeID:     646,
			ArticleNumber: "234",
			Cost:          401,
			Price:         1101,
			Description:   "Description of new addition 1",
			ExternalID:    "EXT-123",
			IsActive:      Ptr(true),
		},
		AdditionType:      1,
		SalaryCode:        "1234",
		SalaryCodeForHour: "1234",
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/admin/articles/addition", request.Path)
	assert.JSONEq(t, `{
		"isActive": true,
		"additionType": 1,
		"salaryCode": "1234",
		"salaryCodeForHour": "1234",
		"articleNumber": "234",
		"cost": 401,
		"description": "Description of new addition 1",
		"externalId": "EXT-123",
		"name": "New Addition 1",
		"price": 1101,
		"unitId": 817,
		"vatCodeId": 646
	}`, string(request.Body))

	assert.Equal(t, 42, created.ID)
	assert.Equal(t, ArticleAdditionTypeNumber, created.AdditionType)
	assert.Equal(t, "1234", created.SalaryCode)
	assert.Equal(t, "234", created.ArticleNumber)
	assert.Equal(t, 401.0, created.Cost)
	assert.Equal(t, "EXT-123", created.ExternalID)
	assert.True(t, created.IsActive)
	assert.Equal(t, "New Addition 1", created.Name)
	assert.Equal(t, 1101.0, created.Price)
	assert.Equal(t, NameID{ID: 123, Name: "Unit 1"}, created.Unit)
	assert.Equal(t, 646, created.VatCode.ID)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, created.VatCode.CreatedBy)
	assert.Equal(t, NameID{ID: 123, Name: "John Doe"}, created.CreatedBy)
	assert.Nil(t, created.UpdatedBy)
	assert.Equal(t, "2026-10-09 09:40:57", created.Created.String())
	assert.Nil(t, created.Updated)
}

func TestArticleInputMinimalSendsOnlyRequiredFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWithArticleFixture(t, http.StatusCreated, "articles_addition.json"))

	_, err := client.Articles.CreateAddition(context.Background(), AdditionArticleInput{
		ArticleInput: ArticleInput{Name: "Minimal", UnitID: 1, VatCodeID: 2, IsActive: Ptr(false)},
	})

	require.NoError(t, err)
	assert.JSONEq(t, `{"name": "Minimal", "unitId": 1, "vatCodeId": 2, "isActive": false}`, string(handler.last().Body))
}

func TestArticlesCreatePerKindSendsKindPathAndKeys(t *testing.T) {
	minimal := ArticleInput{Name: "Article", UnitID: 1, VatCodeID: 2}
	tests := []struct {
		name    string
		fixture string
		create  func(*ArticlesService) error
		path    string
		body    string
	}{
		{
			name:    "allowance",
			fixture: "articles_addition.json",
			create: func(s *ArticlesService) error {
				_, err := s.CreateAllowance(context.Background(), AllowanceArticleInput{ArticleInput: minimal, SalaryCode: "S", SalaryCodeForHour: "SH"})
				return err
			},
			path: "/v1/admin/articles/allowance",
			body: `{"name": "Article", "unitId": 1, "vatCodeId": 2, "salaryCode": "S", "salaryCodeForHour": "SH"}`,
		},
		{
			name:    "equipment",
			fixture: "articles_equipment.json",
			create: func(s *ArticlesService) error {
				_, err := s.CreateEquipment(context.Background(), EquipmentArticleInput{ArticleInput: minimal, EquipmentTypeID: 33, Model: "Model 23", SerialNumber: "SN-123", YearOfPurchase: 2019})
				return err
			},
			path: "/v1/admin/articles/equipment",
			body: `{"name": "Article", "unitId": 1, "vatCodeId": 2, "equipmentTypeId": 33, "model": "Model 23", "serialNumber": "SN-123", "yearOfPurchase": 2019}`,
		},
		{
			name:    "material",
			fixture: "articles_material.json",
			create: func(s *ArticlesService) error {
				_, err := s.CreateMaterial(context.Background(), MaterialArticleInput{ArticleInput: minimal, ArticleGroupID: 123, EAN: "1234567890123", IsFavorite: true, Manufacturer: "Toyota", ManufacturerArticleNumber: "123-123", SupplierID: 5})
				return err
			},
			path: "/v1/admin/articles/material",
			body: `{"name": "Article", "unitId": 1, "vatCodeId": 2, "articleGroupId": 123, "ean": "1234567890123", "isFavorite": true, "manufacturer": "Toyota", "manufacturerArticleNumber": "123-123", "supplierId": 5}`,
		},
		{
			name:    "material minimal keeps required isFavorite",
			fixture: "articles_material.json",
			create: func(s *ArticlesService) error {
				_, err := s.CreateMaterial(context.Background(), MaterialArticleInput{ArticleInput: minimal})
				return err
			},
			path: "/v1/admin/articles/material",
			body: `{"name": "Article", "unitId": 1, "vatCodeId": 2, "isFavorite": false}`,
		},
		{
			name:    "time",
			fixture: "articles_addition.json",
			create: func(s *ArticlesService) error {
				_, err := s.CreateTime(context.Background(), TimeArticleInput{ArticleInput: minimal, SalaryCode: "S", SalaryCodeForHour: "SH"})
				return err
			},
			path: "/v1/admin/articles/time",
			body: `{"name": "Article", "unitId": 1, "vatCodeId": 2, "salaryCode": "S", "salaryCodeForHour": "SH"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, handler := newRecordingClient(t, respondWithArticleFixture(t, http.StatusCreated, tc.fixture))

			require.NoError(t, tc.create(client.Articles))

			request := handler.last()
			assert.Equal(t, http.MethodPost, request.Method)
			assert.Equal(t, tc.path, request.Path)
			assert.JSONEq(t, tc.body, string(request.Body))
		})
	}
}

func TestArticlesUpdatePerKindSendsPutToKindPath(t *testing.T) {
	minimal := ArticleInput{Name: "Article", UnitID: 1, VatCodeID: 2}
	tests := []struct {
		name   string
		update func(*ArticlesService) (int, error)
		path   string
	}{
		{"addition", func(s *ArticlesService) (int, error) {
			updated, err := s.UpdateAddition(context.Background(), 42, AdditionArticleInput{ArticleInput: minimal, AdditionType: 2})
			return articleIDOf(updated, err, func(a *AdditionArticle) int { return a.ID })
		}, "/v1/admin/articles/addition/42"},
		{"allowance", func(s *ArticlesService) (int, error) {
			updated, err := s.UpdateAllowance(context.Background(), 42, AllowanceArticleInput{ArticleInput: minimal})
			return articleIDOf(updated, err, func(a *AllowanceArticle) int { return a.ID })
		}, "/v1/admin/articles/allowance/42"},
		{"equipment", func(s *ArticlesService) (int, error) {
			updated, err := s.UpdateEquipment(context.Background(), 42, EquipmentArticleInput{ArticleInput: minimal})
			return articleIDOf(updated, err, func(a *EquipmentArticle) int { return a.ID })
		}, "/v1/admin/articles/equipment/42"},
		{"material", func(s *ArticlesService) (int, error) {
			updated, err := s.UpdateMaterial(context.Background(), 42, MaterialArticleInput{ArticleInput: minimal})
			return articleIDOf(updated, err, func(a *MaterialArticle) int { return a.ID })
		}, "/v1/admin/articles/material/42"},
		{"time", func(s *ArticlesService) (int, error) {
			updated, err := s.UpdateTime(context.Background(), 42, TimeArticleInput{ArticleInput: minimal})
			return articleIDOf(updated, err, func(a *TimeArticle) int { return a.ID })
		}, "/v1/admin/articles/time/42"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, handler := newRecordingClient(t, respondWithArticleFixture(t, http.StatusOK, "articles_addition.json"))

			id, err := tc.update(client.Articles)

			require.NoError(t, err)
			assert.Equal(t, 42, id)
			request := handler.last()
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, tc.path, request.Path)
			assert.Equal(t, "Article", request.jsonBody(t)["name"])
		})
	}
}

func TestArticlesDecodeKindSpecificResponses(t *testing.T) {
	client := newTestClient(t, routes(t, map[string]http.HandlerFunc{
		"PUT /v1/admin/articles/equipment/42": serveFixture(t, "articles_equipment.json"),
		"PUT /v1/admin/articles/material/42":  serveFixture(t, "articles_material.json"),
	}))

	equipment, err := client.Articles.UpdateEquipment(context.Background(), 42, EquipmentArticleInput{})
	require.NoError(t, err)
	material, err := client.Articles.UpdateMaterial(context.Background(), 42, MaterialArticleInput{})
	require.NoError(t, err)

	assert.Equal(t, 33, equipment.EquipmentTypeID)
	assert.Equal(t, "Model 23", equipment.Model)
	assert.Equal(t, "SN-123", equipment.SerialNumber)
	assert.Equal(t, 2019, equipment.YearOfPurchase)
	assert.Equal(t, NameID{ID: 817, Name: "Unit 1"}, equipment.Unit)
	assert.Equal(t, 56, material.ArticleGroupID)
	assert.Equal(t, "1234567890123", material.EAN)
	assert.True(t, material.IsFavorite)
	assert.Equal(t, "Toyota", material.Manufacturer)
	assert.Equal(t, "123-123", material.ManufacturerArticleNumber)
	assert.Equal(t, 123, material.SupplierID)
	assert.Equal(t, "2026-10-11 09:40:57", material.Updated.String())
}

func respondWithArticleFixture(t *testing.T, statusCode int, name string) http.HandlerFunc {
	fixture := loadFixture(t, name)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(fixture)
	}
}

func articleIDOf[T any](resource *T, err error, id func(*T) int) (int, error) {
	if err != nil {
		return 0, err
	}
	return id(resource), nil
}
