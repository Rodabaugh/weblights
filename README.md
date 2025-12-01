# weblights

HTTP Controlled WS2811 Lights on Raspberry Pi. Built using go, tmepl, and HTMX.

## Why

This project was create to enable control over the lights in my cubicle and on my Chrismas tree. 

## About the Project

The entire project is run on a Raspberry Pi. When on your local network, requests are fulfilled nearly instantly. For public access, you can setup a VPN between the Pi and a VPS running a reverse proxy. With this setup, you can also enable SSL using certbot. This adds some latency, but is still fairly responsive.

## Technical Features
 
- The interface is served using templ, with HTMX for some interactivity
- Presets are stored in a PostgreSQL database, running on the Pi
- Presets support a dynamic number of colors
- While an admin mode does not exist, you can manually update the protected column in the presets table to write protect a preset
- Attempts to update the status of the lights are stored in the logs table, as well as if the request was successful

# Install

## Prerequisites

### Parts
1. A Raspberry Pi running Raspberry Pi OS Lite
2. WS2811 Addressable LED lights (ensure they that are the 5V model to work with the Pi)
3. A 5V power supply with enough amperage to power your lights.
4. The required misc parts to connect everything

### Install Tools

Connect to your pi, and run the following commands.
```
sudo apt install golang git cmake
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Install rpi_ws281x

The rpi_ws281x library is required to enable support for the ws2811 lights. Use these steps to build the library, or read the official docs.

``` bash
git clone https://github.com/jgarff/rpi_ws281x
cd rpi_ws281x/
mkdir build
cmake -D BUILD_SHARED=OFF -D BUILD_TEST=ON ..
sudo make install
```

## Clone the Project

``` bash
cd ~
git clone https://github.com/Rodabaugh/weblights
cd weblights
```

## Setup DB

### Install PostgreSQL

``` bash
sudo apt install postgresql
sudo passwd postgres
# Enter password for postgres user
sudo -u postgres psql
```

### Create the DB and Setup the User

``` sql
CREATE DATABASE weblights;
\c weblights
ALTER USER postgres PASSWORD '(YOUR DB PASSWORD HERE)';
exit
```

### Build the Database Schema using Goose

``` bash
cd sql/schema/
~/go/bin/goose postgres "postgres://postgres:(YOUR DB PASSWORD)@localhost:5432/weblights" up
cd ../..
```

## Configure your Enviroment

Next, create a .env file in the project dir and configure as shown below.

``` bash
nano .env
```

Be sure to enter the number of LEDs that you have, and your database creds.

```
PLATFORM=prod
NUM_LEDS=50 (or your number of LEDs)
DB_URL="postgres://postgres:(YOUR DB PASSWORD)@localhost:5432/weblights"
```

## Build and Install
``` bash
make build
sudo make install
```

## Test

Check to see the service is running using ```systemctl status weblights.service```

On another device go to http://(The IP for your Pi):8080. You should see the weblights UI.

## Public access

From here, I recommend setting up a wireguard VPN between your Pi and a public VPS. Then you can configure an nginx reverse proxy on the VPS and certbot, to provide public access to the weblights UI. This is out of scope for this guide, but information on this is easily found.