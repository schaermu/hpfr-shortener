<script lang="ts">
    import { createForm } from 'svelte-forms-lib'
    import * as yup from 'yup'

    let result = null

    const { form, errors, touched, isValid, isModified, isSubmitting, handleChange, handleSubmit, handleReset } = createForm({
        initialValues: {
            url: ''
        },
        validationSchema: yup.object().shape({
            url: yup.string().required().url()
        }),
        onSubmit: async values => {
            try {
                const res = await fetch('/api/shorten', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ url: values.url })
                }).then(res => res.json())
                
                result = res.short_url
                handleReset()
            } catch (err) {
                console.error(err)
            }
        }
    })
</script>

<style>    
.shorten-frm {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 2em 2em;
}

@media screen and (max-width: 768px) {
    .shorten-frm {
        grid-template-columns: 1fr;
    }
}

</style>

<div class="shorten-frm">
    <div>
        hpfr.ch is an URL shortener service. As of Jan 2022, it is completely open source and <a href="https://github.com/schaermu/hpfr-shortener" target="_blank">work in progress</a>.
    </div>
    <div>
        <form on:submit={handleSubmit}>
            <input id="url" name="url" placeholder="URL to shorten"
                on:keyup={handleChange}
                bind:value={$form.url}
                aria-invalid={!$touched.url ? null : $errors.url ? true : false}>
            <button type="submit" disabled={!$isValid || !$isModified} aria-busy={$isSubmitting}>Shorten</button>
            <pre>{result}</pre>
        </form>
    </div>
</div>
