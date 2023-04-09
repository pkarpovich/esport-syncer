# Esport Match Calendar

This repository contains a Golang application to fetch match data for a specific esports team and generate an iCalendar file that can be imported into various calendar applications.

## Features

- Fetches match data (team names, dates, tournament names, and results) for a specific esports team using the default provider (DotaProvider, which fetches Dota 2 matches from ggscore.com).
- Generates an iCalendar (.ics) file containing the match events.
- Allows you to set the refresh interval for the iCalendar events.
- Runs a simple HTTP server to serve the generated iCalendar file.

## Docker

This project can be run inside a Docker container. To do so, use the provided `Dockerfile` and `docker-compose.yml`.

### Environment Variables

The following environment variables can be set to configure the behavior of the `esport-syncer` service:

- `TEAM_ID`: The ID of the team for which you want to fetch the matches. (Default: Team Spirit ID)
- `CALENDAR_NAME`: The name of the generated iCalendar. (Default: "Esport matches")
- `CALENDAR_COLOR`: The color of the calendar events. (Default: "red")
- `CALENDAR_REFRESH_INTERVAL`: The refresh interval for the calendar events. (Default: "P1D")
- `PORT`: The port on which the server will listen. (Default: "1710")

## Customization

- You can change the default calendar name, default calendar color, default refresh interval, or default port by setting the respective environment variables.
- The default provider for fetching match data is `DotaProvider`. You can create your own provider by implementing the `Provider` interface and modifying the `main.go` file accordingly.

## Default Provider: DotaProvider
The default provider for fetching match data is DotaProvider. This provider fetches Dota 2 match data for a specific team from ggscore.com.

To use a different provider or create your own, implement the Provider interface in the providers package and modify the main.go file accordingly.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open a pull request or report any issues.

1. Fork the repository
2. Create a new branch for your feature or bugfix (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add your feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
