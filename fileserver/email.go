package fileserver

import (
	"crypto/md5"
	"fmt"
	"log"
)

// TOOD: send the email
func (server *Server) sendEmail(email string, id uint64) {

	hash := md5.Sum([]byte(fmt.Sprintf("%d", id)))
	hashString := string(fmt.Sprintf("%x", hash))

	log.Printf("%s\n", hashString)
}
