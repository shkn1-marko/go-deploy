# go-deploy

A small deployment server for Linux. It listens for GitHub push webhooks, runs a build and deploy script for the pushed repository, and sends the result to [StatRep](https://github.com/shkn1-marko/status-report) as a push notification.

Each deployment is a folder with a `build.sh` and a `deploy.sh`. When GitHub reports a push, go-deploy runs both scripts and sends the outcome, along with the error and script output if something failed, to the registered phone through Firebase Cloud Messaging (FCM).

## Commands

The binary is called `gdep`.

| Command              | Description                                                         |
| -------------------- | ------------------------------------------------------------------- |
| `gdep make <name>`   | Create a deployment with empty `build.sh` and `deploy.sh` scripts   |
| `gdep remove <name>` | Delete a deployment and its scripts                                 |
| `gdep list`          | List all deployments                                                |
| `gdep deploy <name>` | Run a deployment by hand (no push notification is sent)             |
| `gdep start`         | Start the webhook server on port `9091`                             |

## StatRep integration

### Device registration

StatRep registers itself by calling:

```
POST /register-device
X-GDEP-Signature-256: sha256=<hex HMAC-SHA256 of the body using GDEP_REGISTER_SECRET>

{"installationID":"<FCM installation ID>"}
```

go-deploy keeps one device at a time. A new registration replaces the previous one, and the ID is saved to `GDEP_DEVICE_FILE` so it survives restarts.

### Deploy results

After every webhook-triggered deploy, go-deploy sends an FCM data message to the registered device:

| Field          | Description                                                     |
| -------------- | --------------------------------------------------------------- |
| `name`         | Deployment name                                                 |
| `buildStatus`  | `OK` or `FAILED`                                                |
| `deployStatus` | `OK`, `FAILED`, or `SKIPPED` (skipped when the build fails)     |
| `timestamp`    | Unix time in seconds                                            |
| `cause`        | Error from the failed script, only sent on failure              |
| `output`       | Script output, only sent on failure (last 3000 bytes)           |

## Setup

### 1. Build

Requires Go 1.27 or newer.

```bash
go build -o gdep .
```

### 2. Set environment variables

| Variable               | Description                                                      |
| ---------------------- | ---------------------------------------------------------------- |
| `GDEP_WEBHOOK_SECRET`  | Secret shared with the GitHub webhook                            |
| `GDEP_REGISTER_SECRET` | Secret shared with StatRep for device registration               |
| `GDEP_DEVICE_FILE`     | Path to the JSON file that stores the registered device          |
| `FCM_CREDENTIALS_FILE` | Path to your Firebase service account JSON key                   |

The Firebase project has to be the same one StatRep uses.

### 3. Create a deployment

```bash
gdep make my-repo
```

This creates `/etc/gdep/deployments/my-repo/` with two scripts. Fill them in with whatever your project needs, for example pulling the latest code in `build.sh` and restarting a service in `deploy.sh`.

- The deployment name **must match the GitHub repository name**. Pushes to repositories without a matching deployment are ignored.
- Both scripts run inside the deployment folder and must finish within 5 minutes.
- `deploy.sh` only runs if `build.sh` succeeds.

### 4. Add the GitHub webhook

In the repository's **Settings → Webhooks**, add a webhook with:

- **Payload URL:** `http://your-server:9091/webhook`
- **Content type:** `application/json`
- **Secret:** the value of `GDEP_WEBHOOK_SECRET`
- **Events:** just the push event

### 5. Start the server

```bash
gdep start
```
