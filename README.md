# Weather App

This is a command-line weather app written in Go that retrieves the current temperature and weather conditions of a city using the OpenWeatherMap API.

## Prerequisites

Before running the app, make sure you have the following:

- Go installed on your machine
- An API key from OpenWeatherMap. You can sign up for a free API key on their website.

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/thecrazygm/gweather.git
   ```

2. Change into the project directory:

   ```
   cd gweather
   ```


3. (Optional) Create a `.env` file and add your OpenWeatherMap API key (or set the environment variable directly):

   ```
   OPENWEATHER_API_KEY=your-api-key
   ```

   If you prefer not to use a `.env` file, set the API key in your shell instead:

   ```bash
   export OPENWEATHER_API_KEY=your-api-key
   ```

4. Build the app:

   ```
   go build
   ```

## Usage

To use the app, run the following command:

```
./gweather [city-name]
```

Replace `[city-name]` with the name of the city you want to get the weather for.

### Options

- `-m, --metric`: Display temperature in Celsius (default is Fahrenheit)

Note: API requests will time out after 10 seconds if the server does not respond.

### Examples

Get weather in Fahrenheit:
```
./gweather "New York"
```

Get weather in Celsius:
```
./gweather -m "London"
```

The app will display the current temperature and weather conditions (e.g., Clear, Cloudy, Rain) for the specified city.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
