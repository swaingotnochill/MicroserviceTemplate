// Sample Handler

package handlers

import (
	"log"
	"net/http"
)

type Bye struct {
	l *log.Logger
}

func NewBye(l *log.Logger) *Bye {
	return &Bye{l}
}

func (b *Bye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("This is Bye Handler..")
	rw.Write([]byte("Byee..."))

}
