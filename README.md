# Bsati's Discord Bot

This project contains my private discord bot that is mainly used for managing birthdays of guild members.
The bot can be extended to different kinds of interactions and serves as a building block for general purpose bots.

## Installation

Clone the repository, add a `.env` file containg a `BOT_TOKEN=your_secret_bot_token`, and then execute
```bash
docker-compose build
docker-compose up
```
in the root directory to have a running instance. Additionally you have to adjust the `command` section for the `migrate` component in `docker-compose.yaml` and change the fromVersion (first argument) to "0". This is needed to initialize the database schema.
 
