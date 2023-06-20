Gin Golang Boilerplate
==============

Introduction
------------

This project is a Boilerplate for Restful-API using Gin-Gonic and GORM for faster development by providing several common functions, such as:
- Service repository for every database models
- JWT Auth
- Logging
- Swagger Documentation
- Env loader
- Response Templating
- Migration & Seeder
- Etc.

Prerequisites
-------------

Make sure you have the following installed on your machine:

*   Go (version 1.20 or higher)
*   Node.js (for the `dev` target in the Makefile with Nodemon)

Installation
------------

1.  Clone the repository:

    git clone https://github.com/dimasbagusussilo/gin-golang-boilerplate.git

2.  Change into the project directory:

    cd gin-golang-boilerplate

3. Install project dependencies:

   make deps

4. Rename file ```.env.example``` to ```.env```

Usage
-----

### Build

To compile the packages, run the following command:

    make build

### Run

To build and run the project in development mode, use the following command:

    make run

This will compile the packages and start the application.

### Clean

To clean the project and remove previous builds, run:

    make clean

### Development Mode

To run the project in development mode using nodemon (requires Node.js), use the following command:

    make dev

This will automatically restart the application whenever changes are made.

### Generate API Documentation V1

To generate API documentation using Swag, run:

    make generatedocs1

This will create the documentation in the `docs/v1/` directory.

### Additional Migration

To generate additional database migration file run:

    make migrationcreate

This will create the migration file in the `database/migrations/` directory.

We can run migration with [Golang Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
we can install with this command (mac):

    brew install golang-migrate

Then we can run up migration with (we need setup credentials in `Makefile`): 
    
    make migrateup

Same as up migration, we can run down migration with:

    make migratedown

Contributors
------------

- Dimas Bagus Susilo  <dimasbagussusilo@gmail.com>

Contributing
------------

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

License
-------

This project is licensed under the [MIT License](LICENSE).# gin-golang-boilerplate
