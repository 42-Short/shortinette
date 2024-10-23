#!/bin/sh

echo "Installing Docker..."

for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done

sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y

sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo docker run hello-world

echo "Installing sqlite3..."
sudo apt update && sudo apt upgrade
sudo apt install sqlite3

echo "Installing tmux..."
sudo apt install tmux

echo "Creating 'Short' user..."
USER="Short"
PWD=$(openssl rand -base64 12)

sudo useradd -m "$USER"
echo "$USER:$PWD" | sudo chpasswd
sudo usermod -aG sudo "$USER"
sudo usermod -aG docker "$USER"

echo "User $USER created, password: $PWD"