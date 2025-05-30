package conversions

import(
	"fmt"
	"encoding/hex"
)

func RGBToHexString(r uint32, g uint32, b uint32)(string){
	bytes := []byte{byte(r), byte(g), byte(b)}
	return hex.EncodeToString(bytes)
}

func StringToRGBValues(colour string)(uint8, uint8, uint8, error){
	var r uint8
	var g uint8
	var b uint8
	var err error
	var dec []byte

	if len(colour) != 7{
		return r,g,b,fmt.Errorf("RGB colour string length should` be 7")
	}

	if colour[0] != '#' {
		return r,g,b,fmt.Errorf("RGB colour string should start with '#'")
	}

	//######### Debug only #######################
	//fmt.Printf("%s %s %s\n",colour[1:3], colour[3:5], colour[5:7])
	//#############################################

	dec, err = hex.DecodeString(colour[1:3])
	r = dec[0]

	dec, err = hex.DecodeString(colour[3:5])
	g = dec[0]

	dec, err = hex.DecodeString(colour[5:7])
	b = dec[0]


	//######### Debug only #######################
	//fmt.Printf("%d %d %d\n",r,g,b)
	//#############################################

	return r,g,b,err
}
