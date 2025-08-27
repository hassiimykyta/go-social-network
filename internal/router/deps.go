package router

import appdb "go-rest-chi/internal/db"

type Deps struct {
	DB *appdb.SQL
}
