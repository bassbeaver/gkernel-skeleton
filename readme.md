Golang web application skeleton based on [bassbeaver/gkernel](https://github.com/bassbeaver/gkernel) framework.

This is simple web application with three pages: public index page, login page, private page.

To serve static files nginx is used. As sessions store Redis is used. Nginx and Redis run as Docker containers so you need
Docker and Docker-compose to run this example. 

How to start application:
1. Copy `config.yaml.example` to `config.yaml`
2. Run `docker-compose up -d` from `./docker` folder.
3. Build and run Go application. Application will be available on http://localhost:50080