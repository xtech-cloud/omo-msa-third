package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Carousel struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Owner  string       `json:"owner" bson:"owner"`
	Status uint8        `json:"status" bson:"status"`
	Quotes []*QuoteInfo `json:"quotes" bson:"quotes"`
}

type QuoteInfo struct {
	Type    uint8  `json:"type" bson:"type"`
	UID     string `json:"uid" bson:"uid"`
	Title   string `json:"title" bson:"title"`
	Asset   string `json:"asset" bson:"asset"`
	Updated int64  `json:"updatedAt" bson:"updatedAt"`
}

func CreateCarousel(info *Carousel) error {
	_, err := insertOne(TableCarousel, info)
	return err
}

func GetCarouselNextID() uint64 {
	num, _ := getSequenceNext(TableCarousel)
	return num
}

func RemoveCarousel(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Carousel uid is empty ")
	}
	_, err := removeOne(TableCarousel, uid, operator)
	return err
}

func GetCarousel(uid string) (*Carousel, error) {
	if len(uid) < 2 {
		return nil, errors.New("db activity uid is empty of GetCarousel")
	}
	result, err := findOne(TableCarousel, uid)
	if err != nil {
		return nil, err
	}
	model := new(Carousel)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCarouselBy(owner string) (*Carousel, error) {
	if len(owner) < 2 {
		return nil, errors.New("db owner is empty of GetCarouselBy")
	}
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	result, err := findOneBy(TableCarousel, filter)
	if err != nil {
		return nil, err
	}
	model := new(Carousel)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCarouselCount() int64 {
	num, _ := getTotalCount(TableCarousel)
	return num
}

func GetCarouselsByOwner(owner string) ([]*Carousel, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableCarousel, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Carousel, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Carousel)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCarouselStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCarousel, uid, msg)
	return err
}

func UpdateCarouselTargets(uid, operator string, list []*QuoteInfo) error {
	msg := bson.M{"quotes": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCarousel, uid, msg)
	return err
}

func AppendCarouselQuote(uid string, quote *QuoteInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"quotes": quote}
	_, err := appendElement(TableCarousel, uid, msg)
	return err
}

func SubtractCarouselQuote(uid, asset string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"quotes": bson.M{"asset": asset}}
	_, err := removeElement(TableCarousel, uid, msg)
	return err
}
