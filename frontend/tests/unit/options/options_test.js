import "../fixtures";
import * as options from "../../../src/options/options";
import {
  AccountTypes,
  Colors, DefaultLocale,
  Expires, FallbackLocale, FeedbackCategories,
  FindLanguage,
  FindLocale, Gender, Intervals,
  ItemsPerPage,
  MapsAnimate,
  MapsStyle, Orientations,
  PhotoTypes,
  RetryLimits, SetDefaultLocale, StartPages, ThumbFilters, Thumbs, ThumbSizes, Timeouts
} from "../../../src/options/options";

let chai = require("chai/chai");
let assert = chai.assert;

describe("options/options", () => {
  it("should get timezones", () => {
    const timezones = options.TimeZones();
    assert.equal(timezones[0].ID, "Local");
    assert.equal(timezones[0].Name, "Local");
    assert.equal(timezones[1].ID, "UTC");
    assert.equal(timezones[1].Name, "UTC");
  });

  it("should get days", () => {
    const Days = options.Days();
    assert.equal(Days[0].text, "01");
    assert.equal(Days[30].text, "31");
  });

  it("should get years", () => {
    const Years = options.Years();
    const currentYear = new Date().getUTCFullYear();
    assert.equal(Years[0].text, currentYear);
  });

  it("should get indexed years", () => {
    const IndexedYears = options.IndexedYears();
    assert.equal(IndexedYears[0].text, "2021");
  });

  it("should get months", () => {
    const Months = options.Months();
    assert.equal(Months[5].text, "June");
  });

  it("should get short months", () => {
    const MonthsShort = options.MonthsShort();
    assert.equal(MonthsShort[5].text, "06");
  });

  it("should get languages", () => {
    const Languages = options.Languages();
    assert.equal(Languages[0].value, "en");
  });

  it("should set default locale", () => {
    assert.equal(DefaultLocale, "en");
    SetDefaultLocale("de");
    assert.equal(DefaultLocale, "de");
    SetDefaultLocale("en");
  });

  it("should return default when no locale is provided", () => {
    assert.equal(FindLanguage("").value, "en");
  });

  it("should return default locale is smaller than 2", () => {
    assert.equal(FindLanguage("d").value, "en");
  });

  it("should return default locale", () => {
    assert.equal(FindLanguage("xx").value, "en");
  });

  it("should return correct locale", () => {
    assert.equal(FindLanguage("de").value, "de");
    assert.equal(FindLanguage("de").text, "Deutsch");
    assert.equal(FindLanguage("de_AT").value, "de");
    assert.equal(FindLanguage("de_AT").text, "Deutsch");
    assert.equal(FindLanguage("zh-tw").value, "zh_TW");
    assert.equal(FindLanguage("zh-tw").text, "繁體中文");
    assert.equal(FindLanguage("zh+tw").value, "zh_TW");
    assert.equal(FindLanguage("zh+tw").text, "繁體中文");
    assert.equal(FindLanguage("zh_AT").value, "zh");
    assert.equal(FindLanguage("zh_AT").text, "简体中文");
    assert.equal(FindLanguage("ZH_TW").value, "zh_TW");
    assert.equal(FindLanguage("ZH_TW").text, "繁體中文");
    assert.equal(FindLanguage("zH-tW").value, "zh_TW");
    assert.equal(FindLanguage("zH-tW").text, "繁體中文");
  });

  it("should return default locale", () => {
    assert.equal(FindLocale("xx"), "en");
    assert.equal(FindLocale(""), "en");
  });

  it("should return fallback locale", () => {
    assert.equal(FallbackLocale(), "en");
  });

  it("should return items per page", () => {
    assert.equal(ItemsPerPage()[0].value, 10);
  });

  it("should return start page options", () => {
    let features = {
      "account": true,
      "albums": true,
      "archive": true,
      "delete": true,
      "download": true,
      "edit": true,
      "estimates": true,
      "favorites": true,
      "files": true,
      "folders": true,
      "import": true,
      "labels": true,
      "library": true,
      "logs": true,
      "calendar": true,
      "moments": true,
      "people": true,
      "places": true,
      "private": true,
      "ratings": true,
      "reactions": true,
      "review": true,
      "search": true,
      "services": true,
      "settings": true,
      "share": true,
      "upload": true,
      "videos": true
    };
    assert.equal(StartPages(features).length, 12);
    assert.equal(StartPages(features)[5].value, "people");
    assert.equal(StartPages(features)[5].props.disabled, false);
    features = {
      "account": true,
      "albums": true,
      "archive": true,
      "delete": true,
      "download": true,
      "edit": true,
      "estimates": true,
      "favorites": true,
      "files": true,
      "folders": true,
      "import": true,
      "labels": true,
      "library": true,
      "logs": true,
      "calendar": false,
      "moments": true,
      "people": false,
      "places": true,
      "private": true,
      "ratings": true,
      "reactions": true,
      "review": true,
      "search": true,
      "services": true,
      "settings": true,
      "share": true,
      "upload": true,
      "videos": true
    };
    assert.equal(StartPages(features).length, 12);
    assert.equal(StartPages(features)[5].value, "people");
    assert.equal(StartPages(features)[5].props.disabled, true);

  });

  it("should return animation options", () => {
    assert.equal(MapsAnimate()[1].value, 2500);
  });

  it("should return photo types", () => {
    assert.equal(PhotoTypes()[0].value, "image");
    assert.equal(PhotoTypes()[1].value, "raw");
  });

  it("should return map styles", () => {
    let styles = MapsStyle(true);
    assert.include(styles[styles.length - 1].value, "low-resolution");
    styles = MapsStyle(false);
    assert.notInclude(styles[styles.length - 1].value, "low-resolution");
  });

  it("should return timeouts", () => {
    assert.equal(Timeouts()[1].value, "high");
  });

  it("should return retry limits", () => {
    assert.equal(RetryLimits()[1].value, 1);
  });

  it("should return intervals", () => {
    assert.equal(Intervals()[0].text, "Never");
    assert.equal(Intervals()[1].text, "1 hour");
  });

  it("should return expiry options", () => {
    assert.equal(Expires()[0].text, "Never");
    assert.equal(Expires()[1].text, "After 1 day");
  });

  it("should return colors", () => {
    assert.equal(Colors()[0].Slug, "purple");
  });

  it("should return feedback categories", () => {
    assert.equal(FeedbackCategories()[0].value, "feedback");
  });

  it("should return thumb sizes", () => {
    assert.equal(ThumbSizes()[1].value, "fit_720");
  });

  it("should return thumb filters", () => {
    assert.equal(ThumbFilters()[0].value, "blackman");
  });

  it("should return gender", () => {
    assert.equal(Gender()[2].value, "other");
  });

  it("should return orientations", () => {
    assert.equal(Orientations()[1].text, "90°");
  });

  it("should return service account type options", () => {
    assert.equal(AccountTypes()[0].value, "webdav");
    assert.equal(AccountTypes().length, 1);
  });

});
