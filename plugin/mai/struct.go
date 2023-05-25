package mai

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type maidx struct {
	AdditionalRating int `json:"additional_rating"`
	Charts           struct {
		Dx []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"dx"`
		Sd []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"sd"`
	} `json:"charts"`
	Nickname string      `json:"nickname"`
	Plate    string      `json:"plate"`
	Rating   int         `json:"rating"`
	UserData interface{} `json:"user_data"`
	UserId   interface{} `json:"user_id"`
	Username string      `json:"username"`
}

type chun struct {
	Nickname string  `json:"nickname"`
	Rating   float64 `json:"rating"`
	Records  struct {
		B30 []struct {
			Cid        int     `json:"cid"`
			Ds         float64 `json:"ds"`
			Fc         string  `json:"fc"`
			Level      string  `json:"level"`
			LevelIndex int     `json:"level_index"`
			LevelLabel string  `json:"level_label"`
			Mid        int     `json:"mid"`
			Ra         float64 `json:"ra"`
			Score      int     `json:"score"`
			Title      string  `json:"title"`
		} `json:"b30"`
		R10 []struct {
			Cid        int     `json:"cid"`
			Ds         float64 `json:"ds"`
			Fc         string  `json:"fc"`
			Level      string  `json:"level"`
			LevelIndex int     `json:"level_index"`
			LevelLabel string  `json:"level_label"`
			Mid        int     `json:"mid"`
			Ra         float64 `json:"ra"`
			Score      int     `json:"score"`
			Title      string  `json:"title"`
		} `json:"r10"`
	} `json:"records"`
	Username string `json:"username"`
}

// HandleMaiDataByUsingText is a function that handle the data from Mai,but using text to output (.
func HandleMaiDataByUsingText(handlejson []byte) (text string) {
	var mai maidx
	_ = json.Unmarshal(handlejson, &mai)
	getUserName := mai.Nickname
	getUserRating := mai.Rating
	geDXtLength := len(mai.Charts.Dx) // DX 2022
	getSDLength := len(mai.Charts.Sd) // old.
	// (Player : MoeMagicMango) | Rating : 4141 + Additional Rating | 称号 : 世界の果て |
	// DX 2022 Music PLay Count : 0 | SD Music Play Count : 0
	/*
		- 1. MusicName (LevelLabel+Level - dx) (Achievement+Rate)  { (FC) (FS) if fs existed not show FC }} |  DXRating: (RA)
		example:
		- 1. 宿星審判 (Expert 12 - 12.8 ) (90.7542% AA) (FC) | DXRating: 105
		...
		Generated by Lucy(HiMoYoBOt),code with lazy.
	*/
	formatHeader := "(Player : " + getUserName + ") | Rating : " + strconv.Itoa(getUserRating) + " + " + strconv.Itoa(mai.AdditionalRating) + " | 称号 : " + mai.Plate + " |\n"
	formatEnd := "\nGenerated by Lucy(HiMoYoBOT),code with lazy."
	var setSongLength, mainText, numList, showFcs string
	for i := 0; i < geDXtLength; i++ {
		numList = strconv.Itoa(i + 1)
		if mai.Charts.Dx[i].Fc != "" {
			if mai.Charts.Dx[i].Fs != "" {
				showFcs = "FS"
			} else {
				showFcs = "FC"
			}
		}
		setSongLength = mai.Charts.Dx[i].Title
		if len(setSongLength) >= 30 {
			setSongLength = setSongLength[:29] + "..."
		}
		DXRating := computeRa(mai.Charts.Dx[i].Ds, mai.Charts.Dx[i].Achievements)
		mainText += fmt.Sprintf("- %v. %v (%v %v - %v) (%v Achievement %s ) (%v) | DXRating: %d \n", numList, setSongLength, mai.Charts.Dx[i].LevelLabel, mai.Charts.Dx[i].Level, mai.Charts.Dx[i].Ds, mai.Charts.Dx[i].Achievements, mai.Charts.Dx[i].Rate, showFcs, DXRating)
	}
	for i := 0; i < getSDLength; i++ {
		numList = strconv.Itoa(geDXtLength + i + 1)
		if mai.Charts.Sd[i].Fc != "" {
			if mai.Charts.Sd[i].Fs != "" {
				showFcs = "FS"
			} else {
				showFcs = "FC"
			}
		}
		setSongLength = mai.Charts.Sd[i].Title
		if len(setSongLength) >= 25 {
			setSongLength = setSongLength[:24] + "..."
		}
		DXRating := computeRa(mai.Charts.Sd[i].Ds, mai.Charts.Sd[i].Achievements)
		mainText += fmt.Sprintf("- %v. %v (%v %v - %v) (%v Achievement %v ) (%v) | DXRating: %d \n", numList, setSongLength, mai.Charts.Sd[i].LevelLabel, mai.Charts.Sd[i].Level, mai.Charts.Sd[i].Ds, mai.Charts.Sd[i].Achievements, mai.Charts.Sd[i].Rate, showFcs, DXRating)
	}
	text = formatHeader + mainText + formatEnd
	return text
}

func computeRa(ds, achievement float64) int {
	baseRa := 22.4
	if achievement < 50 {
		baseRa = 7.0
	} else if achievement < 60 {
		baseRa = 8.0
	} else if achievement < 70 {
		baseRa = 9.6
	} else if achievement < 75 {
		baseRa = 11.2
	} else if achievement < 80 {
		baseRa = 12.0
	} else if achievement < 90 {
		baseRa = 13.6
	} else if achievement < 94 {
		baseRa = 15.2
	} else if achievement < 97 {
		baseRa = 16.8
	} else if achievement < 98 {
		baseRa = 20.0
	} else if achievement < 99 {
		baseRa = 20.3
	} else if achievement < 99.5 {
		baseRa = 20.8
	} else if achievement < 100 {
		baseRa = 21.1
	} else if achievement < 100.5 {
		baseRa = 21.6
	}

	return int(math.Floor(ds * (math.Min(100.5, achievement) / 100) * baseRa))
}
