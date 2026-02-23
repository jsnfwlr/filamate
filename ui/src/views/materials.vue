<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import { useMaterialsStore } from '../stores/materials'


const materialsStore = useMaterialsStore()

const columns: QTableColumn[] = [
  {
    name: 'label',
    required: true,
    label: 'Material',
    align: 'left',
    field: 'label',
    sortable: true,
    style: 'width: 20%',
  },
  {
    name: 'description',
    required: true,
    label: 'Description',
    align: 'left',
    field: 'description',
    sortable: true,
    style: 'width: 40%',
  },
  {
    name: 'class',
    required: true,
    label: 'Class',
    align: 'left',
    field: 'class',
    sortable: true,
    style: 'width: 20%',
  },
  {
    name: 'special',
    required: true,
    label: 'Special',
    align: 'left',
    field: 'special',
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

var materials = ref(materialsStore.sorted)

onMounted(async () => {
  await materialsStore.find()
  materials.value = materialsStore.sorted
})

const editRowData = ref<any>({})

function editRow(id: number) {
  editRowData.value = materialsStore.findByID(id)
}

function saveMaterial() {
  if (editRowData.value.id === undefined || editRowData.value.id === null) {
    materialsStore.create(editRowData.value).then(() => {
      resetEdit()
    })
  } else {
    materialsStore.update(editRowData.value.id, editRowData.value).then(() => {
      resetEdit()
    })
  }
}

function deleteMaterial(id: number) {
  if (confirm("Are you sure you want to delete this material?")) {
    materialsStore.kill(id)
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
    <q-table dark flat bordered binary-state-sort :rows="materials" :columns="columns" row-key="id" virtual-scroll
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
            <div v-if="col.name === 'label'">{{ props.row.label }}</div>
            <div v-if="col.name === 'class'">{{ props.row.class }}</div>
            <div v-if="col.name === 'description'">{{ props.row.description }}</div>
            <div v-else-if="col.name === 'special'">{{ props.row.special ? 'Yes' : 'No' }}</div>
            <div v-else-if="col.name === 'actions'">
              <q-btn flat color="primary" icon="mdi-pencil" size="xs" @click="editRow(props.row.id)" />
              <q-btn flat color="red" icon="mdi-delete" size="xs" @click="deleteMaterial(props.row.id)" />
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <div class="form">
      <q-form v-if="editRowData != null" @submit="saveMaterial(); editRowData = null" @reset="resetEdit()">
        <div class="text-h6 q-mb-md">{{ editRowData.id != null ? 'Edit material' : 'Add new material' }}</div>
        <div><q-input dark v-model="editRowData.label" label="Material name"
            hint="Name of the material: PLA+, PLA-HS, ABS, etc" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']" /></div>
        <div><q-input dark v-model="editRowData.class" label="Material class"
            hint="Class of the material: PLA, PETG, etc" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']" /></div>
        <div><q-input dark v-model="editRowData.description" label="Material description"
            hint="Description of the material" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']" /></div>
        <div><q-toggle label="Special" v-model="editRowData.special" /></div>
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
