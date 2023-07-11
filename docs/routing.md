# Routing in Our Application

This document provides an overview of how routing works within our application, from the point of route organization to executing user role-based handlers.

## Route Organization

Routes in our application are organized within `service.go`. This organization strategy is typically based on different application processes (such as picking, shipping, etc.) or underlying data models (like warehouses).

## Authorization Middleware

Before the execution of any handlers, an authorization middleware is invoked. This middleware reads the JSON Web Token (JWT) sent in the request, authorizes the user, and places the user into the request context. The user context can then be used throughout the rest of the handlers.

## User Router Handlers

Once a route is hit and a user is authorized, the relevant handler is executed. These handlers are primarily located within `user_router.go`, situated in the base of the `handlers` folder. Here, handlers are organized in the same manner as the routes.

Note that certain endpoints, such as webhooks for third-party connections, have their specific handlers separate from `user_router.go`.

## User Role-based Handlers

Inside the handlers in `user_router.go`, the user context is retrieved. Depending on the `user_type`, different actions are performed:

- If the `user_type` does not have permission to execute the requested action, the handler will return an error response.
- Otherwise, it will execute another handler specific to that `user_type`.

We maintain two sets of handlers within two packages: `OrganizationHandlers` and `ClientHandlers`:

- `organization_admin` and `organization_user` roles execute handlers from the `OrganizationHandlers` package. Occasionally, there may be different handlers for the `admin` and `user` roles.
- `client_admin` and `client_user` roles use handlers from the `ClientHandlers` package. As with the organization roles, there can sometimes be different handlers for the `admin` and `user` roles.

By structuring our handlers this way, we ensure a clear separation of roles and responsibilities, and that each user type can only perform the actions they are authorized for.
