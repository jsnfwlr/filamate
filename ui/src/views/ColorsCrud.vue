<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, onMounted, computed } from 'vue'

import { useColorsStore } from '../stores/colors'
import type { Color } from '../stores/colors'


const colorsStore = useColorsStore()

const columns: QTableColumn[] = [
  {
    name: 'label',
    required: true,
    label: 'Color',
    align: 'left',
    field: 'label',
    sortable: true,
    style: 'width: 35%',
  },
  {
    name: 'hex',
    required: true,
    label: 'Hex',
    align: 'left',
    field: 'hex',
    sortable: true,
    style: 'width: 20 %',
  },
  {
    name: 'alias',
    label: 'Alias',
    align: 'left',
    field: 'alias',
    sortable: true,
    style: 'width: 35%',
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

const colors = ref(colorsStore.sorted)

onMounted(() => {

  colorsStore.find().then(() => {
    colors.value = colorsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load color data: " + error.message)
  })
})

const editRowData = ref<Color>({} as Color)

function editRow(id: number) {
  editRowData.value = colorsStore.findByID(id)
}

function saveColor() {
  if (editRowData.value.id === undefined || editRowData.value.id === null) {
    colorsStore.create(editRowData.value).then(() => {
      resetEdit()
    })
  } else {
    colorsStore.update(editRowData.value.id, editRowData.value).then(() => {
      resetEdit()
    })
  }
}

function deleteColor(id: number) {
  if (confirm("Are you sure you want to delete this color?")) {
    colorsStore.kill(id)
    if (editRowData.value.id === id) {
      resetEdit()
    }

  }
}

function resetEdit() {
  editRowData.value = {} as Color
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
    <q-table dark flat bordered binary-state-sort :rows="colors" :columns="columns" row-key="id" virtual-scroll
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
            <div v-if="col.name === 'hex'">
              <div class="hex" :style="{ backgroundColor: props.row.hex }"></div> {{ props.row.hex }}
            </div>
            <div v-if="col.name === 'alias'">{{ props.row.alias }}</div>
            <div v-else-if="col.name === 'actions'">
              <q-btn flat color="primary" icon="mdi-pencil" size="xs" @click="editRow(props.row.id)" />
              <q-btn flat color="red" icon="mdi-delete" size="xs" @click="deleteColor(props.row.id)" />
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <div class="form">
      <q-form v-if="editRowData != null" @submit="saveColor()" @reset="resetEdit()">
        <div class="text-h6 q-mb-md">{{ editRowData.id != null ? 'Edit color' : 'Add new color' }}</div>
        <div><q-input dark v-model="editRowData.label" label="Color name"
            hint="Name of the color: Black, Teal, Orange, etc" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']" /></div>
        <div>
          <q-input dark v-model="editRowData.alias" label="Color alias" hint="Other names for the color" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']"/>
        </div>
        <div>
          <q-input dark v-model="editRowData.hex" label="Color hex" hint="Hex code for the color" lazy-rules
            :rules="[val => val && val.length > 0 || 'Please type something']">
            <template v-slot:append>
              <q-icon name="colorize" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-color dark no-header-tabs v-model="editRowData.hex" />
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>

        </div>
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

.hex {
  width: 16px;
  height: 16px;
  display: inline-block;
  margin-right: 8px;
  border: 1px solid #ffffff77;
}
</style>
