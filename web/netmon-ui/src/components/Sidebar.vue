<template>
    <div class="container">
      <div class="sidebar">
        <div class="sidebar-container">
            <div class="sidebar-nav" @click="goToPage('/')">Dashboard</div>
            <div class="sidebar-nav" @click="goToPage('/packets')">Packets</div>
            <div class="sidebar-nav" @click="goToPage('/hosts')">Hosts</div>
        </div>
      </div>
      <div class="content">
        <RouterView />
      </div>
    </div>
</template>

<script>
import { RouterView } from 'vue-router'

export default {
  name: "Sidebar",
  methods: {
    goToPage(page) {
       this.$router.push(page);
    },
  },
};
</script>
  
<style scoped>
  .container {
    display: flex;
    /* border: 2px solid red; */
  }
  
  .sidebar {
    width: 200px;
    /* height: 100%; */
    min-height: 100vh;
    background-color: #f0f0f0;
    padding: 2rem;
    /* border: 2px solid red; */
  }

  .sidebar-container {
    overflow-y: auto;
    position: fixed;
    /* top: 0;
    left: 0; */
  }

  .sidebar-nav {
    margin-bottom: 0.5rem;
    cursor: pointer;
  }

  .sidebar-nav:hover {
    border-bottom: 2px solid gray;
  }
  
  .content {
    flex: 1;
    padding: 1rem;
    /* border: 2px solid red; */
    overflow-y: auto;
  }
</style>