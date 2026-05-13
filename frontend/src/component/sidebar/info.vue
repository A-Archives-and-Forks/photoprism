<template>
  <div class="p-sidebar-info bg-background metadata" :class="{ 'hide-edit-pencils': hideEditPencils }">
    <v-toolbar density="comfortable" color="background">
      <v-btn :icon="$isRtl ? 'mdi-chevron-left' : 'mdi-chevron-right'" :title="$gettext('Close')" @click.stop="close()"></v-btn>
      <v-toolbar-title>{{ $gettext(`Information`) }}</v-toolbar-title>
    </v-toolbar>
    <div v-if="model.UID">
      <v-list nav slim tile density="compact" bg-color="background" class="metadata__list mt-1">
        <!-- Title -->
        <v-list-item
          v-if="editingField === 'title' || model.Title || isEditable"
          :class="['metadata__item', { clickable: editingField !== 'title' && (isEditable || model.Title) }]"
          @click.stop="onTextRowClick('title', model.Title)"
        >
          <v-text-field
            v-if="editingField === 'title'"
            :ref="setInlineEditorRef"
            v-model="photo.Title"
            :rules="[textRule]"
            density="compact"
            hide-details="auto"
            autocomplete="off"
            class="meta-inline-edit meta-inline-title"
            @keydown.enter.stop.prevent="confirmField"
            @keydown.escape.stop.prevent="cancelEditing"
            @blur="onInlineFieldBlur"
          ></v-text-field>
          <div v-else-if="model.Title" class="text-subtitle-2 meta-title">{{ model.Title }}</div>
          <div v-else class="meta-add-prompt" @click.stop="startEditing('title')">{{ $pgettext("Photo", "Add a Title") }}</div>
          <template v-if="isEditable" #append>
            <p-sidebar-inline-toolbar :editing="editingField === 'title'" @confirm="confirmField" @start="startEditing('title')" />
          </template>
        </v-list-item>

        <!-- Caption -->
        <v-list-item
          v-if="editingField === 'caption' || model.Caption || isEditable"
          :class="['metadata__item', { clickable: editingField !== 'caption' && (isEditable || model.Caption) }]"
          @click.stop="onTextRowClick('caption', model.Caption)"
        >
          <v-textarea
            v-if="editingField === 'caption'"
            :ref="setInlineEditorRef"
            v-model="photo.Caption"
            density="compact"
            auto-grow
            :max-rows="6"
            hide-details="auto"
            autocomplete="off"
            class="meta-inline-edit meta-inline-caption"
            @keydown.escape.stop.prevent="cancelEditing"
            @blur="onInlineFieldBlur"
          ></v-textarea>
          <!-- eslint-disable-next-line vue/no-v-html -- captionHtml is encode-then-sanitized via $util.sanitizeHtml($util.encodeHTML(raw)); see captionHtml() computed -->
          <div v-else-if="model.Caption" class="text-body-2 meta-caption meta-scrollable" v-html="captionHtml"></div>
          <div v-else class="meta-add-prompt" @click.stop="startEditing('caption')">{{ $gettext("Add a Caption") }}</div>
          <template v-if="isEditable" #append>
            <p-sidebar-inline-toolbar :editing="editingField === 'caption'" @confirm="confirmField" @start="startEditing('caption')" />
          </template>
        </v-list-item>

        <v-divider v-if="editingField === 'title' || editingField === 'caption' || model.Title || model.Caption || isEditable" class="my-3"></v-divider>

        <v-list-item
          v-if="fileInfo"
          v-tooltip="fileTypeName"
          :prepend-icon="fileIcon"
          :title="fileInfo"
          :subtitle="fileName"
          :lines="fileName ? 'two' : 'one'"
          :class="['metadata__item', 'meta-file', { clickable: !!fileName }]"
          @click.stop.prevent="fileName && $util.copyText(fileName)"
        >
        </v-list-item>

        <v-divider v-if="fileName || fileInfo" class="my-3"></v-divider>

        <v-list-item
          v-tooltip="$gettext(`Taken`)"
          :title="formatTime(model)"
          prepend-icon="mdi-calendar"
          :class="['metadata__item', { clickable: isEditable }]"
          @click.stop="openDateTimeDialog"
        >
          <template v-if="isEditable" #append>
            <v-btn
              icon="mdi-pencil-outline"
              density="compact"
              variant="plain"
              size="x-small"
              class="meta-inline-pencil"
              :title="$gettext('Edit')"
              @click.stop="openDateTimeDialog"
            ></v-btn>
          </template>
        </v-list-item>

        <v-list-item
          v-if="!restrictedRole && (cameraInfo || isEditable)"
          v-tooltip="$gettext('Camera')"
          :title="cameraInfo || $gettext('Unknown')"
          prepend-icon="mdi-camera"
          :class="['metadata__item', { clickable: isEditable }]"
          @click.stop="openCameraDialog"
        >
          <template v-if="isEditable" #append>
            <v-btn
              icon="mdi-pencil-outline"
              density="compact"
              variant="plain"
              size="x-small"
              class="meta-inline-pencil"
              :title="$gettext('Edit')"
              @click.stop="openCameraDialog"
            ></v-btn>
          </template>
        </v-list-item>

        <v-list-item
          v-if="!restrictedRole && lensInfo"
          v-tooltip="$gettext('Lens')"
          :title="lensInfo"
          prepend-icon="mdi-camera-iris"
          :class="['metadata__item', { clickable: isEditable }]"
          @click.stop="openCameraDialog"
        ></v-list-item>

        <template v-if="locationRowVisible">
          <!-- v-divider class="my-3"></v-divider -->
          <v-list-item
            v-tooltip="$gettext('Location')"
            :title="locationTitle"
            :subtitle="locationSubtitle"
            :lines="locationSubtitle ? 'two' : 'one'"
            prepend-icon="mdi-map-marker"
            :class="['metadata__item', 'meta-location', { 'meta-coordinates': model.Lat && model.Lng, 'clickable': locationRowClickable }]"
            @click.stop.prevent="onLocationRowClick"
          >
            <template v-if="isEditable && featPlaces" #append>
              <v-btn
                icon="mdi-pencil-outline"
                density="compact"
                variant="plain"
                size="x-small"
                class="meta-inline-pencil meta-inline-pencil--location"
                :title="$gettext('Edit')"
                @click.stop.prevent="openLocationDialog"
              ></v-btn>
            </template>
          </v-list-item>
          <v-list-item v-if="featPlaces && model.Lat && model.Lng" class="mx-0 px-0">
            <p-map :latlng="[model.Lat, model.Lng]" :animate-duration="0" :marker-clickable="isEditable" @marker-clicked="openLocationDialog"></p-map>
          </v-list-item>
        </template>

        <template v-if="!restrictedRole && featPeople && (people.length > 0 || isEditable)">
          <v-divider class="my-3"></v-divider>
          <v-list-item class="metadata__item">
            <div class="text-subtitle-2">{{ $gettext("People") }}</div>
            <template v-if="isEditable || people.length > 0" #append>
              <!--
                Per-role toggle:
                - Editable users see the pencil / pencil-off "Edit Faces" toggle.
                  Draw mode is a strict superset of display mode for editable
                  users (boxes + names visible, plus drag-to-create and
                  click-to-remove), so a separate display-mode toggle adds
                  no value.
                - Non-editable users see the eye / eye-off "Show face markers"
                  toggle. Display mode is read-only.

                Both modes hide the lightbox chrome, pause video / slideshow
                on entry, and gate the conflict-only keyboard shortcuts —
                see lightbox.vue isShortcutDisabledInFaceMarkerMode.
              -->
              <v-btn
                v-if="isEditable"
                :icon="addingMarker ? 'mdi-pencil-off-outline' : 'mdi-pencil-outline'"
                density="compact"
                variant="plain"
                size="x-small"
                class="meta-faces-edit"
                :class="{ 'is-active': addingMarker }"
                :title="addingMarker ? $gettext('Done') : $gettext('Edit Faces')"
                :disabled="markersBusy"
                @mousedown.prevent
                @click.stop="onToggleFaceMarkerDraw"
              ></v-btn>
              <v-btn
                v-else
                :icon="markersVisible ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
                density="compact"
                variant="plain"
                size="x-small"
                class="meta-markers-toggle"
                :class="{ 'is-active': markersVisible }"
                :title="markersVisible ? $gettext('Hide face markers') : $gettext('Show face markers')"
                :disabled="markersBusy"
                @mousedown.prevent
                @click.stop="onToggleFaceMarkerMode"
              ></v-btn>
            </template>
          </v-list-item>
          <v-list-item v-for="m in people" :key="m.UID || m.CropID" :data-marker-uid="m.UID" class="metadata__item metadata__person-row">
            <template #prepend>
              <img
                :src="m.thumbnailUrl('tile_160')"
                :alt="m.Name"
                class="meta-person__avatar"
                :class="{ clickable: m.Name && m.SubjUID }"
                @click.stop="m.Name && m.SubjUID ? navigateToPerson(m) : null"
              />
            </template>
            <v-combobox
              v-if="isEditable"
              :model-value="markerInputValue(m.UID)"
              :search="markerInputSearch(m.UID)"
              :items="knownPeople"
              item-title="Name"
              item-value="Name"
              :placeholder="$gettext('Name')"
              :menu-props="markerMenuProps"
              :list-props="chipListProps"
              :readonly="markersBusy || !!m.SubjUID"
              :rules="[markerNameRule]"
              return-object
              hide-no-data
              hide-details="auto"
              single-line
              open-on-clear
              append-icon=""
              :menu-icon="null"
              density="compact"
              bg-color="background"
              autocomplete="off"
              class="meta-inline-edit meta-inline-marker"
              :class="{ 'meta-inline-marker--named': m.SubjUID }"
              @update:model-value="(v) => onPickPerson(m, v)"
              @update:search="(v) => setMarkerInputValue(m.UID, v)"
              @keydown.enter.stop.prevent="confirmMarkerName(m, 'enter')"
              @keydown.escape.stop.prevent="cancelMarkerName(m)"
              @blur="confirmMarkerName(m, 'blur')"
              @click.stop
            >
              <template v-if="m.SubjUID" #append-inner>
                <v-btn
                  icon="mdi-eject"
                  density="compact"
                  variant="plain"
                  size="x-small"
                  class="meta-marker-eject"
                  :title="$gettext('Remove Name')"
                  :disabled="markersBusy"
                  @mousedown.prevent
                  @click.stop="onEjectMarker(m)"
                ></v-btn>
              </template>
            </v-combobox>
            <v-list-item-title v-else-if="m.Name" class="meta-person__name">{{ m.Name }}</v-list-item-title>
            <v-list-item-title v-else class="meta-person__name meta-person__unnamed">{{ $gettext("Unknown") }}</v-list-item-title>
          </v-list-item>
        </template>

        <template v-if="!restrictedRole && (labels.length > 0 || isEditable)">
          <v-divider class="my-3"></v-divider>
          <v-list-item class="metadata__item meta-labels">
            <div class="text-subtitle-2">{{ $gettext("Labels") }}</div>
            <template v-if="isEditable && chipState.labels.removals.length > 0" #append>
              <p-sidebar-inline-toolbar :editing="true" :can-undo="true" @confirm="confirmLabels" @undo="undoChipRemovals('labels')" />
            </template>
          </v-list-item>
          <v-list-item v-if="visibleLabels.length > 0" class="metadata__item metadata__chips meta-labels">
            <div class="d-flex flex-wrap ga-1">
              <span
                v-for="l in visibleLabels"
                :key="l.Label.UID"
                tabindex="0"
                class="meta-chip meta-chip--primary"
                @click.stop.prevent="onChipActivate('labels', l)"
                @keydown.enter.stop.prevent="onChipActivate('labels', l)"
                @keydown.delete.stop.prevent="onChipDelete('labels', l)"
              >
                {{ l.Label.Name }}
                <v-icon
                  v-if="isEditable"
                  icon="mdi-close-circle"
                  size="x-small"
                  class="ml-1 meta-chip__remove"
                  :title="$gettext('Remove')"
                  @click.stop.prevent="onChipDelete('labels', l)"
                ></v-icon>
              </span>
            </div>
          </v-list-item>
          <v-list-item v-if="isEditable" class="metadata__item meta-labels">
            <v-combobox
              :key="chipState.labels.key"
              v-model="chipState.labels.input"
              v-model:search="chipState.labels.search"
              :items="chipState.labels.options"
              item-title="Name"
              item-value="Name"
              return-object
              :placeholder="$gettext('Select or create labels')"
              hide-details
              hide-no-data
              single-line
              append-icon=""
              :menu-icon="null"
              density="compact"
              bg-color="background"
              autocomplete="off"
              :menu-props="chipMenuProps"
              :list-props="chipListProps"
              class="meta-inline-edit"
              @focus="loadChipOptions('labels')"
              @update:model-value="onLabelSelected"
              @keydown.enter.stop.prevent="onLabelEnter"
              @keydown.escape.stop.prevent="onChipEscape('labels')"
            ></v-combobox>
          </v-list-item>
        </template>

        <template v-if="!restrictedRole && (albums.length > 0 || isEditable)">
          <v-divider class="my-3"></v-divider>
          <v-list-item class="metadata__item meta-albums">
            <div class="text-subtitle-2">{{ $gettext("Albums") }}</div>
            <template v-if="isEditable && chipState.albums.removals.length > 0" #append>
              <p-sidebar-inline-toolbar :editing="true" :can-undo="true" @confirm="confirmAlbums" @undo="undoChipRemovals('albums')" />
            </template>
          </v-list-item>
          <v-list-item v-if="visibleAlbums.length > 0" class="metadata__item metadata__chips meta-albums">
            <div class="d-flex flex-wrap ga-1">
              <span
                v-for="a in visibleAlbums"
                :key="a.UID"
                tabindex="0"
                class="meta-chip meta-chip--primary"
                @click.stop.prevent="onChipActivate('albums', a)"
                @keydown.enter.stop.prevent="onChipActivate('albums', a)"
                @keydown.delete.stop.prevent="onChipDelete('albums', a)"
              >
                {{ a.Title }}
                <v-icon
                  v-if="isEditable"
                  icon="mdi-close-circle"
                  size="x-small"
                  class="ml-1 meta-chip__remove"
                  :title="$gettext('Remove')"
                  @click.stop.prevent="onChipDelete('albums', a)"
                ></v-icon>
              </span>
            </div>
          </v-list-item>
          <v-list-item v-if="isEditable" class="metadata__item meta-albums">
            <v-combobox
              :key="chipState.albums.key"
              v-model="chipState.albums.input"
              v-model:search="chipState.albums.search"
              :items="chipState.albums.options"
              item-title="Title"
              item-value="Title"
              return-object
              :placeholder="$gettext('Select or create albums')"
              hide-details
              hide-no-data
              single-line
              append-icon=""
              :menu-icon="null"
              density="compact"
              bg-color="background"
              autocomplete="off"
              :menu-props="chipMenuProps"
              :list-props="chipListProps"
              class="meta-inline-edit"
              @focus="loadChipOptions('albums')"
              @update:model-value="onAlbumSelected"
              @keydown.enter.stop.prevent="onAlbumEnter"
              @keydown.escape.stop.prevent="onChipEscape('albums')"
            ></v-combobox>
          </v-list-item>
        </template>

        <template v-if="showDetailsSection">
          <v-divider class="my-3"></v-divider>
          <v-list-item
            v-for="f in detailsFields"
            v-show="shouldShowFieldRow(f)"
            :key="f.key"
            v-tooltip="f.label"
            :prepend-icon="f.icon"
            :class="['metadata__item', `meta-${f.key}`, { clickable: editingField !== f.key && (isEditable || f.read(photo)) }]"
            @click.stop="onTextRowClick(f.key, f.read(photo))"
          >
            <v-textarea
              v-if="editingField === f.key"
              :ref="setInlineEditorRef"
              :model-value="f.read(photo)"
              :placeholder="f.label"
              :rules="[textRule]"
              density="compact"
              auto-grow
              hide-details="auto"
              autocomplete="off"
              class="meta-inline-edit"
              :class="`meta-inline-${f.key}`"
              @update:model-value="(v) => f.write(photo, v)"
              @keydown.escape.stop.prevent="cancelEditing"
              @blur="onInlineFieldBlur"
            ></v-textarea>
            <div v-else-if="f.read(photo)" class="text-body-2 meta-scrollable">{{ f.read(photo) }}</div>
            <div v-else class="meta-add-prompt" @click.stop="startEditing(f.key)">{{ f.label }}</div>
            <template v-if="isEditable" #append>
              <p-sidebar-inline-toolbar :editing="editingField === f.key" @confirm="confirmField" @start="startEditing(f.key)" />
            </template>
          </v-list-item>
        </template>

        <template v-for="f in textFields" :key="f.key">
          <template v-if="!restrictedRole && shouldShowFieldRow(f)">
            <v-divider class="my-3"></v-divider>
            <v-list-item class="metadata__item" :class="`meta-${f.key}`">
              <div class="text-subtitle-2">{{ f.label }}</div>
              <template v-if="isEditable" #append>
                <p-sidebar-inline-toolbar :editing="editingField === f.key" @confirm="confirmField" @start="startEditing(f.key)" />
              </template>
            </v-list-item>
            <v-list-item
              :class="['metadata__item', `meta-${f.key}`, { clickable: editingField !== f.key && (isEditable || f.read(photo)) }]"
              @click.stop="onTextRowClick(f.key, f.read(photo))"
            >
              <v-textarea
                v-if="editingField === f.key"
                :ref="setInlineEditorRef"
                :model-value="f.read(photo)"
                :placeholder="f.label"
                density="compact"
                auto-grow
                hide-details="auto"
                autocomplete="off"
                class="meta-inline-edit"
                :class="`meta-inline-${f.key}`"
                @update:model-value="(v) => f.write(photo, v)"
                @keydown.escape.stop.prevent="cancelEditing"
                @blur="onInlineFieldBlur"
              ></v-textarea>
              <!-- eslint-disable-next-line vue/no-v-html -- f.htmlValue references a sanitized computed (e.g. notesHtml) — encode-then-sanitize via $util.sanitizeHtml($util.encodeHTML(raw)). -->
              <div v-else-if="f.display === 'html' && fieldHtml(f)" class="text-body-2 meta-scrollable" :class="`meta-${f.key}`" v-html="fieldHtml(f)"></div>
              <div v-else-if="f.display !== 'html' && f.read(photo)" class="text-body-2 meta-scrollable" :class="`meta-${f.key}`">
                {{ f.read(photo) }}
              </div>
              <div v-else class="meta-add-prompt" @click.stop="startEditing(f.key)">{{ f.label }}</div>
            </v-list-item>
          </template>
        </template>
      </v-list>
    </div>
    <p-meta-datetime-dialog :visible="dateTimeDialog" :photo="photo" @close="dateTimeDialog = false" @confirm="confirmDateTime"></p-meta-datetime-dialog>
    <p-meta-camera-dialog :visible="cameraDialog" :photo="photo" @close="cameraDialog = false" @confirm="confirmCamera"></p-meta-camera-dialog>
    <p-meta-location-dialog
      :visible="locationDialog"
      :latlng="[photo ? Number(photo.Lat) || 0 : 0, photo ? Number(photo.Lng) || 0 : 0]"
      @close="locationDialog = false"
      @confirm="confirmLocation"
    ></p-meta-location-dialog>
    <p-confirm-dialog
      :visible="discardDialog.visible"
      icon="mdi-alert-circle-outline"
      :text="$gettext('Discard unsaved changes?')"
      :action="$gettext('Discard')"
      @close="onDiscardCancel"
      @confirm="onDiscardConfirm"
    ></p-confirm-dialog>
    <p-confirm-dialog
      :visible="addNameDialog.visible"
      icon="mdi-account-plus"
      :icon-size="42"
      :text="addNameDialog.name ? $gettext('Add %{s}?', { s: addNameDialog.name }) : $gettext('Add person?')"
      @close="onAddNameCancel"
      @confirm="onAddNameConfirm"
    ></p-confirm-dialog>
  </div>
</template>

<script>
import { DateTime } from "luxon";
import * as formats from "options/formats";
import { $faceMarkers } from "common/face-markers";

import * as media from "common/media";
import typeaheadCache from "common/typeahead-cache";
import { Album } from "model/album";
import PMap from "component/map.vue";
import PMetaDatetimeDialog from "component/meta/datetime/dialog.vue";
import PMetaCameraDialog from "component/meta/camera/dialog.vue";
import PMetaLocationDialog from "component/meta/location/dialog.vue";
import PConfirmDialog from "component/confirm/dialog.vue";
import PSidebarInlineToolbar from "component/sidebar/inline-toolbar.vue";

export default {
  name: "PSidebarInfo",
  components: {
    PMap,
    PMetaDatetimeDialog,
    PMetaCameraDialog,
    PMetaLocationDialog,
    PConfirmDialog,
    PSidebarInlineToolbar,
  },
  props: {
    // UID of the photo currently shown in the parent lightbox. Drives the
    // sidebar lifecycle (re-fetching markers, resetting inline edits) when
    // the user navigates between slides. All other parent state is read
    // through `view` (see data() below) — this matches the pattern used
    // by component/photo/edit/labels.vue.
    uid: {
      type: String,
      default: "",
    },
  },
  emits: ["close", "toggle-face-marker-mode", "toggle-face-marker-draw", "eject-marker", "reload-markers", "naming-started"],
  data() {
    return {
      // Live reactive handle to the parent lightbox's $data, captured once at
      // mount via `$view.getData()`. The lightbox calls `$view.enter(this)`
      // before the sidebar mounts (see lightbox.vue:showDialog), so this is
      // populated by the time data() runs. Mutations through this.view.X
      // write through to the parent and don't trigger vue/no-mutating-props.
      view: this.$view.getData(),
      // Reactive handle to the shared face-marker state singleton
      // (`common/face-markers.js`). Drives the eye / pencil toggles,
      // the inline naming flow, and the markersBusy gates. Sidebar
      // emits transition requests via `toggle-face-marker-mode` /
      // `toggle-face-marker-draw` / `eject-marker` / `reload-markers`;
      // the lightbox is the policy owner and performs the actual writes.
      // Marker removal lives on the face-marker overlay (click an
      // unnamed marker in edit mode → inline confirm pill) — not in
      // the sidebar — so the per-row `mdi-close` button that used to
      // sit inside the combobox was retired.
      faceMarkers: $faceMarkers,
      actions: [],
      featPeople: this.$config.feature("people"),
      featPlaces: this.$config.feature("places"),
      textRule: (v) => !v || v.length <= this.$config.get("clip") || this.$gettext("Text too long"),
      dateTimeDialog: false,
      cameraDialog: false,
      locationDialog: false,
      // When true, the inline pencil edit buttons are hidden via CSS
      // (`.hide-edit-pencils` rule in `css/lightbox.css`). Row-level
      // click handlers still route to `startEditing` / `open*Dialog`,
      // so the pencils are a redundant affordance that can be toggled
      // off without affecting reachability. Save (`meta-inline-confirm`)
      // and Undo (`meta-inline-undo`) buttons stay visible because they
      // commit pending state rather than enter edit mode.
      hideEditPencils: true,
      editingField: null,
      editOriginal: null,
      // Per-field combobox state. The combobox/autocomplete row stays
      // mounted whenever the section is editable (no pencil-to-edit
      // gesture for chips), so each section needs its own input/search
      // scratch refs and per-field force-remount key.
      //
      // - input/search: Vuetify v-model and v-model:search bindings.
      //   `search` doubles as the "typed-but-not-yet-Enter" detector
      //   for hasPendingEdit().
      // - key: incremented on Enter to force-remount the combobox so
      //   stale dropdown state clears alongside the input.
      // - options: typeahead suggestions populated lazily from the
      //   shared typeaheadCache (common/typeahead-cache.js); shape
      //   matches the v-combobox/v-autocomplete item-title bindings.
      // - removals: pending Label.ID / Album.UID removals committed
      //   by the toolbar ✓. Additions take an instant-save path
      //   (addLabelImmediate / addAlbumImmediate → Photo model
      //   methods), so they never enter chipState.
      chipState: {
        labels: { input: null, search: "", key: 0, options: [], removals: [] },
        albums: { input: null, search: "", key: 0, options: [], removals: [] },
      },
      markerDrafts: {},
      markerNameRule: (v) => !v || v.length <= this.$config.get("clip") || this.$gettext("Text too long"),
      markerMenuProps: {
        openOnFocus: true,
        closeOnContentClick: true,
        maxHeight: 260,
        class: "meta-inline-menu",
      },
      chipMenuProps: {
        class: "meta-inline-menu",
      },
      // Forwarded to the inner v-list of the combobox/autocomplete dropdown.
      // density="compact" on the input itself only sizes the trigger field —
      // list-props is the documented way to size the menu items themselves.
      chipListProps: {
        density: "compact",
      },
      discardDialog: {
        visible: false,
        resolver: null,
      },
      // Stores the marker UID (not the transient Marker instance returned
      // by getMarkers, see P1-8). The live Marker is re-derived from
      // photo.getMarkers(true) at commit time so a slide-nav between
      // open + Add/Cancel can't write through a stale reference.
      addNameDialog: {
        visible: false,
        markerUid: "",
        name: "",
      },
    };
  },
  computed: {
    // Aliases for parent-owned reactive state. These read through `view` so
    // every existing template/script reference (this.model.X, this.photo.X, etc.)
    // keeps working without churn. Mutations are explicit: write to
    // `this.view.photo.X` / `this.view.model.X`, never to `this.photo` / `this.model`.
    model() {
      return this.view?.model;
    },
    photo() {
      return this.view?.photo;
    },
    canEdit() {
      return Boolean(this.view?.canEdit && this.view?.contextAllowsEdit);
    },
    collection() {
      return this.view?.collection;
    },
    context() {
      return this.view?.context;
    },
    // Derived from the shared face-marker state singleton
    // (`common/face-markers.js`): null = no overlay, 'display' = read-
    // only markers, 'draw' = Add Face. The sidebar template still binds
    // to `markersVisible` / `addingMarker` booleans so the eye / + icon
    // logic stays compact.
    markersVisible() {
      return this.faceMarkers.active;
    },
    addingMarker() {
      return this.faceMarkers.isDraw;
    },
    markersBusy() {
      return this.faceMarkers.busy;
    },
    newMarkerUid() {
      return this.faceMarkers.pendingNameMarkerUid;
    },
    isEditable() {
      return this.canEdit && this.photo && this.photo.Details && !this.restrictedRole;
    },
    restrictedRole() {
      return this.$session.isSidebarRestricted();
    },
    captionHtml() {
      const raw = this.photo?.Caption ?? this.model?.Caption;
      if (!raw) return "";
      return this.$util.sanitizeHtml(this.$util.encodeHTML(raw));
    },
    notesHtml() {
      if (!this.photo?.Details?.Notes) return "";
      return this.$util.sanitizeHtml(this.$util.encodeHTML(this.photo.Details.Notes));
    },
    cameraInfo() {
      if (!this.photo) return "";
      // Backend returns the "Unknown" placeholder camera (CameraID=1,
      // Camera={Make:"", Model:"Unknown"}) when no EXIF camera is set, and
      // formatCamera() happily renders that as " Unknown". Suppress it so
      // the read-only sidebar doesn't surface an empty camera row.
      const hasRealCamera =
        (this.photo.CameraID && this.photo.CameraID > 1) ||
        (this.photo.CameraMake && this.photo.CameraMake.trim()) ||
        (this.photo.CameraModel && this.photo.CameraModel.trim() && this.photo.CameraModel !== "Unknown");
      if (!hasRealCamera) return "";
      // Suppress "Unknown, ISO 100"-style rows when only ISO/exposure are set.
      if (!this.$util.formatCamera(this.photo.Camera, this.photo.CameraID, this.photo.CameraMake, this.photo.CameraModel, false)) return "";
      const info = this.photo.getCameraInfo();
      return info !== this.$gettext("Unknown") ? info : "";
    },
    lensInfo() {
      if (!this.photo) return "";
      const hasLens =
        (this.photo.LensID && this.photo.LensID > 1) || this.photo.LensMake || this.photo.LensModel || this.photo.Lens?.Model || this.photo.Lens?.Make;
      if (!hasLens) return "";
      const info = this.photo.getLensInfo();
      return info !== this.$gettext("Unknown") ? info : "";
    },
    exifInfo() {
      if (!this.photo) return "";
      return this.photo.getExifInfo();
    },
    people() {
      if (!this.photo) return [];
      return this.photo.getMarkers(true);
    },
    // Sorted, locale-aware copy of $config.values.people for the marker
    // combobox dropdown. Mirrors the Labels/Albums sort baked into
    // loadChipOptions (L12): localeCompare with `undefined` locale +
    // base sensitivity + numeric collation. Reading from $config keeps the
    // list reactive to WS `people.{created,updated,deleted}` events that
    // are already handled inside common/config.js. A future shared
    // people-cache module (F1) will move this logic out of the component;
    // until then the per-consumer sort matches the Labels/Albums pattern.
    knownPeople() {
      const values = this.$config && this.$config.values;
      if (!values || !Array.isArray(values.people)) return [];
      return values.people
        .filter((p) => p && p.Name)
        .slice()
        .sort((a, b) => (a.Name || "").localeCompare(b.Name || "", undefined, { sensitivity: "base", numeric: true }));
    },
    labels() {
      if (!this.photo?.Labels) return [];
      return this.photo.Labels.filter((l) => l.Label && l.Label.Name && l.Uncertainty < 100);
    },
    albums() {
      if (!this.photo?.Albums) return [];
      return this.photo.Albums.filter((a) => a.Title && !a.Private);
    },
    // Visible chip lists — `labels` / `albums` minus anything currently
    // marked for removal in `chipState`. The chip-row wrapper gates on
    // these so it disappears once the user has soft-removed every chip in
    // the section (otherwise the wrapper would render as an empty box).
    // Undo restores the chips by clearing `chipState.<field>.removals`,
    // which makes these computeds repopulate reactively.
    visibleLabels() {
      return this.labels.filter((l) => !this.isChipPendingRemoval("labels", l?.Label?.ID));
    },
    visibleAlbums() {
      return this.albums.filter((a) => !this.isChipPendingRemoval("albums", a?.UID));
    },
    subject() {
      return this.photo?.Details?.Subject || "";
    },
    artist() {
      return this.photo?.Details?.Artist || "";
    },
    copyright() {
      return this.photo?.Details?.Copyright || "";
    },
    license() {
      return this.photo?.Details?.License || "";
    },
    keywords() {
      return this.photo?.Details?.Keywords || "";
    },
    // Single source of truth for inline-text fields. Each entry knows how to
    // read/write its raw value, what label to render (tooltip, placeholder,
    // add-prompt), and whether the display branch should treat the value as
    // sanitized HTML (Caption, Notes) or plain text (everything else).
    // detailsFields/textFields below select subsets for the two visual layouts.
    fieldRegistry() {
      return {
        title: {
          key: "title",
          label: this.$pgettext("Photo", "Title"),
          read: (p) => p?.Title,
          write: (p, v) => {
            if (p) p.Title = v;
          },
          display: "text",
        },
        caption: {
          key: "caption",
          label: this.$gettext("Caption"),
          read: (p) => p?.Caption,
          write: (p, v) => {
            if (p) p.Caption = v;
          },
          display: "html",
          htmlValue: "captionHtml",
        },
        subject: {
          key: "subject",
          label: this.$gettext("Subject"),
          icon: "mdi-text-box-outline",
          read: (p) => p?.Details?.Subject,
          write: (p, v) => {
            if (p?.Details) p.Details.Subject = v;
          },
          display: "text",
        },
        artist: {
          key: "artist",
          label: this.$gettext("Artist"),
          icon: "mdi-palette",
          read: (p) => p?.Details?.Artist,
          write: (p, v) => {
            if (p?.Details) p.Details.Artist = v;
          },
          display: "text",
        },
        copyright: {
          key: "copyright",
          label: this.$gettext("Copyright"),
          icon: "mdi-copyright",
          read: (p) => p?.Details?.Copyright,
          write: (p, v) => {
            if (p?.Details) p.Details.Copyright = v;
          },
          display: "text",
        },
        license: {
          key: "license",
          label: this.$gettext("License"),
          icon: "mdi-license",
          read: (p) => p?.Details?.License,
          write: (p, v) => {
            if (p?.Details) p.Details.License = v;
          },
          display: "text",
        },
        keywords: {
          key: "keywords",
          label: this.$gettext("Keywords"),
          read: (p) => p?.Details?.Keywords,
          write: (p, v) => {
            if (p?.Details) p.Details.Keywords = v;
          },
          display: "text",
        },
        notes: {
          key: "notes",
          label: this.$gettext("Notes"),
          read: (p) => p?.Details?.Notes,
          write: (p, v) => {
            if (p?.Details) p.Details.Notes = v;
          },
          display: "html",
          htmlValue: "notesHtml",
        },
      };
    },
    detailsFields() {
      return ["subject", "artist", "copyright", "license"].map((k) => this.fieldRegistry[k]);
    },
    textFields() {
      return ["keywords", "notes"].map((k) => this.fieldRegistry[k]);
    },
    showDetailsSection() {
      if (this.restrictedRole) return false;
      if (this.isEditable) return true;
      return this.detailsFields.some((f) => Boolean(f.read(this.photo)));
    },
    placeName() {
      if (!this.photo) return "";
      return this.photo.locationInfo() || "";
    },
    altitude() {
      if (!this.photo || !this.photo.Altitude) return "";
      return this.photo.Altitude + " m";
    },
    // Returns the lat/lng (shortened, with optional altitude) for the
    // combined Location row. Restricted users see only the lat/lng so
    // altitude isn't leaked through the sidebar.
    coordinatesLine() {
      if (!this.model?.Lat || !this.model?.Lng) return "";
      const coords = this.model.getLatLngShort();
      if (this.altitude && !this.restrictedRole) {
        return `${coords}\u2002${this.altitude}`;
      }
      return coords;
    },
    // True when the combined Location row should render — i.e., the row
    // has coordinates to display, OR the user is non-restricted and has
    // a place name / can edit a missing location.
    locationRowVisible() {
      if (this.model?.Lat && this.model?.Lng) return true;
      if (this.restrictedRole) return false;
      if (this.placeName) return true;
      return this.isEditable && this.featPlaces;
    },
    // Returns the merged row's title: the place name when allowed (non-
    // restricted, available), otherwise the coordinates line so the row
    // never renders empty when the v-if gate has admitted it.
    locationTitle() {
      if (!this.restrictedRole && this.placeName) return this.placeName;
      if (this.coordinatesLine) return this.coordinatesLine;
      return this.$gettext("Unknown");
    },
    // Returns the merged row's subtitle: the coordinates line, but only
    // when the title already shows the place name (so we don't render
    // the coordinates twice on a restricted / no-placeName row).
    locationSubtitle() {
      if (this.restrictedRole) return null;
      if (this.placeName && this.coordinatesLine) return this.coordinatesLine;
      return null;
    },
    // True when the combined Location row has any click action: editing
    // the location, copying the coordinates, or copying the place name.
    // Drives the .clickable cursor class on the row.
    locationRowClickable() {
      if (this.isEditable && this.featPlaces) return true;
      if (this.model?.Lat && this.model?.Lng) return true;
      return !this.restrictedRole && !!this.placeName;
    },
    // Returns the user-facing file path. For video, Live, and Animated
    // photos the primary file is the generated JPEG cover (used for
    // indexing and thumbnails), not the media file the user uploaded —
    // surface the underlying media file's Name so the sidebar shows the
    // .mp4 / .mov / .gif instead of "...mp4.jpg". The cards view uses the
    // same originalFile() routing via Photo.getOriginalName().
    // Returns `null` (not `""`) for restricted sessions and the
    // no-data state so the merged file row's `:subtitle` binding skips
    // rendering an empty `<v-list-item-subtitle>` element — Vuetify
    // gates the subtitle on `props.subtitle != null`, so an empty
    // string would still render an empty slot. Keeping the gate in the
    // computed lets the template stay free of restricted-role checks.
    fileName() {
      if (this.restrictedRole || !this.photo) return null;
      if (typeof this.photo.originalFile === "function") {
        const original = this.photo.originalFile();
        if (original && original !== this.photo && original.Name) return original.Name;
      }
      if (this.photo.FileName) return this.photo.FileName;
      const primary = typeof this.photo.primaryFile === "function" ? this.photo.primaryFile() : null;
      return primary?.Name || null;
    },
    fileInfo() {
      if (this.photo) {
        switch (this.photo.Type) {
          case media.Video:
          case media.Live:
          case media.Animated:
            return this.photo.getVideoInfo();
          case media.Vector:
          case media.Document:
            return this.photo.getVectorInfo();
          default:
            return this.photo.getImageInfo();
        }
      }
      // Fallback for restricted roles: Thumb.getTypeInfo() produces
      // format, megapixels, and dimensions from the viewer endpoint data.
      if (this.model && typeof this.model.getTypeInfo === "function") {
        return this.model.getTypeInfo();
      }
      return "";
    },
    fileIcon() {
      switch (this.photo?.Type || this.model?.Type) {
        case media.Raw:
          return "mdi-raw";
        case media.Video:
          return "mdi-video";
        case media.Live:
          return "mdi-play-circle-outline";
        case media.Animated:
          return "mdi-file-gif-box";
        case media.Vector:
          return "mdi-vector-polyline";
        case media.Document:
          return "mdi-file-pdf-box";
        default:
          return "mdi-image-outline";
      }
    },
    // Localized media type label for the file row's tooltip. Falls back
    // to the generic "File" label so the tooltip never reads as empty
    // when the photo's Type is unknown or the row is rendered for a
    // restricted-role model that doesn't expose Type.
    fileTypeName() {
      const type = this.photo?.Type || this.model?.Type;
      return this.$util.typeName(type, this.$gettext("File"));
    },
  },
  watch: {
    people: {
      immediate: true,
      handler(markers) {
        this.syncMarkerDrafts(Array.isArray(markers) ? markers : []);
      },
    },
    newMarkerUid(uid) {
      if (!uid) return;
      this.$nextTick(() => this.focusMarkerInput(uid));
    },
  },
  mounted() {
    // Warm the typeahead options for editable sessions so the combobox
    // dropdown is populated by the time the user focuses the input.
    // The shared cache (common/typeahead-cache.js) deduplicates concurrent
    // callers, so a sidebar mount during an open batch-edit session adds
    // no extra network round-trips.
    if (this.isEditable && !this.restrictedRole) {
      this.loadChipOptions("labels");
      this.loadChipOptions("albums");
    }
  },
  methods: {
    close() {
      this.$emit("close");
    },
    // openDateTimeDialog mounts the date-and-time-edit dialog when the
    // session is editable. No-op otherwise, so callers (row @click and
    // pencil button) don't need to inline the gate.
    openDateTimeDialog() {
      if (!this.isEditable) return;
      this.dateTimeDialog = true;
    },
    // openCameraDialog mounts the camera-and-lens-edit dialog when the
    // session is editable. Shared by the Camera and Lens row icons and
    // the Camera pencil button.
    openCameraDialog() {
      if (!this.isEditable) return;
      this.cameraDialog = true;
    },
    // openLocationDialog mounts the location-edit dialog when the
    // session is editable and the `places` feature is enabled. Shared by
    // the Location row pencil and the row click in edit mode.
    openLocationDialog() {
      if (!this.isEditable || !this.featPlaces) return;
      this.locationDialog = true;
    },
    // onLocationRowClick dispatches the combined Location row's click:
    // (1) edit mode -> open the location dialog; (2) read-only with
    // coordinates -> copy them to the clipboard so they paste into
    // mapping tools; (3) read-only with only a place name -> copy the
    // place name. Mirrors the row-level click semantics used by the
    // text-metadata rows but adds a coordinates-first preference under
    // (2) so the existing copy-coords gesture survives the merge.
    onLocationRowClick() {
      if (this.isEditable && this.featPlaces) {
        this.openLocationDialog();
      } else if (this.model?.Lat && this.model?.Lng) {
        this.model.copyLatLng();
      } else if (this.placeName) {
        this.$util.copyText(this.placeName);
      }
    },
    // onTextRowClick routes a sidebar text-row click. In edit mode it
    // enters the inline editor for the given field (Title, Caption,
    // Subject, Artist, Copyright, License, Keywords, Notes). In
    // read-only mode it copies the displayed value to the clipboard so
    // users can grab metadata without pinch-selecting. No-op when the
    // row is already in edit mode or the value is empty.
    onTextRowClick(field, value) {
      if (this.editingField === field) return;
      if (this.isEditable) {
        this.startEditing(field);
      } else if (value) {
        this.$util.copyText(value);
      }
    },
    getFieldValue(field) {
      const f = this.fieldRegistry[field];
      if (!f) return "";
      const v = f.read(this.photo);
      return v == null ? "" : v;
    },
    setFieldValue(field, value) {
      const f = this.fieldRegistry[field];
      if (!f || !this.view?.photo) return;
      f.write(this.view.photo, value);
    },
    // Function ref shared by every inline editor. Vue invokes it with the
    // mounted component on mount and null on unmount; since each editor is
    // gated by a unique `editingField === '<key>'`, only one is mounted at
    // a time, so the latest non-null call always identifies the active one.
    setInlineEditorRef(el) {
      if (el) this._inlineEditorEl = el;
      else if (!this.editingField) this._inlineEditorEl = null;
    },
    fieldHtml(f) {
      if (!f || f.display !== "html" || !f.htmlValue) return "";
      return this[f.htmlValue] || "";
    },
    shouldShowFieldRow(f) {
      if (!f) return false;
      if (this.editingField === f.key) return true;
      if (this.isEditable) return true;
      if (f.display === "html") return Boolean(this.fieldHtml(f));
      return Boolean(f.read(this.photo));
    },
    startEditing(field) {
      if (this.editingField) {
        this.cancelEditing();
      }

      this.editingField = field;
      this.editOriginal = this.getFieldValue(field);
      this._editStartedAt = Date.now();

      this.$nextTick(() => {
        const editor = this._inlineEditorEl;
        if (editor && typeof editor.focus === "function") editor.focus();
      });
    },
    // Eye-icon click handler (non-editable users only — display mode
    // is read-only). Emits `toggle-face-marker-mode` so the lightbox
    // can flip between `null` and `FaceMarkerDisplay`. Gates only on
    // `markersBusy` so an in-flight marker write doesn't race a mode
    // toggle; intentionally does NOT gate on `isEditable` because
    // display mode doesn't require edit permission.
    onToggleFaceMarkerMode() {
      if (this.markersBusy) return;
      this.$emit("toggle-face-marker-mode");
    },
    // Pencil-icon click handler (editable users only — draw mode adds
    // drag-to-create + click-to-remove). Emits `toggle-face-marker-draw`
    // so the lightbox can flip between `null` and `FaceMarkerDraw`.
    onToggleFaceMarkerDraw() {
      if (!this.isEditable || this.markersBusy) return;
      this.$emit("toggle-face-marker-draw");
    },
    onEjectMarker(marker) {
      if (!this.isEditable || this.markersBusy || !marker || !marker.SubjUID) return;
      // When the user types a fresh name and clicks ⏏ on the eject icon:
      // clearSubject flips SubjUID, the combobox re-renders editable, and
      // the implicit @blur from the re-render would commit the typed name
      // we just rejected. The timestamp lets confirmMarkerName bail when
      // an icon click triggered the unmount (P1-7).
      this._lastDestructiveMarkerActionAt = Date.now();
      this.$emit("eject-marker", marker);
    },
    // Combobox can bind either the typed string or the selected subject object.
    unwrapMarkerName(value) {
      return typeof value === "object" && value !== null ? value.Name || "" : value || "";
    },
    // Reconciles the local markerDrafts map with the latest markers from
    // photo.getMarkers(true). Fires from the `people` watcher on every
    // photo-cache mutation. The `editing` flag is the guard for P1-3:
    // when the user is actively typing into a marker's combobox, an
    // unrelated WS-driven photo update (vision worker, label change,
    // etc.) re-runs this method but must NOT overwrite the typed text.
    // confirmMarkerName / cancelMarkerName / onAddNameConfirm clear the
    // flag after committing or discarding the edit.
    syncMarkerDrafts(markers) {
      const seen = new Set();
      markers.forEach((m) => {
        if (!m || !m.UID) return;
        seen.add(m.UID);
        const original = m.Name || "";
        const existing = this.markerDrafts[m.UID];
        if (!existing) {
          this.markerDrafts[m.UID] = { original, current: original, editing: false };
        } else if (existing.original !== original) {
          existing.original = original;
          if (!existing.editing) {
            existing.current = original;
          }
        }
      });
      Object.keys(this.markerDrafts).forEach((uid) => {
        if (!seen.has(uid)) delete this.markerDrafts[uid];
      });
    },
    markerInputValue(uid) {
      const d = this.markerDrafts[uid];
      return d ? d.current : "";
    },
    markerInputSearch(uid) {
      return this.unwrapMarkerName(this.markerInputValue(uid));
    },
    // Records typed text from the v-combobox into the per-marker draft and
    // flips `editing` to true so a concurrent WS update can't snap the
    // input back to the backend value mid-keystroke (P1-3). The flag is
    // cleared by commitMarkerName / cancelMarkerName / onAddName* once
    // the edit reaches a settled state.
    setMarkerInputValue(uid, value) {
      if (!uid) return;
      if (!this.markerDrafts[uid]) {
        this.markerDrafts[uid] = { original: "", current: value, editing: true };
      } else {
        this.markerDrafts[uid].current = value;
        this.markerDrafts[uid].editing = true;
      }
    },
    focusMarkerInput(uid) {
      if (!uid) return;
      this.$emit("naming-started");
      this.$nextTick(() => {
        const input = this.$el && this.$el.querySelector(`[data-marker-uid="${uid}"] input`);
        if (input) input.focus();
      });
    },
    // Match a typed name against knownPeople case-insensitively so the backend
    // doesn't create a duplicate subject when the input only differs in case.
    // Locale `undefined` defers to the user's active locale so case-folding
    // works for Turkish dotted/dotless i, German ß, Cyrillic, Hebrew, etc.
    findKnownPerson(name) {
      if (!name) return null;
      return this.knownPeople.find((p) => p && p.Name && p.Name.localeCompare(name, undefined, { sensitivity: "base" }) === 0) || null;
    },
    // Resolves a marker by UID from the current photo. Used by the
    // Add-name dialog confirm path so a stale Marker reference held in
    // addNameDialog can't write through to a marker that no longer
    // exists on this photo (P1-8). Returns null when the marker has been
    // removed (e.g. rejected) or the slide moved to a different photo.
    findMarker(uid) {
      if (!uid || !this.photo || typeof this.photo.getMarkers !== "function") return null;
      return this.photo.getMarkers(true).find((m) => m && m.UID === uid) || null;
    },
    // Commits typed text from the per-marker draft. Source "enter" fires
    // from the keyboard, "blur" from focus loss; the latter routes through
    // an Add-name confirmation when the marker is still unnamed and the
    // typed name doesn't match an existing subject.
    //
    // Gating (P1-6 + P1-7):
    //   - `markersBusy`: another marker mutation is in flight; bail.
    //   - `marker.Invalid`: marker was rejected; bail before re-committing.
    //   - destructive-action timestamp: an × / ⏏ icon was clicked within
    //     the last 200ms (which destroys this row); bail.
    confirmMarkerName(marker, source = "enter") {
      if (!marker || !marker.UID) return;
      if (this.markersBusy) return;
      if (marker.Invalid) return;
      if (this._lastDestructiveMarkerActionAt && Date.now() - this._lastDestructiveMarkerActionAt < 200) return;
      const draft = this.markerDrafts[marker.UID];
      if (!draft) return;
      const name = this.unwrapMarkerName(draft.current).trim();
      const original = (draft.original || "").trim();

      if (!name || name === original) {
        // Reaching here means the user blurred without changing anything
        // (or restored the original). The draft is settled; clear the
        // editing flag so a concurrent WS update can re-sync.
        draft.editing = false;
        return;
      }
      if (typeof marker.setName !== "function") return;

      const match = this.findKnownPerson(name);

      // Blur without Enter on an unnamed marker → ask before committing a new
      // name. Skip the dialog if the person already exists (match) or if the
      // marker is already named (eject/rename path) — both are unambiguous.
      if (source === "blur" && !marker.SubjUID && !match) {
        this.addNameDialog = { visible: true, markerUid: marker.UID, name };
        return;
      }

      this.commitMarkerName(marker, match, name);
    },
    commitMarkerName(marker, match, name) {
      const draft = this.markerDrafts[marker.UID];
      if (!draft) return;

      if (match) {
        marker.Name = match.Name;
        marker.SubjUID = match.UID;
      } else {
        marker.Name = name;
      }

      // Lock the draft to the saved name so a parallel people-reload watcher
      // tick doesn't snap the input back to the old value mid-request.
      draft.original = marker.Name || "";
      draft.current = marker.Name || "";
      draft.editing = false;

      marker
        .setName()
        .then(() => {
          this.$emit("reload-markers", marker);
        })
        .catch(() => {
          this.$notify.error(this.$gettext("Failed to save name"));
        });
    },
    onPickPerson(marker, value) {
      if (!marker || !value || typeof value !== "object" || !value.Name) return;
      this.setMarkerInputValue(marker.UID, value.Name);
      this.confirmMarkerName(marker, "enter");
    },
    // Resolves the dialog's stored markerUid against the live photo. If
    // the marker has been rejected, navigated away from, or otherwise
    // disappeared, the commit is dropped silently — the dialog already
    // closed, and surfacing an error for a self-resolved state would just
    // confuse the user.
    onAddNameConfirm() {
      const { markerUid, name } = this.addNameDialog;
      this.addNameDialog = { visible: false, markerUid: "", name: "" };
      if (!markerUid || !name) return;
      const marker = this.findMarker(markerUid);
      if (!marker) return;
      this.commitMarkerName(marker, this.findKnownPerson(name), name);
    },
    onAddNameCancel() {
      const { markerUid } = this.addNameDialog;
      this.addNameDialog = { visible: false, markerUid: "", name: "" };
      const draft = markerUid ? this.markerDrafts[markerUid] : null;
      if (draft) {
        draft.current = draft.original || "";
        draft.editing = false;
      }
    },
    cancelMarkerName(marker) {
      if (!marker || !marker.UID) return;
      const draft = this.markerDrafts[marker.UID];
      if (!draft) return;
      draft.current = draft.original;
      draft.editing = false;
      // Blur the marker's own input so the user gets a visual cue the edit
      // was dropped; @blur re-fires confirmMarkerName but it's a no-op
      // now that current === original AND editing is false. Scoped to the
      // marker's input (P1-9) rather than document.activeElement so an
      // unrelated focused element isn't blurred by mistake.
      const input = this.$el && this.$el.querySelector(`[data-marker-uid="${marker.UID}"] input`);
      if (input && typeof input.blur === "function") input.blur();
    },
    resetInlineEdits() {
      if (this.editingField) this.cancelEditing();
      this.resetChipState();
      this.clearChipInput();
      Object.keys(this.markerDrafts).forEach((uid) => {
        const d = this.markerDrafts[uid];
        if (d) {
          d.current = d.original;
          d.editing = false;
        }
      });
      if (this.addNameDialog && this.addNameDialog.visible) {
        this.addNameDialog = { visible: false, markerUid: "", name: "" };
      }
    },
    // Inline text fields (title/caption/subject/...) are excluded on purpose:
    // onInlineFieldBlur() auto-commits them before any navigation source can
    // fire, so they can never have pending state at nav time. Chip-section
    // removals (`chipState.<field>.removals`) ARE counted here because the
    // user can see and toggle them, but `confirmDiscardPending` auto-commits
    // them before checking this — by the time the dialog gate runs they're
    // already gone. The remaining staged inputs that DO open the dialog are
    // marker drafts, typed-but-uncommitted combobox text, and the open
    // Add-name confirmation.
    hasPendingEdit() {
      for (const uid of Object.keys(this.markerDrafts)) {
        const d = this.markerDrafts[uid];
        if (!d) continue;
        if (this.unwrapMarkerName(d.current).trim() !== (d.original || "").trim()) return true;
      }
      // Pending chip removals (staged via × icon) and typed-but-uncommitted
      // text in the always-visible combobox both count as pending. Pressing
      // Enter would fire the instant-save path (addLabelImmediate /
      // addAlbumImmediate), but until then the characters live only in
      // chipState.<field>.search and would silently vanish on navigation.
      if (Object.values(this.chipState).some((s) => s.removals.length || (s.search || "").trim() !== "")) return true;
      // An open Add-name confirmation for an unnamed marker is also pending
      // input until the user picks Add or Cancel.
      if (this.addNameDialog && this.addNameDialog.visible) return true;
      return false;
    },
    // Fire-and-forget commit of any pending chip removals. Mirrors the
    // inline-text auto-commit on blur: the user's intent (clicking ×) is
    // honored on navigation/close instead of being silently discarded.
    // Each Photo.removeLabel / Photo.removeFromAlbum call captures
    // `this.photo` at the time of invocation, so the response patches the
    // original Photo instance even if the slide has changed by the time the
    // promise resolves. The catch path surfaces `$notify.error` and is
    // shared with the manual ✓ Confirm path through confirmLabels /
    // confirmAlbums.
    autoCommitChipRemovals() {
      if (this.chipState.labels.removals.length) {
        this.confirmLabels();
      }
      if (this.chipState.albums.removals.length) {
        this.confirmAlbums();
      }
    },
    // Async guard used by the lightbox before closing / hiding / navigating.
    // Returns a Promise<boolean>: true = safe to proceed, false = user
    // canceled. Pending chip removals auto-commit BEFORE the dialog gate,
    // so the discard prompt only fires for state the user could plausibly
    // still want to keep (marker drafts, typed combobox text, the Add-name
    // dialog) — never for chip × clicks, which are deliberate and final.
    confirmDiscardPending() {
      this.autoCommitChipRemovals();
      if (!this.hasPendingEdit()) return Promise.resolve(true);
      if (this.discardDialog.visible && this.discardDialog.resolver) {
        // Another request is already waiting on the dialog; reuse it.
        return new Promise((resolve) => {
          const prev = this.discardDialog.resolver;
          this.discardDialog.resolver = (ok) => {
            prev(ok);
            resolve(ok);
          };
        });
      }
      return new Promise((resolve) => {
        this.discardDialog.resolver = resolve;
        this.discardDialog.visible = true;
      });
    },
    onDiscardConfirm() {
      this.discardDialog.visible = false;
      const r = this.discardDialog.resolver;
      this.discardDialog.resolver = null;
      this.resetInlineEdits();
      if (r) r(true);
    },
    onDiscardCancel() {
      this.discardDialog.visible = false;
      const r = this.discardDialog.resolver;
      this.discardDialog.resolver = null;
      if (r) r(false);
    },
    confirmField() {
      if (!this.photo || !this.canEdit) {
        this.editingField = null;
        return;
      }

      const field = this.editingField;
      this.editingField = null;
      this.editOriginal = null;

      if (!this.photo.wasChanged()) {
        return;
      }

      // The inline-edit binding (v-model="photo.X") already mutated the photo
      // optimistically; on success sync the matching Thumb fields, on
      // failure roll back so the user sees the title/caption revert and
      // gets a notification instead of a silent ghost edit.
      this.photo
        .update()
        .then(() => {
          if ((field === "title" || field === "caption") && this.view?.model) {
            this.view.model.Title = this.view.photo.Title;
            this.view.model.Caption = this.view.photo.Caption;
          }
        })
        .catch(() => {
          this.photo.rollback();
          this.$notify.error(this.$gettext("Failed to save changes"));
        });
    },
    cancelEditing() {
      // Guard: clicking a pencil icon triggers blur on the previous field before focus lands
      // on the new input, which would immediately cancel the edit we just started.
      if (this._editStartedAt && Date.now() - this._editStartedAt < 200) {
        return;
      }

      if (this.editingField && this.editOriginal !== null) {
        this.setFieldValue(this.editingField, this.editOriginal);
      }

      this.editingField = null;
      this.editOriginal = null;
      this._editStartedAt = null;
    },
    // Blur handler for inline text fields (title/caption/subject/artist/
    // copyright/license/keywords/notes). Commits the edit instead of
    // silently reverting so the user does not lose their typing when
    // they click away, swipe to the next photo, or press the nav arrow.
    // Escape still cancels explicitly via cancelEditing().
    onInlineFieldBlur() {
      if (this._editStartedAt && Date.now() - this._editStartedAt < 200) {
        return;
      }
      if (!this.editingField) return;
      this.confirmField();
    },
    formatTime(model) {
      if (!model || !model.TakenAtLocal) {
        return this.$gettext("Unknown");
      }

      // Always parse as UTC to avoid time shifts
      const dateTime = DateTime.fromISO(model.TakenAtLocal, { zone: "UTC" });

      if (model.TimeZone && model.TimeZone !== "Local" && model.TimeZone !== "UTC") {
        // We use the real timezone just for display, but don't shift the time (prevents double timezone offset as backend already applied it)
        return dateTime.setZone(model.TimeZone, { keepLocalTime: true }).toLocaleString(formats.DATETIME_MED_TZ);
      } else {
        return dateTime.toLocaleString(formats.DATETIME_MED);
      }
    },
    openInNewTab(route) {
      if (!route) return;
      const resolved = this.$router ? this.$router.resolve(route) : null;
      const href = resolved?.href || "";
      if (!href || typeof window === "undefined" || typeof window.open !== "function") return;
      window.open(href, "_blank", "noopener,noreferrer");
    },
    navigateToLabel(label) {
      if (!label) return;
      const slug = label.CustomSlug || label.Slug;
      if (!slug) return;
      this.openInNewTab({ name: "browse", query: { q: "label:" + slug } });
    },
    navigateToAlbum(album) {
      if (!album || !album.UID) return;
      this.openInNewTab({ name: "album", params: { album: album.UID, slug: "view" } });
    },
    navigateToPerson(marker) {
      if (!marker) return;
      if (marker.SubjUID) {
        this.openInNewTab({ name: "browse", query: { q: "subject:" + marker.SubjUID } });
      } else if (marker.Name) {
        this.openInNewTab({ name: "browse", query: { q: "person:" + marker.Name } });
      }
    },
    // Pulls the typeahead suggestions from the shared module-scope
    // cache (`common/typeahead-cache.js`). Fired on combobox @focus —
    // cheap when the cache is warm (returns the same array reference)
    // and refreshes after WS-driven evictions (`labels.updated` /
    // `albums.updated` / `config.updated`) without per-component
    // subscriptions. Errors are swallowed so a transient network hiccup
    // never blocks the editor.
    //
    // The cache returns whatever order the backend emitted (which is
    // not reliably alphabetical even when search?order=name is set),
    // so we sort client-side via locale-aware localeCompare. This
    // also keeps Hebrew/Cyrillic libraries readable, where byte-order
    // sort would not match the user's expectation.
    loadChipOptions(field) {
      const collator = (key) => (a, b) => (a[key] || "").localeCompare(b[key] || "", undefined, { sensitivity: "base", numeric: true });
      if (field === "labels") {
        typeaheadCache
          .getLabels()
          .then((models) => {
            this.chipState.labels.options = models.map((l) => ({ Name: l.Name, UID: l.UID })).sort(collator("Name"));
          })
          .catch(() => {});
      } else if (field === "albums") {
        typeaheadCache
          .getAlbums()
          .then((models) => {
            // Map to plain {Title, UID} objects so v-combobox doesn't
            // try to track the rich Album model instance internally —
            // its reactive metadata (getters, methods, _request slots)
            // breaks v-combobox's input handling and the user can't
            // type. Mirrors the labels mapping above.
            this.chipState.albums.options = models.map((a) => ({ Title: a.Title, UID: a.UID })).sort(collator("Title"));
          })
          .catch(() => {});
      }
    },
    // Clears the typed text and selection for one combobox. The key
    // bump force-remounts the v-combobox / v-autocomplete so any stale
    // dropdown state (a half-rendered no-data row, a tracked input
    // value Vuetify retained after the model went null) goes with it.
    clearChipInput(field) {
      if (!field) {
        // Legacy callers without a field argument clear both —
        // cancelEditing() takes this path during inline-text rollback.
        Object.keys(this.chipState).forEach((f) => this.clearChipInput(f));
        return;
      }
      const state = this.chipState[field];
      if (!state) return;
      state.input = null;
      state.search = "";
      state.key++;
    },
    // Esc inside a chip combobox clears the typed text and the staged
    // pending removals for that field, then drops focus from the input.
    // Matches the inline-text Esc semantic (revert to baseline) without
    // crossing into editingField (chip sections no longer have one).
    onChipEscape(field) {
      const state = this.chipState[field];
      if (!state) return;
      state.removals = [];
      this.clearChipInput(field);
    },
    // Generic chip-state helpers. Field is "labels" or "albums"; the key is
    // whatever uniquely identifies a chip in that field's domain (Label.ID
    // for labels, Album.UID for albums).
    isChipPendingRemoval(field, key) {
      const state = this.chipState[field];
      return Boolean(state && key != null && state.removals.includes(key));
    },
    togglePendingChipRemoval(field, key) {
      const state = this.chipState[field];
      if (!state || key == null) return;
      const idx = state.removals.indexOf(key);
      if (idx >= 0) {
        state.removals.splice(idx, 1);
      } else {
        state.removals.push(key);
      }
    },
    // Clears all pending removals for one chip section in a single click.
    // Wired to the Undo icon in the section toolbar; restores soft-removed
    // chips by emptying `chipState.<field>.removals`, which makes the
    // `visibleLabels` / `visibleAlbums` computeds repopulate reactively.
    undoChipRemovals(field) {
      const state = this.chipState[field];
      if (!state) return;
      state.removals = [];
    },
    resetChipState() {
      Object.values(this.chipState).forEach((s) => {
        s.removals = [];
      });
    },
    // Click + Enter behavior on a primary chip: navigate to the related
    // label / album page in both read-only and editable contexts — the chip
    // acts like a link. Removal lives on the × icon's own click handler and
    // the keyboard Delete / Backspace path via `onChipDelete`. The two chip
    // shapes differ: labels are wrapped (`{ Label: { ID, ... } }`) while
    // albums come through directly (`{ UID, ... }`).
    onChipActivate(field, item) {
      if (!item) return;
      if (field === "labels") return this.navigateToLabel(item.Label);
      if (field === "albums") return this.navigateToAlbum(item);
    },
    // Removal entry point: wired to the chip's × icon click and to the
    // chip's keyboard Delete / Backspace handler. No-op outside edit mode
    // so read-only chips behave as plain links; in edit mode it toggles
    // pending removal so the chip disappears from `visibleLabels` /
    // `visibleAlbums` until the user clicks Undo or auto-commit runs.
    onChipDelete(field, item) {
      if (!item || !this.isEditable) return;
      const key = field === "labels" ? item?.Label?.ID : item.UID;
      this.togglePendingChipRemoval(field, key);
    },
    // Validates `rawName` and, when valid, fires `Photo.addLabel(name)`
    // immediately. The model method chains `setValues(r.data)` so the new
    // label appears as a real primary chip on `this.photo.Labels` as soon
    // as the response lands — there's no transient pending-add chip. On
    // rejection the caller sees a $notify.error and the chip never appears.
    // Returns true when a save was triggered (caller clears the input).
    addLabelImmediate(rawName) {
      if (!this.photo) return false;
      const name = (rawName || "").trim();
      if (!name) return false;
      if (name.length > this.$config.get("clip")) {
        this.$notify.error(this.$gettext("Name too long"));
        return false;
      }
      const norm = this.$util.normalizeTitle(name);
      if (!norm) return false;
      // Already on the photo? Skip the API call.
      if (this.labels.some((l) => this.$util.normalizeTitle(l?.Label?.Name) === norm)) return false;
      // Match against the system-wide label list — if a normalized match
      // exists, send the canonical existing-label name so the backend
      // doesn't create a near-duplicate (e.g. typed `Hello Cat` reuses an
      // existing `hello-cat` label) and the user sees the canonical casing.
      const existing = this.chipState.labels.options.find((l) => this.$util.normalizeTitle(l.Name) === norm);
      const finalName = existing ? existing.Name : name;
      this.photo.addLabel(finalName).catch(() => {
        this.$notify.error(this.$gettext("Failed to save changes"));
      });
      return true;
    },
    albumTitleConflicts(norm) {
      if (!norm) return true;
      return this.albums.some((a) => this.$util.normalizeTitle(a?.Title) === norm);
    },
    // Validates `album` and, when valid, fires `Photo.addToAlbum(uid)`
    // immediately. The model method evicts the LRU cache and refetches
    // the canonical photo so `this.photo.Albums` repopulates with the
    // saved state — no transient pending-add chip. Caller in onAlbumEnter
    // wraps brand-new albums in `Album.save()` first so a UID exists.
    addAlbumImmediate(album) {
      if (!this.photo) return false;
      if (!album || typeof album !== "object" || !album.UID) return false;
      const title = (album.Title || "").trim();
      if (!title) return false;
      if (title.length > this.$config.get("clip")) {
        this.$notify.error(this.$gettext("Name too long"));
        return false;
      }
      if (this.albums.some((a) => a.UID === album.UID)) return false;
      const norm = this.$util.normalizeTitle(title);
      if (this.albumTitleConflicts(norm)) return false;
      this.photo.addToAlbum(album.UID).catch(() => {
        this.$notify.error(this.$gettext("Failed to save changes"));
      });
      return true;
    },
    onLabelSelected(value) {
      if (!value || typeof value !== "object" || !value.Name) return;
      this.addLabelImmediate(value.Name);
      this.clearChipInput("labels");
    },
    // Read the typed name from the per-field search ref. The ev.target
    // fallback guards against Vuetify clearing `search` on the same Enter
    // keystroke we handle, which would otherwise drop the pending addition.
    pendingChipName(field, ev) {
      const search = this.chipState[field]?.search;
      if (search) return search;
      const target = ev && ev.target ? ev.target : null;
      return target && typeof target.value === "string" ? target.value : "";
    },
    // Enter inside the Labels combobox. Mirrors onAlbumEnter's structure:
    // empty → no-op; too-long → notify and leave the typed text so the
    // user can fix it; otherwise hand off to addLabelImmediate (which
    // does its own normalize + canonicalize + already-on-photo dedup) and
    // ALWAYS clear the input + bump the key so the menu closes — even
    // when the label was already on the photo and no API call fired. The
    // earlier `if (addLabelImmediate(...)) { clear }` shape left the
    // input populated and the menu open in the already-on-photo case,
    // which felt unresolved compared to the Albums combobox.
    onLabelEnter(ev) {
      const search = this.pendingChipName("labels", ev).trim();
      if (!search) return;

      if (search.length > this.$config.get("clip")) {
        this.$notify.error(this.$gettext("Name too long"));
        return;
      }

      this.addLabelImmediate(search);
      this.clearChipInput("labels");
    },
    // Confirms pending REMOVALS via Photo.removeLabel — additions take the
    // instant-save path (addLabelImmediate) and never reach this method.
    confirmLabels() {
      if (!this.photo) return;

      const state = this.chipState.labels;
      const removals = state.removals.slice();
      state.removals = [];

      // photo.removeLabel chains .then((r) => this.setValues(r.data)) (see
      // model/photo.js), so a successful response repopulates
      // this.photo.Labels with the backend-provided list. The websocket
      // photos.updated subscriber additionally evicts the cached entry via
      // evictCachedFromEntities, so the next read after navigation
      // rehydrates from GET /photos/:uid. confirmAlbums takes a different
      // path because album mutations only emit albums.updated (not
      // photos.updated) — see Photo.removeFromAlbum / addToAlbum for the
      // explicit evict+refind. The asymmetry is intentional.
      if (removals.length) {
        Promise.all(removals.map((id) => this.photo.removeLabel(id))).catch(() => {
          this.$notify.error(this.$gettext("Failed to save changes"));
        });
      }
    },
    // Confirms pending REMOVALS via Photo.removeFromAlbum — additions take
    // the instant-save path (addAlbumImmediate) and never reach this method.
    confirmAlbums() {
      if (!this.photo) return;

      const state = this.chipState.albums;
      const removals = state.removals.slice();
      state.removals = [];

      // Photo.removeFromAlbum owns the evict+refind dance per call, so the
      // sidebar's this.photo.Albums reflects the saved state without an
      // extra Photo.evictCache + find here. See model/photo.js for the
      // contract; the per-call extra GET is acceptable for the typical
      // 1-2 chip removals at a time.
      if (removals.length) {
        Promise.all(removals.map((uid) => this.photo.removeFromAlbum(uid))).catch(() => {
          this.$notify.error(this.$gettext("Failed to save changes"));
        });
      }
    },
    onAlbumSelected(value) {
      // v-combobox emits update:model-value transiently while the user
      // types free text (the model can momentarily flip to a string or
      // a Title-only stub before settling). Bail silently on anything
      // that isn't a real album object — clearing the input here would
      // bump chipState.albums.key, force-remount the v-combobox, and
      // kill focus mid-keystroke. Free-text Enter is committed via
      // onAlbumEnter, which owns the canonical clear path.
      if (!value || typeof value !== "object" || !value.UID) return;
      this.addAlbumImmediate(value);
      this.clearChipInput("albums");
    },
    onAlbumEnter(ev) {
      const search = this.pendingChipName("albums", ev).trim();
      if (!search) return;

      if (search.length > this.$config.get("clip")) {
        this.$notify.error(this.$gettext("Name too long"));
        return;
      }

      const norm = this.$util.normalizeTitle(search);
      if (!norm) {
        this.clearChipInput("albums");
        return;
      }

      const options = this.chipState.albums.options;

      // Normalized exact-match against the full known-albums list first.
      // normalizeTitle ignores case and converts every punctuation character
      // to whitespace, so `Hello Cat`, `hello-cat`, `hello,cat`, and
      // `hello.CAT` all resolve to the same canonical album. This mirrors
      // the Labels validation pipeline. Substring fuzzy matching is
      // intentionally NOT applied here — typing `test` must not silently
      // merge into an existing `LRUTEST-ALBUM-…`. Users pick partial
      // matches via the dropdown (click or arrow-key + Enter on a
      // highlighted item, which fires `onAlbumSelected`).
      const exactMatch = options.find((a) => this.$util.normalizeTitle(a.Title) === norm);
      if (exactMatch) {
        this.onAlbumSelected(exactMatch);
        return;
      }

      // Skip the API round-trip if a normalized title clash already exists
      // among the photo's current albums.
      if (this.albumTitleConflicts(norm)) {
        this.clearChipInput("albums");
        return;
      }

      const album = new Album({ Title: search });

      album
        .save()
        .then(() => {
          if (album.UID && this.addAlbumImmediate(album)) {
            options.push(album);
          }
        })
        .catch(() => {})
        .finally(() => {
          this.clearChipInput("albums");
        });
    },
    confirmDateTime(data) {
      this.dateTimeDialog = false;

      const photo = this.view?.photo;
      if (!photo || !photo.UID || !this.canEdit) return;

      photo.Day = data.Day;
      photo.Month = data.Month;
      photo.Year = data.Year;
      photo.TimeZone = data.TimeZone;

      const localDate = photo.localDate(data.time);
      if (!localDate.isValid) return;

      const isoTime =
        localDate.toISO({
          suppressMilliseconds: true,
          includeOffset: false,
        }) + "Z";

      photo.TakenAtLocal = isoTime;

      if (photo.currentTimeZoneUTC()) {
        photo.TakenAt = isoTime;
      }

      photo
        .update()
        .then(() => {
          if (!this.view?.model) return;
          this.view.model.TakenAtLocal = photo.TakenAtLocal;
          this.view.model.TimeZone = photo.TimeZone;
        })
        .catch(() => {
          photo.rollback();
          this.$notify.error(this.$gettext("Failed to save changes"));
        });
    },
    confirmCamera(data) {
      this.cameraDialog = false;

      const photo = this.view?.photo;
      if (!photo || !photo.UID || !this.canEdit) return;

      photo.CameraID = data.CameraID;
      photo.LensID = data.LensID;
      photo.Iso = data.Iso;
      photo.Exposure = data.Exposure;
      photo.FNumber = data.FNumber;
      photo.FocalLength = data.FocalLength;

      // photo.update() resets __originalValues only on success; on rejection
      // the snapshot still holds the pre-mutation state, so rollback() puts
      // every field back without per-field bookkeeping here.
      photo.update().catch(() => {
        photo.rollback();
        this.$notify.error(this.$gettext("Failed to save changes"));
      });
    },
    confirmLocation(data) {
      this.locationDialog = false;

      const photo = this.view?.photo;
      if (!photo || !photo.UID || !this.canEdit) return;

      photo.Lat = data.lat;
      photo.Lng = data.lng;
      photo.PlaceSrc = "manual";

      if (data.location?.country) {
        photo.Country = data.location.country;
      }

      photo
        .update()
        .then(() => {
          if (!this.view?.model) return;
          this.view.model.Lat = photo.Lat;
          this.view.model.Lng = photo.Lng;
        })
        .catch(() => {
          photo.rollback();
          this.$notify.error(this.$gettext("Failed to save changes"));
        });
    },
  },
};
</script>
