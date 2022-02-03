import { render, screen, waitFor } from '@testing-library/svelte'
import userEvent from '@testing-library/user-event'
import ApiClient from './api';

jest.mock('./api')
jest.mock('./utils')

import ShortenForm from './ShortenForm.svelte'

describe('ShortenForm', () => {
    const MockedApiClient = jest.mocked(ApiClient, true);
  
    beforeEach(() => {
        jest.clearAllMocks()
    });

    test('it has an input field for the url', () => {
        render(ShortenForm);
        expect(() => screen.getByRole('textbox', { name: /url/i })).not.toThrow()
    })

    test('it does not accept invalid input', async () => {
        const testValue = 'INVALID_URL'
        render(ShortenForm);

        const inputField = screen.getByRole('textbox', { name: /url/i })
        userEvent.type(inputField, testValue);
        
        expect(inputField).toHaveValue(testValue)
        expect(screen.getByRole('button', {name: /shorten/i})).toHaveAttribute('disabled')
    })
    
    test('it does accept valid input', async () => {
        const testValue = 'http://www.foobar.org'
        render(ShortenForm);

        const inputField = screen.getByRole('textbox', { name: /url/i })
        await userEvent.type(inputField, testValue, { delay: 10 });
        
        expect(inputField).toHaveValue(testValue)
        expect(screen.getByRole('button', {name: /shorten/i})).not.toHaveAttribute('disabled')
    })

    test('it sends valid input data to the backend on submit', async () => {
        const testValue = 'http://www.foobar.org'
        render(ShortenForm);

        const inputField = screen.getByRole('textbox', { name: /url/i })
        const button = screen.getByRole('button', {name: /shorten/i})

        await userEvent.type(inputField, testValue, { delay: 10 })
        button.click()


        await waitFor(() => {
            expect(screen.getByText('http://hpfr.ch/foo_bar_code')).toBeInTheDocument()
        })
    })
})