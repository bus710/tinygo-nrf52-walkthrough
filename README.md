# Tinygo on nRF52 Walkthrough

This repo briefly introduces how to build Go based BLE app and flash to nRF52 devices.

<br/><br/>

## References

- Getting started: https://tinygo.org/getting-started/linux/
- nRF52832 doc: https://tinygo.org/microcontrollers/pca10040/
- nRF52832 API: https://tinygo.org/microcontrollers/machine/pca10040/
- nRF52840 doc: https://tinygo.org/microcontrollers/pca10056/
- nRF52840 API: https://tinygo.org/microcontrollers/machine/pca10056/
- TinyGo BLE: https://github.com/tinygo-org/bluetooth

<br/>


## Prerequisites

- Ubuntu 20.04 (or newer) installed X86_64 PC
- Go SDK 1.14 (or newer)
- Git and other typical utils for development
- VSCODE
- nRF52832 (= PCA10040, s132) or nRF52840 (= PCA10056, s140) DK board

<br/>

## Install prerequisites

```sh
$ sudo apt install \
    libncurses5 \
    openocd \
    minicom
```

## Install TinyGo

```sh
# Download the latest version
$ cd ~/Downloads
$ curl --silent https://github.com/tinygo-org/tinygo/releases/latest \
    | grep -oP '(?<=/v)[^">]+' \
    | xargs -I '{}' wget https://github.com/tinygo-org/tinygo/releases/download/v'{}'/tinygo_'{}'_amd64.deb 

# Install and check - the version may vary
$ sudo dpkg -i tinygo_*.deb
$ tinygo version
tinygo version 0.16.0 linux/amd64 (using go version go1.14.4 and LLVM version 10.0.1)
```

<br/>

## Install JLink tools

Unlike Adafruit's Feather nRF52840 or ItsyBitsy-nRF52840 (those have UF2 bootloader and SoftDevice already), nRF52 DK boards have a built-in JLink. To utilize the JLink, some tools (nrfjprog, mergehex and friends) should be installed.

NordicSemi provides files that include the tools:
- https://www.nordicsemi.com/Software-and-tools/Development-Tools/nRF-Command-Line-Tools
- Latest version: https://www.nordicsemi.com/-/media/Software-and-other-downloads/Desktop-software/nRF-command-line-tools/sw/Versions-10-x-x/10-12-1/nRFCommandLineTools10121Linuxamd64tar.gz

```sh
# Download
$ cd ~/Downloads
$ mkdir nrf-tools
$ wget https://www.nordicsemi.com/-/media/Software-and-other-downloads/Desktop-software/nRF-command-line-tools/sw/Versions-10-x-x/10-12-1/nRFCommandLineTools10121Linuxamd64tar.gz
$ tar xvf nRFCommandLineTools10121Linuxamd64tar.gz -C nrf-tools

# Install 
$ cd nrf-tools
$ sudo dpkg -i JLink_Linux_V688a_x86_64.deb
$ sudo dpkg -i nRF-Command-Line-Tools_10_12_1_Linux-amd64.deb
```

<br/>

## Build examples

First, connect a target DK board to the host.

Then, follow below commands to build an example and flash the target
```sh
# Clone
$ git clone https://github.com/tinygo-org/bluetooth.git
$ cd bluetooth

# Build an example - this is not really needed
$ tinygo build -o heartrate -target=pca10040-s132v6 ./examples/heartrate

# Flash SoftDevice - do this only once
$ nrfjprog -f nrf52 --eraseall
$ nrfjprog -f nrf52 --program s132_nrf52_6.1.1/s132_nrf52_6.1.1_softdevice.hex

# Flash the example built
$ tinygo flash -target=pca10040-s132v6 ./examples/heartrate
```

NordicSemi's Connect mobile app provides extremely useful features to check BLE communications. 

<br/>

## Test from RPI

- https://elinux.org/RPi_Bluetooth_LE
- https://elinux.org/images/3/32/Doing_Bluetooth_Low_Energy_on_Linux.pdf
- https://www.argenox.com/library/bluetooth-low-energy/using-raspberry-pi-ble/
- https://developpaper.com/the-basic-method-of-using-bluetooth-function-in-linux-system/
- https://www.makeuseof.com/manage-bluetooth-linux-with-bluetoothctl/
- https://punchthrough.com/creating-a-ble-peripheral-with-bluez/
- https://punchthrough.com/android-ble-guide/

Some packages should be installed:
```sh
$ sudo apt install -y \
    bluez \
    bluetooth \
    bluez-tools
$ sudo apt install -y \
    libusb-dev \
    libreadline-dev \
    libglib2.0-dev \
    libudev-dev \
    libdbus-1-dev \
    libical-dev
$ btmon -v
$ sudo systemctl enable bluetooth.service
$ sudo systemctl start bluetooth.service
$ sudo usermod -aG bluetooth $USER
```

Some commands might be helpful to scan, connect, pair, and bond:
```sh
# hciconfig, hcitool, and gatttool are low level commands

$ hciconfig
$ sudo hciconfig hci0 up
$ sudo hcitool lescan -i hci0
$ gatttool -I -b MAC-address
[BT MAC][LE]> char-read-uuid 00002a00-0000-1000-8000-00805f9b34fb

# rfcomm and rfkill are more of wireless system control in general
# useful to see if BT is disabled

$ rfkill list
$ rfkill unblock

# btmgmt and bluetoothctl are high level commands

$ btmgmt
[mgmt] \# info
[mgmt] \# select hci0
[hci0] \# power up
[hci0] \# info

$ bluetoothctl 
[bluetooth] \# show
[bluetooth] \# scan on
[bluetooth] \# discoverable on
[bluetooth] \# pair MAC
[bluetooth] \# connect MAC
```

<br/>

## Setup VSCODE for further work

Go to the example and open with VSCODE:
```sh
$ cd examples/heartrate
$ code .
```

Press CTRL+SHIFT+B. 

Then, a dialog will appear 
- "No build task to run found. Configure build task..." => press enter
- "Create tasks.json from template" => press enter
- Select Others

Attach below snippet:
```json
{
    "version": "2.0.0",
    "type": "shell",
    "echoCommand": true,
    "cwd": "${workspaceFolder}",
    "tasks": [
        {
            "label": "echo",
            "type": "shell",
            "command": "echo Hello"
        },
        {
            "label": "Build",
            "command": "tinygo flash -target=pca10040-s132v6 .",
            "group": {
                "kind": "build",
                "isDefault": true
            },
            "problemMatcher": [
                "$go"
            ]
        },
    ]
}
```

Next time when CTRL+SHIFT+B is pressed, tinygo will build and flash the project to the target board.

<br/>

## Debug log

The example code has some lines with the **println** function and the messages go to the USB-serial connection with the host. To monitor the messages:
```sh
$ sudo minicom -b 115200 \
    -o -D /dev/ttyACM0 # the port name can be different! 
```

<br/>
