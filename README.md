**github-proxy**

A simple REST API that interacts with github to manage repositories and pull requests.

---

## **Features**

- **Health Check**: Verify service availability.
- **Repository Management**: List, create, or delete repositories.
- **Pull Request Insights**: List pull requests for a specific repository.

---

## **Getting Started**

### **Run the Server Locally**

To start the service locally, use Docker Compose:

```bash
docker compose up -d
```

This will expose the REST API at `http://localhost:8080` with the following endpoints:

#### **Endpoints**
| Method | Endpoint                       | Description                                        |
|--------|--------------------------------|----------------------------------------------------|
| **GET**    | `/health`                      | Check service status.                              |
| **GET**    | `/repository`                  | List all repositories for the token user.          |
| **POST**   | `/repository`                  | Create a new repository.                           |
| **DELETE** | `/repository/{name}`           | Delete an existing repository.                     |
| **GET**    | `/pull-request/{owner}/{repo}` | List pull requests for a given `owner/repository`. |

---

### **Environment Configuration**

Before starting the service, configure the required environment variables:

```bash
GH_AUTH_TOKEN=<your GitHub access token with read and write permissions>
SRV_PORT=8080
```

- **`GH_AUTH_TOKEN`**: GitHub access token (requires `repo` scope for accessing repositories).
- **`SRV_PORT`**: The port on which the service runs (default is `8080`).

---

## **CI/CD Pipeline**

The pipeline is designed with four jobs:

1. **Lint**: runs static code analysis.
2. **Test**: executes Gherkin defined unit tests for the service.
3. **Build**:
    - builds a Docker image.
    - publishes the image to the ghcr.
    - signs the image using Cosign with a private key stored in repo secrets.
4. **Deploy**:
    - creates a Minikube cluster for deployment.
    - installs a policy manager (Kyverno) and applies a custom image validation policy.
    - deploys the built and signed Docker image to the Minikube cluster.

---

## **Possible improvements**

- deepen the reach of the unit tests as they only cover success scenarios
- better assert a successful deployment on minikube by exposing the deployment with a service and making api requests
- refactor http handlers as they share common behaviour which is replicated
---
