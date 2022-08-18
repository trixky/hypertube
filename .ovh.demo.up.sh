# ============= clone project
echo "################################ clone project"
git clone https://github.com/trixky/hypertube.git

# ============= docker
echo "################################ update repos"
sudo apt-get update -y

echo "################################ install basic tools"
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# https://docs.docker.com/engine/install/ubuntu/

echo "################################ install docker # 1"
sudo mkdir -p /etc/apt/keyrings

echo "### 5"
echo "################################ install docker # 2"
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo "################################ install docker # 3"
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "################################ install docker # 4"
sudo apt-get update -y

echo "################################ install docker # 5"
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# https://docs.docker.com/engine/install/linux-postinstall/

# sudo groupadd docker

echo "################################ configure docker # 1"
sudo usermod -aG docker $USER
echo "################################ configure docker # 2"
sudo gpasswd -a $USER docker

# ============= env

cd hypertube

echo "################################ copy .env"
cp .env.example .env; \
cp postgres/.env.example postgres/.env; \
cp redis/.env.example redis/.env; \
cp client/.env.example client/.env; \
cp api-position/.env.example api-position/.env; \
cp api-streaming/.env.example api-streaming/.env; \
cp api-auth/.env.example api-auth/.env; \
cp api-scrapper/.env.example api-scrapper/.env; \
cp api-media/.env.example api-media/.env; \
cp api-picture/.env.example api-picture/.env; \
cp api-user/.env.example api-user/.env; \
cp streaming-proxy/.env.example streaming-proxy/.env; \
cp tmdb-proxy/.env.example tmdb-proxy/.env

echo "################################ source"
source env.sh

# ============= start project
echo "################################ start docker grp"
newgrp docker << END
echo "################################ start docker compose project in demo mode"
docker compose -f docker-compose.demo.yml up -d
END