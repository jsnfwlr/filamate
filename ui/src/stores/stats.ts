import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { mande } from "mande";
import type { Spool } from './spools';
import type { Color } from './colors';


const multiStatAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/stats') : mande('/api/stats')
const multiChartAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/chart') : mande('/api/chart')

export interface UsageStat {
    color: string
    material: string
    used: number
    ordered: number
}

export interface StorageStat {
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


export const useUsageStatsStore = defineStore('usageStats', () => {

    const sorted = ref<UsageStat[]>([])
    // const editError = ref<string[]>([])

    async function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

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

    async function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

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

export const useStorageChartStore = defineStore('storageChart', () => {

    const sorted = ref<StorageChart>({} as StorageChart)

    async function find(sortBy: string = "label", sortDir: string = "asc") {
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
