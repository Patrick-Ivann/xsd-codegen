package xmlgen

import (
	"strings"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/xmlgen/mocks"
	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestGenerateElement(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")

	schema := &model.XSDSchema{
		Elements: []model.XSDElement{
			{Name: "test", Type: "xsd:string"},
		},
	}
	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	if elem.Tag != "test" {
		t.Errorf("expected tag 'test', got %s", elem.Tag)
	}
}

func TestGenerateDocument(t *testing.T) {

	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")

	schema := &model.XSDSchema{
		Elements: []model.XSDElement{
			{
				Name: "root",
				ComplexType: &model.XSDComplexType{
					Sequence: &model.XSDSequence{
						Elements: []model.XSDElement{
							{Name: "child", Type: "xsd:string"},
						},
					},
				},
			},
		},
	}

	doc := etree.NewDocument()
	root := GenerateElement(schema, schema.Elements[0], mockGen)
	doc.SetRoot(root)

	out, err := doc.WriteToString()
	if err != nil || !strings.Contains(out, "<child>") {
		t.Error("Expected output with <child> tag")
	}
}

func TestGenerateElementWhenElementHasTnsComplexType(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")
	mockGen.On("Generate", "xs:int", (*model.XSDRestriction)(nil)).Return(31)

	schema := &model.XSDSchema{
		ComplexTypes: []model.XSDComplexType{
			{
				Name: "PersonType",
				Sequence: &model.XSDSequence{
					Elements: []model.XSDElement{
						{Name: "FirstName", Type: "xs:string"},
						{Name: "LastName", Type: "xs:string"},
					},
				},
				Attrs: []model.XSDAttribute{
					{Name: "id", Type: "xs:int", Fixed: "69"},
				},
			},
		},
		Elements: []model.XSDElement{
			{Name: "Person", Type: "tns:PersonType"},
		},
	}

	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	assert.Equal(t, "Person", elem.Tag)
	assert.NotNil(t, elem.SelectElement("FirstName"))
	assert.NotNil(t, elem.SelectElement("LastName"))
	assert.NotEmpty(t, elem.SelectAttrValue("id", ""))
}

func TestGenerateElementWhenElementHasTnsSimpleType(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:int", (*model.XSDRestriction)(nil)).Return(69)

	schema := &model.XSDSchema{
		SimpleTypes: []model.XSDSimpleType{
			{
				Name: "AgeType",
				Restriction: &model.XSDRestriction{
					Base: "xs:int",
				},
			},
		},
		Elements: []model.XSDElement{
			{Name: "Age", Type: "tns:AgeType"},
		},
	}

	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	assert.Equal(t, "Age", elem.Tag)
	assert.NotEmpty(t, elem.Text())
}

func TestGenerateElementWhenElementHasPrimitiveType(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")

	schema := &model.XSDSchema{
		Elements: []model.XSDElement{
			{Name: "Title", Type: "xs:string"},
		},
	}

	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	assert.Equal(t, "Title", elem.Tag)
	assert.NotEmpty(t, elem.Text())
}

func TestGenerateElementWhenElementHasInlineComplexType(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")

	schema := &model.XSDSchema{
		Elements: []model.XSDElement{
			{
				Name: "Book",
				ComplexType: &model.XSDComplexType{
					Sequence: &model.XSDSequence{
						Elements: []model.XSDElement{
							{Name: "Title", Type: "xs:string"},
							{Name: "Author", Type: "xs:string"},
						},
					},
				},
			},
		},
	}

	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	assert.Equal(t, "Book", elem.Tag)
	assert.NotNil(t, elem.SelectElement("Title"))
	assert.NotNil(t, elem.SelectElement("Author"))
}

func TestGenerateElementWhenElementHasInlineSimpleType(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:int", (*model.XSDRestriction)(nil)).Return(93)

	schema := &model.XSDSchema{
		Elements: []model.XSDElement{
			{
				Name: "Rating",
				SimpleType: &model.XSDSimpleType{
					Restriction: &model.XSDRestriction{
						Base: "xs:int",
					},
				},
			},
		},
	}

	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	assert.Equal(t, "Rating", elem.Tag)
	assert.NotEmpty(t, elem.Text())
}

func TestGenerateElementWhenElementIsAReference(t *testing.T) {

	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")

	schema := &model.XSDSchema{
		Elements: []model.XSDElement{
			{
				Name: "Original",
				Type: "xs:string",
			},
			{
				Name: "Alias",
				Ref:  "tns:Original",
			},
		},
	}

	elem := GenerateElement(schema, schema.Elements[1], mockGen)
	assert.Equal(t, "Original", elem.Tag)
	assert.NotEmpty(t, elem.Text())
}

func TestGenerateElementWhenComplexTypeHasChoice(t *testing.T) {
	mockGen := new(mocks.MockValueGenerator)
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")
	mockGen.On("Generate", "xs:string", (*model.XSDRestriction)(nil)).Return("MockTitle")

	schema := &model.XSDSchema{
		ComplexTypes: []model.XSDComplexType{
			{
				Name: "ContactType",
				Choice: &model.XSDChoice{
					Elements: []model.XSDElement{
						{Name: "Email", Type: "xs:string"},
						{Name: "Phone", Type: "xs:string"},
					},
				},
			},
		},
		Elements: []model.XSDElement{
			{
				Name: "Contact",
				Type: "tns:ContactType",
			},
		},
	}

	elem := GenerateElement(schema, schema.Elements[0], mockGen)
	assert.Equal(t, "Contact", elem.Tag)
	children := elem.ChildElements()
	assert.Len(t, children, 1)

	child := children[0]
	assert.Contains(t, []string{"Email", "Phone"}, child.Tag)
	assert.NotEmpty(t, child.Text())
}
