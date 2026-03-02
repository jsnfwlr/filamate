<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import { useStoresStore } from '../stores/stores'
import { useBrandsStore } from '../stores/brands'


const storesStore = useStoresStore()
const brandsStore = useBrandsStore()

const columns: QTableColumn[] = [
  {
    name: 'label',
    required: true,
    label: 'Brand',
    align: 'left',
    field: 'label',
    sortable: true,
    style: 'width: 40%',
  },
  {
    name: 'active',
    required: true,
    label: 'Active',
    align: 'left',
    field: 'active',
    sortable: true,
    style: 'width: 10%',
  },
  {
    name: 'store',
    label: 'Store',
    align: 'left',
    field: 'store',
    sortable: true,
    style: 'width: 40%',
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

const brands = ref(brandsStore.sorted)
const stores = ref(storesStore.sorted)

onMounted(async () => {
  await storesStore.find()
  await brandsStore.find()
  brands.value = brandsStore.sorted
  stores.value = storesStore.sorted
})

const editRowData = ref<any>({})

function editRow(id: number) {
  editRowData.value = brandsStore.findByID(id)
}

function saveBrand() {
  if (editRowData.value.id === undefined || editRowData.value.id === null) {
    brandsStore.create(editRowData.value).then( () => {
      resetEdit()
    })
  } else {
    brandsStore.update(editRowData.value.id, editRowData.value).then( () => {
      resetEdit()
    })
  }
}

function deleteBrand(id: number) {
  if (confirm("Are you sure you want to delete this brand?")) {
    brandsStore.kill(id)
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
</script>

<template>
  <div class="row">
    <q-table dark flat bordered binary-state-sort :rows="brands" :columns="columns" row-key="id" virtual-scroll :rows-per-page-options="[0]" v-model:pagination="pagination" class="grid sticky-header">
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">{{ col.label }}</q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props" :style="editRowData != null && editRowData.id === props.row.id ? 'background-color: #ffffff12' : ''">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <div v-if="col.name === 'label'">{{ props.row.label }}</div>
            <div v-else-if="col.name === 'active'">{{ props.row.active ? 'Yes' : 'No' }}</div>
            <div v-else-if="col.name === 'store'"><a v-if="props.row.store_id != null" :href="storesStore.findByID(props.row.store_id)?.url" target="_blank" >{{ storesStore.findByID(props.row.store_id)?.label }}</a></div>
            <div v-else-if="col.name === 'actions'">
              <q-btn flat color="primary" icon="mdi-pencil" size="xs"  @click="editRow(props.row.id)" />
              <q-btn flat color="red" icon="mdi-delete" size="xs"  @click="deleteBrand(props.row.id)" />
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <div class="form">
      <q-form v-if="editRowData != null" @submit="saveBrand(); editRowData = null" @reset="resetEdit()">
      <div class="text-h6 q-mb-md">{{ editRowData.id != null ? 'Edit brand' : 'Add new brand' }}</div>
      <div><q-input dark v-model="editRowData.label" label="Brand name" hint="Name of the brand: Creality, eSun, etc" lazy-rules :rules="[ val => val && val.length > 0 || 'Please type something']"/></div>
      <div><q-toggle label="Active" v-model="editRowData.active" /></div>
      <div><q-select label="Store" clearable dark v-model="editRowData.store_id" :options="stores" option-label="label" option-value="id" map-options emit-value /> </div>
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
</style>
