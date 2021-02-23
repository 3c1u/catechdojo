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

	characters, err := pickCharacters(times)
	if err != nil {
		return
	}

	for i := 0; i < times; i++ {
		character := UserCharacter{
			UserID:      userID,
			CharacterID: characters[i].ID,
		}
		userCharacters = append(userCharacters, character)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Create(&userCharacters).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	for i := 0; i < times; i++ {
		userCharacters[i].Character = characters[i]
	}

	return
}

func pickCharacters(times int) (characters []Character, err error) {
	var character Character

	for i := 0; i < times; i++ {
		// FIXME: ループ内でのクエリ実行（timesは十分小さいとはいえ...）
		character, err = pickCharacter()
		if err != nil {
			characters = nil
			return
		}

		characters = append(characters, character)
	}

	return
}

func pickCharacter() (character Character, err error) {
	// TODO: ガチャが膨れ上がった時にtotal_rate, sum_rateはキャッシュしたほうが良い気がする
	err = db.Raw(`with
				characters_with_sum as (select
					id,
					name,
					rate,
					sum(rate) over(order by id) as sum_rate,
					sum(rate) over() as total_rate
				from characters),
				random_rate as (select rand() as random_rate)
			select id, name from
				characters_with_sum, random_rate
			where
				sum_rate > total_rate * random_rate
			order by rate desc, rand()
			limit 1`).Scan(&character).Error

	return
}
