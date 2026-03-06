<script setup lang="ts">
import { ref, onMounted } from 'vue'


import UsageStats from '../components/UsageStats.vue'
import StorageStats from '../components/StorageStats.vue'
import RatingStats from '../components/RatingStats.vue'
import StorageChart from '../components/StorageChart.vue'
import MaterialChart from '../components/MaterialChart.vue'

import { useStorageChartStore } from '../stores/stats'
import { useMaterialChartStore } from '../stores/stats'
import type { ChartData } from 'chart.js'

const storageChartStore = useStorageChartStore()
const materialChartStore = useMaterialChartStore()

const storageChartEmptied = ref<(number | [number, number] | null)[]>([])
const storageChartAdded = ref<(number | [number, number] | null)[]>([])
const storageChartStored = ref<(number | [number, number] | null)[]>([])
const storageChartLabels = ref<Array<string>>([])
const materialChartData = ref<ChartData<'pie'>>({
  labels: [],
  datasets: []
})

onMounted(async () => {
  await storageChartStore.find()
  storageChartLabels.value = storageChartStore.sorted.labels
  storageChartAdded.value = storageChartStore.sorted.added
  storageChartStored.value = storageChartStore.sorted.stored
  storageChartEmptied.value = storageChartStore.sorted.emptied

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
      <StorageChart :emptied="storageChartEmptied" :added="storageChartAdded" :labels="storageChartLabels"
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
    <div class="col-4 q-pr-lg">
      <div class="text-h6 q-mb-md">
        Quality Ratings by Material/Brand
        <q-icon name="info" left size="xs" color="#111" id="RatingStats" />
        <q-tooltip target="#RatingStats" anchor="center right" self="center left" class="bg-black">
          An overview of your quality ratings by material and brand.
        </q-tooltip>
      </div>
      <RatingStats />
    </div>
  </div>
</template>

<style></style>
