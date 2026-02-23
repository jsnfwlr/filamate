import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {mande} from "mande";

const multiStatAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/stats') : mande('/api/stats')

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
