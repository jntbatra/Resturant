# Restaurant Management API - Postman Collection

This repository contains Postman collection and environment files for testing the Restaurant Management API.

## Files

- `Restaurant_API.postman_collection.json` - Complete Postman collection with all API endpoints
- `Restaurant_API.postman_environment.json` - Environment variables for the API
- `Restaurant_API_Test_Data.json` - Test data file for collection runner iterations

## Setup

1. **Import the Collection:**
   - Open Postman
   - Click "Import" button
   - Select "File" tab
   - Choose `Restaurant_API.postman_collection.json`

2. **Import the Environment:**
   - Click "Import" button
   - Select "File" tab
   - Choose `Restaurant_API.postman_environment.json`
   - Select the "Restaurant API Environment" from the environment dropdown

3. **Configure Environment Variables:**
   - Update the `base_url` variable if your API is running on a different host/port
   - Update database variables if needed for your setup

## API Endpoints Covered

### Sessions Management
- List all sessions and active sessions
- Create, read, update, delete sessions
- Change session table
- Get sessions by table

### Tables Management
- List, create, read, delete tables

### Menu Management
- List menu items (with pagination and filtering)
- Create, read, update, delete menu items
- Get menu items by category

### Category Management
- List, create, update, delete categories
- Get category by name or ID
- Get category ID by name

### Order Management
- List, create, update orders
- Add items to orders
- Get order items
- Get orders by session
- Get order items by session

## Testing

Each request includes automated tests that:
- Verify correct HTTP status codes
- Validate response structure
- Set environment variables for subsequent requests
- Ensure data integrity

## Usage Flow

The collection follows a logical workflow for testing the complete restaurant API:

1. **Create Table** - Sets up a table for dining
2. **Create Category** - Creates a menu category
3. **Create Menu Item** - Adds an item to the menu under the category
4. **Create Session** - Starts a dining session for the table
5. **Create Order** - Creates an order for the session
6. **Add Item to Order** - Adds menu items to the order
7. **Get Order Items** - Verifies items were added
8. **Update Order Status** - Changes order status (pending → preparing → served)
9. **Get Orders by Session** - Lists all orders for a session
10. **Update Session Status** - Changes session status (active → completed)
11-13. **List endpoints** - Verify all data was created correctly

## Running the Collection

### Option 1: Manual Testing (Individual Requests)

1. Import the collection and environment as described above
2. Run requests individually in order
3. Each request sets collection variables for the next request

### Option 2: Automated Collection Runner

1. Open Postman and import the collection
2. Click on the collection name in the sidebar
3. Click the "Run" button to open Collection Runner
4. Select the environment
5. **Important:** Select the data file `Restaurant_API_Test_Data.json`
6. Set iterations to match the number of data rows (5)
7. Click "Run Restaurant Management API - Workflow"

The runner will:
- Execute all 13 requests for each iteration
- Use different data from each row in the test data file
- Create multiple tables, categories, menu items, sessions, and orders
- Validate each step with automated tests
- Show console logs for created resource IDs

## What Gets Tested

Each request includes automated tests that:
- Verify correct HTTP status codes (200, 201, 204)
- Validate response structure and required fields
- Set collection variables for subsequent requests
- Log created resource information to console
- Ensure data integrity across the workflow

## Expected Results

After running the collection:
- Multiple tables, categories, menu items, sessions, and orders will be created
- All relationships will be properly established
- Status updates will work correctly
- List endpoints will show the accumulated data
- Console logs will show all created resource IDs

## Troubleshooting

- Ensure your API server is running on `http://localhost:8080`
- Check the Postman console for detailed error messages
- Verify database is properly set up and migrations are run
- If tests fail, check the API logs for server-side errors

### Table Management
- List, create, get, delete tables

### API Documentation
- Access to Swagger UI

## Running Tests

### Individual Requests
- Select any request in the collection
- Click "Send" to execute
- View response in the response panel

### Collection Runner
To run all requests automatically:

1. Click "Runner" button in Postman
2. Select "Restaurant API" collection
3. Choose "Restaurant API Environment"
4. Optionally select the test data file: `Restaurant_API_Test_Data.json`
5. Set iterations (number of times to run)
6. Click "Run Restaurant API"

### Manual Test Flow

For a complete test scenario:

1. **Setup Tables:**
   - Create tables using "Create Table"

2. **Create Session:**
   - Create a dining session for a table

3. **Manage Menu:**
   - Create categories
   - Create menu items

4. **Place Orders:**
   - Create an order for the session
   - Add menu items to the order
   - Update order status

5. **Complete Session:**
   - Update session status to completed

## Environment Variables

The collection uses the following variables:

- `base_url`: API base URL (default: http://localhost:8080)
- `table_id`: Table ID for operations
- `session_id`: Session ID for operations
- `order_id`: Order ID for operations
- `menu_item_id`: Menu item ID for operations
- `category_id`: Category ID for operations
- Various other variables for request bodies

## Notes

- Make sure the API server is running before testing
- Database should be set up and migrations applied
- Some requests require IDs from previous operations (use environment variables to store them)
- The collection includes proper request bodies and headers for each endpoint