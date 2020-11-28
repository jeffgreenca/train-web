package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/train/:speed", updateTrain)
	r.Run(":3000")
}

func updateTrain(c *gin.Context) {
	speed := c.Param("speed")
	if len(speed) != 2 {
		c.Status(http.StatusBadRequest)
		return
	}
	for i := 0; i < 2; i++ {
		n, err := strconv.Atoi(speed[i : i+1])
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if n < 0 || n > 9 {
			c.Status(http.StatusBadRequest)
			return
		}

	}

	err := writeTrain(speed)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(200)
}

func writeTrain(speed string) error {
	// conn, err := net.Dial("udp", "10.0.1.213:1011")
	conn, err := net.Dial("udp", "localhost:1011")
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Fprintf(conn, speed+"\n")

	// recv
	p := make([]byte, 4)
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return err
	}
	if string(p) != "OK!\n" {
		return errors.New("unexpected response")
	}
	return nil
}
