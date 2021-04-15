package bullet

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BulletServer struct {
	db *gorm.DB
}

func NewBulletServer(db *gorm.DB) BulletServer {
	return BulletServer{
		db: db,
	}
}

type bulletHandler struct {
	db *gorm.DB
}

type bulletContent struct {
	Bullet string `json:"bullet"`
}

func (h *bulletHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.FormValue("id")
		db := h.db
		bullets := make([]*Bullet, 0)
		if id == "0" {
			db = db.Limit(1).Order("id DESC")
		} else {
			db = db.Where("id > ?", id).Order("id")
		}
		db.Find(&bullets)
		b, _ := json.Marshal(&bullets)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
	if r.Method == "POST" {
		n, _ := ioutil.ReadAll(r.Body)
		var bullet bulletContent
		json.Unmarshal(n, &bullet)

		entity := Bullet{
			Content:  bullet.Bullet,
			SentTime: time.Now(),
		}
		h.db.Save(&entity)

	}

}

func (s *BulletServer) Serve(bulletListener net.Listener) {
	http.Handle("/bullet", &bulletHandler{db: s.db})
	http.Serve(bulletListener, nil)
}
