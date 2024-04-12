# GoLang API with MySQL, Docker, JWT, and Testing

 This project embodies a robust backend architecture adhering to Clean Architecture principles, ensuring modularity, testability, and maintainability.

## Deployment

Our API is live on an EC2 instance. Access the Swagger documentation [here](http://ec2-18-195-89-34.eu-central-1.compute.amazonaws.com/api/v1/swagger/index.html).

## Postman Documentation

Interact effortlessly with our API using Postman. Refer to the [Postman Documentation](https://documenter.getpostman.com/view/19898564/2sA3Bhfabi) and [Request Collection](https://api.postman.com/collections/19898564-424a85cf-5089-411a-ae51-e60d2d33df88?access_key=PMAT-01HV9J2MRXNZ6X03PDB0363WB1).

## Deployment

Here's how to get started with the API:

1. **Environment Setup**: Copy `.env.example` and create `.env`, adjusting variables as necessary.
2. **Initialize Docker**: Run the command `docker-compose up`.
3. **Swagger Documentation**: Generate Swagger documentation:
   ```bash
   make swagger
   ```
4. **Database Setup**: Utilize Docker to create the required database. Apply migrations with after creating docker:
   ```bash
   make migrate-up
   ```

5. **Testing**: Run comprehensive tests with:
   ```bash
   make test
   ```

6. **Local Deployment**: Launch the API locally (without Docker) using, in this scenario user needs to create db manually:
   ```bash
   make run
   ```

7. **Logging and Metrics**: Access Prometheus metrics via `/metrics` and view logs on `/logs`.

## Additional Insights
- **Backend Architecture**: Our API is architected following Clean Architecture principles, emphasizing separation of concerns and testability.
  
- **GoDoc Support**: While GoDoc support is integrated, it requires initialization for enhanced documentation readability. Can start godoc with ```godoc -http=localhost:8080```


