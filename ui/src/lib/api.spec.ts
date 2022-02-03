import fetchMock from 'jest-fetch-mock'
import ApiClient from './api';

describe('ApiClient', () => {
    beforeEach(() => {
        fetchMock.resetMocks();
    })

    test('it does return a json object', async () => {
        const mockRes = { short_url: 'https://hpfr.ch/Gk9Fj9' }
        fetchMock.mockResponseOnce(JSON.stringify(mockRes))

        const res = await new ApiClient().shortenUrl('http://foobar.org')

        expect(res).toEqual(mockRes)
    })

    test('it does return an error on network failure', async () => {
        const mockRes = new Error('request timeout')
        fetchMock.mockRejectedValue(mockRes)

        await expect(new ApiClient().shortenUrl('http://foobar.org')).rejects.toEqual(mockRes.message)
    })

    test('it does return an error object on certain error status codes', async () => {
        const mockRes = { message: 'validation error' }
        fetchMock.mockResponseOnce(JSON.stringify(mockRes), { status: 400 });

        expect.assertions(1)
        await expect(new ApiClient().shortenUrl('http://foobar.org')).rejects.toEqual(mockRes.message)
    })
})