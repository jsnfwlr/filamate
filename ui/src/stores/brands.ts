import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {mande} from "mande";

const multiBrandAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/brands') : mande('/api/brands')
const singleBrandAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/brand') : mande('/api/brand')

export interface Brand {
  id: number | null
  label: string
  active: boolean
  store_id: number | null
}



export const useBrandsStore = defineStore('brands', () => {

    const sorted = ref<Brand[]>([])
    // const editError = ref<string[]>([])

    function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

        return multiBrandAPI.get<Array<Brand>>({ query: { "sort_by": sortBy } }).then((results: Brand[]) => {
            sorted.value = results
        })
    }

    function update(id: number, record: Brand) {
        return singleBrandAPI.patch<Brand>(id, record).then((result: Brand) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        })

    }

    function create(record: Brand) {
        return multiBrandAPI.post<Brand>(record).then( (result: Brand) => {
            sorted.value.push( result)
        })
    }

    function kill(brandID: number) {
        return singleBrandAPI.delete(brandID).then(() => {
            const idx = indexOfID(brandID)
            sorted.value.splice(idx, 1)
        })
    }

    const count = computed( () => {
        return sorted.value.length
    })

    function findByID(id: number): Brand {
        const found = sorted.value.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                label: "",
                active: false,
                store_id: null,
            }
        }
        return {
            id: found.id,
            label: found.label,
            active: found.active,
            store_id: found.store_id,
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
