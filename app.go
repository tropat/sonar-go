package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func main() {
	e := echo.New()

	db, err := gorm.Open(sqlite.Open("zad04.db"), &gorm.Config{})
	if err != nil {
		panic("Blad polaczenia z baza danych")
	}
	db.AutoMigrate(&Produkt{}, &Koszyk{}, &Kategoria{})

	e.GET("/produkty", func(c echo.Context) error {
		var produkty []Produkt
		if err := db.Find(&produkty).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Brak produktow")
		}
		return c.JSON(http.StatusOK, produkty)
	})

	e.GET("/produkty/:id", func(c echo.Context) error {
		id := c.Param("id")
		var produkt Produkt
		if err := db.First(&produkt, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Brak podanego produktu")
		}
		return c.JSON(http.StatusOK, produkt)
	})

	e.POST("/produkty", func(c echo.Context) error {
		produkt := new(Produkt)
		if err := c.Bind(produkt); err != nil {
			return err
		}
		db.Create(&produkt)

		return c.JSON(http.StatusCreated, produkt)
	})

	e.PUT("/produkty/:id", func(c echo.Context) error {
		id := c.Param("id")
		var produkt Produkt
		if err := db.First(&produkt, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Brak podanego produktu")
		}
		if err := c.Bind(&produkt); err != nil {
			return err
		}
		db.Save(&produkt)
		return c.JSON(http.StatusOK, produkt)
	})

	e.DELETE("/produkty/:id", func(c echo.Context) error {
		id := c.Param("id")
		var produkt Produkt
		if err := db.First(&produkt, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Brak podanego produktu")
		}
		db.Delete(&produkt)
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/koszyk/:id", func(c echo.Context) error {
		id := c.Param("id")
		var koszyk Koszyk
		if err := db.Preload("Produkty").First(&koszyk, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, "Koszyk nie znaleziony")
		}
		return c.JSON(http.StatusOK, koszyk)
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

	e.Logger.Fatal(e.Start(":8080"))
}
