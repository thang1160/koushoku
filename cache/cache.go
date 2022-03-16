package cache

import "time"

var Archives *LRU
var Taxonomies *LRU
var Templates *LRU

var Users *LRU
var Favorites *LRU
var Cache *LRU

const defaultExpr = time.Hour
const templateExpr = 5 * time.Minute

func Init() {
	Archives = New(4096, defaultExpr)
	Taxonomies = New(4096, defaultExpr)
	Templates = New(4096, templateExpr)

	Users = New(2048, defaultExpr)
	Favorites = New(2048, defaultExpr)
	Cache = New(512, defaultExpr)
}
