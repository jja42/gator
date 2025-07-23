# Gator - Boot.Dev Project

This program fetches and processes RSS feeds using Go and PostgreSQL.

## Prerequisites

Make sure you have both **Go** and **PostgreSQL** installed before proceeding.

## Installation

Install  by running:
```bash
go install github.com/jja42/gator@latest
```

## Setup
You must create a config file in your home directory named ".gatorconfig.json"
This file must have the following structure:
{
  "db_url": "protocol://username:password@host:port/database",
  "current_user_name": "your-username"
}

## Use
Run the program by using the prefix gator
Use the following commands to interact with the program
ie: gator [command] [arguments] (Remove Brackets when you fill in the blanks)

Commands
- register [username] : sets your current username to the username and creates a user in the database
- login [username] : sets your current username to the username provided that it exists in the database
- reset: resets the users table which cascades unto other tables, effectively reseting the database
- users: returns a list of registered users
- addfeed [name] [url]: adds a feed to the database. also creates a feedfollow in the database between the current user and the feed
- feeds: returns all feeds in the database
- follow [url]: retrieves a feed from the database by url and then creates a feedfollow between the current user and the feed
- following: using the current user, returns a list of all feeds followed by referencing the feedfollows table
- unfollow [url]: retrieves a feed from the database by url and then deletes any feedfollow between the current user and the feed
- agg: aggregates posts from feeds into the posts table, creating new records. marks feeds and goes to oldest marked feed next to pull data.
- browse [limit (optional)]: retrieves a number of posts (limit) from feeds that the current user is following. if limit is not provided, defaults to 2
