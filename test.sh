#!/bin/bash

url="http://localhost:8080"

produkt='{"Nazwa": "nowy produkt", "Cena": 20, "KategoriaID": 1}'

echo "Tworzenie nowego produktu..."
nowy_produkt=$(curl -X POST -H "Content-Type: application/json" -d "$produkt" "$url/produkty")
echo "Nowy produkt: $nowy_produkt"