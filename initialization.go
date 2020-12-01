package main

import (
	"crypto/aes"
	"crypto/cipher"
	_ "image/png"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/opentype"
)

//fonts
var (
	FontBig    font.Face
	FontMedium font.Face
	FontSmall  font.Face
)

func init() {
	tt, err := opentype.Parse(gomedium.TTF)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	FontMedium, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	FontBig, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    70,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	FontSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	TitleTextX = getCentredPosForText(TitleTextValue, FontBig)
	PressSpaceTextX = getCentredPosForText(PressSpaceTextValue, FontSmall)
	ControlsTextX = getCentredPosForText(ControlsTextValue, FontSmall)
	GameOverTextX = getCentredPosForText(GameOverTextValue, FontMedium)
	NewRecordTextX = getCentredPosForText(NewRecordTextValue, FontMedium)

}

//Cipher using for record file
var Cipher cipher.Block

func init() {
	var err error
	Cipher, err = aes.NewCipher(SecretKey)
	if err != nil {
		log.Fatal(err)
	}
}
