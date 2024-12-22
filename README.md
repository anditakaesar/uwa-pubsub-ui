# UWA PubSub UI

UWA PubSub UI is a web-based user interface for managing and monitoring Pub/Sub messaging systems. It provides an intuitive way to visualize and interact with your Pub/Sub Emulator.

## Summary

This application allows users to:
- Monitor Pub/Sub topics and subscriptions
- Publish messages to topics (TBD)
- View and manage messages in subscriptions (TBD)
- Configure and manage Pub/Sub settings (TBD)

## Using the Docker Image

To run the UWA PubSub UI using Docker, follow these steps:

1. **Pull the Docker image:**
    ```sh
    docker pull anditakaesar/uwa-pubsub-ui:latest
    ```

2. **Run the Docker container:**
    ```sh
    docker run -d \
      -p 8080:8080 \
      --name uwa-pubsub-ui \
      -e ProjectID=your_project_id \
      -e PUBSUB_EMULATOR_HOST=your_emulator_host \
      anditakaesar/uwa-pubsub-ui:latest
    ```

3. **Environment Variables:**

    - `ProjectID`: The Local Cloud project ID for your Pub/Sub setup.
    - `PUBSUB_EMULATOR_HOST`: The host address of the Pub/Sub emulator (if using an emulator).

4. **Access the UI:**
    Open your web browser and navigate to `http://localhost:8080` to access the UWA PubSub UI.

## Environment Variables

- `ProjectID`: (Required) The ID of your Local Cloud project.
- `PUBSUB_EMULATOR_HOST`: (Optional) The address of the Pub/Sub emulator, if you are using one.

Make sure to replace `your_project_id` and `your_emulator_host` with your actual project ID and emulator host address.
