Golang web application skeleton based on [bassbeaver/gkernel](https://github.com/bassbeaver/gkernel) framework.

This is simple web application with three pages: public index page, login page, private page.

To serve static files nginx is used. For sessions store Redis is used. Nginx and Redis run as Docker containers so you need
Docker and Docker-compose to run this example. 

How to start application:
1. Copy `config.yaml.example` to `config.yaml`
2. Run `docker-compose up -d` from `./docker` folder.
3. Build and run Go application. Application will be available on http://localhost:50080

Login / Password for the private page is: `login1 / password1`

More detailed documentation: [https://bassbeaver.github.io/gkernel-docs/](https://bassbeaver.github.io/gkernel-docs/)