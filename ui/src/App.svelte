<script lang="ts">
  import { Router, Route } from 'svelte-routing';
  import New from './routes/New.svelte';
  import Stats from './routes/Stats.svelte'

  export let url = '';
</script>

<style>
footer {
  margin-top: 5rem;
  text-align: center;
  font-size: 0.9rem;
}
</style>

<header class="container">
  <hgroup>
    <h1>hpfr.ch</h1>
    <h2>url shortener.</h2>
  </hgroup>
</header>
<main class="container">
  <Router url="{url}">
    <Route path="/"><New /></Route>
    <Route path="/:code" component={Stats}></Route>
  </Router>
</main>
<footer class="container">
  <small>made with ❤️ by schaermu &dash; &copy; 2022</small>
</footer>