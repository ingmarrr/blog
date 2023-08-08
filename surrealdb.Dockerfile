version: '3'

services:
  db:
    image: surrealdb/surrealdb:latest
    restart: always
    volumes: 
      - db:/var/lib/surrealdb/data
    ports:
      - 8000:8000
    command: start --user $DB_USER --pass $DB_PASSWORD file:/blogdata/blog.db
    env_file:
      - .env
volumes:
  db:
    driver: local

