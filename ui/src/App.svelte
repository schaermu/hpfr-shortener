<script lang="ts">
  import ShortenForm from './lib/ShortenForm.svelte'
</script>

<style>
.frm {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 2em 2em;
}

@media screen and (max-width: 768px) {
    .frm {
        grid-template-columns: 1fr;
    }
}

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
  <div class="frm">
    <div>
        hpfr.ch is an URL shortener service. It is completely open source, <a href="https://github.com/schaermu/hpfr-shortener" target="_blank">work in progress</a> and comes without any guarantees 😇
    </div>
    <div>
      <ShortenForm />
    </div>
  </div>
</main>
<footer class="container">
  <small>made with ❤️ by schaermu &dash; &copy; 2022</small>
</footer>