# Newsletter service in golang

Cloud service template build in golang.

* docker
* AWS lightsail deployment
* AWS SQS
* tests
* postgres
* templates migration
* newsletter signup
* message queue

https://www.golang.dk/courses/build-cloud-apps-in-go

For deployment provide a 'containers.json' file:

```
{
  "app": {
    "image": "",
    "environment": {
      "LOG_ENV": "production",
      "HOST": "",
      "PORT": "8080",
      "DB_USER": "canvas",
      "DB_PASSWORD": "{{your db password}}",
      "DB_HOST": "{{your db host}}",
      "DB_NAME": "canvas",
      "BASE_URL": "{{your base URL}}",
      "POSTMARK_TOKEN": "{{your postmark token}}",
      "MARKETING_EMAIL_ADDRESS": "{{your marketing email address}}",
      "TRANSACTIONAL_EMAIL_ADDRESS": "{{your transactional email address}}",
      "AWS_ACCESS_KEY_ID": "{{the aws access key ID from the cloudformation output}}",
      "AWS_SECRET_ACCESS_KEY": "{{the aws secret access key from the cloudformation output}}",
      "ADMIN_PASSWORD": "{{your admin password}}"
    },
    "ports": {
      "8080": "HTTP"
    }
  }
}

```
