Ukoni Client Roadmap

This document outlines the phased implementation plan for the Ukoni web client. It is designed to align with the API roadmap while prioritizing user-centric workflows and **production-grade engineering standards**.

**Tech Stack Strategy:**
*   **Framework:** React + TypeScript (Vite)
*   **Styling:** Tailwind CSS (Utility-first, responsive by default)
*   **State Management:** React Query (Server state) + Context/Zustand (Client state)
*   **Testing:** Vitest (Unit), React Testing Library (Component), Playwright (E2E)
*   **Quality:** ESLint, Prettier, Husky (Pre-commit hooks)

⸻

Phase 0 – Foundation & Engineering Rigor

Goals
Establish a robust, production-ready project structure with CI/CD checks, authentication, and design system foundations.

Tasks
- [ ] Initialize React + TypeScript (Vite) project
- [ ] Setup **Tailwind CSS** with a custom theme configuration (colors, typography)
- [ ] Setup **Vitest & React Testing Library** for unit/component tests
- [ ] Setup **Playwright** for E2E testing
- [ ] Configure ESLint, Prettier, and Husky (pre-commit checks)
- [ ] Setup Routing (React Router) with **Error Boundaries**
- [ ] Implement API Client (Axios/Fetch) with robust error handling and interceptors
- [ ] Implement Sign Up & Sign In (JWT handling, secure storage)
- [ ] Create Reusable UI Components (Button, Input, Card) using Tailwind
- [ ] Setup **Layout Skeleton** (Responsive Sidebar/Navbar)

Milestone
A secure, tested, and styled "Hello World" environment is ready for feature work.

⸻

Phase 1 – Household Context (The Scope)

Goals
Establish the tenancy scope. Since data is **scoped to an inventory/household**, this context must be established before managing content.

Tasks
- [ ] Create/Select Household (Inventory) flow upon login
- [ ] Implement `InventoryProvider` to manage the active scope globally
- [ ] Dashboard View: High-level summary of the active household
- [ ] Manage Memberships (View current user's role)

Milestone
User is authenticated and anchored to a specific Inventory context.

⸻

Phase 2 – Scoped Product Catalog

Goals
Manage the definition of products **within the current inventory**.
*Note: Canonical products are inventory-specific, not global.*

Tasks
- [ ] List Canonical Products (Paginated, Searchable)
- [ ] Create Canonical Product (Name, Category) - *Optimistic Updates*
- [ ] Edit/Delete Canonical Product
- [ ] View Product Details
- [ ] **Integration Test:** Verify products from Inventory A do not appear in Inventory B

Milestone
Users can build their household's specific dictionary of items.

⸻

Phase 3 – Inventory Management

Goals
Manage the physical stock (Variants) linked to the catalog.

Tasks
- [ ] List Inventory Items (Inventory Products) with "Low Stock" indicators
- [ ] View Details (Quantity, Unit, Variant info)
- [ ] Manual "Add to Inventory" (linking to Canonical Product)
- [ ] Implement **Virtualization** for long inventory lists (Performance)

Milestone
Users have a real-time view of their stock.

⸻

Phase 4 – Shopping Lists & Planning

Goals
Plan purchases based on inventory needs.

Tasks
- [ ] Create Shopping Lists
- [ ] Add Items to List (Search scoped Canonical Products)
- [ ] Manage List Items (Notes, Preferred Outlet)
- [ ] **Smart Suggestions:** Suggest items based on low inventory (future proofing)
- [ ] Interactive "Shopping Mode" UI (Mobile optimized, large tap targets)

Milestone
Digital planning replaces paper lists.

⸻

Phase 5 – Transactions & Loop Closure

Goals
Execute purchases and automatically update inventory.

Tasks
- [ ] Manage Sellers & Outlets
- [ ] Transaction Wizard: Convert List -> Transaction
- [ ] Handle Substitutions (UI to swap "Planned" for "Bought")
- [ ] **Form Validation:** Complex validation for multi-step transaction forms
- [ ] Verify Inventory increments automatically

Milestone
The "Shopping Cycle" is complete.

⸻

Phase 6 – Consumption

Goals
Track usage to close the loop on inventory counts.

Tasks
- [ ] Log Consumption Event (Source: Manual/Recipe)
- [ ] "Quick Actions" on Inventory List (Swipe to consume)
- [ ] Visual Feedback for stock reduction

Milestone
Inventory reflects reality.

⸻

Phase 7 – Household Management & Collaboration

Goals
Multiplayer features.

Tasks
- [ ] Generate Invitation Link/Token
- [ ] Invitation Acceptance UI
- [ ] Manage Members (RBAC: Admin vs Member)

Milestone
Multiple users collaborate on the same inventory.

⸻

Phase 8 – Production Hardening & Analytics

Goals
Ensure the app is robust, accessible, and performant.

Tasks
- [ ] **Accessibility Audit:** Ensure WCAG compliance (ARIA labels, keyboard nav)
- [ ] **Performance Tuning:** Code splitting, lazy loading routes
- [ ] **Monitoring:** Setup Sentry for frontend error tracking
- [ ] Analytics Dashboard: Purchase frequency, Waste estimation

Milestone
Staff-level polish and operational excellence.
