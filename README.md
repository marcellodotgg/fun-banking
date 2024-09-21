# Fun Banking

The official Fun Banking rewrite with Go, HTMX, and SQLite. This inspiration came from needed
server-side rendering to boost performance and SEO, as well as needing to keep our
cost down so that we can be in business; thus the "boring" architecture.

# Getting Started

1. Create your `.env`

```.env
DATABASE_URL=fun_banking.db
JWT_SECRET=something_secret
EMAIL_USERNAME=bytebury@gmail.com
EMAIL_PASSWORD=PASSWORD
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
WEBSITE_URL=http://localhost:8080
COOKIE_URL=localhost:8080 # or .fun-banking.com
BUILD_HASH=some_generated_hash
GOOGLE_ANALYTICS_ID=G-XXXXXXXX
```

# Local Environment

We recommend that you use Air. It's a dev server the will listen to your files.

```shell
air init
```

Make sure that you update this in your `.air.toml` file

```diff
+ cmd = "go build -o ./tmp/main cmd/fun-banking/main.go"
```

Now, you can just run the following command and your server will automatically launch in watch mode.

```shell
./dev.sh
```

# Deployment

Deployment is as easy as running the docker file via the `deploy.sh` script

```shell
./deploy.sh
```
