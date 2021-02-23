package db

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// Character describes a model of a character.
type Character struct {
	gorm.Model
	Name     string
	RarityID int
	Rarity   Rarity
}

// Rarity describes a model of a rarity of a characcter.
type Rarity struct {
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

var rarity []Rarity
var charactersByRarityID map[int][]Character
var invalidateAt time.Time

func cacheGachaRarity() error {
	if db == nil {
		return fmt.Errorf("no database connection")
	}

	db.Order("rate").Find(&rarity)

	return nil
}

func cacheGachaCharacters() error {
	if db == nil {
		return fmt.Errorf("no database connection")
	}

	var charactersList []Character
	db.Find(&charactersList)

	charactersByRarityID = map[int][]Character{}

	for i := 0; i < len(charactersList); i++ {
		item := charactersList[i]
		charactersByRarityID[item.RarityID] = append(charactersByRarityID[item.RarityID], item)
	}

	return nil
}

func cacheData() error {
	// already cached
	if rarity != nil && charactersByRarityID != nil && invalidateAt.After(time.Now()) {
		return nil
	}

	invalidateAt = time.Now().Add(10 * time.Minute)

	fmt.Println("cache gacha data...")

	if err := cacheGachaRarity(); err != nil {
		return err
	}
	if err := cacheGachaCharacters(); err != nil {
		return err
	}

	fmt.Println("cache gacha data done")
	return nil
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

// DrawGacha draws gacha for given times.
func DrawGacha(userID uint, times int) (userCharacters []UserCharacter, err error) {
	if db == nil {
		err = fmt.Errorf("no database connection")
		return
	}

	if err = cacheData(); err != nil {
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
	for i := 0; i < times; i++ {
		rate := rand.Float64()
		idx := 0
		for j := 0; j < len(rarity); j++ {
			if rate < rarity[j].Rate {
				break
			}
			idx = j
		}

		rarityID := int(rarity[idx].ID)
		charactersByRarity := charactersByRarityID[rarityID]
		characterIdx := rand.Int() % len(charactersByRarity)

		characters = append(characters, charactersByRarity[characterIdx])
	}

	return
}

func pickCharacter() (character Character, err error) {
	err = db.Raw(`
	with
		rarity as (select rand() as rarity),
		characters_with_rarity as (
			select
				characters.id as id,
				characters.name as name,
				rate
			from characters
			join rarities
			on characters.rarity_id = rarities.id
		)
	select id, name, rate, rarity from
		characters_with_rarity, rarity
	where rate < rarity
	order by rate desc, rand()
	limit 1
	`).Scan(&character).Error

	return
}
