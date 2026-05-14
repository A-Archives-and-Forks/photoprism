// Module-scope cache for the labels and albums typeahead lists used by
// the sidebar info panel, the batch-edit dialog, and the edit-dialog
// labels tab. Each consumer would otherwise re-fetch the same dataset
// on mount, costing repeated GET /api/v1/{labels,albums}?count=<cap>
// round-trips for the same browser session.
//
// API: getLabels() / getAlbums() each return a Promise<Array> of raw
// model instances. Concurrent callers share the same in-flight promise
// so exactly one request fires for any number of consumers. WS-driven
// invalidation evicts on labels.updated / labels.deleted / albums.updated
// / albums.deleted; the next read after invalidation triggers a fresh
// fetch. Album deletion currently only emits config.updated (no
// dedicated channel), so we also evict albums on config.updated to
// catch that case — the cost is one extra fetch per unrelated config
// change, which is bounded.
//
// The cache stays a client-side preload. Server-side debounced
// typeahead (search-as-you-type) is the right answer once libraries
// genuinely exceed the cap and lives in its own future proposal; this
// module would become its orchestrator at that point.
import Album from "model/album";
import Label from "model/label";
import $event, { subscribeEntityActions } from "common/event";

// Pragmatic ceiling shared by every consumer. Power users with more
// than CAP labels or albums see a console.warn and a truncated list;
// the long-term answer for those libraries is server-side debounced
// typeahead, not raising the cap further.
export const CAP = 5000;

const state = {
  labels: { data: null, fetch: null },
  albums: { data: null, fetch: null },
};

function evict(field) {
  const slot = state[field];
  if (!slot) return;
  slot.data = null;
  slot.fetch = null;
}

function fetchLabels() {
  return Label.search({ count: CAP, order: "name", all: true }).then((resp) => {
    const models = Array.isArray(resp?.models) ? resp.models : [];
    if (models.length === CAP) {
      console.warn(`Label.search returned ${CAP} results — list may be truncated.`);
    }
    return models;
  });
}

function fetchAlbums() {
  return Album.search({ count: CAP, order: "name", type: "album" }).then((resp) => {
    const models = Array.isArray(resp?.models) ? resp.models : [];
    if (models.length === CAP) {
      console.warn(`Album.search returned ${CAP} results — list may be truncated.`);
    }
    return models;
  });
}

function get(field, fetcher) {
  const slot = state[field];
  if (slot.data) return Promise.resolve(slot.data);
  if (slot.fetch) return slot.fetch;
  slot.fetch = fetcher()
    .then((data) => {
      slot.data = data;
      slot.fetch = null;
      return data;
    })
    .catch((err) => {
      slot.fetch = null;
      throw err;
    });
  return slot.fetch;
}

// Public surface — call-site agnostic. Consumers map the returned
// model arrays to whatever shape they need at the boundary.
export const typeaheadCache = {
  getLabels: () => get("labels", fetchLabels),
  getAlbums: () => get("albums", fetchAlbums),
  evictLabels: () => evict("labels"),
  evictAlbums: () => evict("albums"),
  clear: () => {
    evict("labels");
    evict("albums");
  },
};

// One hierarchical subscriber per entity namespace, filtered to the
// standard mutation verbs by subscribeEntityActions. The cache is
// permissive — the action only matters as a "something changed"
// signal, so we ignore payload and just evict.
//
// Backend emit sites at the time of writing:
//
//   labels.created  — `event.EntitiesCreated("labels", …)` from
//                     `entity/label.go` FirstOrCreateLabel.
//   labels.updated  — `PublishLabelEvent(StatusUpdated, …)` from
//                     `internal/api/labels.go` (rename, like, unlike).
//   labels.deleted  — `event.EntitiesDeleted("labels", …)` from
//                     `batch_labels.go`.
//   albums.created  — `event.PublishUserEntities("albums", EntityCreated, …)`
//                     from `entity/album.go`. The WS writer strips the
//                     `user.<uid>.` prefix so the client receives
//                     `albums.created`.
//   albums.updated  — same entity path with EntityUpdated, plus
//                     `PublishAlbumEvent(StatusUpdated, …)` from REST
//                     handlers (albums.go, links.go, import.go,
//                     users_upload.go).
//   albums.deleted  — `event.EntitiesDeleted("albums", …)` from
//                     `entity/album.go` Album.Delete.
//
// A future `labels.edited` / `albums.edited` (published via
// event.EntitiesEdited) automatically joins this list without a code
// change here, as does any new entity-mutation verb added to
// ENTITY_MUTATIONS. Non-mutation channels on the same namespace
// (today: none — the existing convention namespaces those elsewhere,
// e.g. `count.labels` lives under `count.`) are no-ops.
subscribeEntityActions("labels", () => evict("labels"));
subscribeEntityActions("albums", () => evict("albums"));

// config.updated is outside the albums/labels namespace but the album
// DELETE handler also calls UpdateClientConfig(), historically the
// only way the client learned about album removal. Now that
// `entity/album.go` emits `albums.deleted` directly the dependency
// is belt-and-braces, but keeping the subscription preserves
// coverage for any future config-touching mutation we don't anticipate.
$event.subscribe("config.updated", () => evict("albums"));

// Drop both lists on logout so user A's labels/albums cannot be
// served to user B inside the same tab. Mirrors Photo.clearCache()'s
// session.logout path in common/session.js (via deleteData → reset).
$event.subscribe("session.logout", () => {
  evict("labels");
  evict("albums");
});

export default typeaheadCache;
