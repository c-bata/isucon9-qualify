package main

import (
	"errors"
	"sync"

	"github.com/c-bata/measure"
	"github.com/jmoiron/sqlx"
)

var (
	categories   = make(map[int]Category, 43)
	categoriesmu sync.RWMutex
)

func initCategories() error {
	var categoryList []int
	err := dbx.Select(&categoryList, "SELECT id FROM categories")
	if err != nil {
		return err
	}

	categoriesmu.Lock()
	defer categoriesmu.Unlock()
	for _, cid := range categoryList {
		c, err := getCategoryByIDFromDB(dbx, cid)
		if err != nil {
			return err
		}
		categories[c.ID] = c
	}
	return nil
}

func getCategoryByIDFromDB(q sqlx.Queryer, categoryID int) (category Category, err error) {
	defer measure.Start("get_category_by_id").Stop()
	err = sqlx.Get(q, &category, "SELECT * FROM `categories` WHERE `id` = ?", categoryID)
	if category.ParentID != 0 {
		parentCategory, err := getCategoryByID(q, category.ParentID)
		if err != nil {
			return category, err
		}
		category.ParentCategoryName = parentCategory.CategoryName
	}
	return category, err
}

func getCategoryByID(q sqlx.Queryer, categoryID int) (category Category, err error) {
	defer measure.Start("get_category_by_id").Stop()

	categoriesmu.RLock()
	defer categoriesmu.RUnlock()
	if c, ok := categories[categoryID]; ok {
		return c, nil
	}
	return Category{}, errors.New("category is not found")
}

func getParentCategory(q sqlx.Queryer, base *Category) error {
	defer measure.Start("get_parent_category").Stop()

	categoriesmu.RLock()
	defer categoriesmu.RUnlock()
	if c, ok := categories[base.ParentID]; ok {
		base.ParentCategoryName = c.CategoryName
		return nil
	}
	return errors.New("category is not found")
}
