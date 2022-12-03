# QR Generator

Microservice for generating QR codes for products.

Available endpoints:
- `/qr`: returns a png image of a QR code describing a products (currently lists number of products in the database)

Branches:
- `main`: Contains stable, tagged releases
- `dev`: Contains latest development version

## Setup/installation

To run the microservice using Docker Compose run `make compose`.

To see other available options run `make help`.

## License

Multimo is licensed under the [GNU AGPLv3 license](LICENSE).
