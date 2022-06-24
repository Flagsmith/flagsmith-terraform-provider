package flagsmith

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"math/big"
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
	}
	// TODO: Implement other types
	panic("unsupported FeatureStateValue type")
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
	}
	// TODO: Implement other types
	panic("unsupported FeatureStateValue type")
	return FeatureStateValue{}
	// if clientResponse.Type == "unicode"{

	// }
	// return FeatureStateValue{
	// 	Type: types.String(clientResponse.Type),
	// 	StringValue: types.String(clientResponse.StringValue),
	// 	IntegerValue: types.Number(clientResponse.IntegerValue),
	// 	BooleanValue: types.Bool(clientResponse.BooleanValue),
	// }
}

type June struct {
	ID types.Number `tfsdk:"id"`
}

type flagResourceData struct {
	ID                types.Number       `tfsdk:"id"`
	Enabled           types.Bool         `tfsdk:"enabled"`
	FeatureStateValue *FeatureStateValue `tfsdk:"feature_state_value"`
	Feature           types.Number       `tfsdk:"feature"`
	Environment       types.Number       `tfsdk:"environment"`
	FeatureName       types.String       `tfsdk:"feature_name"`
	EnvironmentKey    types.String       `tfsdk:"environment_key"`
}

func (f *flagResourceData) ToClientFS(featureStateID int64) *flagsmithapi.FeatureState {
	intFeature, _ := f.Feature.Value.Int64()
	intEnvironment, _ := f.Environment.Value.Int64()
	return &flagsmithapi.FeatureState{
		ID:                featureStateID,
		Enabled:           f.Enabled.Value,
		FeatureStateValue: f.FeatureStateValue.ToClientFSV(),
		Feature:           intFeature,
		Environment:       intEnvironment,
	}
}

// Generate a new flagResourceData from a client FeatureState
func MakeFlagResourceDataFromClientFS(clientFS *flagsmithapi.FeatureState) flagResourceData {
	fsValue := MakeFeatureStateValueFromClientFSV(clientFS.FeatureStateValue)
	return flagResourceData{
		ID:                types.Number{Value: big.NewFloat(float64(clientFS.ID))},
		Enabled:           types.Bool{Value: clientFS.Enabled},
		FeatureStateValue: &fsValue,
		Feature:           types.Number{Value: big.NewFloat(float64(clientFS.Feature))},
		Environment:       types.Number{Value: big.NewFloat(float64(clientFS.Environment))},
	}
}
