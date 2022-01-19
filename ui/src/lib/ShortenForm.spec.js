import { fireEvent, render, screen } from '@testing-library/svelte'
import ShortenForm from './ShortenForm.svelte'

describe('ShortenForm', () => {
    test('it has an input field for the url', () => {
        render(ShortenForm);
        expect(() => screen.getByRole('textbox', { name: /url/i })).not.toThrow()
    })

    test('it does not accept invalid input', async () => {
        const testValue = 'INVALID_URL'
        render(ShortenForm);

        const inputField = screen.getByRole('textbox', { name: /url/i })
        await fireEvent.change(inputField, { target: { value: testValue } });
        
        expect(inputField).toHaveValue(testValue)
        expect(screen.getByText('Shorten')).toHaveAttribute('disabled')
    })
    
    test('it does accept valid input', async () => {
        const testValue = 'http://www.foobar.org'
        render(ShortenForm);

        const inputField = screen.getByRole('textbox', { name: /url/i })
        await fireEvent.change(inputField, { target: { value: testValue } });
        
        expect(inputField).toHaveValue(testValue)
        expect(screen.getByRole('button', {name: /shorten/i})).toHaveAttribute('disabled', '')
    })
})