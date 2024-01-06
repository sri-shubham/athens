# API

This package contains REST api handlers for all entities and search

About Implementation:
- The CRUD operations use Generics to save time.
- All Handlers are bound by Interfaces for loose coupling and easy swapping out of services if need be
- The setupRoutes.go file maps URLs to handlers

