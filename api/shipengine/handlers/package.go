package ShipengineHandlers

import (
	"context"
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	shipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func CreatePackage(ctx context.Context, p shipengineModels.Package) (*shipengineModels.PackageResponse, error) {
	packageResponse, errorResponse, err := ShipengineAPI.CreatePackage(ctx, p)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error creating package: %v", errorResponse.Errors[0].Message)
	}

	return packageResponse, nil
}
