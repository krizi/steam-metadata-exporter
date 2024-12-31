```yaml
api_key: yyy
user_id: xxx
features:
  playtime:
    enabled: true
    schedule: "15s"
    api_url: "https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/"
  owned_games:
    enabled: true
    schedule: "2h"
    api_url: "https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/"
log_level: "info"
metrics_port: ":8080"
```