<script lang="ts">
import ShortenForm from 'src/components/ShortenForm.svelte'
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
</style>
  
<div class="frm">
    <div>
        hpfr.ch is an URL shortener service. It is completely open source, <a href="https://github.com/schaermu/hpfr-shortener" target="_blank">work in progress</a> and comes without any guarantees 😇
    </div>
    <div>
        <ShortenForm />
    </div>
</div>