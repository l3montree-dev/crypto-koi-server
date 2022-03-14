package models

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type Cryptogotchi struct {
	Base
	Name    *string   `json:"name" gorm:"type:varchar(255);default:null"`
	OwnerId uuid.UUID `json:"owner" gorm:"type:char(36); not null"`

	IsValidNft bool `json:"isValidNft" gorm:"default:false"`

	PredictedDeathDate time.Time  `json:"-" gorm:"not null"`
	LastFed            *time.Time `json:"-" gorm:"default:null"`
	// values between 100 and 0.
	Food float64 `json:"food" gorm:"default:100"`
	// drain per minute
	FoodDrain float64 `json:"foodDrain" gorm:"default:0.5"`
	// the id of the token - might be changed in the future.
	// mapping to the event struct.
	Events    []Event    `json:"events"`
	GameStats []GameStat `json:"game_stats" gorm:"foreignKey:cryptogotchi_id"`
	Active    bool       `json:"released" gorm:"default:true"`
	// the timestamp of the current snapshot stored inside the database.
	// in most cases this equals the LastFeed value. - nevertheless to build the struct a bit more
	// future proof, we store the timestamp of the snapshot in the database as a separate column.
	// currently this affects only the food value.
	SnapshotValid time.Time `json:"-" gorm:"not null"`
	Rank          int       `json:"rank" gorm:"default:-1"`
}

func (c *Cryptogotchi) ToOpenseaNFT(baseUrl string) (OpenseaNFT, error) {
	uintStr, err := util.UuidToUint256(c.Id.String())
	koi := cryptokoi.NewKoi(c.Id.String())
	attributes := koi.GetAttributes()
	if err != nil {
		return OpenseaNFT{}, err
	}

	state := "Alive"
	if !c.IsAlive() {
		state = "Dead"
	}

	backgroundColor := util.Shade(attributes.PrimaryColor, -20)
	if util.IsDark(attributes.PrimaryColor) {
		backgroundColor = util.Shade(attributes.PrimaryColor, 20)
	}

	return OpenseaNFT{
		Name:            *c.Name,
		Image:           baseUrl + "v1/images/" + uintStr.String(),
		BackgroundColor: util.ConvertColor2HexWithoutHash(backgroundColor),

		Attributes: []OpenseaNFTAttribute{
			{
				TraitType:   "Birthday",
				DisplayType: DateDisplayType,
				Value:       c.CreatedAt.Unix(),
			},
			{
				TraitType: "State",
				Value:     state,
			},
			{
				TraitType: "Primary Color",
				Value:     util.ConvertColor2Hex(attributes.PrimaryColor),
			},
			{
				TraitType: "Body Color",
				Value:     util.ConvertColor2Hex(attributes.BodyColor),
			},
			{
				TraitType: "Fin Color",
				Value:     util.ConvertColor2Hex(attributes.FinColor),
			},
			{
				TraitType:   "Pattern Quantity",
				DisplayType: NumberDisplayType,
				Value:       len(attributes.BodyImages) + len(attributes.FinImages) + len(attributes.HeadImages),
			},
			{
				TraitType: "Species",
				Value:     attributes.KoiType,
			},
		}}, nil
}

// make sure to only call this function after the food value has been updated.
func (c *Cryptogotchi) PredictNewDeathDate() time.Time {
	return time.Now().Add(time.Duration(c.Food/c.FoodDrain) * time.Minute)
}

func (c *Cryptogotchi) IsAlive() bool {
	return c.PredictedDeathDate.After(time.Now())
}

// the cryptogotchi has a few time dependentant state variables.
// this function updates the state variables according to the provided time
// returns if the cryptogotchi is still alive
func (c *Cryptogotchi) ProgressUntil(nextTime time.Time) (bool, time.Time) {
	if c.PredictedDeathDate.Before(nextTime) {
		c.Food = 0
		return false, c.PredictedDeathDate
	}
	// calculate the time difference
	timeDiffMinutes := nextTime.Sub(c.SnapshotValid).Minutes()

	// calculate the food value
	// just calculating the current state does not change the predicted death date.
	c.Food -= timeDiffMinutes * c.FoodDrain
	// update the snapshot validity since we mutated the cryptogotchi
	c.SnapshotValid = nextTime
	return true, time.Time{}
}

func (c *Cryptogotchi) GetNextFeedingTime() time.Time {

	if c.LastFed == nil {
		return time.Now()
	}

	return (*c.LastFed).Add(config.TIME_BETWEEN_FEEDINGS)
}

func NewCryptogotchi(user *User) Cryptogotchi {
	return Cryptogotchi{OwnerId: user.Id}
}
