<script setup lang="ts">
import type { RatingStat } from '../stores/stats'
import type { Color } from '../stores/colors'
import type { ChartData } from 'chart.js'
import type { UsageStat } from '../stores/stats'

import { ref, onMounted, computed } from 'vue'

import UsageStats from '../components/UsageStats.vue'
import StorageStats from '../components/StorageStats.vue'
import RatingStats from '../components/RatingStats.vue'
import StorageChart from '../components/StorageChart.vue'
import MaterialChart from '../components/MaterialChart.vue'

import { useStorageStatsStore } from '../stores/stats'
import { useBrandsStore } from '../stores/brands'
import { useMaterialsStore } from '../stores/materials'
import { useUsageStatsStore } from '../stores/stats'


import { useStorageChartStore, useMaterialChartStore, useRatingStatsStore } from '../stores/stats'
import { useColorsStore } from '../stores/colors'

const usageStatsStore = useUsageStatsStore()
const storageChartStore = useStorageChartStore()
const materialChartStore = useMaterialChartStore()
const ratingStatsStore = useRatingStatsStore()

const colorsStore = useColorsStore()
const storageStatsStore = useStorageStatsStore()
const brandsStore = useBrandsStore()
const materialsStore = useMaterialsStore()

const usage = ref<UsageStat[]>([])
const storage = ref(storageStatsStore.sorted)
const brands = ref(brandsStore.sorted)
const materials = ref(materialsStore.sorted)
const colors = ref<Color[]>([])
const ratings = ref<RatingStat[]>([])
const storageChartEmptied = ref<(number | [number, number] | null)[]>([])
const storageChartAdded = ref<(number | [number, number] | null)[]>([])
const storageChartStored = ref<(number | [number, number] | null)[]>([])
const storageChartLabels = ref<string[]>([])
const materialChartData = ref<ChartData<'pie'>>({
  labels: [],
  datasets: []
})

onMounted(() => {
  storageChartStore.find().then(() => {
    storageChartLabels.value = storageChartStore.sorted.labels
    storageChartAdded.value = storageChartStore.sorted.added
    storageChartStored.value = storageChartStore.sorted.stored
    storageChartEmptied.value = storageChartStore.sorted.emptied
  }).catch((error) => {
    errors.value.push("Failed to load storage chart data: " + error.message)
  })

  colorsStore.find().then(() => {
    colors.value = colorsStore.sorted
  }).catch((error) => {
    errors.value.push("Failed to load color data: " + error.message)
  })

  materialChartStore.find().then(() => {
    materialChartData.value = materialChartStore.sorted
  }).catch((error) => {
    errors.value.push("Failed to load material chart data: " + error.message)
  })

  ratingStatsStore.find().then(() => {
    ratings.value = ratingStatsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load rating chart data: " + error.message)
  })

  storageStatsStore.find().then(() => {
    storage.value = storageStatsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load storage table data: " + error.message)
  })

  brandsStore.find().then(() => {
    brands.value = brandsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load brands data: " + error.message)
  })

  materialsStore.find().then(() => {
    materials.value = materialsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load materials data: " + error.message)
  })

  usageStatsStore.find().then(() => {
    usage.value = usageStatsStore.sorted
  }).catch((error) => {
   errors.value.push("Failed to load usage table data: " + error.message)
  })
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
      <UsageStats :usage="usage" />
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
      <StorageStats :storage="storage" :brands="brands" :materials="materials" />
    </div>
    <div class="col-4 q-pr-lg">
      <div class="text-h6 q-mb-md">
        Quality Ratings by Material/Brand
        <q-icon name="info" left size="xs" color="#111" id="RatingStats" />
        <q-tooltip target="#RatingStats" anchor="center right" self="center left" class="bg-black">
          An overview of your quality ratings by material and brand.
        </q-tooltip>
      </div>
      <RatingStats :ratings="ratings" :colors="colors" />
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

<style></style>
