from fabric import task
from invoke import run

@task
def devdb(c):
    run("docker run --name iris_postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=iris -e POSTGRES_DB=iris -p 0.0.0.0:5432:5432 -d postgres:11")