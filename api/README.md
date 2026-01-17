# ukoni api
This is the backend for Ukoni, a household inventory management app.
This API will be built in Golang and will support all the core functionality for the Ukoni app.


## Data model
```mermaid
erDiagram
    USERS {
        uuid id PK
        string email
        string name
        string password_hash
        datetime created_at
        datetime deleted_at
    }

    INVENTORIES {
        uuid id PK
        string name
        uuid owner_user_id FK
        datetime created_at
        datetime deleted_at
    }

    INVENTORY_MEMBERSHIPS {
        uuid id PK
        uuid inventory_id FK
        uuid user_id FK
        string role "admin | editor | viewer"
        datetime invited_at
        datetime deleted_at
    }

    INVITATIONS {
        uuid id PK
        uuid inventory_id FK
        string email
        string role
        uuid invited_by_user_id FK
        string status "pending | accepted | revoked | expired"
        datetime created_at
        datetime accepted_at
        datetime expires_at
    }

    CANONICAL_PRODUCTS {
        uuid id PK
        string name
        text description
        datetime created_at
        datetime deleted_at
    }

    PRODUCTS {
        uuid id PK
        uuid canonical_product_id FK
        string brand
        string name
        text description
        uuid category_id FK
        datetime created_at
        datetime deleted_at
    }

    PRODUCT_VARIANTS {
        uuid id PK
        uuid product_id FK
        string variant_name
        string sku
        string unit
        datetime deleted_at
    }

    PRODUCT_CATEGORIES {
        uuid id PK
        string name
        uuid parent_category_id FK
        datetime deleted_at
    }

    SELLERS {
        uuid id PK
        string name
        string type "chain | independent | online"
        datetime created_at
        datetime deleted_at
    }

    OUTLETS {
        uuid id PK
        uuid seller_id FK
        string name
        string channel "physical | online"
        string address "nullable"
        string website_url "nullable"
        datetime created_at
        datetime deleted_at
    }

    INVENTORY_PRODUCTS {
        uuid id PK
        uuid inventory_id FK
        uuid product_variant_id FK
        decimal quantity
        datetime last_updated
        datetime deleted_at
    }

    TRANSACTIONS {
        uuid id PK
        uuid inventory_id FK
        uuid outlet_id FK
        uuid created_by_user_id FK
        datetime transaction_date
        decimal total_amount
        datetime deleted_at
    }

    TRANSACTION_ITEMS {
        uuid id PK
        uuid transaction_id FK
        uuid product_variant_id FK
        decimal quantity
        decimal price_per_unit
        datetime deleted_at
    }

    CONSUMPTION_EVENTS {
        uuid id PK
        uuid inventory_id FK
        uuid canonical_product_id FK
        uuid created_by_user_id FK
        decimal quantity_consumed "nullable"
        string unit "nullable"
        string note
        datetime consumed_at
        datetime deleted_at
    }

    ACTIVITY_LOGS {
        uuid id PK
        uuid inventory_id FK
        uuid user_id FK
        string action
        json metadata
        datetime created_at
    }

    SHOPPING_LISTS {
        uuid id PK
        uuid inventory_id FK
        string name
        datetime created_at
        datetime last_updated_at
        datetime deleted_at
    }

    SHOPPING_LIST_ITEMS {
        uuid id PK
        uuid shopping_list_id FK
        string target_type "canonical_product | product_variant"
        uuid target_id
        uuid preferred_outlet_id "nullable"
        datetime created_at
        datetime deleted_at
    }

    %% Relationships
    USERS ||--o{ INVENTORY_MEMBERSHIPS : participates_in
    INVENTORIES ||--o{ INVENTORY_MEMBERSHIPS : has_members
    USERS ||--o{ INVENTORIES : owns

    INVENTORIES ||--o{ INVITATIONS : has
    USERS ||--o{ INVITATIONS : sends

    CANONICAL_PRODUCTS ||--o{ PRODUCTS : groups
    PRODUCT_CATEGORIES ||--o{ PRODUCTS : categorises
    PRODUCT_CATEGORIES ||--o{ PRODUCT_CATEGORIES : parent_of

    PRODUCTS ||--o{ PRODUCT_VARIANTS : has

    INVENTORIES ||--o{ INVENTORY_PRODUCTS : contains
    PRODUCT_VARIANTS ||--o{ INVENTORY_PRODUCTS : tracked_as

    INVENTORIES ||--o{ TRANSACTIONS : records
    USERS ||--o{ TRANSACTIONS : creates
    SELLERS ||--o{ OUTLETS : operates
    OUTLETS ||--o{ TRANSACTIONS : involved_in

    TRANSACTIONS ||--o{ TRANSACTION_ITEMS : contains
    PRODUCT_VARIANTS ||--o{ TRANSACTION_ITEMS : purchased_as

    SHOPPING_LISTS ||--o{ SHOPPING_LIST_ITEMS : contains
    SHOPPING_LIST_ITEMS ||--o{ TRANSACTION_ITEMS : fulfilled_by
    OUTLETS ||--o{ SHOPPING_LIST_ITEMS : preferred_source

    INVENTORIES ||--o{ CONSUMPTION_EVENTS : records
    CANONICAL_PRODUCTS ||--o{ CONSUMPTION_EVENTS : consumed
    USERS ||--o{ CONSUMPTION_EVENTS : logs

    INVENTORIES ||--o{ ACTIVITY_LOGS : logs
    USERS ||--o{ ACTIVITY_LOGS : performs
```

## Features
### User management
#### Authentication
Users can sign up with email and password as well as log in via thier set passwords. We should be able to remember user devices, just to reduce the length of time before they log in again.

#### Invitations
A user can invite a user to an inventory. They can invite existing users, in which case the existing users get a notification to accept the invite. The user can also invite new users, in which case those emails get a notification to sign up. Once they sign up, they can then accept the invitation to join an inventory.

### Product Management
#### Products
These are the items that are bought and consumed in the household.

Product Variants are variations on a specific product. E.g Sainsbury's Olive Oil 1 Litre and Sainsbury's Olive Oil 3 litres are variants of a Sainsbury's Olive Oil product.

A canonical product refers to a generic instance of an item independent of brand or variety. For example, olive oil is a canonical product that can have specific brands such as Sainsbury's Olive Oil or Tesco's Olive Oil.

These are useful as a way to interact with products without necessarily caring about the specific brand, e.g we bought 3L of olive oil over the last week.

#### Product Categories
These are a useful way to group related products together. E.g Seasonings & Condiments or Baking.

### Sellers & Outlets
A seller is the business entity that a product was purchased from. This seller can have one or more outlets. The outlet is the place the actual purchase was made and could be a physical store or could be online.

### Transactions
These are as implied. A transaction is made up of multiple transaction items which themselves record how much of a product variant was bought and at how much.

### Shopping Lists
A shopping list reflects intent to purchase some items. We should add shopping list items to the list which are linked either to a product variant or to a canonical product.

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Postgres Client (optional, for manual inspection)

### Running the Project

1. **Start the Database**
   ```bash
   docker-compose up -d
   ```
   This starts Postgres on localhost:5432 with:
   - User: `etin`
   - Password: `etin`
   - DB: `ukoni`

2. **Run Migrations**
   ```bash
   DATABASE_URL="postgres://etin:etin@localhost:5432/ukoni?sslmode=disable" go run cmd/migrate/main.go up
   ```

3. **Run the Server**
   ```bash
   export DATABASE_URL="postgres://etin:etin@localhost:5432/ukoni?sslmode=disable"
   go run cmd/api/main.go
   # OR
   ./bin/api
   ```
   Server listens on port 8080.

3. **Run the Seeder (Optional)**
   ```bash
   go run cmd/seeder/main.go
   ```
   Creates a user: `test@example.com` / `password123`.

## Current Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/signup` | Create a new user |
| POST | `/login` | Get JWT token |

### Example Curl

**Signup:**
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com", "password": "securepassword"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "securepassword"}'
```
