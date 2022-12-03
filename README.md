# QR Generator

Microservice for generating QR codes for products.

Available endpoints:
- `/live`: Liveliness check
- `/ready`: Readiness check
- `/generate`: Generates a QR code for a product

Branches:
- `main`: Contains stable, tagged releases
- `dev`: Contains latest development version

## Setup/installation

To run the microservice using Docker Compose run `make compose`.

To see other available options run `make help`.

## License

Multimo is licensed under the [GNU AGPLv3 license](LICENSE).
