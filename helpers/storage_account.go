package helpers

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	aztables "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

// Connect to Azure Tables and retrieve all entities
func GetEntitiesFromStorage(tablesClient *aztables.Client) []Entity {
	var entities []Entity
	log.Println("Getting entities from storage")
	// Create query to retrieve all entities
	listPager := tablesClient.NewListEntitiesPager(nil)
	for listPager.More() {
		response, err := listPager.NextPage(context.Background())
		if err != nil {
			log.Printf("error getting next page: %v", err)
		}
		for _, entityBytes := range response.Entities {
			var entity Entity
			err := json.Unmarshal(entityBytes, &entity)
			if err != nil {
				if strings.Contains(err.Error(), "Entity.Reprocess") {
					// Ignore this specific error, as it is expected when the entity has the Reprocess property
				} else {
					log.Printf("error unmarshalling entity: %v", err)
				}
			}

			entities = append(entities, entity)

		}
	}
	return entities
}

func MarkEntityForReprocessing(tablesClient *aztables.Client, entity Entity, reprocess bool) {
	// If uploading to github fails, mark entity for reprocessing on next run.
	// This will be done by updating the entity in Azure Tables with a new property
	// called "reprocess" set to true.
	entity.Reprocess = reprocess

	entityBytes, err := json.Marshal(entity)
	if err != nil {
		log.Printf("error marshalling entity: %v", err)
	}

	updateOptions := aztables.UpdateEntityOptions{
		UpdateMode: aztables.UpdateModeReplace,
	}

	tablesClient.UpdateEntity(context.Background(), entityBytes, &updateOptions)

	log.Printf("Entity marked for reprocessing: ", entity)
}
