-- Create the 'parking' and 'app schemas for organizing related tables and functions.
CREATE SCHEMA app;
CREATE SCHEMA parking;

-- Set the transaction isolation level and timezone to ensure data consistency and uniform time representation.
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
SET TIMEZONE='UTC';



