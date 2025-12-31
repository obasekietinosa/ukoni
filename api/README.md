# ukoni api
This is the backend for Ukoni, a household inventory management app.
This API will be built in Golang and will support all the core functionality for the Ukoni app.


## Data model
```mermaid
erDiagram
    CANONICAL_PRODUCT {
        uuid id PK
        string name
        string category
        datetime created_at
    }

    PRODUCT_VARIANT {
        uuid id PK
        uuid canonical_product_id FK
        string brand
        string size
        string unit
        decimal typical_price
    }

    SHOPPING_LIST {
        uuid id PK
        string name
        datetime created_at
    }

    SHOPPING_LIST_ITEM {
        uuid id PK
        uuid shopping_list_id FK
        string target_type  "canonical_product | product_variant"
        uuid target_id
        decimal requested_quantity
        string notes
    }

    TRANSACTION {
        uuid id PK
        datetime occurred_at
        string retailer
    }

    TRANSACTION_ITEM {
        uuid id PK
        uuid transaction_id FK
        uuid product_variant_id FK
        uuid shopping_list_item_id FK "nullable"
        decimal quantity
        decimal price_paid
    }

    CANONICAL_PRODUCT ||--o{ PRODUCT_VARIANT : "has"
    SHOPPING_LIST ||--o{ SHOPPING_LIST_ITEM : "contains"

    SHOPPING_LIST_ITEM ||--o{ TRANSACTION_ITEM : "fulfilled by"
    TRANSACTION ||--o{ TRANSACTION_ITEM : "contains"

    PRODUCT_VARIANT ||--o{ TRANSACTION_ITEM : "purchased as"
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
