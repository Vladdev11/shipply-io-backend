# User, Client, and Organization Relationships

This document outlines the relationships between Users, Clients, and Organizations within our software.

## Organizations and Clients

### Organizations

Organizations represent a Third-Party Logistics provider (3PL). An Organization can own multiple Clients.

### Clients

Clients represent customers of a 3PL. Each Client must belong to exactly one Organization.

## User Roles

User roles determine the level of access and the type of actions a user can perform. They are classified into `organization_admin`, `organization_user`, `client_admin`, and `client_user` based on their association with either an Organization or a Client.

### Organization Roles

#### Organization Admin

The `organization_admin` role belongs to Organizations. This user has additional privileges compared to the `organization_user` role and has access to information from all Clients that belong to their Organization.

#### Organization User

The `organization_user` role also belongs to Organizations. This role has fewer privileges than the `organization_admin` but has access to any information from any Client that belongs to their Organization.

### Client Roles

#### Client Admin

The `client_admin` role belongs to a Client. This user has more privileges than the `client_user` role and has access to all information that belongs to the Client they are associated with.

#### Client User

The `client_user` role belongs to a Client. This user has fewer privileges than the `client_admin` but has access to any information that belongs to the Client they are associated with.
