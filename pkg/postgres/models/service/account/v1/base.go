package servicemodelv1

import (
	sharedmodel "github.com/vyolayer/vyolayer/pkg/postgres/models/shared"
)

type (
	UUID       = sharedmodel.UUID
	TimeStamps = sharedmodel.TimeStampsWithSoftDelete
)
