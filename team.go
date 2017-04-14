package main

import (
	"net/http"
	//"strconv"
	//"fmt"

	"github.com/labstack/echo"
	//"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	//"github.com/gocraft/dbr/dialect"
)

type (
	dataJSON struct {
		Domain_name    string `json:"domain_name"`
		Google_user_id int    `json:"google_user_id"`
	}

	team struct {
		Team_id     int    `db:team_id`
		Domain_name string `db:"domain_name"`
		Team_name   string `db:"team_name"`
		Created_at  string `db:"created_at"`
	}

	teamJSON struct {
		Domain_name string `json:"domain_name"`
		Team_name   string `json:"team_name"`
		Created_at  string `json:"created_at"`
	}

	user struct {
		User_id        int    `db:"user_id"`
		Team_id        int    `db:"team_id"`
		Google_user_id int    `db:"google_user_id"`
		Name           string `db:"name"`
		Account        string `db:"account"`
		Email          string `db:"email"`
		Status         int    `db:"status"`
		Created_at     string `db:"created_at"`
	}

	channel struct {
		Channel_id int    `db:"channel_id"`
		User_id    int    `db:"user_id"`
		Name       string `db:"name"`
		Type       int    `db:"type"`
	}

	channelJSON struct {
		Channel_id int    `json:"channel_id"`
		User_id    int    `json:"user_id"`
		Name       string `json:"name"`
		Type       int    `json:"type"`
	}

	create_channelJSON struct {
		Team_id string `json:"team_id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
	}

	responseData struct {
		Channels []channel `json:"channels"`
	}
)

var (
	table1  = "team"
	table2  = "user"
	table3  = "channel"
	table4  = "channel_access"
	seq     = 1
	conn, _ = dbr.Open("mysql", "root:@localhost:3306", nil)
	sess    = conn.NewSession(nil)
)

//----------
// Handlers
//----------

func selectChannel(c echo.Context) error {
	var m team
	var u user
	var a []channel

	i := new(dataJSON)
	if err := c.Bind(u); err != nil {
		return err
	}

	sess.Select("*").From(table1).Where("domain_name = ?", i.Domain_name).Load(&m)
	//	return c.JSON(http.StatusOK, m)

	sess.Select("*").From(table2).Where("team_id = ? AND google_user_id = ?", m.Team_id, i.Google_user_id).Load(&u)

	sess.Select("channel_id, user_id, table3.name, table3.type").From(table4).LeftJoin(table3, "table4.channel_id = table3.channel_id").Where("user_id = ?", u.User_id).Load(&a)

	response := new(responseData)
	response.Channels = a
	return c.JSON(http.StatusOK, response)
}

func createChannel(c echo.Context) error {
	u := new(create_channelJson)
	if err := c.Bind(u); err != nil {
		return err
	}
	sess.InsertInto(table3).Columns("team_id", "name", "type").Values(u.Team_id, u.Name, u.Type).Exec()

	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes

	e.POST("/login/", selectChannel)
	e.POST("/create_channel/", createChannel)
	e.Logger.Fatal(e.Start(":3000"))
}
