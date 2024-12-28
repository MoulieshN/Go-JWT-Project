package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MoulieshN/Go-JWT-Project.git/config"
	"github.com/MoulieshN/Go-JWT-Project.git/repository"
)

func Init(logCtx context.Context, port string) {
	config := config.GetConfig()

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=%v",
			config.MySQL.Username,
			config.MySQL.Password,
			config.MySQL.Hostname,
			config.MySQL.Port,
			config.MySQL.DBName,
			config.MySQL.ParseTime,
		),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	if err := db.PingContext(dbCtx); err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	repo := repository.NewRepository(db)

	r := NewRoutes(logCtx, repo)

	r.Run(":" + port)
}
