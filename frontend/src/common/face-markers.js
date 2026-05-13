/*

Copyright (c) 2018 - 2025 PhotoPrism UG. All rights reserved.

    This program is free software: you can redistribute it and/or modify
    it under Version 3 of the GNU Affero General Public License (the "AGPL"):
    <https://docs.photoprism.app/license/agpl>

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    The AGPL is supplemented by our Trademark and Brand Guidelines,
    which describe how our Brand Assets may be used:
    <https://www.photoprism.app/trademark>

Feel free to send an email to hello@photoprism.app if you have questions,
want to support our work, or just want to say hello.

Additional information can be found in our Developer Guide:
<https://docs.photoprism.app/developer-guide/>

*/

import { reactive } from "vue";
import { FaceMarkerDisplay, FaceMarkerDraw, isFaceMarkerMode } from "options/face-marker";

// FaceMarkers carries the reactive UI state for the face-marker overlay
// and the surrounding controls (sidebar eye-toggle, + button, inline
// name editor). State that depends on a specific photo (the marker
// array itself) lives on the Photo model and is read via
// `photo.getMarkers(true)` at render time; this singleton only carries
// state that is global to "is the overlay active right now".
//
// Fields:
//
//   mode  — null | FaceMarkerDisplay | FaceMarkerDraw. `null` means the
//           overlay is hidden; the two truthy values mirror the
//           constants in `options/face-marker.js`.
//   busy  — true while a marker-mutating API call is in flight. Sidebar
//           buttons gate on this to prevent overlapping reject / eject /
//           name-set requests.
//   pendingNameMarkerUid — UID of a freshly-created marker whose name
//           input should auto-focus after the overlay's save settles.
//           The sidebar consumes this to drive the inline naming UI and
//           emits `naming-started` to reset it back to "".
//
// The lightbox is the policy owner (it gates writes on
// `shouldShowEditButton()` etc.) and writes through the singleton; the
// sidebar reads from the singleton and emits sidebar→lightbox events
// for state transitions it wants to request.
export class FaceMarkers {
  constructor() {
    this.mode = null;
    this.busy = false;
    this.pendingNameMarkerUid = "";
  }

  // active reports whether any face-marker mode (display or draw) is on.
  get active() {
    return !!this.mode;
  }

  // display reports whether the overlay is in read-only display mode.
  get isDisplay() {
    return this.mode === FaceMarkerDisplay;
  }

  // isDraw reports whether the overlay is in draw mode.
  get isDraw() {
    return this.mode === FaceMarkerDraw;
  }

  // setMode flips the state machine to the specified mode (or `null` to
  // clear it). Invalid values are ignored.
  setMode(mode) {
    if (mode === null || isFaceMarkerMode(mode)) {
      this.mode = mode;
    }
  }

  // display enters read-only display mode (eye toggle on).
  display() {
    this.mode = FaceMarkerDisplay;
  }

  // draw enters draw mode (+ button / Add face).
  draw() {
    this.mode = FaceMarkerDraw;
  }

  // exit clears the mode to `null`; the lightbox watcher reacts to this
  // by tearing down the overlay (and leaves any paused playback paused
  // — playback is NOT resumed on exit, by design).
  exit() {
    this.mode = null;
  }

  // setBusy toggles the in-flight lock; sidebar buttons gate on this.
  setBusy(b) {
    this.busy = !!b;
  }

  // setPendingNameMarkerUid records the UID of a marker whose name
  // input should auto-focus. Pass "" to clear it.
  setPendingNameMarkerUid(uid) {
    this.pendingNameMarkerUid = typeof uid === "string" ? uid : "";
  }

  // reset returns every field to its default. Called when the lightbox
  // closes or the user navigates to a different photo so a fresh slide
  // never inherits stale state.
  reset() {
    this.mode = null;
    this.busy = false;
    this.pendingNameMarkerUid = "";
  }
}

// $faceMarkers is the shared singleton consumed by `PLightbox` and
// `PSidebarInfo`. Components import it directly; it is not installed as
// a Vue plugin because only two consumers need it today (and explicit
// imports keep the dependency graph greppable).
export const $faceMarkers = reactive(new FaceMarkers());

export default $faceMarkers;
