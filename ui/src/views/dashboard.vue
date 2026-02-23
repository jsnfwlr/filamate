<script setup lang="ts">
import { ref } from 'vue'

import Usage from '../components/Usage.vue'
import Storage from '../components/Storage.vue'

</script>

<template>
  <div class="row">
    <div class="col-3 q-mr-md">
      <div class="text-h6 q-mb-md">Usage</div>
      <Usage />
    </div>
    <div class="col-3 q-ml-md">
      <div class="text-h6 q-mb-md">Storage</div>
      <Storage />
    </div>
  </div>
</template>

<style></style>
