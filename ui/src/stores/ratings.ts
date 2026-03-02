import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {mande} from "mande";

const multiRatingAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/ratings') : mande('/api/ratings')
const singleRatingAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/rating') : mande('/api/rating')

export interface Rating {
    id: number
    spool_id: number
    created_at: Date
    updated_at: Date
    rating: number
}

export interface NewRating {
    spool_id: number
    rating: number
}



export const useRatingStore = defineStore('ratings', () => {

    const sorted = ref<Rating[]>([])
    // const editError = ref<string[]>([])

    async function find(brandID: number | null = null, materialID: number | null = null) {
        if (brandID !== null && materialID !== null) {

            await multiRatingAPI.get<Array<Rating>>({ query: { "material_id": materialID, "brand_id": brandID } }).then((results: Rating[]) => {
                sorted.value = results
            }).catch(err => {
                alert("find: " + err)
            })
        } else {
            await multiRatingAPI.get<Array<Rating>>().then((results: Rating[]) => {
                sorted.value = results
            }).catch(err => {
                alert("find: " + err)
            })
        }

        return sorted.value
    }

    async function update(id: number, record: Rating) {
        await singleRatingAPI.patch<Rating>(id, record).then((result: Rating) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        })

    }

    async function create(record: NewRating) {
        multiRatingAPI.post<Rating>(record).then( (result: Rating) => {
            sorted.value.push( result)
        })
    }

    async function kill(ratingID: number) {
        singleRatingAPI.delete(ratingID).then(() => {
            const idx = indexOfID(ratingID)
            sorted.value.splice(idx, 1)
            // }).catch(err => {
            //     addNewStudentErrors.value = err.body
        })
    }

    const count = computed( () => {
        return sorted.value.length
    })

    function findByID(id: number): Rating | undefined {
        const found = sorted.value.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: 0,
                rating: 0,
                spool_id: 0,
                created_at: new Date(),
                updated_at: new Date(),
            }
        }
        return {
            id: found.id,
            rating: found.rating,
            spool_id: found.spool_id,
            created_at: found.created_at,
            updated_at: found.updated_at,
        }
    }
    function findBySpoolID(id: number): Rating {
        const found = sorted.value.find((b: { spool_id: number | null }) => b.spool_id === id)
        if (found === null || found === undefined) {
            return {
                id: 0,
                rating: 0,
                spool_id: 0,
                created_at: new Date(),
                updated_at: new Date(),
            }
        }
        return {
            id: found.id,
            rating: found.rating,
            spool_id: found.spool_id,
            created_at: found.created_at,
            updated_at: found.updated_at,
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
        findBySpoolID,
        indexOfID,
    }

})
