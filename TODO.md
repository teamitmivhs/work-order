# TODO List - IT Work Order System

This document outlines the current tasks and future improvements for the IT Work Order System.

## High Priority

-   [ ] **Fix Frontend JavaScript Errors:**
    -   Address `ReferenceError: Cannot access 'updateCurrentUser' before initialization`.
    -   Resolve `ReferenceError: workOrders is not defined`.
    -   Fix `Uncaught ReferenceError: workOrders is not defined at initializeWorkOrderTimers`.
    -   Review and refactor `src/static/assets/script.js` to fix scoping issues and improve code quality.
    -   Fix the 'POST' 500 Error when creating orders.

-   [ ] **Production-Ready Frontend:**
    -   Replace Tailwind CSS CDN with a proper build process (e.g., PostCSS or Tailwind CLI) to remove development warnings and optimize for production.
    -   Minify JavaScript and CSS for production.

## Backend and Infrastructure

-   [ ] **Integrate MySQL in Docker Compose:**
    -   Add a PostgreSQL service to `src/docker-compose.yml`.
    -   Ensure the Go backend connects to the PostgreSQL container.
    -   Provide a `init.sql` script to create the necessary tables. The existing SQL files (`dbwoit_executors.sql`, `dbwoit_members.sql`, `dbwoit_orders.sql`, `dbwoit_safetychecklist.sql`) should be reviewed and used.

-   [ ] **Complete User Authentication:**
    -   Implement the backend logic for user registration and login (`POST /api/register` and `POST /api/login`).
    -   Secure the relevant API routes with authentication middleware.
    -   Connect the `login.html` and `register.html` pages to the backend API.
    -   Implement session management (e.g., JWT).

## Medium Priority

-   [ ] **Implement Backend Tests:**
    -   Write unit tests for the controllers, repositories, and models.
    -   Write integration tests for the API endpoints.

-   [ ] **Implement Frontend Tests:**
    -   Write tests for the main frontend functionalities in `src/static/assets/script.js`.

-   [ ] **Improve Error Handling:**
    -   Improve error handling in the Go backend.
    -   Display user-friendly error messages on the frontend.

## Future Enhancements

-   [ ] **Email/Whatsapp Notifications:**
    -   Integrate a notification service to send alerts for new work orders or status changes.

-   [ ] **Advanced Reporting & Analytics:**
    -   Expand the Kaizen page with more detailed analytics and reporting features.

-   [ ] **Mobile-Friendly UI:**
    -   Further improve the responsiveness of the UI for mobile devices.
