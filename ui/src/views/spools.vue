<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, reactive, onMounted } from 'vue'
import { Temporal } from 'temporal-polyfill'

import { useSpoolsStore } from '../stores/spools'
import { useLocationsStore } from '../stores/locations'
import { useMaterialsStore } from '../stores/materials'
import { useColorsStore } from '../stores/colors'
import { useBrandsStore } from '../stores/brands'
import { useStoresStore } from '../stores/stores'


const spoolsStore = useSpoolsStore()
const locationsStore = useLocationsStore()
const materialsStore = useMaterialsStore()
const colorsStore = useColorsStore()
const brandsStore = useBrandsStore()
const storesStore = useStoresStore()

const columns: QTableColumn[] = [
    {
    name: 'brand',
    required: true,
    label: 'Brand/Store',
    align: 'left',
    field: 'brand',
    sortable: true,
    style: 'width: 15%',
  },
  {
    name: 'material',
    required: true,
    label: 'Material',
    align: 'left',
    field: 'material',
    sortable: true,
    style: 'width: 10%',
  },

  {
    name: 'colors',
    required: true,
    label: 'Color',
    align: 'left',
    field: 'colors',
    sortable: true,
    style: 'width: 10%',
  },
  {
    name: 'location',
    required: true,
    label: 'Location',
    align: 'left',
    field: 'location',
    sortable: true,
    style: 'width: 10%',
  },

  {
    name: 'current_weight',
    required: true,
    label: 'Weight',
    align: 'left',
    field: 'current_weight',
    sortable: true,
    style: 'width: 10%',
  },

  {
    name: 'price',
    required: true,
    label: 'Price',
    align: 'left',
    field: 'price',
    sortable: true,
    style: 'width: 10%',
  },
  {
    name: 'empty',
    required: true,
    label: 'Empty',
    align: 'left',
    field: 'empty',
    sortable: true,
    style: 'width: 5%',
  },
  {
    name: 'created_at',
    required: true,
    label: 'Created At',
    align: 'left',
    field: 'created_at',
    sortable: true,
    style: 'width: 10%',
  },
  {
    name: 'updated_at',
    required: true,
    label: 'Updated At',
    align: 'left',
    field: 'updated_at',
    sortable: true,
    style: 'width: 10%',
  },
  {
    name: 'actions',
    label: 'Actions',
    align: 'center',
    field: 'label',
    style: 'width: 10%',
    format: val => `${val}`,
  }
]

var spools = ref(spoolsStore.sorted)
var locations = ref(locationsStore.sorted)
var materials = ref(materialsStore.sorted)
var colors = ref(colorsStore.sorted)
var brands = ref(brandsStore.sorted)
var stores = ref(storesStore.sorted)

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
})

const editRowData = ref<any>({})

function editRow(id: number) {
  editRowData.value = spoolsStore.findByID(id)
}

function saveSpool() {
  if (editRowData.value.id === undefined || editRowData.value.id === null) {
    spoolsStore.create(editRowData.value).then(() => {
      resetEdit()
    })
  } else {
    spoolsStore.update(editRowData.value.id, editRowData.value).then(() => {
      resetEdit()
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

function resetEdit() {
  editRowData.value = {}
}

const pagination = ref({
  rowsPerPage: 0
})

function DateTime(date: string): string {

  let d = new Date(date)

  let dt = Temporal.PlainDateTime.from(date)
  let s = dt.toString()
  s = s.replace('T', ' ')
  return s.split('.')[0] as string
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

              <div v-else-if="col.name === 'empty'">{{ props.row.empty ? 'Yes' : 'No' }}</div>

            <div v-else-if="col.name === 'actions'">
              <q-btn flat color="primary" icon="mdi-pencil" size="xs" @click="editRow(props.row.id)" />
              <q-btn flat color="red" icon="mdi-delete" size="xs" @click="deleteSpool(props.row.id)" />
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <div class="form">
      <q-toggle label="Show Empty" v-model="spoolsStore.showEmpty" @update:model-value="toggleEmpty"/>
      <hr />
      <q-form v-if="editRowData != null" @submit="saveSpool(); editRowData = null" @reset="resetEdit()">
        <div class="text-h6 q-mb-md">{{ editRowData.id != null ? 'Edit spool' : 'Add new spool' }}</div>
        <div><q-select label="Brand" dark v-model="editRowData.brand" :options="brands" option-label="label" option-value="id" map-options emit-value /> </div>
        <div><q-select label="Store" dark v-model="editRowData.store" :options="stores" option-label="label" option-value="id" map-options emit-value /> </div>
        <div><q-select label="Material" dark v-model="editRowData.material" :options="materials" option-label="label" option-value="id" map-options emit-value /> </div>
        <div><q-select label="Location" dark v-model="editRowData.location" :options="locations" option-label="label" option-value="id" map-options emit-value /> </div>
        <div><q-input dark v-model="editRowData.current_weight" label="Current weight" type="number" /></div>
        <div><q-input dark v-model="editRowData.combined_weight" label="Combined weight" type="number" /></div>
        <div><q-input dark v-model="editRowData.weight" label="Total weight" type="number" /></div>
        <div><q-input dark v-model="editRowData.price" label="Price" type="number" /></div>
        <div>
          <q-select label="Colors" clearable dark v-model="editRowData.colors" :options="colors" option-label="label" option-value="id" map-options emit-value multiple />

        </div>

        <div><q-toggle label="Empty" v-model="editRowData.empty" /></div>
        <div class="q-mt-md">
          <q-btn label="Save" icon="mdi-content-save" type="submit" color="primary" />
          <q-btn label="Cancel" icon="mdi-undo" type="reset" color="secondary" class="q-ml-sm" />
        </div>
      </q-form>

    </div>
  </div>
</template>

<style scoped>
.row {
  display: flex;
  flex-direction: row;

}

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
    background-color: #222222;
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
}

.hex {
  width: 16px;
  height: 16px;
  display: inline-block;
  margin-right: 8px;
  border: 1px solid #ffffff77;
}
</style>
