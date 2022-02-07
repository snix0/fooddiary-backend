# Overview

Backend for FoodDiary application written in Golang (using Gin)

# Instructions

Deployment is managed fully by docker-compose. Install docker-compose and docker if you do not already have it.

There is a dependency with the repo at `fooddiary-frontend`.
Clone the project `fooddiary-frontend`. Then build it with `yarn install` to generate the `build` folder with a production build of the frontend.
Then run `docker build -t fdfrontend .` at the project root to build the docker image for the frontend.

Run `docker-compose up` at the root of the project directory to fully deploy this application.

After docker-compose brings up all the containers, visit the proxy container at port 80 to view the application
