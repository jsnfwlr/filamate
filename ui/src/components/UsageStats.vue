<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import { useUsageStatsStore } from '../stores/stats'


const usageStatsStore = useUsageStatsStore()

const columns: QTableColumn[] = [
  {
    name: 'color',
    required: true,
    label: 'Color',
    align: 'left',
    field: 'color',
    sortable: true,
    style: 'width: 25%',
  },
  {
    name: 'material',
    required: true,
    label: 'Material',
    align: 'left',
    field: 'material',
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
    name: 'ordered',
    required: true,
    label: 'Purchased',
    align: 'left',
    field: 'ordered',
    sortable: true,
    style: 'width: 25%',
  }
]

const usage = ref(usageStatsStore.sorted)

onMounted(async () => {
  await usageStatsStore.find()
  usage.value = usageStatsStore.sorted
})

const editRowData = ref<any | null>(null)

const pagination = ref({
  rowsPerPage: 0
})



export interface RowClass {
  id: number
  class: string
}

const rowClasses = ref<Array<RowClass>>([])

function rowClassFn(row: any): string {
  for (let i = 0; i < rowClasses.value.length; i++) {
    if (rowClasses.value !== undefined && rowClasses.value !== undefined && rowClasses.value[i]?.id === row.id) {
      return rowClasses.value[i]?.class as string
    }
  }

  if (rowClasses.value.length % 2 === 0) {
    rowClasses.value.push({ id: row.id, class: 'odd' })
    return 'odd' // length of 1 means index
  }

  rowClasses.value.push({ id: row.id, class: 'even' })
  return 'even'
}

function resetRowClasses() {
  rowClasses.value = []
}

</script>

<template>
  <div class="row">
    <q-table dark flat bordered binary-state-sort :rows="usage" :columns="columns" row-key="id" virtual-scroll :rows-per-page-options="[0]" v-model:pagination="pagination" class="grid sticky-header" :table-row-class-fn="rowClassFn" @update:pagination="resetRowClasses">
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">{{ col.label }}</q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props"
          :style="editRowData != null && editRowData.id === props.row.id ? 'background-color: #ffffff12' : ''">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <div v-if="col.name === 'color'">{{ props.row.color }}</div>
            <div v-if="col.name === 'material'">{{ props.row.material }}</div>
            <div v-if="col.name === 'used'">{{ props.row.used }}</div>
            <div v-if="col.name === 'ordered'">{{ props.row.ordered }}</div>

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
   width: 100%;
    max-height: calc(100vh - (98px + 3rem));
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

.even {
  background-color: #44444411;
}
</style>
