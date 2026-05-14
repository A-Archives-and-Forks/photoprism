// True if debug and/or trace logs should be recorded.
const debug = window.__CONFIG__?.debug;
const trace = window.__CONFIG__?.trace;

// Use global variable to configure pubsub.js, a dependency-free publish/subscribe event hub:
// https://sahadar.github.io/pubsub/#installation
window.pubsub = {
  separator: ".",
  recurrent: true,
  async: true,
  log: trace,
};

// Import pubsub.js, see https://www.npmjs.com/package/pubsub-js.
import * as PubSub from "pubsub-js";

// Use $event as a generic alias for publishing and subscribing to events.
const $event = PubSub;

// Log all events when running in trace log mode, and log config events in debug mode.
// Event names are displayed in blue so that they are easy to recognize.
if (trace) {
  $event.subscribeAll((ev, data) => {
    if (data) {
      console.debug(`%c${ev}`, "background: transparent; color: #9FA8DA; font-weight: normal;", data);
    } else {
      console.debug(`%c${ev}`, "background: transparent; color: #9FA8DA; font-weight: normal;");
    }
  });
} else if (debug) {
  $event.subscribe("config", (ev, data) => {
    if (data) {
      console.debug(`%c${ev}`, "background: transparent; color: #9FA8DA; font-weight: normal;", data);
    } else {
      console.debug(`%c${ev}`, "background: transparent; color: #9FA8DA; font-weight: normal;");
    }
  });
}

// Action verbs the backend's entity-event helpers publish — kept
// in sync with EntityUpdated / EntityCreated / EntityDeleted /
// EntityArchived / EntityRestored / EntityEdited in
// internal/event/publish_entities.go. Frozen so call sites can't
// mutate the shared default; pass a different Set explicitly when
// a caller needs a narrower scope.
export const ENTITY_MUTATIONS = Object.freeze(new Set(["created", "updated", "deleted", "archived", "restored", "edited"]));

// Subscribes to every <namespace>.<action> event whose action is
// in `actions`. Mirrors the page/photos.vue onUpdate switch
// pattern at the cache layer: one hierarchical subscriber per
// namespace, filtered by action inside the handler. Future
// entity-mutation verbs join via one edit to ENTITY_MUTATIONS;
// non-mutation events on the same namespace (a hypothetical
// `<ns>.viewed`, etc.) stay no-ops. Returns the pubsub-js
// subscription token so callers can pass it to $event.unsubscribe.
export function subscribeEntityActions(namespace, handler, actions = ENTITY_MUTATIONS) {
  return $event.subscribe(namespace, (ev, data) => {
    const action = typeof ev === "string" ? ev.split(".")[1] || "" : "";
    if (actions.has(action)) {
      handler(ev, data, action);
    }
  });
}

// Export $event to publish and subscribe to events.
export default $event;
