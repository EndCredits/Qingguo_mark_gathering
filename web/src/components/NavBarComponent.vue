<script setup lang="ts">
import { ref } from "vue";

const mainItems = ref([
  { itemName: "Gathering", itemURL: "/" },
  { itemName: "About", itemURL: "/about" },
]);
</script>

<template>
  <div class="md:flex md:flex-row md:box-content md:p-6 md:m-3 m-3 border-4 rounded-xl shadow-xl items-center flex flex-col">
    <div class="md:flex md:flex-row md:justify-center">
      <div class="flex">
        <img src="../assets/hebei-university-of-genshin-impact-black-no-blank.png" class="img_vue flex md:my-auto my-4 h-10 md:h-12" />
      </div>
    </div>
    <div class="md:flex md:flex-row md:inset-y-0 md:ml-auto md:mr-4 flex flex-row my-2">
      <div class="md:mx-4 md:my-auto mx-2 my-4" v-for="item in mainItems" :key="item.itemName">
        <router-link
          class="border-2 p-2 border-gray-300 rounded-2xl drop-shadow-md hover:drop-shadow-2xl"
          :to="{ path: item.itemURL }"
          >{{ item.itemName }}</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.img_hebust {
  height: 48px;
}

.img_hebust:hover {
  filter: drop-shadow(0 0 0.75rem #083f8e);
}
</style>