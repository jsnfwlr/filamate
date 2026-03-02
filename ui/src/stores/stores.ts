import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {mande} from "mande";

const multiStoreAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/stores') : mande('/api/stores')
const singleStoreAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/store') : mande('/api/store')

export interface Store {
  id: number | null
  label: string
  url: string
}



export const useStoresStore = defineStore('stores', () => {

    const sorted = ref<Store[]>([])
    // const editError = ref<string[]>([])

    async function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

        await multiStoreAPI.get<Array<Store>>({ query: { "sort_by": sortBy } }).then((results: Store[]) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    async function update(id: number, record: Store) {
        await singleStoreAPI.patch<Store>(id, record).then((result: Store) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        })

    }

    async function create(record: Store) {
        multiStoreAPI.post<Store>(record).then( (result: Store) => {
            sorted.value.push( result)
        })
    }

    async function kill(storeID: number) {
        singleStoreAPI.delete(storeID).then(resp => {
            const idx = indexOfID(storeID)
            sorted.value.splice(idx, 1)
            // }).catch(err => {
            //     addNewStudentErrors.value = err.body
        })
    }

    const count = computed( () => {
        return sorted.value.length
    })

    function findByID(id: number): Store | undefined {
        const found = sorted.value.find((f: { id: number | null }) => f.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                label: "",
                url: "",
            }
        }
        return {
            id: found.id,
            label: found.label,
            url: found.url,
        }
    }

    function indexOfID(id: number): number {
        return sorted.value.findIndex((b: { id: number | null }) => b.id === id)
    }

    return {
        sorted,
        count,
        find,
        update,
        create,
        kill,

        findByID,
        indexOfID,
    }

})
