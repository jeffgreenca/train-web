package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
)

type config struct {
	HttpPort  string `env:"HTTP_PORT" envDefault:"3000"`
	TrainHost string `env:"TRAIN_HOST" envDefault:"10.0.1.213"`
	TrainPort string `env:"TRAIN_PORT" envDefault:"1011"`
}

var cfg = config{}

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
}

func main() {

	r := gin.Default()
	r.POST("/train/:speed", updateTrain)
	r.Run(fmt.Sprintf(":%s", cfg.HttpPort))
}

func updateTrain(c *gin.Context) {
	speed := c.Param("speed")
	// validate input
	if len(speed) != 2 {
		c.Status(http.StatusBadRequest)
		return
	}
	err := checkRange(speed[0:1])
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = checkRange(speed[1:2])
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// update train
	err = writeTrain(speed)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(200)
}

func checkRange(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if n < 0 || n > 9 {
		return errors.New("out of range")
	}
	return nil
}

func writeTrain(speed string) error {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", cfg.TrainHost, cfg.TrainPort))
	// conn, err := net.Dial("udp", "localhost:1011")
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
	if string(p[:2]) != "OK" {
		fmt.Printf("unexpected response: %v", string(p))
		return errors.New("unexpected response")
	}
	return nil
}
