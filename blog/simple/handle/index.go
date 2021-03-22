package handle

import (
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/conf"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func handleIndex(c *gin.Context) {
	spage := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(spage)
	nums := len(NewPosts)
	allPages := nums / 5
	if nums%5 != 0 {
		allPages = nums/5 + 1
	}
	posts := NewPosts
	if (page * 5) < nums {
		posts = NewPosts[(page-1)*5 : page*5]
	} else {
		posts = NewPosts[(page-1)*5 : nums]
	}

	tabs := make([]int, allPages+2)
	if (page - 1) == 0 {
		tabs[0] = 1
	} else {
		tabs[0] = page - 1
	}
	for i := 1; i <= allPages; i++ {
		tabs[i] = i
	}
	if page+1 <= allPages {
		tabs[allPages+1] = page + 1
	} else {
		tabs[allPages+1] = 1
	}
	conf.Stat.ToCnt++
	//ip := c.ClientIP()
	c.HTML(http.StatusOK, "index.html", gin.H{"post": posts, "items": conf.Item.Items, "about": conf.Abt,
		"newcmts": NewCmts, "newart": NewArts, "hotart": HotArts, "vistcnt": conf.Stat.ToCnt, "curpage": page,
		"tabs": tabs})
}
