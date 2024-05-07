package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ProduktyEndpoint = "/produkty/:id"
var BrakProduktuMessage = "Brak podanego produktu"

type Produkt struct {
	gorm.Model
	Nazwa       string
	Cena        float64
	KategoriaID uint
}

type Koszyk struct {
	gorm.Model
	Produkty []*Produkt `gorm:"many2many:koszyk_produkty;"`
}

type Kategoria struct {
	gorm.Model
	Nazwa string
}

func initializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("zad04.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Produkt{}, &Koszyk{}, &Kategoria{}); err != nil {
		return nil, err
	}
	return db, nil
}

func setupGetRouts(e *echo.Echo, db *gorm.DB) {
	e.GET("/produkty", func(c echo.Context) error {
		var produkty []Produkt
		if err := db.Find(&produkty).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Brak produktow")
		}
		return c.JSON(http.StatusOK, produkty)
	})

	e.GET(ProduktyEndpoint, func(c echo.Context) error {
		id := c.Param("id")
		var produkt Produkt
		if err := db.First(&produkt, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, BrakProduktuMessage)
		}
		return c.JSON(http.StatusOK, produkt)
	})

	e.GET("/koszyk/:id", func(c echo.Context) error {
		id := c.Param("id")
		var koszyk Koszyk
		if err := db.Preload("Produkty").First(&koszyk, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Koszyk nie znaleziony")
		}
		return c.JSON(http.StatusOK, koszyk)
	})
}

func setupRoutes(e *echo.Echo, db *gorm.DB) {

	setupGetRouts(e, db)

	e.POST("/produkty", func(c echo.Context) error {
		produkt := new(Produkt)
		if err := c.Bind(produkt); err != nil {
			return err
		}
		db.Create(&produkt)

		return c.JSON(http.StatusCreated, produkt)
	})

	e.PUT(ProduktyEndpoint, func(c echo.Context) error {
		id := c.Param("id")
		var produkt Produkt
		if err := db.First(&produkt, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, BrakProduktuMessage)
		}
		if err := c.Bind(&produkt); err != nil {
			return err
		}
		db.Save(&produkt)
		return c.JSON(http.StatusOK, produkt)
	})

	e.DELETE(ProduktyEndpoint, func(c echo.Context) error {
		id := c.Param("id")
		var produkt Produkt
		if err := db.First(&produkt, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, BrakProduktuMessage)
		}
		db.Delete(&produkt)
		return c.NoContent(http.StatusNoContent)
	})

	e.PUT("/koszyk/:id", func(c echo.Context) error {
		id := c.Param("id")

		var koszyk Koszyk
		if err := db.First(&koszyk, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Brak podanego koszyka")
		}

		var nowyProdukt Produkt
		if err := c.Bind(&nowyProdukt); err != nil {
			return err
		}

		koszyk.Produkty = append(koszyk.Produkty, &nowyProdukt)

		if err := db.Save(&koszyk).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, koszyk)
	})
}

func main() {
	e := echo.New()

	db, err := initializeDB()
	if err != nil {
		panic("Blad polaczenia z baza danych: " + err.Error())
	}

	setupRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}
