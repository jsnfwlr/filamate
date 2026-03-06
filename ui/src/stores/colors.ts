import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {mande} from "mande";

const multiColorAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/colors') : mande('/api/colors')
const singleColorAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/color') : mande('/api/color')

export interface Color {
  id: number | null
  label: string
  hex: string
  alias: string | null
}



export const useColorsStore = defineStore('colors', () => {

    const sorted = ref<Color[]>([])
    // const editError = ref<string[]>([])

    function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

        return multiColorAPI.get<Array<Color>>({ query: { "sort_by": sortBy } }).then((results: Color[]) => {
            sorted.value = results
        })
    }

    function update(id: number, record: Color) {
        return singleColorAPI.patch<Color>(id, record).then((result: Color) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        })

    }

    function create(record: Color) {
        return multiColorAPI.post<Color>(record).then( (result: Color) => {
            sorted.value.push(result)
        })
    }

    function kill(colorID: number) {
        return singleColorAPI.delete(colorID).then(() => {
            const idx = indexOfID(colorID)
            sorted.value.splice(idx, 1)
            // }).catch(err => {
            //     addNewStudentErrors.value = err.body
        })
    }

    const count = computed( () => {
        return sorted.value.length
    })

    function findByID(id: number): Color {
        const found = sorted.value.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                label: "",
                hex: "",
                alias: null,
            }
        }
        return {
            id: found.id,
            label: found.label,
            hex: found.hex,
            alias: found.alias,
        }
    }

    function filterByID(colors: Color[], id: number): Color {
        const found = colors.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                label: "",
                hex: "",
                alias: null,
            }
        }
        return {
            id: found.id,
            label: found.label,
            hex: found.hex,
            alias: found.alias,
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
        filterByID,
        indexOfID,
    }

})
