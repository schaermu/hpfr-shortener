import { render, screen, waitFor } from '@testing-library/svelte'
import userEvent from '@testing-library/user-event'
import ApiClient from 'src/lib/api';

jest.mock('src/lib/api')
jest.mock('src/lib/utils')

import ShortenForm from './ShortenForm.svelte'
import { MOCK_SHORT_URL } from 'src/lib/__mocks__/api';

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
            expect(screen.getByText(MOCK_SHORT_URL)).toBeInTheDocument()
        })
    })

    test('it renders backend errors on the page', async () => {
        ApiClient.prototype.shortenUrl = jest.fn().mockRejectedValueOnce('invalid input data');

        const testValue = 'http://www.foobar.org'
        render(ShortenForm);

        const inputField = screen.getByRole('textbox', { name: /url/i })
        const button = screen.getByRole('button', {name: /shorten/i})

        await userEvent.type(inputField, testValue, { delay: 10 })
        button.click()

        await waitFor(() => {
            expect(screen.getByText('invalid input data')).toBeInTheDocument()
        })
    })
})