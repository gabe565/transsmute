## transsmute



### Synopsis

Build RSS feeds for websites that don't provide them.

```
transsmute [flags]
```

### Options

```
      --address string                 Listen address (default ":3000")
      --docker-enabled                 Docker API enabled (default true)
      --dockerhub-password string      DockerHub password for private repos
      --dockerhub-username string      DockerHub username for private repos
      --ghcr-app-id int                GitHub app ID
      --ghcr-installation-id int       GitHub installation ID
      --ghcr-password string           GitHub personal access token for ghcr.io
      --ghcr-private-key string        GitHub app private key
      --ghcr-private-key-path string   GitHub app private key file path
      --ghcr-username string           GitHub username for ghcr.io
  -h, --help                           help for transsmute
      --kemono-enabled                 Kemono API enabled (default true)
      --kemono-hosts stringToString    Kemono API hosts, where the key is the URL prefix and the value is the host (default [kemono=kemono.su])
  -v, --version                        version for transsmute
      --youtube-enabled                YouTube API enabled (default true)
      --youtube-key string             YouTube API key
```

