<template>
  <div class="p-tab p-tab-photo-labels">
    <v-form ref="form" class="p-form p-form--table p-form-photo-labels" validate-on="invalid-input" accept-charset="UTF-8" tabindex="-1" @submit.prevent>
      <div class="form-body">
        <div class="form-controls">
          <v-row dense align="start">
            <v-col cols="0" sm="2" class="form-thumb">
              <div>
                <img :alt="view?.model.Title" :src="view?.model.thumbnailUrl('tile_500')" class="clickable" @click.stop.prevent.exact="openPhoto()" />
              </div>
            </v-col>
            <v-col cols="12" sm="10" class="d-flex flex-column ga-4">
              <div
                :class="$vuetify.display.smAndDown ? 'v-table--density-compact' : 'v-table--density-comfortable'"
                class="v-table v-table--has-top v-table--hover v-data-table elevation-0 edit-table list-view"
              >
                <div class="v-table__wrapper">
                  <table>
                    <thead>
                      <tr>
                        <th class="v-data-table__td v-data-table-column--align-left v-data-table__th" colspan="1" rowspan="1">
                          <div class="v-data-table-header__content">
                            <span>{{ $gettext(`Label`) }}</span>
                          </div>
                        </th>
                        <th class="v-data-table__td v-data-table-column--align-left v-data-table__th" colspan="1" rowspan="1">
                          <div class="v-data-table-header__content">
                            <span>{{ $gettext(`Source`) }}</span>
                          </div>
                        </th>
                        <th class="v-data-table__td v-data-table-column--align-center v-data-table__th" colspan="1" rowspan="1">
                          <div class="v-data-table-header__content">
                            <span>{{ $gettext(`Confidence`) }}</span>
                          </div>
                        </th>
                        <th class="v-data-table__td v-data-table-column--align-center v-data-table__th" colspan="1" rowspan="1">
                          <div class="v-data-table-header__content">
                            <span>{{ $gettext(`Action`) }}</span>
                          </div>
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="label in view.model.Labels" :key="label.LabelID" class="label result">
                        <td class="text-start">
                          {{ label.Label.Name }}
                          <!--                  TODO: add this dialog later-->
                          <!--                  <v-dialog class="p-inline-edit" @save="renameLabel(props.item.Label)">-->
                          <!--                    {{ props.item.Label.Name }}-->
                          <!--                    <template #input>-->
                          <!--                      <v-text-field v-model="props.item.Label.Name" :rules="[nameRule]" :label="$gettext('Name')" color="surface-variant" class="input-rename background-inherit elevation-0" single-line autofocus variant="solo" hide-details></v-text-field>-->
                          <!--                    </template>-->
                          <!--                  </v-dialog>-->
                        </td>
                        <td class="text-start">
                          {{ sourceName(label.LabelSrc) }}
                        </td>
                        <td class="text-center">{{ 100 - label.Uncertainty }}%</td>
                        <td class="text-center">
                          <v-btn
                            v-if="disabled"
                            icon
                            density="comfortable"
                            variant="text"
                            :ripple="false"
                            class="action-view"
                            title="Search"
                            @click.stop.prevent="searchLabel(label.Label)"
                          >
                            <v-icon color="surface-variant">mdi-magnify</v-icon>
                          </v-btn>
                          <v-btn
                            v-else-if="(label.LabelSrc === 'manual' && label.Uncertainty < 100) || (label.LabelSrc === 'batch' && label.Uncertainty === 0)"
                            icon
                            density="comfortable"
                            variant="text"
                            :ripple="false"
                            class="action-delete"
                            title="Delete"
                            @click.stop.prevent="removeLabel(label.Label)"
                          >
                            <v-icon color="surface-variant">mdi-delete</v-icon>
                          </v-btn>
                          <v-btn
                            v-else-if="label.Uncertainty < 100"
                            icon
                            density="comfortable"
                            variant="text"
                            :ripple="false"
                            class="action-remove"
                            title="Remove"
                            @click.stop.prevent="removeLabel(label.Label)"
                          >
                            <v-icon color="surface-variant">mdi-minus</v-icon>
                          </v-btn>
                          <v-btn
                            v-else
                            icon
                            density="comfortable"
                            variant="text"
                            :ripple="false"
                            class="action-on"
                            title="Activate"
                            @click.stop.prevent="activateLabel(label.Label)"
                          >
                            <v-icon color="surface-variant">mdi-plus</v-icon>
                          </v-btn>
                        </td>
                      </tr>
                      <tr v-if="!disabled" class="label result">
                        <td class="text-start">
                          <!-- No autofocus: v-combobox auto-opens its menu on
                            focus, so autofocusing here would pop the dropdown
                            the moment the Labels tab renders (and lay the menu
                            out before the tab geometry settles, which
                            mispositions it). The user clicks the input when
                            they want to type. -->
                          <v-combobox
                            v-model="newLabelModel"
                            v-model:search="newLabel"
                            :items="labelOptions"
                            item-title="Name"
                            item-value="Name"
                            return-object
                            :rules="[nameRule]"
                            color="surface-variant"
                            autocomplete="off"
                            single-line
                            flat
                            variant="plain"
                            hide-details
                            hide-no-data
                            append-icon=""
                            :list-props="{ density: 'compact' }"
                            class="input-label ma-0 pa-0"
                            @focus="loadLabelOptions"
                            @update:model-value="onLabelSelected"
                            @keydown.enter.stop.prevent="addLabel"
                          ></v-combobox>
                        </td>
                        <td class="text-start">
                          {{ sourceName("manual") }}
                        </td>
                        <td class="text-center">100%</td>
                        <td class="text-center">
                          <v-btn icon density="comfortable" variant="text" :ripple="false" title="Add" class="p-photo-label-add" @click.stop.prevent="addLabel">
                            <v-icon color="surface-variant">mdi-plus</v-icon>
                          </v-btn>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </v-col>
          </v-row>
        </div>
      </div>
    </v-form>
  </div>
</template>

<script>
import Thumb from "model/thumb";
import typeaheadCache from "common/typeahead-cache";

export default {
  name: "PTabPhotoLabels",
  props: {
    uid: {
      type: String,
      default: "",
    },
  },
  emits: ["close"],
  data() {
    return {
      view: this.$view.getData(),
      disabled: !this.$config.feature("edit"),
      config: this.$config.values,
      readonly: this.$config.get("readonly"),
      selected: [],
      newLabel: "",
      newLabelModel: null,
      // Typeahead suggestions sourced lazily from the shared cache
      // (common/typeahead-cache.js) on first focus. Cache de-dupes
      // across the sidebar combobox and the batch-edit dialog.
      labelOptions: [],
      listColumns: [
        {
          title: this.$gettext("Label"),
          key: "",
          sortable: false,
          align: "left",
        },
        {
          title: this.$gettext("Source"),
          key: "LabelSrc",
          sortable: false,
          align: "left",
        },
        {
          title: this.$gettext("Confidence"),
          key: "Uncertainty",
          sortable: false,
          align: "center",
        },
        {
          title: this.$gettext("Action"),
          key: "",
          sortable: false,
          align: "center",
        },
      ],
      nameRule: (v) => v.length <= this.$config.get("clip") || this.$gettext("Name too long"),
    };
  },
  methods: {
    refresh() {},
    sourceName(s) {
      return this.$util.sourceName(s);
    },
    removeLabel(label) {
      if (!label || !this.view?.model) {
        return;
      }

      const name = label.Name;

      this.view.model.removeLabel(label.ID).then(() => {
        this.$notify.success("removed " + name);
      });
    },
    addLabel() {
      const typed = (this.newLabel || "").trim();
      if (!typed || !this.view?.model) {
        return;
      }

      // Apply the same canonical-match dedup the sidebar uses for L3:
      // typing `Hello Cat` resolves to an existing `hello-cat` label so
      // the backend isn't asked to create a near-duplicate. normalizeTitle
      // ignores case, strips punctuation, and treats `+`/`_`/`-` as
      // space.
      const norm = this.$util.normalizeTitle ? this.$util.normalizeTitle(typed) : typed.toLowerCase();
      let finalName = typed;
      if (norm) {
        const existing = this.labelOptions.find((l) =>
          this.$util.normalizeTitle ? this.$util.normalizeTitle(l.Name) === norm : (l.Name || "").toLowerCase() === norm
        );
        if (existing) {
          finalName = existing.Name;
        }
      }

      this.view.model.addLabel(finalName).then(() => {
        this.$notify.success("added " + finalName);
        this.newLabel = "";
        this.newLabelModel = null;
      });
    },
    // Selecting an existing label from the dropdown commits via the same
    // canonical-name path as addLabel — keeps the chip name consistent
    // with what's stored server-side.
    onLabelSelected(value) {
      if (!value || typeof value !== "object" || !value.Name) {
        return;
      }
      this.newLabel = value.Name;
      this.addLabel();
    },
    // Pulls suggestions from the shared cache on first focus. Cheap on
    // re-focus (cache hit) and refreshes after WS-driven evictions.
    loadLabelOptions() {
      typeaheadCache
        .getLabels()
        .then((models) => {
          this.labelOptions = models.map((l) => ({ Name: l.Name, UID: l.UID }));
        })
        .catch(() => {});
    },
    activateLabel(label) {
      if (!label || !this.view?.model) {
        return;
      }

      this.view.model.activateLabel(label.ID);
    },
    // TODO: add this dialog later
    // renameLabel(label) {
    //   if (!label) {
    //     return;
    //   }
    //
    //   this.view.model.renameLabel(label.ID, label.Name);
    // },
    searchLabel(label) {
      this.$router.push({ name: "all", query: { q: "label:" + label.Slug } }).catch(() => {});
      this.$emit("close");
    },
    openPhoto() {
      if (!this.view?.model) {
        return;
      }

      this.$lightbox.openModels(Thumb.fromPhotos([this.view.model]), 0);
    },
  },
};
</script>
