<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import type { NewRating } from '../stores/ratings'
import { ref, onMounted, computed } from 'vue'
// import { Temporal } from 'temporal-polyfill'

import { useSpoolsStore } from '../stores/spools'
import { useLocationsStore } from '../stores/locations'
import { useMaterialsStore } from '../stores/materials'
import { useColorsStore } from '../stores/colors'
import { useBrandsStore } from '../stores/brands'
import { useRatingStore as useRatingsStore } from '../stores/ratings'
import { useStoresStore } from '../stores/stores'



const spoolsStore = useSpoolsStore()
const locationsStore = useLocationsStore()
const materialsStore = useMaterialsStore()
const colorsStore = useColorsStore()
const brandsStore = useBrandsStore()
const storesStore = useStoresStore()
const ratingsStore = useRatingsStore()

const spools = ref(spoolsStore.sorted)
const locations = ref(locationsStore.sorted)
const materials = ref(materialsStore.sorted)
const colors = ref(colorsStore.sorted)
const brands = ref(brandsStore.sorted)
const stores = ref(storesStore.sorted)
const ratings = ref(ratingsStore.sorted)

const tab = ref('browse')
const formTabName = ref('New Spool')
const editRowData = ref<any>({})
const rateSpoolData = ref<any>(null)
const pagination = ref({
  rowsPerPage: 0
})

const columns = computed<Array<QTableColumn>>(() => {
  const cols: Array<QTableColumn> = []
  cols.push({
    name: 'brand',
    required: true,
    label: 'Brand/Store',
    align: 'left',
    field: 'brand',
    sortable: true,
    sort: (a, b, rowA, rowB) => {
      const labelA = brandsStore.findByID(rowA?.brand)?.label ?? String(a ?? '')
      const labelB = brandsStore.findByID(rowB?.brand)?.label ?? String(b ?? '')
      return String(labelA).localeCompare(String(labelB))
    },
    style: 'width: clamp(200px, 15%, 260px);',
  })

  cols.push({
    name: 'material',
    required: true,
    label: 'Material',
    align: 'left',
    field: 'material',
    sortable: true,
    sort: (a, b, rowA, rowB) => {
      const labelA = materialsStore.findByID(rowA?.material)?.label ?? String(a ?? '')
      const labelB = materialsStore.findByID(rowB?.material)?.label ?? String(b ?? '')
      return String(labelA).localeCompare(String(labelB))
    },
    style: 'width: clamp(130px, 10%, 170px);',
  })

  cols.push({
    name: 'colors',
    required: true,
    label: 'Color',
    align: 'left',
    field: 'colors',
    sortable: true,
    sort: (a, b, rowA, rowB) => {
      const labelA = colorsStore.findByID(rowA?.colors[0])?.label ?? String(a ?? '')
      const labelB = colorsStore.findByID(rowB?.colors[0])?.label ?? String(b ?? '')
      return String(labelA).localeCompare(String(labelB))
    },
    style: 'width: clamp(130px, 10%, 170px);',
  })

  cols.push({
    name: 'location',
    required: true,
    label: 'Location',
    align: 'left',
    field: 'location',
    sortable: true,
    sort: (a, b, rowA, rowB) => {
      const labelA = locationsStore.findByID(rowA?.location)?.label ?? String(a ?? '')
      const labelB = locationsStore.findByID(rowB?.location)?.label ?? String(b ?? '')
      return String(labelA).localeCompare(String(labelB))
    },
    style: 'width: clamp(130px, 10%, 170px);',
  })

  cols.push({
    name: 'current_weight',
    required: true,
    label: 'Weight',
    align: 'left',
    field: 'current_weight',
    sortable: true,
    style: 'width: clamp(130px, 10%, 170px);',
  })

  cols.push({
    name: 'price',
    required: true,
    label: 'Price',
    align: 'left',
    field: 'price',
    sortable: true,
    style: 'width: clamp(130px, 11%, 170px);',
  })

  if (spoolsStore.showEmpty) {
    cols.push({
      name: 'emptied_at',
      required: true,
      label: 'Empty',
      align: 'left',
      field: 'emptied_at',
      sortable: true,
      style: 'width: clamp(140px, 8%, 150px);',
    })
    cols.push({
      name: 'rating',
      required: false,
      label: 'Rating',
      align: 'left',
      field: 'rating',
      sortable: true,
      rawSort: (a, b, rowA, rowB) => {
        const numA = ratingsStore.findBySpoolID(rowA?.id)?.rating ?? 0
        const numB = ratingsStore.findBySpoolID(rowB?.id)?.rating ?? 0

        return Number(numA) == Number(numB) ? 0 : Number(numA) > Number(numB) ? -1 : 1
      },
      style: 'width: clamp(65px, 5%, 85px);',
    })

  }

  cols.push({
    name: 'created_at',
    required: true,
    label: 'Created At',
    align: 'left',
    field: 'created_at',
    sortable: true,
    style: 'width: clamp(140px, 8%, 150px);',
  })

  cols.push({
    name: 'updated_at',
    required: true,
    label: 'Updated At',
    align: 'left',
    field: 'updated_at',
    sortable: true,
    style: 'width: clamp(140px, 8%, 150px);',
  })

  cols.push({
    name: 'actions',
    label: 'Actions',
    align: 'center',
    field: 'label',
    style: 'width: clamp(65px, 5%, 85px);',
    format: val => `${val}`,
  })

  return cols
})


onMounted(async () => {

  await locationsStore.find()
  locations.value = locationsStore.sorted

  await materialsStore.find()
  materials.value = materialsStore.sorted

  await colorsStore.find()
  colors.value = colorsStore.sorted

  await brandsStore.find()
  brands.value = brandsStore.sorted

  await storesStore.find()
  stores.value = storesStore.sorted

  await spoolsStore.find()
  spools.value = spoolsStore.sorted

  await ratingsStore.find()
  ratings.value = ratingsStore.sorted
})


function editRow(id: number) {
	  editRowData.value = spoolsStore.findByID(id)
	  tab.value = 'form'
	  formTabName.value = `Edit Spool`
}

function editRating(id: number) {
  rateSpoolData.value = { spool_id: id, rating: 0 }
  tab.value = 'form'
  formTabName.value = `Rate Spool`
}

function saveSpool() {

  if (editRowData.value.id === undefined || editRowData.value.id === null) {

    spoolsStore.create(editRowData.value).then(() => {
      if (editRowData.value.empty && (editRowData.value.emptied_at === undefined || editRowData.value.emptied_at === null)) {
        editRating(spoolsStore.lastCreatedID as number)
        editRowData.value = null
      } else {
        resetEdit()
      }
    })



  } else {
    spoolsStore.update(editRowData.value.id, editRowData.value).then(() => {
      if (editRowData.value.empty && (editRowData.value.emptied_at === undefined || editRowData.value.emptied_at === null)) {
        editRating(editRowData.value.id as number)
        editRowData.value = null
      } else {
        resetEdit()
      }
    })

  }
}

function deleteSpool(id: number) {
  if (confirm("Are you sure you want to delete this spool?")) {
	  spoolsStore.kill(id)
	  if (editRowData.value.id === id) {
	    resetEdit()
	  }

  }
}

function saveRating() {
  const rating: NewRating = {
    spool_id: rateSpoolData.value.spool_id,
    rating: rateSpoolData.value.rating,
  }

  ratingsStore.create(rating).then(() => {
    resetEdit()
  })
}

function resetEdit() {
  editRowData.value = {}
  rateSpoolData.value = null
  formTabName.value = 'New Spool'
}

// function DateTime(date: string): string {

//   let d = new Date(date)

//   let dt = Temporal.PlainDateTime.from(date)
//   let s = dt.toString()
//   s = s.replace('T', ' ')
//   return s.split('.')[0] as string
// }

function DateTime(date: string): string {
  const d = new Date(date)
  return d.toISOString().replace('T', ' ').split('.')[0] as string
}

async function toggleEmpty() {
  spoolsStore.toggleEmpty()

  await spoolsStore.find()
  spools.value = spoolsStore.sorted
}


</script>

<template>
  <div class="row">
    <q-table dark flat bordered binary-state-sort :rows="spools" :columns="columns" row-key="id" virtual-scroll
	    :rows-per-page-options="[0]" v-model:pagination="pagination" class="grid sticky-header">
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">{{ col.label }}</q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props"
	        :style="editRowData != null && editRowData.id === props.row.id ? 'background-color: #ffffff12' : ''">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <div v-if="col.name === 'brand'">
              <a v-if="spoolsStore.storeBrandLinks(props.row.id).brand_link !== ''" :href="spoolsStore.storeBrandLinks(props.row.id).brand_link">{{ brandsStore.findByID(props.row.brand)?.label }}</a><span v-else>{{ brandsStore.findByID(props.row.brand)?.label }}</span>

              <a v-if="spoolsStore.storeBrandLinks(props.row.id).store_link !== '' && spoolsStore.storeBrandLinks(props.row.id).store_link !== 'skip'" :href="spoolsStore.storeBrandLinks(props.row.id).store_link">{{ storesStore.findByID(props.row.store)?.label }}</a><span v-else-if="spoolsStore.storeBrandLinks(props.row.id).store_link !== 'skip'">{{ storesStore.findByID(props.row.store)?.label }}</span>
            </div>
            <div v-if="col.name === 'material'">
	            {{ materialsStore.findByID(props.row.material)?.label }}
            </div>
            <div v-if="col.name === 'colors'">
              <div v-for="color in props.row.colors" :key="color" class="q-mr-sm">
                <div class="hex" :style="{ backgroundColor: colorsStore.findByID(color)?.hex }"></div> {{ colorsStore.findByID(color)?.label }}
              </div>
            </div>
            <div v-if="col.name === 'location'">{{ locationsStore.findByID(props.row.location)?.label }}</div>
            <div v-if="col.name === 'current_weight'">{{ props.row.current_weight }} / {{ props.row.weight }} ({{ props.row.combined_weight }})</div>

            <div v-if="col.name === 'price'">{{ props.row.price }}</div>
            <div v-if="col.name === 'created_at'">
	            {{ DateTime(props.row.created_at) }}
            </div>
            <div v-if="col.name === 'updated_at'">
	            {{ DateTime(props.row.updated_at) }}
            </div>

            <div v-else-if="col.name === 'emptied_at' && spoolsStore.showEmpty">{{ props.row.empty ? DateTime(props.row.emptied_at) : '' }}</div>
            <div v-else-if="col.name === 'rating' && spoolsStore.showEmpty">{{   ratingsStore.findBySpoolID(props.row.id)?.rating > 0 ? ratingsStore.findBySpoolID(props.row.id)?.rating + '/5': '' }}</div>

            <div v-else-if="col.name === 'actions'">
              <q-btn dense flat color="primary" icon="mdi-pencil" size="xs" @click="editRow(props.row.id)" />
              <q-btn dense flat color="red" icon="mdi-delete" size="xs" @click="deleteSpool(props.row.id)" />
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <div class="form">
      <q-tabs v-model="tab" dense dark active-color="primary" indicator-color="primary" align="justify" narrow-indicator>
        <q-tab name="browse" label="Browse" />
        <q-tab name="form" :label="formTabName" />
      </q-tabs>
      <q-tab-panels dark v-model="tab">
        <q-tab-panel dark name="browse">
          <q-toggle label="Show Empty" v-model="spoolsStore.showEmpty" @update:model-value="toggleEmpty" left-label />
        </q-tab-panel>
        <q-tab-panel dark name="form">
          <q-form v-if="rateSpoolData !== null && editRowData === null" @submit="saveRating();" @reset="resetEdit()">
            <div>
              <q-rating v-model="rateSpoolData.rating" size="3.5em" color="yellow-5" icon="star" />
            </div>
            <div class="q-mt-md">
              <q-btn label="Save" icon="mdi-content-save" type="submit" color="primary" />
              <q-btn label="Cancel" icon="mdi-undo" type="reset" color="secondary" class="q-ml-sm" />
            </div>
          </q-form>
          <q-form v-if="rateSpoolData === null && editRowData !== null" @submit="saveSpool();" @reset="resetEdit()">
            <div class="text-h6 q-mb-md">{{ editRowData.id != null ? 'Edit spool' : 'Add new spool' }}</div>
            <div><q-select label="Brand" dark v-model="editRowData.brand" :options="brands" option-label="label" option-value="id" map-options emit-value hint="The brand of the filament on the spool" /></div>
            <div><q-select label="Store" dark v-model="editRowData.store" :options="stores" option-label="label" option-value="id" map-options emit-value hint="The store the filament was purchased from" /></div>
            <div><q-select label="Material" dark v-model="editRowData.material" :options="materials" option-label="label" option-value="id" map-options emit-value hint="The material type of the filament on the spool" /></div>
            <div><q-select label="Location" dark v-model="editRowData.location" :options="locations" option-label="label" option-value="id" map-options emit-value hint="Where the spool is currently located" /></div>
            <div><q-input dark v-model="editRowData.current_weight" label="Current weight" type="number" hint="The current combined weight of the filament and the spool" /></div>
            <div><q-input dark v-model="editRowData.combined_weight" label="Combined weight" type="number" hint="The weight of the spool and the filament" /></div>
            <div><q-input dark v-model="editRowData.weight" label="Weight" type="number" hint="The marketed weight of the filament on the spool" /></div>
            <div><q-input dark v-model="editRowData.price" label="Price" type="number" hint="The price paid for the spool of filament" /></div>
            <div>
              <q-select label="Colors" clearable dark v-model="editRowData.colors" :options="colors" option-label="label" option-value="id" map-options emit-value multiple hint="The color(s) of the filament on the spool" />
            </div>

            <div><q-toggle label="Empty" v-model="editRowData.empty" left-label /></div>
            <div class="q-mt-md">
              <q-btn label="Save" icon="mdi-content-save" type="submit" color="primary" />
              <q-btn label="Cancel" icon="mdi-undo" type="reset" color="secondary" class="q-ml-sm" />
            </div>
          </q-form>
        </q-tab-panel>

      </q-tab-panels>

    </div>
  </div>
</template>

<style scoped>

.form {
  width: 400px;
  margin-left: 16px;
}

.grid {
  flex: 1 1 calc(100% - 416px);
}


.sticky-header {

  /* height or max-height is important */
  & {
	  height: calc(100vh - (50px + 3rem));
  }

  .q-table__top,
  .q-table__bottom,
  thead tr:first-child th {
	  /* bg color is important for th; just specify one */
	  background-color: #111111;
  }

  thead tr th {
	  position: sticky;
	  z-index: 1;
  }

  thead tr:first-child th {
	  top: 0
  }

  /* this is when the loading indicator appears */
  &.q-table--loading thead tr:last-child th {
	  /* height of all previous header rows */
	  top: 48px;
  }

  /* prevent scrolling behind sticky top row on focus */
  tbody {
	  /* height of all previous header rows */
	  scroll-margin-top: 48px;
  }
  tbody.q-virtual-scroll__content {
	  tr.q-tr:nth-child(odd) {
	    background-color: #44444411;
	  }
  }
}

.q-table__bottom {
  background-color: #111111;
}

.hex {
  width: 16px;
  height: 16px;
  display: inline-block;
  margin-right: 8px;
  border: 1px solid #ffffff77;
}
</style>
