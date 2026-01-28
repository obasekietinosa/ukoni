Ukoni Client Roadmap

This document outlines the phased implementation plan for the Ukoni web client. It is designed to align with the API roadmap while prioritizing user-centric workflows.

⸻

Phase 0 – Foundation & Authentication

Goals
Establish the project structure and allow users to sign up/in and establish a session.

Tasks
- [ ] Initialize React + TypeScript (Vite) project
- [ ] Setup Routing (React Router)
- [ ] Setup API Client (Axios/Fetch with interceptors)
- [ ] Implement Sign Up (creates User + Auth Token)
- [ ] Implement Sign In
- [ ] Persistent Session Management
- [ ] Basic Layout (Header, Navigation)

Milestone
A user can log in and see a secure home screen.

⸻

Phase 1 – Core Product Catalog

Goals
Allow users to define "what" items are (Canonical Products). Corresponds to API Phase 1.

Tasks
- [ ] List Canonical Products
- [ ] Create Canonical Product (Name, Category)
- [ ] Edit/Delete Canonical Product
- [ ] View Product Details

Milestone
Users can manage the global dictionary of products.

⸻

Phase 2 – Household & Inventory Basics

Goals
Allow users to see what they have. Corresponds to API Phase 5 and initial setup.

Tasks
- [ ] Create/Select Household (Inventory)
- [ ] List Inventory Items (Inventory Products)
- [ ] View details of Inventory Items (Quantity, Variant info)
- [ ] Manual "Add to Inventory" (for initial population)

Milestone
Users can view and manually manage their current stock.

⸻

Phase 3 – Shopping Lists

Goals
Plan purchases. Corresponds to API Phase 3.

Tasks
- [ ] Create Shopping Lists
- [ ] Add Items to List (Search Canonical Products)
- [ ] Manage List Items (Notes, Preferred Outlet)
- [ ] Mark items as "Done" or delete

Milestone
Users can replace paper lists with the app.

⸻

Phase 4 – Sellers, Outlets & Transactions

Goals
Execute purchases and update inventory. Corresponds to API Phase 2 & 4.

Tasks
- [ ] Manage Sellers & Outlets (Create/List)
- [ ] "Go Shopping" Mode: Turn a List into a Transaction
- [ ] Record Transaction (Select Items, Price, Outlet)
- [ ] Handle Substitutions (Planned X, bought Y)
- [ ] Verify Inventory updates after transaction

Milestone
The loop is closed: Planning -> Purchasing -> Inventory.

⸻

Phase 5 – Consumption

Goals
Track usage. Corresponds to API Phase 6.

Tasks
- [ ] Log Consumption Event (Select Product, Quantity, Source)
- [ ] "Quick Eat" from Inventory list

Milestone
Inventory counts reflect actual usage.

⸻

Phase 6 – Household Management

Goals
Multiplayer mode. Corresponds to API Phase 8.

Tasks
- [ ] Generate Invitation Link/Token
- [ ] UI for Accepting Invitation
- [ ] Manage Members (List/Kick)

Milestone
Multiple users can manage the same household.

⸻

Phase 7 – Polish & Analytics

Goals
UX improvements and insights.

Tasks
- [ ] Dashboard/Home view (Recent activity, Low stock)
- [ ] Mobile responsiveness check
- [ ] Derived stats (Purchase frequency)

Milestone
The app feels complete and provides insight beyond raw data.
