# Use the official PostgreSQL image as the base image
FROM postgres:latest

# Set environment variables for PostgreSQL
ENV POSTGRES_USER="postgres"
ENV POSTGRES_PASSWORD="mypassword"
# ENV POSTGRES_DB="astella"

# Copy any custom SQL scripts to initialize the database
COPY ./conf/database/v1.0.0.sql /docker-entrypoint-initdb.d/

# Expose the PostgreSQL port
EXPOSE 5432

CMD ["postgres"]