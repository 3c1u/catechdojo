package db

import (
	"fmt"

	"gorm.io/gorm"
)

// Character describes a model of a character.
type Character struct {
	gorm.Model
	Name string
	Rate float64
}

// UserCharacter describes a relation between a user and a character.
type UserCharacter struct {
	gorm.Model
	UserID      uint
	User        User
	CharacterID uint
	Character   Character
}

// EnumerateUserCharacters enumerates characters that a user has.
func EnumerateUserCharacters(userID uint) (characters []UserCharacter, err error) {
	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	err = db.Preload("User").Preload("Character").Where("user_id = ?", userID).Find(&characters).Error

	return
}

func DrawGacha(userID uint, times int) (userCharacters []UserCharacter, err error) {
	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	characters, err := PickCharacters(times)
	if err != nil {
		return
	}

	for i := 0; i < times; i++ {
		character := UserCharacter{
			UserID:      userID,
			CharacterID: characters[i].ID,
			Character:   characters[i],
		}
		userCharacters = append(userCharacters, character)
	}

	err = db.Create(&userCharacters).Error
	return
}

func PickCharacters(times int) (characters []Character, err error) {
	var character Character

	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	for i := 0; i < times; i++ {
		// FIXME: ループ内でのクエリ実行（timesは十分小さいとはいえ...）
		character, err = PickCharacter()
		if err != nil {
			characters = nil
			return
		}

		characters = append(characters, character)
	}

	return
}

func PickCharacter() (character Character, err error) {
	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	// TODO: ガチャが膨れ上がった時にtotal_rate, sum_rateはキャッシュしたほうが良い気がする
	db.Raw(`with
				characters_with_sum as (select
					id,
					name,
					rate,
					sum(rate) over(order by id) as sum_rate,
					sum(rate) over() as total_rate
				from characters),
				random_rate as (select rand() random_rate)
			select id, name from
				characters_with_sum, random_rate
			where
				sum_rate > total_rate * random_rate
			order by rate desc, rand()
			limit 1;`).Scan(&character)

	return
}
