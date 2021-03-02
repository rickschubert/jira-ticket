package cucumberfeatureparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFeatures(t *testing.T) {
	features := getFeatures()
	for _, feature := range features {
		assert.GreaterOrEqual(t, len(feature.Name), 1)
		assert.GreaterOrEqual(t, len(feature.Bucket), 1)
	}
}

func TestGetFeatureInfoFromTextMentioningFeatureError(t *testing.T) {
	feat, err := GetFeatureInfoOfTextMentioningFeature("this text doesn'ta\n\ncontain a \nvalid feature")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, feat, feature{
		Name:   "",
		Bucket: "",
	})
}

func TestGetFeatureInfoFromTextMentioningFeatureSuccess(t *testing.T) {
	feat, err := GetFeatureInfoOfTextMentioningFeature("this text does\n\ncontain a featureName three\n\ninbetween")
	assert.Equal(t, err, nil)
	assert.Equal(t, feat, feature{
		Name:   "three",
		Bucket: "bucket_two",
	})
}
