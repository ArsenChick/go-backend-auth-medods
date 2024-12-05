# Simple Go Auth service

## Endpoints

- `/new` - request a new pair of tokens.

  ### Request parameters

  - **Type**: POST
  - **Body**:

    ```json
    {
        "guid": "<valid-guid>"
    }
    ```
  - **Returns**:

    ```json
    {
        "access": "<jwt-access-token>",
        "refresh": "<base64-refresh-token>"
    }
    ```
    or
    ```json
    {
        "message": "<error-message>"
    }
    ```

- `/refresh` - refresh pair of tokens. They are needed to be sent via headers.

  ### Request parameters

  - **Type**: GET
  - **Headers**:

    - *Access-Token*: \<contains access token value\>
    - *Refresh-Token*: \<contains refresh token value\>

  - **Returns**:

    ```json
    {
        "access": "<jwt-access-token>",
        "refresh": "<base64-refresh-token>"
    }
    ```
    or
    ```json
    {
        "message": "<error-message>"
    }
    ```

## How to run

1. Have [Docker Compose](https://docs.docker.com/compose/install/) installed.
2. Run `docker compose up`.