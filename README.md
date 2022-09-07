# Moon

Bootleg Logstash for use with Koiwai and Lnx

## Features

- Sequentially indexes all posts in your DB to Lnx
- Updates modified posts by itself
- Almost ACID

## Usage

- Edit config.example.toml to fit your use case
- Either export the ```MOON_CONFIG``` environment variable to point it to your configuration file or leave it as config.toml in the project root
- Install golang 1.18 or above
- Run ```go build .``` on the project root to build your executable
- Run it
