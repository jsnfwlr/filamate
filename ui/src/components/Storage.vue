<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import { useStorageStatsStore } from '../stores/stats'


const storageStatsStore = useStorageStatsStore()

const columns: QTableColumn[] = [
  {
    name: 'label',
    required: true,
    label: 'Label',
    align: 'left',
    field: 'label',
    sortable: true,
    style: 'width: 25%',
  },
  {
    name: 'max',
    required: true,
    label: 'Max',
    align: 'left',
    field: 'max',
    sortable: true,
    style: 'width: 25%',
  },
  {
    name: 'used',
    required: true,
    label: 'Used',
    align: 'left',
    field: 'used',
    sortable: true,
    style: 'width: 25%',
  },
  {
    name: 'free',
    required: true,
    label: 'Free',
    align: 'left',
    field: 'free',
    sortable: true,
    style: 'width: 25%',
  }
]

var storage = ref(storageStatsStore.sorted)

onMounted(async () => {
  await storageStatsStore.find()
  storage.value = storageStatsStore.sorted
})

const editRowData = ref<any | null>(null)

const pagination = ref({
  rowsPerPage: 0
})
</script>

<template>
  <div class="row">
    <q-table dark flat bordered binary-state-sort :rows="storage" :columns="columns" row-key="id" virtual-scroll
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
            <div v-if="col.name === 'max'">{{ props.row.max }}</div>
            <div v-if="col.name === 'used'">{{ props.row.used }}</div>
            <div v-if="col.name === 'free'">{{ props.row.free }}</div>

          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>

<style scoped>
.row {
  display: flex;
  flex-direction: row;
}


.sticky-header {

  /* height or max-height is important */
  & {
    height: calc(100vh - (98px + 3rem));
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
