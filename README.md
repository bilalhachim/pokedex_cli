# pokedex_cli
# Pokédex CLI

A command-line Pokédex application built in Go that allows you to explore the Pokémon world, discover Pokémon in different locations, and catch them for your collection.

## Features

- 🗺️ **Explore Locations**: Navigate through different areas in the Pokémon world
- 🔍 **Discover Pokémon**: Find which Pokémon inhabit specific locations
- ⚡ **Catch Pokémon**: Attempt to catch Pokémon with a 50% success rate
- 📊 **Inspect Pokémon**: View detailed stats of your caught Pokémon
- 📖 **Personal Pokédex**: Keep track of all your caught Pokémon
- 🚀 **Caching System**: Built-in caching for faster API responses

## Installation

### Prerequisites
- Go 1.19 or higher
- Internet connection (for PokéAPI access)

### Clone and Build
```bash
git clone https://github.com/bilalhachim/pokedex_cli.git
cd pokedex_cli
go build -o pokedexcli
```

### Run
```bash
./pokedexcli
```

## Usage

Once you start the application, you'll see the `Pokedex >` prompt. Here are the available commands:

### Commands

| Command | Description | Usage |
|---------|-------------|-------|
| `help` | Display help message | `help` |
| `exit` | Exit the Pokedex | `exit` |
| `map` | Show next 20 location areas | `map` |
| `mapb` | Show previous 20 location areas | `mapb` |
| `explore` | List Pokémon in a specific area | `explore <area_name>` |
| `catch` | Attempt to catch a Pokémon | `catch <pokemon_name>` |
| `inspect` | View details of a caught Pokémon | `inspect <pokemon_name>` |
| `pokedex` | List all caught Pokémon | `pokedex` |

### Example Session

```
Pokedex > help

Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
sunyshore-city-area
...

Pokedex > explore canalave-city-area
Exploring canalave-city-area
Found Pokemon:
 - tentacool
 - tentacruel
 - staryu
 - magikarp
 - gyarados

Pokedex > catch magikarp

Throwing a Pokeball at magikarp...
magikarp was caught!
You may now inspect it with the inspect command.

Pokedex > inspect magikarp
Name: magikarp
Height: 9
Weight: 100
Stats: 
  -hp: 20
  -attack: 10
  -defense: 55
  -special-attack: 15
  -special-defense: 20
  -speed: 80
Types: 
  - water

Pokedex > pokedex

Your Pokedex:
 - magikarp
```

## Project Structure

```
POKEDEXCLI/
├── .vscode/               # VS Code configuration
├── internal/
│   ├── pokecache_test.go  # Cache tests
│   └── pokecache.go       # Caching implementation
├── go.mod                 # Go module file
├── main.go                # Main application logic
├── pokedexcli             # Compiled binary
├── README.md              # Project documentation
├── repl_test.go           # REPL tests
├── repl.log               # Application logs
└── simple_repl.go         # REPL implementation
```

## Technical Details

### API Integration
- Uses the [PokéAPI](https://pokeapi.co/) for all Pokémon data
- Implements HTTP client with proper error handling
- JSON unmarshaling for API responses

### Caching System
- Built-in cache with 300ms duration
- Reduces API calls and improves performance
- Caches location areas and Pokémon data

### Core Features
- **REPL Interface**: Interactive command-line interface
- **Command Registry**: Modular command system
- **Random Catch Mechanic**: 50% success rate for catching Pokémon
- **Data Persistence**: Maintains caught Pokémon during session

### Data Structures
- `Pokemon`: Complete Pokémon information including stats and types
- `location_areas`: Location data with pagination
- `PokemonEncountersInLocationArea`: Pokémon encounters per location
- `config`: Navigation state management

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Future Enhancements

- [ ] Save/load Pokédex data to file
- [ ] Battle system between caught Pokémon
- [ ] More detailed Pokémon information
- [ ] Configuration options for catch rates
- [ ] Colored terminal output
- [ ] Search functionality for locations/Pokémon

## Dependencies

- Standard Go libraries only
- External API: [PokéAPI v2](https://pokeapi.co/docs/v2)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [PokéAPI](https://pokeapi.co/) for providing the comprehensive Pokémon database
- The Pokémon Company for creating the amazing world of Pokémon

---

**Happy Pokémon hunting!** 🎯
