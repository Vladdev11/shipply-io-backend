# Project Structure

This document outlines the structure of our Warehouse Management System (WMS) developed by Warehance.

The project has a relatively standard structure, with `service.go` serving as the entry point of the program. This file handles routing, server startup, and all other operations required when initializing the program.

We have four primary directories:

1. `Models`: This directory manages anything related to the underlying data structure and any communication with the database.

2. `Handlers`: This directory handles business logic and user-requested actions. It is divided into two subdirectories:

   - `OrganizationHandlers` for organizational users.
   - `ClientHandlers` for client-level users.

3. `API`: This directory manages any third-party integrations. While most third-party integrations have their own subdirectory within `API`, not all do.

4. `Util`: This directory contains miscellaneous functions, constant declarations, and so on.

This structure allows us to maintain a clean separation of concerns, making the codebase easier to navigate and manage.
