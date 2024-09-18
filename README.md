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
   go build -o go-nar main.go
   ```

## Usage

You can use Gonar to check the status of a project's quality gate on SonarQube. Here's how to use it:

```
./go-nar -n "PROJECT_NAME" -h "SONAR_URL" -k "SONAR_TOKEN" -g "GITLAB_URL" -t "GITLAB_TOKEN" -p "PROJECT_ID" -m "MERGE_REQUEST_IID"
```

- <PROJECT_NAME>: The name of the project on SonarQube.
- <SONAR_URL>: The URL of your SonarQube instance.
- <SONAR_LOGIN>: The authentication token for SonarQube.
- <GITLAB_URL>: The URL of your self-hosted or cloud GitLab instance.
- <GITLAB_TOKEN>: The personal access token for GitLab authentication, with api and read_api permissions.
- <PROJECT_ID>: The project ID in your GitLab repository. If using GitLab CI, use the predefined variable $CI_PROJECT_ID.
- <MERGE_REQUEST_IID>: The merge request ID. If you trigger the job only for merge requests, use the predefined variable $CI_MERGE_REQUEST_IID.

Example:

### Bash
```bash
./go-nar -n "PROJECT_NAME" -h "SONAR_URL" -k "SONAR_TOKEN" -g "GITLAB_URL" -t "GITLAB_TOKEN" -p "PROJECT_ID" -m "MERGE_REQUEST_IID"
```

### Docker 
```bash
docker run --rm jdnielss/go-nar:latest -n "PROJECT_NAME" -h "SONAR_URL" -k "SONAR_TOKEN" -g "GITLAB_URL" -t "GITLAB_TOKEN" -p "PROJECT_ID" -m "MERGE_REQUEST_IID"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
