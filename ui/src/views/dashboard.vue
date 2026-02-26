<script setup lang="ts">
import { ref, onMounted } from 'vue'


import UsageStats from '../components/UsageStats.vue'
import StorageStats from '../components/StorageStats.vue'
import StorageChart from '../components/StorageChart.vue'

import { useStorageChartStore } from '../stores/stats'


const storageChartStore = useStorageChartStore()

// const used =      ref([0, 0,  0,  0,  0,  0,  0,  0,  0, -9, -9, -3])
// const purchased = ref([4, 4,  4,  4,  3,  4,  4,  4,  4,  4,  7,  8])
// const stored =    ref([5, 9, 13, 17, 20, 24, 28, 32, 36, 40, 35, 33])
// const labels = ref(['March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December', 'January', 'February'])

const used =      ref<(number | [number, number] | null)[]>([])
const purchased = ref<(number | [number, number] | null)[]>([])
const stored =    ref<(number | [number, number] | null)[]>([])
const labels =    ref<Array<string>>([])

onMounted(async () => {
  await storageChartStore.find()
  labels.value = storageChartStore.sorted.labels
  purchased.value = storageChartStore.sorted.purchased
  stored.value = storageChartStore.sorted.stored
  used.value = storageChartStore.sorted.used
})


</script>

<template>
  <div class="row">
    <div class="col-3 q-pr-lg">
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
    <div class="col-3 q-pr-lg">
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
    <div class="col-6">
      <div class="text-h6 q-mb-md">
        Stock Levels Over Time
        <q-icon name="info" left size="xs" color="#111" id="StorageChart" />
        <q-tooltip target="#StorageChart" anchor="center right" self="center left" class="bg-black">
          Your filament purchases and usage over the last 12 months.
        </q-tooltip>
      </div>
      <StorageChart :used="used" :purchased="purchased" :labels="labels" :stored="stored" />
    </div>
  </div>

</template>

<style></style>
