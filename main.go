package main

import (
	"log"
	"net/http"
	_ "net/url"
	"os"

	_ "github.com/satori/go.uuid"

	"database/sql"

	_ "github.com/lib/pq"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/gin-gonic/gin"
	"fmt"


	_ "io/ioutil"


	_ "github.com/jinzhu/gorm"
	_ "github.com/gin-gonic/gin"
	_ "github.com/aws/aws-sdk-go/private/protocol"
	_ "strconv"
	"strconv"

)


type Member struct {
	ID   int
	ChineseName  string
	EnglishName string
	Email string
	CellPhone string
	Street string
	City string
	State string
	Zip string
	ShortPixName string
	PictureURL string
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func uploadAFile(c *gin.Context) (string, string, error) {
		return "", "", nil
		}

func main() {

	var DATABASE_URL = os.Getenv("DATABASE_URL")

	db, errDB := sql.Open("postgres", DATABASE_URL)
	defer db.Close()

	if errDB != nil {
		log.Fatalf("Error connecting to the DB")
	} else {
		log.Println("Connection is successful!")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")


	router.GET("/members/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "members.create.tmpl.html", nil)
	})

	router.GET("/members/select/:id", func(c *gin.Context) {

		id := c.Param("id")

		rows, err := db.Query("SELECT ID, ChineseName, EnglishName, " +
			" Email, CellPhone, Street, City, State, Zip, " +
				"PictureURL FROM members where ID = $1", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		member := new(Member)

		for rows.Next() {
			err := rows.Scan(&member.ID, &member.ChineseName,
				&member.EnglishName, &member.Email,
					&member.CellPhone, &member.Street,
						&member.City, &member.State,
							&member.Zip, &member.PictureURL)
			if err != nil {
				log.Fatal(err)
			}
		}


		c.HTML(http.StatusOK, "members.select.tmpl.html", member)
	})

	router.POST("/members/create", func(c *gin.Context) {
		EnglishName := c.PostForm("EnglishName")
		ChineseName := c.PostForm("ChineseName")
		Email := c.PostForm("Email")
		CellPhone := c.PostForm("CellPhone")
		Street := c.PostForm("Street")
		City := c.PostForm("City")
		State := c.PostForm("State")
		Zip := c.PostForm("Zip")

		//calling uploadAFile to upload it.
		shortPixName, returnedFile, err := uploadAFile(c)

		if err != nil {
			log.Println("Upload an image file encounters a problem.")
			c.HTML(http.StatusOK, "members.create_error.tmpl.html", err)
		}

		fmt.Println("returned file name is : ", returnedFile)


		_, errInsert := db.
		Exec("INSERT INTO members(ChineseName, EnglishName, " +
			"Email, CellPhone, Street, City, State, Zip, " +
				"ShortPixName, PictureURL) VALUES " +
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
			ChineseName, EnglishName, Email, CellPhone,
				Street, City, State, Zip, shortPixName, returnedFile)

		if errInsert != nil {
			log.Println("DB Insertion is in error.")
			c.HTML(http.StatusOK,
				"members.create_error.tmpl.html", errInsert)
		} else {
			log.Println("DB Insertion successful.")
			rows, err := db.Query("SELECT ID, ChineseName, " +
				"EnglishName, Email, CellPhone, Street, City, State, " +
					"Zip, PictureURL FROM members order by ID DESC")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()


			c.HTML(http.StatusOK, "members.create_ok.tmpl.html", nil)
		}
	})

	router.GET("/members/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		rows, err := db.Query("SELECT ID, ChineseName, EnglishName, " +
			"Email, CellPhone, Street, City, State, zip, " +
			"ShortPixName, PictureURL FROM members where ID = $1", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		member := new(Member)

		for rows.Next() {
			err := rows.Scan(&member.ID, &member.ChineseName, &member.EnglishName,
				&member.Email, &member.CellPhone,
				&member.Street, &member.City,
				&member.State, &member.Zip,
				&member.ShortPixName, &member.PictureURL)
			if err != nil {
				log.Fatal(err)
			}
		}
		c.HTML(http.StatusOK, "members.update.tmpl.html", member)

	})

	router.POST("/members/update/:id", func(c *gin.Context) {
		ID := c.Param("id")

		IDNumber, err1 := strconv.Atoi(ID)
		checkErr(err1)

		EnglishName := c.PostForm("EnglishName")
		ChineseName := c.PostForm("ChineseName")
		Email := c.PostForm("Email")
		CellPhone := c.PostForm("CellPhone")
		Street := c.PostForm("Street")
		City := c.PostForm("City")
		State := c.PostForm("State")
		Zip := c.PostForm("Zip")


		DistanceFromChurch := c.PostForm("DistanceFromChurch")

		// Update
		stmt, err := db.Prepare(
			"update members set EnglishName = $1, ChineseName = $2, " +
				"Email = $3, CellPhone = $4, " +
		        "Street = $5, City = $6, " +
		        "State = $7, Zip = $8, " +
		        "DistanceFromChurch = $9 where ID=$10")
		checkErr(err)
		fmt.Println("update statement is: ", stmt)

		val, err := strconv.ParseFloat(DistanceFromChurch, 32)

		fmt.Println("EnglishName, ChineseName, val, ID are: ", EnglishName, ChineseName, val, ID)

		res, err2 := stmt.Exec(EnglishName, ChineseName, Email, CellPhone,
			Street, City, State, Zip, val, IDNumber)

		checkErr(err2)
		defer stmt.Close()

		rowsAffected, err3 := res.RowsAffected()
		checkErr(err3)
		fmt.Println("rowsAffected is: ", rowsAffected)



		c.HTML(http.StatusOK, "members.update_post.tmpl.html", ID)

	})

	router.GET("/members/delete/:id", func(c *gin.Context) {
		ID := c.Param("id")

		member := new(Member)
		idHolder, err1 := strconv.Atoi(ID)
		member.ID = idHolder

		if err1 != nil {
			panic("ID is not a big integer. Terribly wrong")
		}

		c.HTML(http.StatusOK, "members.delete.tmpl.html", member)

	})

	router.POST("/members/delete/:id", func(c *gin.Context) {
		ID := c.Param("id")

		// Update
		stmt, err := db.Prepare(
			"delete from Members where ID=$1")
		checkErr(err)
		fmt.Println("statement is: ", stmt)
		fmt.Println("ID is: ", ID)

		res, err2 := stmt.Exec(ID)
		checkErr(err2)
		defer stmt.Close()

		rowsAffected, err3 := res.RowsAffected()
		checkErr(err3)
		fmt.Println("rowsAffected is: ", rowsAffected)

		c.HTML(http.StatusOK, "members.delete_post.tmpl.html", ID)

	})

	router.GET("/", func(c *gin.Context) {
		rows, err := db.Query("SELECT ID, ChineseName, " +
			"EnglishName, Email, CellPhone, ShortPixName, " +
				" PictureURL FROM Members order by ID DESC")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		members := make([]*Member, 0)
		for rows.Next() {
			member := new(Member)
			err := rows.Scan(&member.ID, &member.ChineseName, &member.EnglishName,
				&member.Email, &member.CellPhone,
					&member.ShortPixName, &member.PictureURL)
			if err != nil {
				log.Fatal(err)
			}
			members = append(members, member)
		}
		c.HTML(http.StatusOK, "index.tmpl.html", members)

	})

	router.Run(":" + port)
}
