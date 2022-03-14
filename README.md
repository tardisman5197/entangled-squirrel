# entangled-squirrel

Squirrels who connect to each other.

## Building

1. Have [Go](https://go.dev/) installed

    ```
    go version
    ```

    This command will verify that Go is installed.

2. Navigate to the command folder (`/cmd/es`)

    ```
    cd ./cmd/es
    ```

3. Run Go build

    ```
    go build
    ```

    This should create a binary is the same folder with the name `es`.

## Setup

When deploying onto a pi zero w there are a few steps that need to be done.

1. Connect your Pi to your wifi.

    This [guide](https://howchoo.com/g/ndy1zte2yjn/how-to-set-up-wifi-on-your-raspberry-pi-without-ethernet) has most of the steps to do this.

    Basically edit the `wpa_supplicant.conf` file to follow this structure.

    ```
    country=GB # Your 2-digit country code
    ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev

    network={
        ssid="YOUR_NETWORK_NAME"
        psk="YOUR_PASSWORD"
        key_mgmt=WPA-PSK
    }
    ```

2. Ensure the `squirrels.service` file is correct

    An example of the service file can be found in the repository.
    The command being run on start should only have the `-squirrels` flag in use, the rest are for debugging purposses.

3. Move the service file into the correct location and enable the service

    ```
    cp ./squirrel.service /etc/systemd/system/
    ```

    ```
    sudo systemctl daemon-reload
    ```

    ```
    sudo systemctl start squirrel
    ```

    ```
    sudo systemctl enable squirrel
    ```

    These commands should start the squirrel service and allow the program to start on boot.

4. Ensure the list of konwn squirrels are up to date.

    These can be found in `./cmd/es/squirrels.txt`. It should be a list of urls seperated by new lines.

5. Ensure [Duck DNS](https://www.duckdns.org/install.jsp) is setup

    The `duck.sh` script needs to include the correct domain for your squirrel (make sure you are not using the same domain as a differnt squirrel)

After these steps everything should be setup and your squirrel should be accessable.
This can be tested by pressing the right arm of the squirrel and the eyes should flash.

## Usage

*Note: it can take a while to start once plugged in*

Once setup the squirrel should flash its eyes once when your button is pressed.

If a differnt squirrel is pressed the squirrel will flash its eyes 3 times.
