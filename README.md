# snooker-assistant-server
This is the backend of [台球计分助手](https://github.com/superwhys/snooker-assistant)

It is designed with a microservices architecture, before start the servers, you need to install [consul](https://www.consul.io/)

## To install consul
### Mac
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/consul
```

### Linux
```bash
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install consul
```

# Dependency
Before you start the servers, run `go mod download` command to ensure the dependencies package installed.

# Servers
Snooker-Assistant-Server includes the following services in total:
- snooker-assistant-server
- snooker-game-server
- chinese-ball-game-server
