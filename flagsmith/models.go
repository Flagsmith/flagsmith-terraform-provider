package flagsmith

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"math/big"
	"sort"
)

type FeatureStateValue struct {
	Type         types.String `tfsdk:"type"`
	StringValue  types.String `tfsdk:"string_value"`
	IntegerValue types.Number `tfsdk:"integer_value"`
	BooleanValue types.Bool   `tfsdk:"boolean_value"`
}

func (f *FeatureStateValue) ToClientFSV() *flagsmithapi.FeatureStateValue {
	switch f.Type.Value {
	case "unicode":
		return &flagsmithapi.FeatureStateValue{
			Type:        "unicode",
			StringValue: &f.StringValue.Value,
		}
	case "int":
		intValue, _ := f.IntegerValue.Value.Int64()
		return &flagsmithapi.FeatureStateValue{
			Type:         "int",
			IntegerValue: &intValue,
		}
	case "bool":
		return &flagsmithapi.FeatureStateValue{
			Type:         "bool",
			BooleanValue: &f.BooleanValue.Value,
		}
	}
	return nil
}

func MakeFeatureStateValueFromClientFSV(clientFSV *flagsmithapi.FeatureStateValue) FeatureStateValue {
	fsvType := clientFSV.Type
	switch fsvType {
	case "unicode":
		return FeatureStateValue{
			Type:         types.String{Value: fsvType},
			StringValue:  types.String{Value: *clientFSV.StringValue},
			IntegerValue: types.Number{Null: true, Value: nil},
			BooleanValue: types.Bool{Null: true},
		}
	case "int":
		return FeatureStateValue{
			Type:         types.String{Value: fsvType},
			StringValue:  types.String{Null: true},
			IntegerValue: types.Number{Value: big.NewFloat(float64(*clientFSV.IntegerValue))},
			BooleanValue: types.Bool{Null: true},
		}
	case "bool":
		return FeatureStateValue{
			Type:         types.String{Value: fsvType},
			StringValue:  types.String{Null: true},
			IntegerValue: types.Number{Null: true},
			BooleanValue: types.Bool{Value: *clientFSV.BooleanValue},
		}

	}
	return FeatureStateValue{}
}

type FeatureStateResourceData struct {
	ID                types.Number       `tfsdk:"id"`
	Enabled           types.Bool         `tfsdk:"enabled"`
	FeatureStateValue *FeatureStateValue `tfsdk:"feature_state_value"`
	Feature           types.Number       `tfsdk:"feature"`
	Environment       types.Number       `tfsdk:"environment"`
	FeatureName       types.String       `tfsdk:"feature_name"`
	EnvironmentKey    types.String       `tfsdk:"environment_key"`
}

func (f *FeatureStateResourceData) ToClientFS(featureStateID int64, feature int64, environment int64) *flagsmithapi.FeatureState {
	return &flagsmithapi.FeatureState{
		ID:                featureStateID,
		Enabled:           f.Enabled.Value,
		FeatureStateValue: f.FeatureStateValue.ToClientFSV(),
		Feature:           feature,
		Environment:       environment,
	}
}

// Generate a new FeatureStateResourceData from client `FeatureState`
func MakeFeatureStateResourceDataFromClientFS(clientFS *flagsmithapi.FeatureState) FeatureStateResourceData {
	fsValue := MakeFeatureStateValueFromClientFSV(clientFS.FeatureStateValue)
	return FeatureStateResourceData{
		ID:                types.Number{Value: big.NewFloat(float64(clientFS.ID))},
		Enabled:           types.Bool{Value: clientFS.Enabled},
		FeatureStateValue: &fsValue,
		Feature:           types.Number{Value: big.NewFloat(float64(clientFS.Feature))},
		Environment:       types.Number{Value: big.NewFloat(float64(clientFS.Environment))},
	}
}

type MultivariateOption struct {
	Type                        types.String `tfsdk:"type"`
	ID                          types.Number `tfsdk:"id"`
	IntegerValue                types.Number `tfsdk:"integer_value"`
	StringValue                 types.String `tfsdk:"string_value"`
	BooleanValue                types.Bool   `tfsdk:"boolean_value"`
	DefaultPercentageAllocation types.Number `tfsdk:"default_percentage_allocation"`
}

func (m *MultivariateOption) ToClientMultivariateOption() *flagsmithapi.MultivariateOption {
	defaultPercentageAllocation, _ := m.DefaultPercentageAllocation.Value.Float64()
	stringValue := m.StringValue.Value
	booleanValue := m.BooleanValue.Value

	mo := flagsmithapi.MultivariateOption{
		Type:                        m.Type.Value,
		DefaultPercentageAllocation: defaultPercentageAllocation,
	}
	if m.ID.Value != nil {
		moID, _ := m.ID.Value.Int64()
		mo.ID = &moID
	}
	if m.IntegerValue.Value != nil {
		moIntegerValue, _ := m.IntegerValue.Value.Int64()
		mo.IntegerValue = &moIntegerValue
	} else if m.StringValue.Null == false {
		mo.StringValue = &stringValue
	} else if m.BooleanValue.Null == false {
		mo.BooleanValue = &booleanValue
	}

	return &mo
}

type FeatureResourceData struct {
	UUID                types.String          `tfsdk:"uuid"`
	ID                  types.Number          `tfsdk:"id"`
	Name                types.String          `tfsdk:"feature_name"`
	Type                types.String          `tfsdk:"type"`
	Description         types.String          `tfsdk:"description"`
	InitialValue        types.String          `tfsdk:"initial_value"`
	DefaultEnabled      types.Bool            `tfsdk:"default_enabled"`
	IsArchived          types.Bool            `tfsdk:"is_archived"`
	Owners              *[]types.Number       `tfsdk:"owners"`
	MultivariateOptions *[]MultivariateOption `tfsdk:"multivariate_options"`
	ProjectID           types.Number          `tfsdk:"project_id"`
	ProjectUUID         types.String          `tfsdk:"project_uuid"`
}

func (f *FeatureResourceData) ToClientFeature() *flagsmithapi.Feature {
	feature := flagsmithapi.Feature{
		UUID:           f.UUID.Value,
		Name:           f.Name.Value,
		Type:           &f.Type.Value,
		Description:    &f.Description.Value,
		InitialValue:   f.InitialValue.Value,
		DefaultEnabled: f.DefaultEnabled.Value,
		IsArchived:     f.IsArchived.Value,
		ProjectUUID:    f.ProjectUUID.Value,
	}
	if f.ID.Value != nil {
		featureID, _ := f.ID.Value.Int64()
		feature.ID = &featureID
	}
	if f.ProjectID.Value != nil {
		projectID, _ := f.ProjectID.Value.Int64()
		feature.ProjectID = &projectID
	}

	if f.Owners != nil {
		for _, owner := range *f.Owners {
			ownerID, _ := owner.Value.Int64()
			*feature.Owners = append(*feature.Owners, ownerID)

		}
	}
	if f.MultivariateOptions != nil {
		mvOptions := []flagsmithapi.MultivariateOption{}
		for _, mo := range *f.MultivariateOptions {
			mvOptions = append(mvOptions, *mo.ToClientMultivariateOption())
		}
		feature.MultivariateOptions = &mvOptions
	}
	return &feature

}

func MakeFeatureResourceDataFromClientFeature(clientFeature *flagsmithapi.Feature) FeatureResourceData {
	var multivariateOptions []MultivariateOption
	// Sort the multivariate options by ID to ensure state stays consistent
	sort.Slice(*clientFeature.MultivariateOptions, func(i, j int) bool {
		iID := (*clientFeature.MultivariateOptions)[i].ID
		jID := (*clientFeature.MultivariateOptions)[j].ID
		return *iID < *jID
	})
	for _, option := range *clientFeature.MultivariateOptions {
		mvOption := MultivariateOption{
			Type: types.String{Value: option.Type},
			ID:   types.Number{Value: big.NewFloat(float64(*option.ID))},
			IntegerValue:                types.Number{Null: true},
			StringValue:                 types.String{Null: true},
			BooleanValue: 	      types.Bool{Null: true},
			DefaultPercentageAllocation: types.Number{Value: big.NewFloat(option.DefaultPercentageAllocation)},
		}
		if option.IntegerValue != nil {
			mvOption.IntegerValue = types.Number{Value: big.NewFloat(float64(*option.IntegerValue))}
		}

		if option.StringValue != nil {
			mvOption.StringValue = types.String{Value: *option.StringValue}
		}
		if option.BooleanValue != nil {
			mvOption.BooleanValue = types.Bool{Value: *option.BooleanValue}
		}

		multivariateOptions = append(multivariateOptions, mvOption)

	}
	resourceData := FeatureResourceData{
		UUID:                types.String{Value: clientFeature.UUID},
		ID:                  types.Number{Value: big.NewFloat(float64(*clientFeature.ID))},
		Name:                types.String{Value: clientFeature.Name},
		Type:                types.String{Value: *clientFeature.Type},
		DefaultEnabled:      types.Bool{Value: clientFeature.DefaultEnabled},
		IsArchived:          types.Bool{Value: clientFeature.IsArchived},
		InitialValue:        types.String{Value: clientFeature.InitialValue},
		MultivariateOptions: &multivariateOptions,
		ProjectID:           types.Number{Value: big.NewFloat(float64(*clientFeature.ProjectID))},
		ProjectUUID:         types.String{Value: clientFeature.ProjectUUID},
	}
	if clientFeature.Description != nil {
		resourceData.Description = types.String{Value: *clientFeature.Description}
	}
	if clientFeature.Owners != nil {
		for _, owner := range *clientFeature.Owners {
			*resourceData.Owners = append(*resourceData.Owners, types.Number{Value: big.NewFloat(float64(owner))})
		}
	}
	return resourceData
}
