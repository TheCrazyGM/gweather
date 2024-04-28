# Weather App

This is a simple weather app written in Go that retrieves the current temperature of a city using the OpenWeatherMap API.

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

3. Create a `.env` file and add your OpenWeatherMap API key:

   ```
   OPENWEATHER_API_KEY=your-api-key
   ```

4. Build the app:

   ```
   go build
   ```

## Usage

To use the app, run the following command:

```
./gweather <city-name>
```

Replace `<city-name>` with the name of the city you want to get the weather for.

The app will retrieve the current temperature for the specified city and display it in Fahrenheit.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
