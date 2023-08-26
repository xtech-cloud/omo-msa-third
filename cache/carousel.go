package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"time"
)

type CarouselInfo struct {
	Status uint8
	baseInfo
	Owner  string
	Quotes []*nosql.QuoteInfo
}

func (mine *cacheContext) CreateCarousel(owner, operator string, list []*pb.QuoteInfo) (*CarouselInfo, error) {
	db := new(nosql.Carousel)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecommendNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = ""
	db.Status = 0
	db.Owner = owner
	db.Quotes = make([]*nosql.QuoteInfo, 0, len(list))
	for _, info := range list {
		db.Quotes = append(db.Quotes, &nosql.QuoteInfo{Type: uint8(info.Type), Title: info.Title,
			UID: info.Uid, Asset: info.Asset, Updated: info.Updated})
	}
	err := nosql.CreateCarousel(db)
	if err != nil {
		return nil, err
	}
	info := new(CarouselInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCarousel(owner string) (*CarouselInfo, error) {
	db, err := nosql.GetCarouselBy(owner)
	if err != nil {
		return nil, err
	}
	info := new(CarouselInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveCarousel(uid, operator string) error {
	return nosql.RemoveCarousel(uid, operator)
}

func (mine *CarouselInfo) initInfo(db *nosql.Carousel) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Name = db.Name
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Status = db.Status
	mine.Owner = db.Owner
	mine.Quotes = db.Quotes
}

func (mine *CarouselInfo) UpdateQuote(asset, quo, name string, ty uint8) error {
	if !mine.HadQuote(asset) {
		return errors.New("the asset is no existed")
	}
	for i, quote := range mine.Quotes {
		if quote.Asset == asset {
			mine.Quotes[i].UID = quo
			mine.Quotes[i].Title = name
			mine.Quotes[i].Type = ty
		}
	}
	err := nosql.UpdateCarouselTargets(mine.UID, "", mine.Quotes)
	if err == nil {
		return err
	}
	return err
}

func (mine *CarouselInfo) UpdateQuotes(operator string, list []*pb.QuoteInfo) error {
	arr := make([]*nosql.QuoteInfo, 0, len(list))
	for _, quote := range mine.Quotes {
		arr = append(arr, &nosql.QuoteInfo{Type: quote.Type, UID: quote.UID,
			Title: quote.Title, Asset: quote.Asset, Updated: quote.Updated})
	}
	err := nosql.UpdateCarouselTargets(mine.UID, operator, arr)
	if err == nil {
		return err
	}
	mine.Quotes = arr
	return err
}

func (mine *CarouselInfo) HadQuote(asset string) bool {
	for _, quote := range mine.Quotes {
		if quote.Asset == asset {
			return true
		}
	}
	return false
}

func (mine *CarouselInfo) AppendQuote(name, asset, source string, tp uint8) error {
	if mine.HadQuote(asset) {
		return errors.New("the asset is existed")
	}
	tmp := &nosql.QuoteInfo{}
	tmp.Asset = asset
	tmp.Title = name
	tmp.Type = tp
	tmp.UID = source
	tmp.Updated = time.Now().Unix()
	err := nosql.AppendCarouselQuote(mine.UID, tmp)
	if err == nil {
		mine.Quotes = append(mine.Quotes, tmp)
	}
	return err
}

func (mine *CarouselInfo) SubtractQuote(asset string) error {
	if !mine.HadQuote(asset) {
		return nil
	}

	err := nosql.SubtractCarouselQuote(mine.UID, asset)
	if err == nil {
		for i, quote := range mine.Quotes {
			if quote.Asset == asset {
				if i == len(mine.Quotes)-1 {
					mine.Quotes = append(mine.Quotes[:i])
				} else {
					mine.Quotes = append(mine.Quotes[:i], mine.Quotes[i+1:]...)
				}
				break
			}
		}

	}
	return err
}
