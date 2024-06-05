# Coffeecloud User Guide

### Introduction

> The Coffeecloud app provides integration and synchronization between Eliona and Coffeecloud services.

## Overview

This guide provides instructions on configuring, installing, and using the Coffeecloud app to manage resources and synchronize data between Eliona and [CoffeeCloud services](https://www.scanomat.com/coffeecloud/).

## Installation

Install the Coffeecloud app via the Eliona App Store.

## Configuration

The Coffeecloud app requires configuration through Eliona’s settings interface. Below are the general steps and details needed to configure the app effectively.

### Registering the app in Coffeecloud Service

Create an API-Key in the [Coffeecloud dashboard](https://coffeecloud.scanomat.com) to connect the Coffeecloud services from Eliona. All required credentials are listed below in the [configuration section](#configure-the-coffeecloud-app).

### Configure the Coffeecloud app

Configurations can be created in Eliona under `Apps > Coffeecloud > Settings` which opens the app's [Generic Frontend](https://doc.eliona.io/collection/v/eliona-english/manuals/settings/apps). Here you can use the appropriate endpoint with the POST method. Each configuration requires the following data:

| Attribute         | Description                                      |
|-------------------|--------------------------------------------------|
| `username`        | Username. Same as for the Coffeecloud dashboard. |
| `password`        | Password. Same as for the Coffeecloud dashboard. |
| `api_key`         | Generated API-Key (see below)                    |
| `url`             | URL of the Coffeecloud services.                 |
| `refreshInterval` | Interval in seconds for data synchronization.    |
| `requestTimeout`  | API query timeout in seconds.                    |
| `enable`          | Flag to enable or disable this configuration.    |
| `projectIDs`      | List of Eliona project IDs for data collection.  |

Example configuration JSON:

```json
{
  "username": "foo",
  "password": "secret",
  "api_key": "SeCrEt",
  "url": "https://coffeecloud.scanomat.com/rest/login",
  "refreshInterval": 60,
  "requestTimeout": 120,
  "enable": true,
  "projectIDs": [
    "10"
  ]
}
```

## Continuous Asset Creation

Once configured, the app starts Continuous Asset Creation (CAC). Discovered resources are automatically created as assets in Eliona.

The assets are group in Eliona like in the Coffeecloud dashboard. To avoid conflicts, the Global Asset Identifier is a manufacturer's ID prefixed with asset type name as a namespace.

## Additional Features

### Eliona dashboard templates

The app offers a predefined dashboard for ELiona that clearly displays the most important information. YOu can create such a dashboard under `Dashboards > Copy Dashboard > From App > Coffeecloud`.

