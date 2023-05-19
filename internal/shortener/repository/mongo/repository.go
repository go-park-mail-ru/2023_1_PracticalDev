package mongo

import (
	"context"
	"errors"

	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	short "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener"
	gen "github.com/go-park-mail-ru/2023_1_PracticalDev/pkg/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type shortener struct {
	db  *mongo.Database
	log *zap.Logger
}

type shortURL struct {
	OriginalURL string `bson:"original_url"`
	ShortCode   string `bson:"short_code"`
}

var (
	ErrGenerateAttemptsLimitExceeded = errors.New("generate attempts limit exceeded")
)

func NewShortenerRepository(db *mongo.Database, log *zap.Logger) short.ShortenerRepository {
	return &shortener{
		db:  db,
		log: log,
	}
}

func (s *shortener) Get(hash string) (string, error) {
	res := &shortURL{}
	collection := s.db.Collection("urls")
	err := collection.FindOne(context.Background(), bson.M{"short_code": hash}).Decode(res)
	if err == mongo.ErrNoDocuments {
		s.log.Error("failed to find short link", zap.String("ShortCode", hash), zap.Error(err))
		return "", pkgErrors.ErrLinkNotFound
	}

	s.log.Debug("found short link", zap.String("OriginalURL", res.OriginalURL), zap.String("ShortCode", hash))
	return res.OriginalURL, nil
}

func (s *shortener) Create(url string) (string, error) {
	collection := s.db.Collection("urls")
	res := &shortURL{}
	err := collection.FindOne(context.Background(), bson.M{"original_url": url}).Decode(res)
	if err == nil {
		s.log.Debug("short link already exists", zap.String("OriginalURL", url), zap.String("ShortCode", res.ShortCode))
		return res.ShortCode, nil
	}
	count := 10
	for count > 0 {
		sc := gen.GenerateShortCode()
		err := collection.FindOne(context.Background(), bson.M{"short_code": sc}).Decode(res)
		if err == mongo.ErrNoDocuments {
			_, err := collection.InsertOne(context.Background(), &shortURL{
				OriginalURL: url,
				ShortCode:   sc,
			})
			if err != nil {
				s.log.Error("failed to insert short link", zap.String("OriginalURL", url), zap.Error(err))
				return "", err
			}
			s.log.Debug("generated short link", zap.String("OriginalURL", url), zap.String("ShortCode", sc))
			return sc, nil
		}
		if err != nil {
			s.log.Error("failed to lookup short link", zap.String("OriginalURL", url), zap.Error(err))
		}
		count--
	}
	s.log.Error("generate attempts limit exceeded")
	return "", ErrGenerateAttemptsLimitExceeded
}
