import { describe, it, expect, beforeEach } from "vitest";
import { FaceMarkers, $faceMarkers } from "common/face-markers";
import { FaceMarkerDisplay, FaceMarkerDraw } from "options/face-marker";

describe("common/face-markers", () => {
  describe("FaceMarkers (constructor + helpers)", () => {
    it("initializes with the inactive default state", () => {
      const fm = new FaceMarkers();
      expect(fm.mode).toBeNull();
      expect(fm.busy).toBe(false);
      expect(fm.pendingNameMarkerUid).toBe("");
      expect(fm.active).toBe(false);
      expect(fm.isDisplay).toBe(false);
      expect(fm.isDraw).toBe(false);
    });

    it("setMode accepts FaceMarkerDisplay, FaceMarkerDraw, and null", () => {
      const fm = new FaceMarkers();
      fm.setMode(FaceMarkerDisplay);
      expect(fm.mode).toBe(FaceMarkerDisplay);
      expect(fm.active).toBe(true);
      expect(fm.isDisplay).toBe(true);
      expect(fm.isDraw).toBe(false);
      fm.setMode(FaceMarkerDraw);
      expect(fm.mode).toBe(FaceMarkerDraw);
      expect(fm.isDraw).toBe(true);
      expect(fm.isDisplay).toBe(false);
      fm.setMode(null);
      expect(fm.mode).toBeNull();
      expect(fm.active).toBe(false);
    });

    it("setMode ignores invalid values without throwing", () => {
      const fm = new FaceMarkers();
      fm.setMode("garbage");
      expect(fm.mode).toBeNull();
      fm.setMode(123);
      expect(fm.mode).toBeNull();
    });

    it("display / draw / exit are convenience setters", () => {
      const fm = new FaceMarkers();
      fm.display();
      expect(fm.mode).toBe(FaceMarkerDisplay);
      fm.draw();
      expect(fm.mode).toBe(FaceMarkerDraw);
      fm.exit();
      expect(fm.mode).toBeNull();
    });

    it("setBusy coerces to boolean", () => {
      const fm = new FaceMarkers();
      fm.setBusy(true);
      expect(fm.busy).toBe(true);
      fm.setBusy(0);
      expect(fm.busy).toBe(false);
      fm.setBusy("yes");
      expect(fm.busy).toBe(true);
    });

    it("setPendingNameMarkerUid only accepts strings; non-strings clear it", () => {
      const fm = new FaceMarkers();
      fm.setPendingNameMarkerUid("abc123");
      expect(fm.pendingNameMarkerUid).toBe("abc123");
      fm.setPendingNameMarkerUid(null);
      expect(fm.pendingNameMarkerUid).toBe("");
      fm.setPendingNameMarkerUid(42);
      expect(fm.pendingNameMarkerUid).toBe("");
    });

    it("reset returns every field to its default", () => {
      const fm = new FaceMarkers();
      fm.setMode(FaceMarkerDraw);
      fm.setBusy(true);
      fm.setPendingNameMarkerUid("uid1");
      fm.reset();
      expect(fm.mode).toBeNull();
      expect(fm.busy).toBe(false);
      expect(fm.pendingNameMarkerUid).toBe("");
    });
  });

  describe("$faceMarkers (shared singleton)", () => {
    beforeEach(() => {
      $faceMarkers.reset();
    });

    it("is a reactive instance of FaceMarkers", () => {
      // Vue's reactive() proxy preserves the prototype chain, so the
      // shared singleton still recognizes as a FaceMarkers instance.
      expect($faceMarkers).toBeInstanceOf(FaceMarkers);
      expect($faceMarkers.mode).toBeNull();
    });

    it("mode mutations propagate to the active/isDisplay/isDraw getters", () => {
      $faceMarkers.display();
      expect($faceMarkers.active).toBe(true);
      expect($faceMarkers.isDisplay).toBe(true);
      $faceMarkers.draw();
      expect($faceMarkers.isDraw).toBe(true);
      $faceMarkers.exit();
      expect($faceMarkers.active).toBe(false);
    });
  });
});
