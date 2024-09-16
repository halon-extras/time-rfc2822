# Build instructions

```
export HALON_REPO_USER=exampleuser
export HALON_REPO_PASS=examplepass
docker compose -p halon-extras-time-rfc2822 up --build
docker compose -p halon-extras-time-rfc2822 down --rmi local
```