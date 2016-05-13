package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	//	"github.com/robfig/config"
)

type dbConfig struct {
	Dbname     string
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
}

type hostConfig struct {
	TargetDomain1 string
	TargetDomain2 string
}

var (
	db        *sql.DB
	idMapping = map[string]string{}
	dbconf    = dbConfig{}
	hostconf  = hostConfig{}
)

func main() {
	initConf()
	db := initDb()
	r := gin.Default()
	r.Use(gin.Recovery())

	//for pattern1
	r.GET("/detail/:id", func(c *gin.Context) {
		before_id := c.Param("id")

		after_id, ok := idMapping[before_id]
		if !ok {
			after_id = get(db, before_id)
		}
		c.Redirect(http.StatusMovedPermanently, "http://"+hostconf.TargetDomain1+"/detail/"+after_id)
	})

	//for pattern2
	r.GET("/id/:id", func(c *gin.Context) {
		before_id := c.Param("id")

		after_id, ok := idMapping[before_id]
		if !ok {
			after_id = get(db, before_id)
		}
		c.Redirect(http.StatusMovedPermanently, "http://"+hostconf.TargetDomain2+"/id/"+after_id)
	})

	//for cache clear
	r.GET("/clear", func(c *gin.Context) {
		idMapping = map[string]string{}
		c.String(200, "cache cleard")
	})

	if gin.Mode() == gin.ReleaseMode {
		endless.ListenAndServe(":18080", r)
	} else {
		r.Run(":18080")
	}
}

func initConf() {
	//c, _ := config.ReadDefault("config.cfg")
	//dbconf.Dbname, _ = c.String("default", "dbname")
	//dbconf.DbUser, _ = c.String("default", "dbuser")
	//dbconf.DbPassword, _ = c.String("default", "dbpassword")
	//dbconf.DbHost, _ = c.String("default", "dbhost")
	//dbconf.DbPort, _ = c.String("default", "dbport")
	//hostConfig.TargetDomain1, _ = c.String("default", "TargetDomain1")
	//hostConfig.TargetDomain2, _ = c.String("default", "TargetDomain2")

	dbconf.Dbname = "proxy"
	dbconf.DbUser = "postgres"
	dbconf.DbPassword = "postgres"
	dbconf.DbHost = "localhost"
	dbconf.DbPort = "5436"

	hostconf.TargetDomain1 = "localhost"
	hostconf.TargetDomain2 = "127.0.0.1"
}

func initDb() (db *sql.DB) {
	db, err := sql.Open("postgres", "user="+dbconf.DbUser+" password="+dbconf.DbPassword+" dbname="+dbconf.Dbname+" host="+dbconf.DbHost+" port="+dbconf.DbPort+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	return db
}

func get(db *sql.DB, before_id string) (after_id string) {
	err := db.QueryRow(`SELECT after_id FROM redirect WHERE before_id=$1`, before_id).Scan(&after_id)
	if err != nil {
		log.Fatal(err)
	}
	idMapping[before_id] = after_id
	return after_id
}
