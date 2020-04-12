package cucumberfeatureparser

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/mitchellh/go-homedir"
	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/utils"
)

type feature struct {
	Name   string
	Bucket string
}

func filterForFilesWithFeatureExtension(filePaths []string) []string {
	var featureFiles []string
	for _, filePath := range filePaths {
		matched, _ := regexp.Match(`\.feature$`, []byte(filePath))
		if matched {
			featureFiles = append(featureFiles, filePath)
		}
	}
	return featureFiles
}

func getPathsToAllFeatureFiles() []string {
	settings := constants.GetSettings()
	if settings.FeatureFolderPath == "" {
		log.Fatal(fmt.Sprintf("It looks like your settings file, located at %s, doesn't contain the necessary setting which would point to the feature folder of the cucumber journeys", constants.SETTINGS_FILE_RELATIVE_PATH))
	}
	featuresFolderAbsolute, err := homedir.Expand(settings.FeatureFolderPath)
	utils.HandleErrorStrictly(err)
	allFilesInFeaturesFolderDeep := utils.DeepGetAllFilesInFolder(featuresFolderAbsolute)
	return filterForFilesWithFeatureExtension(allFilesInFeaturesFolderDeep)
}

func getFeatures() []feature {
	var features []feature
	allFeaturePaths := getPathsToAllFeatureFiles()
	for _, featurePath := range allFeaturePaths {
		re := regexp.MustCompile(`.+/(.+?)/(.+)(?:\.feature)`)
		results := re.FindAllStringSubmatch(featurePath, -1)
		if len(results[0]) < 3 {
			log.Fatal(fmt.Sprintf("Unable to apply regex to get feature file from path %s", featurePath))
		}
		features = append(features, feature{
			Name:   results[0][2],
			Bucket: results[0][1],
		})
	}
	return features
}

func GetFeatureInfoOfTextMentioningFeature(text string) (feature, error) {
	features := getFeatures()

	var featureToReturn feature
	var errToReturn error

	for _, feature := range features {
		matched, _ := regexp.Match(feature.Name, []byte(text))
		if matched {
			featureToReturn = feature
			break
		}
	}

	if featureToReturn.Name == "" {
		errToReturn = errors.New("We couldn't find any feature in this text, sorry.")
	}

	return featureToReturn, errToReturn
}
