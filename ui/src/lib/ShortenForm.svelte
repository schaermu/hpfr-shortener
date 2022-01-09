<script lang="ts">
    import { createForm } from 'svelte-forms-lib'
    import * as yup from 'yup'

    let result = null

    const { form, errors, touched, isValid, isModified, isSubmitting, handleChange, handleSubmit } = createForm({
        initialValues: {
            url: ''
        },
        validationSchema: yup.object().shape({
            url: yup.string().required().url()
        }),
        onSubmit: async values => {
            return fetch('/api/shorten', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: values.url })
            }).then((res) => {
                res.json().then(json => result = json)
            })
        }
    })
</script>

<style>    
.shorten-frm {
    display: grid;
    grid-template-columns: 70% 30%;
    gap: 0 2em;
}
</style>

<div class="shorten-frm">
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
    <div>
        hpfr.ch is an URL shortener service. As of Jan 2022, it is completely open source and <a href="https://github.com/schaermu/hpfr-shortener" target="_blank">work in progress</a>.
    </div>
</div>
