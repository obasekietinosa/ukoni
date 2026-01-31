# Ukoni Web Client Roadmap

This document outlines the phased implementation plan for the **Ukoni web client**, aligned with the backend API roadmap and the domain model.

It prioritises user workflows, clean architecture, and production-grade tooling.

---

## Tech Stack Strategy

- **Framework:** React + TypeScript (Vite)
- **Styling:** Tailwind CSS
- **State Management:** React Query (server state) + Context/Zustand (client UI state)
- **Networking:** Fetch API wrapped for type safety & interceptors
- **Testing:** Vitest (unit), React Testing Library (components), MSW (mock API), Playwright (E2E)
- **Quality:** ESLint, Prettier, Husky (pre-commit checks)
- **CI/CD:** GitHub Actions (lint, test, build)

---

## Phase 0 — Project Setup & Tooling

**Goals**
- Establish a robust, production-ready client infrastructure.

**Tasks**
- [x] Initialize React + TypeScript (Vite)
- [x] Setup Tailwind CSS with custom theme
- [x] Configure Vitest + React Testing Library
- [x] Setup MSW for API mocking
- [x] Add Playwright for E2E
- [x] Configure ESLint, Prettier, Husky
- [x] Create basic UI atoms (Button, Input, Icons, etc.)

**Milestone**
> All foundational tooling and styles are configured and tested.

---

## Phase 1 — Authentication & App Foundation

**Goals**
- Implement secure authentication and core app skeleton.

**Tasks**
- [x] Routing (React Router) + Error boundaries
- [x] Authentication UI (Sign In, Sign Up)
- [x] API client with auth handling (JWT / session)
- [x] Persistent session & secure storage
- [x] Layout skeleton (Sidebar, Header, responsive nav)

**Milestone**
> Users can log in and view a secure home screen.

---

## Phase 2 — Household Context (Scope)

> Aligns directly with backend inventory scoping.

**Goals**
- Enforce identity + household selection before browsing content.

**Tasks**
- [x] Inventory selection/creation flow on login
- [x] Implement `InventoryProvider` (React context)
- [x] Dashboard showing high-level inventory summary
- [x] Display user role & permission state

**Milestone**
> Authenticated users are scoped to a household/inventory context.

---

## Phase 3 — Product Catalog (Intent vs Execution)

> Maps to backend canonical products + variants endpoints.

**Goals**
- Manage products scoped to the active household.

**Tasks**
- [ ] Canonical product list & details
- [ ] Create/Edit canonical product (name, category)
- [ ] Product variant list per canonical product
- [ ] Create/Edit product variant (brand, size, unit)
- [ ] Product search + filters
- [ ] Integration tests for catalog flows

**Milestone**
> Users can define generic products and their purchasable variants.

**Missing**
- UI forms and integrations for catalog create/edit
- Search/filter with backend support

---

## Phase 4 — Inventory Management

> Tied to backend `inventory_products`.

**Goals**
- Manage current stock of product variants.

**Tasks**
- [ ] List inventory items with key details (quantity, unit)
- [ ] Manual addition/edit of inventory products
- [ ] Low stock indicators
- [ ] Virtualisation for long lists
- [ ] Responsive mobile support

**Milestone**
> Users have a real-time view of inventory levels.

**Missing**
- UI + API integration
- Working low-stock indicators
- Pagination/virtualisation

---

## Phase 5 — Shopping Lists & Planning

**Goals**
- Support flexible, user-driven shopping intent.

**Tasks**
- [ ] Create shopping lists
- [ ] Polymorphic item entry:
  - Add canonical product (generic intent)
  - Add specific product variant
- [ ] Preferred outlet selection (optional)
- [ ] Notes & quantity UI
- [ ] Suggestions based on low inventory
- [ ] Integration tests for list flows

**Milestone**
> Users can plan purchases at both generic and specific levels.

**Missing**
- UI for smart suggestions
- Backend transformation logic support

---

## Phase 6 — Transactions & Fulfilment

> Completes the “shopping cycle” with backend support.

**Goals**
- Execute purchases and update inventory.

**Tasks**
- [ ] Sellers & outlets management UI
- [ ] Transaction wizard (from list → transaction)
- [ ] Fulfilment logic handling:
  - Match list items → transaction items
  - Support substitutions
- [ ] UI to confirm bought vs planned
- [ ] Inventory update feedback
- [ ] Integration & E2E tests

**Milestone**
> Users complete shopping and inventory updates automatically.

**Missing**
- Transaction list and detail screen
- Fulfilment matching UI
- Substitute flows

---

## Phase 7 — Consumption

> Tracks usage and reduces inventory.

**Goals**
- Offer consumption logging and quick actions.

**Tasks**
- [ ] Record consumption event (canonical + optional variant)
- [ ] Quick consume action on inventory list
- [ ] History of consumption events
- [ ] UI feedback for consumption

**Milestone**
> Inventory reflects real usage; consumption history visible.

**Missing**
- UI patterns for consumption input
- Tests for consumption flows

---

## Phase 8 — Household Collaboration

> Same order as backend `invitations` + `memberships`.

**Goals**
- Invite and manage collaborators on inventories.

**Tasks**
- [ ] Generate invite link
- [ ] Invitation acceptance UI
- [ ] Member management (roles & permissions)
- [ ] Real-time UI updates on member changes

**Milestone**
> Multi-user support with roles and collaboration.

**Missing**
- Role-aware UI restrictions
- Real-time sync (optional)

---

## Phase 9 — Production Hardening & Analytics

**Goals**
- Polish, performance, accessibility, and insights.

**Tasks**
- [ ] Accessibility audit (ARIA, keyboard nav)
- [ ] Route lazy loading & code splitting
- [ ] Frontend monitoring (e.g., Sentry)
- [ ] Analytics dashboards:
  - Purchase frequency
  - Substitutions stats
  - Consumption vs inventory

**Milestone**
> App is polished, robust, and insight-ready.

**Missing**
- Analytics visualisations
- Performance metrics

---

## Cross-Cutting & Quality

A few execution priorities across all phases:

### Error & Loading States
- Consistent UI for API failures
- Retry / offline & skeleton loaders

### Testing
- Unit & component tests
- Integration tests with MSW
- E2E Playwright flows

### Accessibility
- Ensure WCAG standards
- Focus management & keyboard nav

### Documentation
- Storybook for UI components
- API contracts surfaced in client code

---

## Timeline Suggestions

| Phase | Weeks |
|-------|-------|
| 0 | 1 |
| 1 | 1–2 |
| 2 | 1 |
| 3 | 2–3 |
| 4 | 2–3 |
| 5 | 3–4 |
| 6 | 3–4 |
| 7 | 2 |
| 8 | 2 |
| 9 | 2–3 |

---

## Notes

This aligns closely with the backend API roadmap:

- Inventory scoping precedes product and transaction flows
- Intent (shopping lists) always precedes reality (transactions)
- Consumption and analytics sit on top of core CRUD workflows

As the backend API evolves, this roadmap should be updated to reflect changes, especially where new API capabilities are introduced (e.g., smart suggestions, analytics endpoints).

---
