import fetchMock from 'jest-fetch-mock'
import ApiClient from './api';

describe('ApiClient', () => {
    beforeEach(() => {
        fetchMock.resetMocks();
    })

    test('it does return a fetch promise', async () => {
        const mockRes = { short_url: 'https://hpfr.ch/Gk9Fj9' }
        fetchMock.mockResponseOnce(JSON.stringify(mockRes))

        const res = await new ApiClient().shortenUrl('http://foobar.org')
        const json = await res.json()

        expect(json).toEqual(mockRes)
    })
})