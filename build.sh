#!/bin/bash
cd backend/
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../build/train-backend-pi
env GOOS=linux go build -o ../build/train-backend
