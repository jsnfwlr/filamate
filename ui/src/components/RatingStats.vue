<script setup lang="ts">
import type { QTableColumn } from 'quasar'
import { ref, onMounted } from 'vue'

import type { RatingStat } from '../stores/stats'
import { useRatingStatsStore } from '../stores/stats'
import { useColorsStore } from '../stores/colors'

const ratingStatsStore = useRatingStatsStore()
const colorsStore = useColorsStore()


const colors = ref(colorsStore.sorted)
const rating = ref(ratingStatsStore.sorted)

const pagination = ref({
  rowsPerPage: 0
})

const rootColumns: QTableColumn[] = [
  {
    name: 'brand',
    required: true,
    label: 'Brand',
    align: 'left',
    field: 'brand',
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
    name: 'rating_count',
    required: true,
    label: '# of Ratings',
    align: 'left',
    field: 'rating_count',
    sortable: true,
    style: 'width: 25%',
  },
  {
    name: 'rating_average',
    required: true,
    label: 'Average Rating',
    align: 'left',
    field: 'rating_average',
    sortable: true,
    style: 'width: 25%',
    format: (val) => val.toFixed(1) + ' out of 5'
  }
]


const detailsColumns: QTableColumn[] = [
  {
    name: 'spool_colors',
    required: true,
    label: 'Color',
    align: 'left',
    field: 'spool_colors',
    sortable: false,
    style: 'width: 18%',
  },
  {
    name: 'spool_weight',
    required: true,
    label: 'Weight',
    align: 'left',
    field: 'spool_weight',
    sortable: false,
    style: 'width: 15%',
  },
  {
    name: 'rating',
    required: true,
    label: 'Rating',
    align: 'left',
    field: 'rating',
    sortable: false,
    style: 'width: 15%',
    format: (val) => val.toFixed(1) + ' out of 5'
  },
  {
    name: 'spool_created_at',
    required: true,
    label: 'Purchased At',
    align: 'left',
    field: 'spool_created_at',
    sortable: false,
    style: 'width: 26%',
  },
  {
    name: 'rating_created_at',
    required: true,
    label: 'Rated At',
    align: 'left',
    field: 'rating_created_at',
    sortable: false,
    style: 'width: 26%',
  }
]

onMounted(async () => {
  await ratingStatsStore.find()
  rating.value = ratingStatsStore.sorted

  await colorsStore.find()
  colors.value = colorsStore.sorted
})

function DateTime(date: string): string {
  const d = new Date(date)
  return d.toISOString().replace('T', ' ').split('.')[0] as string
}

export interface RowClass {
  id: number
  class: string
}

const rowClasses = ref<Array<RowClass>>([])

function rowClassFn(row: RatingStat): string {
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
    <q-table dark flat bordered binary-state-sort :rows="rating" :columns="rootColumns" row-key="id" virtual-scroll :rows-per-page-options="[0]" v-model:pagination="pagination" class="grid sticky-header" :table-row-class-fn="rowClassFn" @update:pagination="resetRowClasses">
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">{{ col.label }}</q-th>
          <q-th auto-width />
        </q-tr>
      </template>
      <template v-slot:body="props">
        <q-tr :props="props" @click="props.expand = !props.expand" :key="props.row.id + '-main'">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">{{ col.value }}</q-td>
          <q-td auto-width><q-icon @click="props.expand = !props.expand" size="xs" :name="props.expand ? 'mdi-chevron-up' : 'mdi-chevron-down'" /></q-td>
        </q-tr>
        <q-tr v-show="props.expand" :props="props" :key="props.row.id + '-expand'">
          <q-td colspan="100%" class="details">
            <q-table dark flat bordered binary-state-sort :rows="props.row.details" :columns="detailsColumns" row-key="label" virtual-scroll :rows-per-page-options="[0]" v-model:pagination="pagination" class="subgrid grid sticky-header">
               <template v-slot:header="props">
                <q-tr :props="props">
                  <q-th v-for="col in props.cols" :key="col.name" :props="props">{{ col.label }}</q-th>
                </q-tr>
              </template>
              <template v-slot:body="props">
                <q-tr :props="props">
                  <q-td v-for="col in props.cols" :key="col.name" :props="props">
                    <div v-if="col.name === 'spool_colors'">
                      <div v-for="color in props.row.spool_colors" :key="color" class="q-mr-sm">
                        <div class="hex" :style="{ backgroundColor: colorsStore.findByID(color)?.hex }"></div> {{ colorsStore.findByID(color)?.label }}
                      </div>
                    </div>
                    <div v-else-if="col.name === 'spool_created_at'">
                      {{ DateTime(col.value) }}
                    </div>
                    <div v-else-if="col.name === 'rating_created_at'">
                      {{ DateTime(col.value) }}
                    </div>
                    <div v-else>{{ col.value }}</div>
                  </q-td>
                </q-tr>
              </template>
            </q-table>
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

td.details {
  padding-left: 5px !important;
  padding-right: 0 !important;
}
.subgrid {
  width: 100%;
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

  tbody tr {
    cursor: pointer;
    padding: 7px 16px;
  }
}

.hex {
  width: 16px;
  height: 16px;
  display: inline-block;
  margin-right: 8px;
  border: 1px solid #ffffff77;
}

.even {
  background-color: #44444411;
}
</style>
