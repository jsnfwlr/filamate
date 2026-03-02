import { defineStore } from 'pinia'
import { ref } from 'vue'
import { mande } from "mande";


const multiStatAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/stats') : mande('/api/stats')
const multiChartAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/chart') : mande('/api/chart')

export interface UsageStat {
    id: number
    color: string
    material: string
    used: number
    ordered: number
}

export interface RatingStat {
    id: number
    brand: string
    material: string
    rating_count: number
    rating_average: number
    details: RatingStatDetails[]
}

export interface RatingStatDetails {
    spool_id: number
    spool_created_at: string
    spool_weight: number
    rating: number
    rating_created_at: string
}
export interface StorageStat {
    id: number
    label: string
    max: number
    used: number
    free: number
    details: StorageStatDetails[]
}

export interface StorageStatDetails {
    material: string
    brand: string
    weight: number
    colors_hex: string[]
    colors_label: string[]
}

export interface StorageChart {
    labels: string[]
    used: (number | [number, number] | null)[]
    purchased: (number | [number, number] | null)[]
    stored: (number | [number, number] | null)[]
}

export interface MaterialChartDatasets {
    labels: string[]
    datasets: dataset[]
}

export interface dataset {
    backgroundColor: string[]
    data: number[]
}

export const useUsageStatsStore = defineStore('usageStats', () => {

    const sorted = ref<UsageStat[]>([])
    // const editError = ref<string[]>([])

    async function find() {
        await multiStatAPI.get<Array<UsageStat>>("usage").then((results: UsageStat[]) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }



    return {
        sorted,
        find,
    }

})

export const useStorageStatsStore = defineStore('storageStats', () => {

    const sorted = ref<StorageStat[]>([])
    // const editError = ref<string[]>([])

    async function find() {
        await multiStatAPI.get<Array<StorageStat>>("storage").then((results: StorageStat[]) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }



    return {
        sorted,
        find,
    }

})

export const useRatingStatsStore = defineStore('ratingStats', () => {

    const sorted = ref<RatingStat[]>([])
    // const editError = ref<string[]>([])

    async function find() {
        await multiStatAPI.get<Array<RatingStat>>("ratings").then((results: RatingStat[]) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    return {
        sorted,
        find,
    }

})

export const useStorageChartStore = defineStore('storageChart', () => {

    const sorted = ref<StorageChart>({} as StorageChart)

    async function find() {
        await multiChartAPI.get<StorageChart>("storage").then((results: StorageChart) => {
            sorted.value = results

        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    return {
        sorted,
        find,
    }

})

export const useMaterialChartStore = defineStore('materialChart', () => {

    const sorted = ref<MaterialChartDatasets>({} as MaterialChartDatasets)

    async function find() {
        await multiChartAPI.get<MaterialChartDatasets>("material").then((results: MaterialChartDatasets) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    return {
        sorted,
        find,
    }

})
