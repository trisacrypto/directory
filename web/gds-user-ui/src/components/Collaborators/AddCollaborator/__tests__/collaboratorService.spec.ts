/* eslint-disable @typescript-eslint/no-unused-vars */
import { getAllCollaborators } from 'modules/dashboard/collaborator/CollaboratorService';
import { collaboratorMockValue } from '../__mocks__';
import axios from 'axios';
import mockedAxios from 'jest-mock-axios';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
const mockAxios = axios as jest.Mocked<typeof axios>;
jest.mock('axios', () => {
    return {

        create: () => {
            return {
                get: jest.fn(),
                post: jest.fn(),
                put: jest.fn(),
                delete: jest.fn(),
                interceptors: {
                    request: { eject: jest.fn(), use: jest.fn() },
                    response: { eject: jest.fn(), use: jest.fn() },
                },
                defaults: {
                    withCredentials: true,
                },
            };
        },

    };
});
describe('CollaboratorService', () => {
    it('should not be called if the service is mocked out', async () => {
        const { data } = collaboratorMockValue;
        axios.get = jest.fn().mockResolvedValue({ data });
        mockedAxios.get.mockReturnValue(collaboratorMockValue);

        await getAllCollaborators();
        // expect(response).toBe(collaboratorMockValue);
        // await expect(getAllCollaborators()).resolves.toEqual(collaboratorMockValue.data);
        // expect(axios.get).toHaveBeenCalledWith('/collaborators');
        expect(mockedAxios.get).toHaveBeenCalledTimes(0);
        // expect(result).toEqual(collaboratorMockValue);
    });
});
