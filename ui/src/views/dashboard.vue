<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'


import UsageStats from '../components/UsageStats.vue'
import StorageStats from '../components/StorageStats.vue'
import StorageChart from '../components/StorageChart.vue'
import MaterialChart from '../components/MaterialChart.vue'

import { useStorageChartStore } from '../stores/stats'
import { useMaterialChartStore } from '../stores/stats'
import type { ChartData } from 'chart.js'
import type { MaterialChartDatasets } from '../stores/stats'

const storageChartStore = useStorageChartStore()
const materialChartStore = useMaterialChartStore()

const storageChartUsed = ref<(number | [number, number] | null)[]>([])
const storageChartPurchased = ref<(number | [number, number] | null)[]>([])
const storageChartStored = ref<(number | [number, number] | null)[]>([])
const storageChartLabels = ref<Array<string>>([])
const materialChartData = ref<ChartData<'pie'>>({
  labels: [],
  datasets: []
})

/*
const materialChartResp = ref<Array<MaterialChartResponse>>([
  {
    label: 'ABS',
    color: '#FF0000',
    value: 2,
    children: [
      {
        label: 'ABS',
        color: '#FF0000',
        value: 2,
        children: [
          {
            label: 'Sunlu',
            color: '#FF0000',
            value: 1,
            children: null
          }
        ]
      }
    ]
  },
  {
    label: 'Nylon',
    color: '#00FF00',
    value: 1,
    children: [
      {
        label: 'PA6CF',
        color: '#00FF00',
        value: 1,
        children: [
          {
            label: 'Sunlu',
            color: '#00FF00',
            value: 1,
            children: null
          }
        ]

      }
    ]
  },
  {
    label: 'PETG',
    color: '#0000FF',
    value: 7,
    children: [
      {
        label: 'Basic PETG',
        color: '#0000FF',
        value: 5,
        children: [
          {
            label: 'Anycubic',
            color: '#0000DD',
            value: 1,
            children: null
          },
          {
            label: 'Slic3d',
            color: '#0000BB',
            value: 1,
            children: null
          }
        ]
      },
      {
        label: 'Hyper PETG',
        color: '#0033FF',
        value: 2,
        children: [
          {
            label: 'Creality',
            color: '#0033DD',
            value: 2,
            children: null
          }
        ]
      }
    ]
  },
  {
    label: 'PLA',
    color: '#FFFF00',
    value: 48,
    children: [
      {
        label: 'Basic PLA',
        color: '#FFFF00',
        value: 20,
        children: [
          {
            label: 'Anycubic',
            color: '#FFFF00',
            value: 12,
            children: null
          },
          {
            label: 'CCDIY',
            color: '#FFFF00',
            value: 1,
            children: null
          },
          {
            label: 'eSun',
            color: '#FFFF00',
            value: 5,
            children: null
          },
          {
            label: 'Slic3d',
            color: '#FFFF00',
            value: 1,
            children: null
          }
        ]
      },
      {
        label: 'HS PLA',
        color: '#FFFF00',
        value: 1,
        children: [
          {
            label: 'CCDIY',
            color: '#FFFF00',
            value: 1,
            children: null
          }
        ]
      },
      {
        label: 'Hyper PLA',
        color: '#FFFF00',
        value: 6,
        children: [
          {
            label: 'Creality',
            color: '#FFFF00',
            value: 6,
            children: null
          }
        ]
      },
      {
        label: 'Matte PLA',
        color: '#FFFF00',
        value: 1,
        children: [
          {
            label: 'Siddament',
            color: '#FFFF00',
            value: 1,
            children: null
          }
        ]
      },
      {
        label: 'PLA+',
        color: '#FFFF00',
        value: 18,
        children: [
          {
            label: '3DFillies',
            color: '#FFFF00',
            value: 5,
            children: null
          },
          {
            label: 'eSun',
            color: '#FFFF00',
            value: 4,
            children: null
          },
          {
            label: 'Sunlu',
            color: '#FFFF00',
            value: 9,
            children: null
          }
        ]
      },
      {
        label: 'Silk PLA',
        color: '#FFFF00',
        value: 1,
        children: [
          {
            label: 'Siddament',
            color: '#FFFF00',
            value: 1,
            children: null
          }
        ]
      },
      {
        label: 'Wood PLA',
        color: '#FFFF00',
        value: 1,
        children: [
          {
            label: 'Sunlu',
            color: '#FFFF00',
            value: 1,
            children: null
          }
        ]
      },
      {
        label: 'TPU',
        color: '#FF00FF',
        value: 1,
        children: [
          {
            label: 'TPU 95A',
            color: '#FF00FF',
            value: 1,
            children: [
              {
                label: 'Elegoo',
                color: '#FF00FF',
                value: 1,
                children: null
              }
            ]
          }
        ]
      }
    ]
  }
])

const materialChartData = computed<ChartData<'pie'>>(() => {
  const data: ChartData<'pie'> = {
    labels: [
      '1', '2', '3',
      '1.1', '1.2', '2.1', '3.1', '3.2',
      '1.1.1', '1.1.2', '1.2.1', '2.1.1', '3.1.1', '3.2.1', '3.3.1', '3.3.2'
    ],
    datasets: [
      {
        backgroundColor: ['#FF7800', '#00FD00', '#0077FF'],
        data: [11, 5, 85]
      },
      {
        backgroundColor: ['#FF3300', '#FFFE00', '#00FE00', '#0022FF', '#00AAFF'],
        data: [7, 4, 5, 40, 45]
      },
      {
        backgroundColor: ['#FF0000', '#FF7700', '#FFFF00', '#00FF00', '#0000FF', '#0044FF', '#0088FF', '#00CCFF'],
        data: [3, 4, 4, 5, 18, 22, 40, 5]
      }

    ]
  }

  return data
})
*/
onMounted(async () => {
  await storageChartStore.find()
  storageChartLabels.value = storageChartStore.sorted.labels
  storageChartPurchased.value = storageChartStore.sorted.purchased
  storageChartStored.value = storageChartStore.sorted.stored
  storageChartUsed.value = storageChartStore.sorted.used

  await materialChartStore.find()
  materialChartData.value = materialChartStore.sorted
})


</script>

<template>
  <div class="row q-mb-lg">
    <div class="col-8 q-pr-lg">
      <div class="text-h6 q-mb-md">
        Stock Levels Over Time
        <q-icon name="info" left size="xs" color="#111" id="StorageChart" />
        <q-tooltip target="#StorageChart" anchor="center right" self="center left" class="bg-black">
          Your filament purchases and usage over the last 12 months.
        </q-tooltip>
      </div>
      <StorageChart :used="storageChartUsed" :purchased="storageChartPurchased" :labels="storageChartLabels"
        :stored="storageChartStored" />
    </div>
    <div class="col-4">
      <div class="text-h6 q-mb-md">
        Spools by Material/Brand
        <q-icon name="info" left size="xs" color="#111" id="StorageChart" />
        <q-tooltip target="#StorageChart" anchor="center right" self="center left" class="bg-black">
          Your filament purchases and usage over the last 12 months.
        </q-tooltip>
      </div>
      <MaterialChart :chartData="materialChartData" />
    </div>
  </div>
  <div class="row">
    <div class="col-4 q-pr-lg">
      <div class="text-h6 q-mb-md">
        Top Filament by Material/Color
        <q-icon name="info" left size="xs" color="#111" id="UsageStats" />
        <q-tooltip target="#UsageStats" anchor="center right" self="center left" class="bg-black">
          Your most ordered & used filaments,<br />
          separated by Material and Color.
        </q-tooltip>
      </div>
      <UsageStats />
    </div>
    <div class="col-4 q-pr-lg">
      <div class="text-h6 q-mb-md">
        Storage Capacity
        <q-icon name="info" left size="xs" color="#111" id="StorageStats" />
        <q-tooltip target="#StorageStats" anchor="center right" self="center left" class="bg-black">
          An overview of your storage capacity and filament locations.<br />
          Click on a row for more details.
        </q-tooltip>
      </div>
      <StorageStats />
    </div>
  </div>
</template>

<style></style>
