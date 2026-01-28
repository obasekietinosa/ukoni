Ukoni Implementation Roadmap

This document outlines a high-fidelity, phased plan for implementing the remaining Ukoni functionality, based on:
	•	The current state of the codebase
	•	The agreed domain model and design principles
	•	The household-manager point of view

Anything not currently implemented is explicitly flagged.

⸻

Guiding Principles (Anchor These Throughout)
	•	Household-centric model: all data is scoped to a household unless explicitly global
	•	Canonical vs Variant separation: intent vs execution is first-class
	•	Derivable state over stored state: avoid flags like is_purchased
	•	Auditability by default: soft deletes + activity log everywhere
	•	Intent ≠ Outcome: shopping lists vs transactions vs consumption are distinct

⸻

Phase 0 – Baseline & Hardening (Completed)

Goals

Ensure the existing system is safe to extend and aligned with the domain model.

Tasks
	•	[x] Audit existing entities vs agreed domain model
	•	[x] Add soft delete support (deleted_at) to core tables
	•	[x] Add activity_log table
	•	actor (user)
	•	household
	•	entity_type / entity_id
	•	action (created, updated, deleted, consumed, etc)
	•	timestamp
	•	[x] Add middleware / hooks to auto-log mutations

Milestone

System is auditable and future changes won’t lose history

⸻

Phase 1 – Product Model Completion (Completed)

Canonical Products

Purpose: represent what something is, independent of brand or size.
	•	[x] canonical_products
	•	[x] id
	•	[x] name (“Rapeseed Oil”)
	•	[x] category (optional, future-proofing)
	•	[x] created_at / updated_at / deleted_at

Product Variants

Represents what was actually purchased.
	•	[x] Ensure variants:
	•	[x] belong to a canonical product
	•	[x] encode size + unit
	•	[x] encode brand / seller-specific info

Milestone

Clear separation between intent (canonical) and execution (variant)

⸻

Phase 2 – Sellers & Outlets (Completed)

Sellers

Represents the business entity.
	•	[x] sellers
	•	[x] id
	•	[x] name (“Lidl”, “Morrisons”)
	•	[x] created_at / deleted_at

Outlets

Represents a concrete place (physical or online).
	•	[x] outlets
	•	[x] id
	•	[x] seller_id
	•	[x] name (optional – “Lidl Tottenham”)
	•	[x] address (physical or URL)
	•	[x] created_at / deleted_at

Explicitly out of scope
	•	Seller inventory
	•	Prices over time (future concern)

Milestone

Transactions can reference where something was bought without modelling store inventory

⸻

Phase 3 – Shopping Lists (Completed)

Shopping Lists
	•	[x] shopping_lists
	•	[x] id
	•	[x] household_id
	•	[x] created_by
	•	[x] created_at
	•	[x] last_actioned_at (suggested name; open to rename)

Shopping List Items

Key design decisions:
	•	No is_purchased
	•	State is derived from linked transaction items
	•	Polymorphic product reference
	•	[x] shopping_list_items
	•	[x] id
	•	[x] shopping_list_id
	•	[x] product_type (canonical | variant)
	•	[x] product_id
	•	[x] preferred_outlet_id (nullable)
	•	[x] notes (optional)
	•	[x] created_at / deleted_at

Linking to Transactions
	•	[x] transaction_items gain optional shopping_list_item_id
	•	[x] Shopping lists do not link directly to transactions

Milestone

Lists model intent, transactions model reality, and substitutions are analysable

⸻

Phase 4 – Transactions & Fulfilment (Completed)

Transactions
	•	[x] Ensure transactions:
	•	[x] belong to a household
	•	[x] optionally reference an outlet
	•	[x] contain transaction items only

Transaction Items
	•	[x] Support substitutions naturally via shopping_list_item link
	•	[x] Support partial fulfilment (multiple transaction items per list item over time)

Milestone

Real-world shopping behaviour is faithfully modelled

⸻

Phase 5 – Inventory & Inventory Products (Completed)

Inventory Products

Represents what the household currently has.
	•	[x] inventory_products
	•	[x] id
	•	[x] household_id (as inventory_id)
	•	[x] product_variant_id
	•	[x] quantity
	•	[x] unit
	•	[x] created_at / deleted_at

Inventory Updates
	•	[x] Created from transaction items
	•	Reduced via consumption events (Pending Phase 6)

Milestone

Inventory reflects reality without manual syncing

⸻

Phase 6 – Consumption Tracking (Completed)

Consumption Events

Loosely modelled by design.
	•	[x] consumption_events
	•	[x] id
	•	[x] household_id
	•	[x] canonical_product_id
	•	[x] quantity (nullable)
	•	[x] unit (nullable)
	•	[x] consumed_at
	•	[x] source (manual, recipe, estimate)

Design Notes
	•	Canonical-level by default
	•	Variant-level attribution can be inferred later

Milestone

Ukoni knows what was consumed and when, without pretending to know too much

⸻

Phase 7 – Units & Conversions (Deferred Logic)

⚠️ Explicitly deferred

Planned approach:
	•	Store quantities + unit as-is
	•	Maintain a unit graph (not strict base units)
	•	Allow probabilistic / best-effort conversions

No schema changes yet.

Milestone

System remains flexible and honest about uncertainty

⸻

Phase 8 – Permissions & Invitations

Invitations (⚠️ missing)
	•	invitations
	•	id
	•	household_id
	•	email
	•	role
	•	token
	•	expires_at
	•	accepted_at

Roles
	•	Manager
	•	Member (future: child, guest)

Milestone

Households can grow safely

⸻

Phase 9 – Analytics & Derivations (Later)

All derived, no new core state:
	•	Purchase frequency per canonical product
	•	Substitution rates
	•	Waste estimation (inventory vs consumption)

⸻

Final Notes

This roadmap intentionally:
	•	Avoids premature optimisation
	•	Prioritises correctness over convenience
	•	Models reality as it is, not as we wish it were

Each phase is independently valuable and shippable.
