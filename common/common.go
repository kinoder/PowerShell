package common

import (
	models "hamkaran_system/bootcamp/final/project/model"
)

var LoginUser = &models.User{Username: "", Password: ""}

var LogHistory = make([]models.LogHistory, 0)
