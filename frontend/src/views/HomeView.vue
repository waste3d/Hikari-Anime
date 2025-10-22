<script setup>
import { ref } from 'vue';
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8081/api/v1';

const searchQuery = ref('');
const results = ref([]);
const isLoading = ref(false);

async function performSearch() {
  if (!searchQuery.value.trim()) return;

  isLoading.value = true;
  results.value = [];

  try {
    const [movieRes, tvRes] = await Promise.all([
      axios.get(`${API_BASE_URL}/movies/search`, { params: { query: searchQuery.value, language: 'ru-RU' } }),
      axios.get(`${API_BASE_URL}/tv/search`, { params: { query: searchQuery.value, language: 'ru-RU' } })
    ]);
    results.value = [...movieRes.data.results, ...tvRes.data.results];
    console.log('Полученные результаты:', results.value);
  } catch (error) {
    console.error("Ошибка при поиске:", error);
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <main class="container">
    <h1 class="title">Hikari</h1>
    <p class="subtitle">Совместный просмотр нового поколения</p>

    <div class="search-box">
      <input 
        type="text" 
        class="search-input" 
        placeholder="Что ищем? Например, 'Дюна'..."
        v-model="searchQuery"
        @keyup.enter="performSearch" 
      />
      <button class="search-button" @click="performSearch">Найти</button>
    </div>

    <div v-if="isLoading" class="loader">
      Загрузка...
    </div>

    <div v-if="results.length > 0 && !isLoading" class="results-grid">
      <div v-for="item in results" :key="item.id" class="result-card">
        <img :src="item.poster_path" :alt="item.title" class="poster" loading="lazy"/>
        <h3 class="item-title" :title="item.title">{{ item.title }}</h3>
      </div>
    </div>
  </main>
</template>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
  text-align: center;
}

.title {
  font-size: 4rem;
  font-weight: 700;
  color: #c8a2c8; /* Светло-фиолетовый */
  margin-bottom: 0;
}

.subtitle {
  font-size: 1.2rem;
  color: #a9a9a9;
  margin-top: 5px;
  margin-bottom: 40px;
}

.search-box {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 50px;
}

.search-input {
  width: 50%;
  padding: 15px 20px;
  font-size: 1.1rem;
  border-radius: 30px;
  border: 2px solid #333;
  background-color: #1a1a1a;
  color: #fff;
  outline: none;
  transition: border-color 0.3s;
}

.search-input:focus {
  border-color: #6e48d1; /* Фиолетовый */
}

.search-button {
  padding: 15px 30px;
  font-size: 1.1rem;
  font-weight: bold;
  color: #fff;
  background: linear-gradient(45deg, #6e48d1, #486ed1); /* Градиент от фиолетового к синему */
  border: none;
  border-radius: 30px;
  cursor: pointer;
  transition: transform 0.2s;
}

.search-button:hover {
  transform: scale(1.05);
}

.loader {
  margin-top: 50px;
  font-size: 1.5rem;
  color: #c8a2c8;
}

.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  text-align: left;
}

.result-card {
  background-color: #181818;
  border-radius: 10px;
  overflow: hidden;
  transition: transform 0.2s;
  cursor: pointer;
}

.result-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(110, 72, 209, 0.2);
}

.poster {
  width: 100%;
  height: 300px;
  object-fit: cover;
  display: block;
  background-color: #222; /* Плейсхолдер, пока картинка грузится */
}

.item-title {
  padding: 15px 10px;
  font-size: 1rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>