package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	ErrUserIDRequired modelError = "models: user ID is required"
	ErrTitleRequired  modelError = "models: title is required"
)

var _ GalleryService = &galleryGorm{}

type galleryValFn func(*Gallery) error

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not null;index"`
	Title  string  `gorm:"not null"`
	Images []Image `gorm:"-"`
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

func (g *galleryGorm) Create(gallery *Gallery) error {
	return g.db.Create(gallery).Error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (g *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := g.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (g *galleryGorm) Update(gallery *Gallery) error {
	return g.db.Save(gallery).Error
}

func (gv *galleryValidator) nonZeroID(gallery *Gallery) error {
	if gallery.ID <= 0 {
		return ErrInvaildID
	}
	return nil
}

func (gv *galleryValidator) Delete(id uint) error {
	var gallery Gallery
	gallery.ID = id
	err := runGalleryValFns(&gallery,
		gv.nonZeroID)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Delete(gallery.ID)
}

func (g *galleryGorm) Delete(id uint) error {
	var gollery Gallery
	gollery.ID = id
	return g.db.Delete(&gollery).Error
}

func (g *galleryGorm) ByUserID(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	db := g.db.Where("user_id = ?", userID)
	if err := db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (g *Gallery) ImagesSplitN(n int) [][]Image {
	ret := make([][]Image, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}

	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}
