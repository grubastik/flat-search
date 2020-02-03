# flat-search

This is a small project for look-up a flat according to the options stored in config file in CZ.

## Installation

1. Build docker-compose
    ```bash
   docker-compose build
   ```

2. Run docker-compose storage
    ```bash
   docker-compose up storage
   ```

3. Run migration
    ```bash
   make migration-up
   ```

4. Run application
    ```bash
   docker-compose up flat-search
   ```

## TODO 
1. Minimize usage of pointers
2. tests
3. move sreality to separate repo
