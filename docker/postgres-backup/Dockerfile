FROM alpine:latest

# Install required packages
RUN apk add --no-cache postgresql-client dcron bash

# Copy backup script and .env file
COPY backup.sh /usr/local/bin/backup.sh
COPY .pgpass .

# Make the script executable
RUN chmod +x backup.sh

# Set up backup directory
RUN mkdir -p /backups

# Copy crontab
COPY crontab /etc/crontabs/root

# Run cron in foreground with log level 8
CMD ["crond", "-f", "-d", "8"]