<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, onMounted, computed } from 'vue'


import type { Location } from '../stores/locations'
import { useLocationsStore } from '../stores/locations'


const locationsStore = useLocationsStore()

const columns: QTableColumn[] = [
  {
    name: 'label',
    required: true,
    label: 'Location',
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
    name: 'printable',
    required: true,
    label: 'Printable',
    align: 'left',
    field: 'printable',
    sortable: true,
    style: 'width: 10%',
  },
  {
    name: 'tally',
    required: true,
    label: 'Tally',
    align: 'left',
    field: 'tally',
    sortable: true,
    style: 'width: 10%',
  },

  {
    name: 'capacity',
    required: true,
    label: 'Capacity',
    align: 'left',
    field: 'capacity',
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

const locations = ref(locationsStore.sorted)

onMounted(() => {
  locationsStore.find().then(() => {
    locations.value = locationsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load location data: " + error.message)
  })
})

const editRowData = ref<Location>({} as Location)

function editRow(id: number) {
  editRowData.value = locationsStore.findByID(id)
}

function saveLocation() {
  if (editRowData.value.id === undefined || editRowData.value.id === null) {
    locationsStore.create(editRowData.value).then(() => {
      resetEdit()
    })
  } else {
    locationsStore.update(editRowData.value.id, editRowData.value).then(() => {
      resetEdit()
    })
  }
}

function deleteLocation(id: number) {
  if (confirm("Are you sure you want to delete this location?")) {
    locationsStore.kill(id)
    if (editRowData.value.id === id) {
      resetEdit()
    }

  }
}

function resetEdit() {
  editRowData.value = {} as Location
}

const pagination = ref({
  rowsPerPage: 0
})
const showErrors = computed({
  get() {
    return errors.value.length > 0
  },
  set(newValue: boolean) {
    if (!newValue) {
      errors.value = []
    }
  }
})
const errors = ref<string[]>([])
</script>

<template>
  <div class="row">
    <q-table dark flat bordered binary-state-sort :rows="locations" :columns="columns" row-key="id" virtual-scroll
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
            <div v-else-if="col.name === 'printable'">{{ props.row.printable ? 'Yes' : 'No' }}</div>
            <div v-else-if="col.name === 'tally'">{{ props.row.tally ? 'Yes' : 'No' }}</div>
            <div v-else-if="col.name === 'capacity'">{{ props.row.capacity }}</div>
            <div v-else-if="col.name === 'actions'">
              <q-btn flat color="primary" icon="mdi-pencil" size="xs" @click="editRow(props.row.id)" />
              <q-btn flat color="red" icon="mdi-delete" size="xs" @click="deleteLocation(props.row.id)" />
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <div class="form">
      <q-form v-if="editRowData != null" @submit="saveLocation()" @reset="resetEdit()">
        <div class="text-h6 q-mb-md">{{ editRowData.id != null ? 'Edit location' : 'Add new location' }}</div>
        <div><q-input dark v-model="editRowData.label" label="Label"
            hint="Name of the location: Box 1, AMS, etc" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']" /></div>
        <div><q-input dark v-model="editRowData.description" label="Description"
            hint="Description of the location" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']" /></div>
        <div><q-toggle label="Printable" v-model="editRowData.printable" /></div>
        <div><q-toggle label="Tally" v-model="editRowData.tally" /></div>
        <div><q-input dark v-model="editRowData.capacity" label="Capacity" hint="Maximum number of items the location can hold, enter 0 for unlimited" type="number" min="0" lazy-rules
            :rules="[val => val != null && val >= 0 || 'Please enter a non-negative number']" /></div>
        <div class="q-mt-md">
          <q-btn label="Save" icon="mdi-content-save" type="submit" color="primary" />
          <q-btn label="Cancel" icon="mdi-undo" type="reset" color="secondary" class="q-ml-sm" />
        </div>
      </q-form>

    </div>
  </div>
  <q-dialog v-model="showErrors" dark>
    <q-card dark>
      <q-card-section>
        <div class="text-h6">Error(s) </div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <div v-for="error in errors" :key="error">{{ error }}</div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="OK" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
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
