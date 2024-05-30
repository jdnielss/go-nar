# Gonar

Gonar is a Go application to check the status of a project's quality gate on SonarQube.

## Installation

To install Gonar, you can use the following steps:

1. Clone this repository:

   ```
   git clone https://github.com/jdnielss/go-nar.git
   ```

2. Build the executable:

   ```
   go build -o gonar main.go
   ```

## Usage

You can use Gonar to check the status of a project's quality gate on SonarQube. Here's how to use it:

```
./gonar -n <PROJECT_NAME> -h <SONAR_URL> -k <SONAR_LOGIN>
```

- `<PROJECT_NAME>`: The name of the project on SonarQube.
- `<SONAR_URL>`: The URL of your SonarQube instance.
- `<SONAR_LOGIN>`: The login token for SonarQube authentication.

Example:

```
./gonar -n myproject -h https://sonar.example.com -k mytoken123
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
